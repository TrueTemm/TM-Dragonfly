// reallib: the historical-gophertunnel half of the verifier. Pin the module to a version, then:
//
//	reallib decsend <send_file>   decode our outgoing bytes with the real pool, report problems
//	reallib encrecv <out_file>    encode zero-value of every real CLIENT packet -> id \t name \t hex
package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// fill populates every field of v with a small non-zero value so that value-dependent branches and nested
// slices are exercised: one element per slice, one entry per map, every Optional present, every bool true.
// Unexported fields (Optional's internals) are reached through unsafe so generics need no special casing.
// Interfaces are left nil — their concrete type cannot be guessed, and the resulting encode panic shows up
// identically for every protocol, so it drops out of the baseline diff.
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

func merged() map[uint32]func() packet.Packet {
	m := map[uint32]func() packet.Packet{}
	for id, f := range packet.NewClientPool() {
		m[id] = f
	}
	for id, f := range packet.NewServerPool() {
		m[id] = f
	}
	return m
}

func main() {
	switch os.Args[1] {
	case "decsend":
		pool := merged()
		f, _ := os.Open(os.Args[2])
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 1<<20), 1<<24)
		for sc.Scan() {
			p := strings.Split(sc.Text(), "\t")
			if len(p) < 4 {
				continue
			}
			name, status, hx := p[0], p[2], p[3]
			if status != "ok" {
				fmt.Printf("  SEND %-34s our encode: %s\n", name, status)
				continue
			}
			var id uint32
			fmt.Sscanf(p[1], "%d", &id)
			nf, ok := pool[id]
			if !ok {
				fmt.Printf("  SEND %-34s id=%d NOT IN REAL POOL\n", name, id)
				continue
			}
			b, _ := hex.DecodeString(hx)
			buf := bytes.NewBuffer(b)
			pk := nf()
			rn := reflect.TypeOf(pk).Elem().Name()
			func() {
				defer func() {
					if e := recover(); e != nil {
						fmt.Printf("  SEND %-34s DECODE PANIC (real=%s): %v\n", name, rn, e)
					}
				}()
				pk.Marshal(protocol.NewReader(buf, 0, false))
				if buf.Len() != 0 {
					fmt.Printf("  SEND %-34s %d TRAILING (real=%s)\n", name, buf.Len(), rn)
				}
				if rn != name {
					fmt.Printf("  SEND %-34s ID MISMATCH: real id=%d is %s\n", name, id, rn)
				}
			}()
		}
	case "encrecv":
		var lines []string
		for id, nf := range packet.NewClientPool() {
			pk := nf()
			if os.Getenv("BLINDVERIFY_FILL") != "" {
				fill(reflect.ValueOf(pk).Elem(), 0)
			}
			name := reflect.TypeOf(pk).Elem().Name()
			buf := bytes.NewBuffer(nil)
			status := "ok"
			func() {
				defer func() {
					if e := recover(); e != nil {
						status = fmt.Sprintf("ENCODE_PANIC:%v", e)
					}
				}()
				pk.Marshal(protocol.NewWriter(buf, 0))
			}()
			lines = append(lines, fmt.Sprintf("%d\t%s\t%s\t%s", id, name, status, hex.EncodeToString(buf.Bytes())))
		}
		sort.Strings(lines)
		os.WriteFile(os.Args[2], []byte(strings.Join(lines, "\n")+"\n"), 0644)
		fmt.Printf("  wrote %d client packets\n", len(lines))
	}
}
