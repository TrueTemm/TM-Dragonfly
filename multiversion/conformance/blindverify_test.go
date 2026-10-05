/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package conformance

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

var finaliseOnce sync.Once

func fill(v reflect.Value, depth int) {
	if depth > 7 || !v.IsValid() {
		return
	}
	if !v.CanSet() {
		v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(1)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(1)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(1.5)
	case reflect.String:
		v.SetString("a")
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			v.SetBytes([]byte{1, 2})
			return
		}
		s := reflect.MakeSlice(v.Type(), 1, 1)
		fill(s.Index(0), depth+1)
		v.Set(s)
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			fill(v.Index(i), depth+1)
		}
	case reflect.Map:
		m := reflect.MakeMap(v.Type())
		k := reflect.New(v.Type().Key()).Elem()
		fill(k, depth+1)
		e := reflect.New(v.Type().Elem()).Elem()
		if e.Kind() == reflect.Interface {
			e.Set(reflect.ValueOf(int32(1)))
		} else {
			fill(e, depth+1)
		}
		m.SetMapIndex(k, e)
		v.Set(m)
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		fill(v.Elem(), depth+1)
	case reflect.Struct:
		if v.Type().String() == "big.Int" || strings.HasSuffix(v.Type().Name(), "Bitset") {
			return
		}
		for i := 0; i < v.NumField(); i++ {
			fill(v.Field(i), depth+1)
		}
	}
}

func verifyProtocols(t *testing.T) []minecraft.Protocol {
	if os.Getenv("TMDRAGONFLY_BLINDVERIFY") == "" {
		t.Skip("blind verification only runs through tools/blindverify/verify.sh")
	}

	finaliseOnce.Do(world.DefaultBlockRegistry.Finalize)
	return protocols()
}

func TestBlindVerifySend(t *testing.T) {
	for _, proto := range verifyProtocols(t) {
		proto := proto
		t.Run(strconv.Itoa(int(proto.ID())), func(t *testing.T) {
			var lines []string
			for _, nf := range packet.NewServerPool() {
				src := nf()
				if os.Getenv("BLINDVERIFY_FILL") != "" {
					fill(reflect.ValueOf(src).Elem(), 0)
				}
				name := reflect.TypeOf(src).Elem().Name()
				status, hx := "ok", ""
				var id uint32
				func() {
					defer func() {
						if e := recover(); e != nil {
							status = fmt.Sprintf("ENCODE_PANIC:%v", e)
						}
					}()
					out := proto.ConvertFromLatest(src, nil)
					if len(out) != 1 {
						status = "skip"
						return
					}
					id = out[0].ID()
					buf := bytes.NewBuffer(nil)
					out[0].Marshal(proto.NewWriter(buf, 0))
					hx = hex.EncodeToString(buf.Bytes())
				}()
				if status == "skip" {
					continue
				}
				lines = append(lines, fmt.Sprintf("%s\t%d\t%s\t%s", name, id, status, hx))
			}
			sort.Strings(lines)
			os.WriteFile(fmt.Sprintf("%s/send_%d.txt", os.Getenv("DUMPDIR"), proto.ID()), []byte(strings.Join(lines, "\n")+"\n"), 0644)
		})
	}
}

func TestBlindVerifyRecv(t *testing.T) {
	for _, proto := range verifyProtocols(t) {
		proto := proto
		t.Run(strconv.Itoa(int(proto.ID())), func(t *testing.T) {
			f, err := os.Open(fmt.Sprintf("%s/recv_%d.txt", os.Getenv("DUMPDIR"), proto.ID()))
			if err != nil {
				t.Skip("no recv file")
			}
			defer f.Close()
			pool := proto.Packets(true)
			sc := bufio.NewScanner(f)
			sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
			for sc.Scan() {
				p := strings.Split(sc.Text(), "\t")
				if len(p) < 4 {
					continue
				}
				id64, _ := strconv.Atoi(p[0])
				id, name, status, hx := uint32(id64), p[1], p[2], p[3]
				if status != "ok" {
					continue
				}
				nf, ok := pool[id]
				if !ok {
					t.Logf("RECV %-34s id=%d NOT IN OUR CLIENT POOL", name, id)
					continue
				}
				b, _ := hex.DecodeString(hx)
				buf := bytes.NewBuffer(b)
				pk := nf()
				ourName := reflect.TypeOf(pk).String()
				func() {
					defer func() {
						if e := recover(); e != nil {
							t.Logf("RECV %-34s DECODE PANIC (ours=%s): %v", name, ourName, e)
						}
					}()
					pk.Marshal(proto.NewReader(buf, 0, false))
					if buf.Len() != 0 {
						t.Logf("RECV %-34s %d TRAILING (ours=%s)", name, buf.Len(), ourName)
					}
					func() {
						defer func() {
							if e := recover(); e != nil {
								t.Logf("RECV %-34s ConvertToLatest PANIC (ours=%s): %v", name, ourName, e)
							}
						}()
						proto.ConvertToLatest(pk, nil)
					}()
				}()
			}
		})
	}
}
