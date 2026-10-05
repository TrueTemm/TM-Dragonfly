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

package blockpalette

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/df-mc/dragonfly/server/world"
)

//go:embed data/p786.txt.gz
var data786 []byte

//go:embed data/p800.txt.gz
var data800 []byte

//go:embed data/p818.txt.gz
var data818 []byte

//go:embed data/p819.txt.gz
var data819 []byte

//go:embed data/p827.txt.gz
var data827 []byte

type palette struct {
	keyToRID map[string]uint32
	air      uint32
	once     sync.Once
	raw      []byte

	cacheMu sync.RWMutex
	cache   map[uint32]uint32

	reverseOnce sync.Once
	reverse     map[uint32]uint32
}

var palettes = map[uint32]*palette{

	685: {},
	686: {},
	712: {},
	729: {},
	748: {},
	766: {},
	776: {},
	786: {raw: data786},
	800: {raw: data800},
	818: {raw: data818},
	819: {raw: data819},
	827: {raw: data827},

	844: {},
	859: {},
	898: {},

	924:  {},
	944:  {},
	975:  {},
	1001: {},
}

func init() { palettes[uint32(protocol.CurrentProtocol)] = &palette{} }

func Supported(protocol uint32) bool {
	_, ok := palettes[protocol]
	return ok
}

func (p *palette) load() {
	p.once.Do(func() {
		p.keyToRID = make(map[string]uint32, 16000)
		p.cache = make(map[uint32]uint32, 4096)
		gz, err := gzip.NewReader(bytes.NewReader(p.raw))
		if err != nil {
			return
		}
		defer gz.Close()
		sc := bufio.NewScanner(gz)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			line := sc.Text()

			tab := -1
			for j := 0; j < len(line); j++ {
				if line[j] == '\t' {
					tab = j
					break
				}
			}
			if tab < 0 {
				continue
			}
			key := line[:tab]
			rid, err := strconv.ParseUint(line[tab+1:], 10, 32)
			if err != nil {
				continue
			}
			if _, exists := p.keyToRID[key]; !exists {
				p.keyToRID[key] = uint32(rid)
			}
		}
		if rid, ok := p.keyToRID["air|"]; ok {
			p.air = rid
		}
	})
}

func Translate(protocol, dragonflyRID uint32) uint32 {
	p, ok := palettes[protocol]
	if !ok {
		return dragonflyRID
	}
	p.load()

	p.cacheMu.RLock()
	if rid, hit := p.cache[dragonflyRID]; hit {
		p.cacheMu.RUnlock()
		return rid
	}
	p.cacheMu.RUnlock()

	name, props, found := world.DefaultBlockRegistry.RuntimeIDToState(dragonflyRID)
	out := p.air
	if found {
		if rid, ok := p.keyToRID[canonicalKey(name, props)]; ok {
			out = rid
		}
	}

	p.cacheMu.Lock()
	p.cache[dragonflyRID] = out
	p.cacheMu.Unlock()
	return out
}

func HashOf(dragonflyRID uint32) uint32 {
	if h, ok := world.DefaultBlockRegistry.RuntimeIDToHash(dragonflyRID); ok {
		return h
	}
	if h, ok := world.DefaultBlockRegistry.RuntimeIDToHash(world.DefaultBlockRegistry.AirRuntimeID()); ok {
		return h
	}
	return dragonflyRID
}

func RuntimeIDFromHash(hash uint32) uint32 {
	if rid, ok := world.DefaultBlockRegistry.HashToRuntimeID(hash); ok {
		return rid
	}
	return world.DefaultBlockRegistry.AirRuntimeID()
}

func ReverseTranslate(protocol, versionRID uint32) uint32 {
	p, ok := palettes[protocol]
	if !ok {
		return versionRID
	}
	p.buildReverse()
	if rid, hit := p.reverse[versionRID]; hit {
		return rid
	}
	return world.DefaultBlockRegistry.AirRuntimeID()
}

func (p *palette) buildReverse() {
	p.reverseOnce.Do(func() {
		p.load()
		count := uint32(world.DefaultBlockRegistry.BlockCount())
		p.reverse = make(map[uint32]uint32, count)
		for rid := uint32(0); rid < count; rid++ {
			v := Translate(p.protocolOf(), rid)
			if _, exists := p.reverse[v]; !exists {
				p.reverse[v] = rid
			}
		}
	})
}

func (p *palette) protocolOf() uint32 {
	for proto, pal := range palettes {
		if pal == p {
			return proto
		}
	}
	return 0
}

func canonicalKey(name string, props map[string]any) string {
	name = strings.TrimPrefix(name, "minecraft:")
	if len(props) == 0 {
		return name + "|"
	}
	parts := make([]string, 0, len(props))
	for k, v := range props {
		parts = append(parts, k+"="+valueString(v))
	}
	sort.Strings(parts)

	buf := make([]byte, 0, len(name)+1+len(parts)*8)
	buf = append(buf, name...)
	buf = append(buf, '|')
	for i, part := range parts {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, part...)
	}
	return string(buf)
}

func valueString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "1"
		}
		return "0"
	case uint8:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	default:
		return ""
	}
}
