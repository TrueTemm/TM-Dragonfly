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
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/df-mc/dragonfly/server/world"
)

//go:embed data/states_*.txt.gz
var statesFS embed.FS

//go:embed data/schema/*.json
var schemaFS embed.FS

var schemaFiles = []struct {
	file string
	to   int
}{
	{"0291_1.21.0.25_beta_to_1.21.20.24_beta.json", ver(1, 21, 20)},
	{"0301_1.21.20.24_beta_to_1.21.30.24_beta.json", ver(1, 21, 30)},
	{"0311_1.21.30.24_beta_to_1.21.40.25_beta.json", ver(1, 21, 40)},
	{"0321_1.21.50.29_beta_to_1.21.60.28_beta.json", ver(1, 21, 60)},
	{"0331_1.21.100.23_beta_to_1.21.110.26_beta.json", ver(1, 21, 110)},
	{"0341_1.26.20_to_1.26.30.json", ver(1, 26, 30)},
}

func ver(major, minor, patch int) int { return major*1e6 + minor*1e3 + patch }

var statePalettes = map[uint32]*statePalette{
	671:  {file: 686, version: ver(1, 20, 80)}, // no 1.20 schema, 1.21.0 shapes
	685:  {file: 686, version: ver(1, 21, 0)},
	686:  {file: 686, version: ver(1, 21, 2)},
	712:  {file: 712, version: ver(1, 21, 20)},
	729:  {file: 729, version: ver(1, 21, 30)},
	748:  {file: 748, version: ver(1, 21, 40)},
	766:  {file: 766, version: ver(1, 21, 50)},
	776:  {file: 776, version: ver(1, 21, 60)},
	786:  {file: 786, version: ver(1, 21, 70)},
	800:  {file: 800, version: ver(1, 21, 80)},
	818:  {file: 818, version: ver(1, 21, 90)},
	819:  {file: 819, version: ver(1, 21, 93)},
	827:  {file: 827, version: ver(1, 21, 100)},
	844:  {file: 844, version: ver(1, 21, 111), incomplete: true},
	859:  {file: 844, version: ver(1, 21, 120), incomplete: true},
	898:  {file: 844, version: ver(1, 21, 130), incomplete: true},
	924:  {file: 924, version: ver(1, 26, 0)},
	944:  {file: 944, version: ver(1, 26, 10)},
	975:  {file: 975, version: ver(1, 26, 20)},
	1001: {file: 1001, version: ver(1, 26, 30)},
	2168: {file: 2168, version: ver(1, 26, 40), incomplete: true},
}

type statePalette struct {
	file    uint32
	version int

	incomplete bool

	once     sync.Once
	toOld    []uint32
	toLatest map[uint32]uint32

	unmatchedOld, unmatchedLatest int
	unmatchedLatestNames          map[string]int
	substituted                   map[string]string
}

func HashFor(protocol, rid uint32) uint32 {
	p, ok := statePalettes[protocol]
	if !ok {
		return HashOf(rid)
	}
	p.build()
	if rid < uint32(len(p.toOld)) {
		return p.toOld[rid]
	}
	return HashOf(rid)
}

func RuntimeIDFor(protocol, hash uint32) uint32 {
	if p, ok := statePalettes[protocol]; ok {
		p.build()
		if rid, ok := p.toLatest[hash]; ok {
			return rid
		}
	}
	return RuntimeIDFromHash(hash)
}

func DowngradeStats(protocol uint32) (unmatchedOld int, unmatchedLatest map[string]int, ok bool) {
	p, ok := statePalettes[protocol]
	if !ok {
		return 0, nil, false
	}
	p.build()
	return p.unmatchedOld, p.unmatchedLatestNames, true
}

var (
	schemasOnce sync.Once
	schemas     []*schema

	latestOnce sync.Once
	latestKeys map[string]uint32

	latestFirst    map[string]uint32
	latestPropKeys map[string]map[string]bool
)

func loadSchemas() {
	schemasOnce.Do(func() {
		for _, sf := range schemaFiles {
			data, err := schemaFS.ReadFile("data/schema/" + sf.file)
			if err != nil {
				panic("blockpalette: " + err.Error())
			}
			s := &schema{to: sf.to}
			if err := json.Unmarshal(data, s); err != nil {
				panic("blockpalette: " + sf.file + ": " + err.Error())
			}
			schemas = append(schemas, s)
		}
	})
}

func loadLatest() {
	latestOnce.Do(func() {
		reg := world.DefaultBlockRegistry
		n := uint32(reg.BlockCount())
		latestKeys = make(map[string]uint32, n)
		latestFirst = make(map[string]uint32, 1200)
		latestPropKeys = make(map[string]map[string]bool, 1200)
		for rid := uint32(0); rid < n; rid++ {
			name, props, ok := reg.RuntimeIDToState(rid)
			if !ok {
				continue
			}
			if k := untypedKey(name, props); latestKeys[k] == 0 {
				latestKeys[k] = rid
			}
			if _, seen := latestFirst[name]; !seen {
				latestFirst[name] = rid
				latestPropKeys[name] = make(map[string]bool, len(props))
				for k := range props {
					latestPropKeys[name][k] = true
				}
			}
		}
	})
}

func (p *statePalette) build() {
	p.once.Do(func() {
		loadSchemas()
		loadLatest()
		reg := world.DefaultBlockRegistry
		n := uint32(reg.BlockCount())
		p.toOld = make([]uint32, n)
		for rid := uint32(0); rid < n; rid++ {
			p.toOld[rid] = HashOf(rid)
		}
		p.toLatest = make(map[uint32]uint32, n)
		mapped := make([]bool, n)

		data, err := statesFS.ReadFile(fmt.Sprintf("data/states_%d.txt.gz", p.file))
		if err != nil {
			panic("blockpalette: " + err.Error())
		}
		gz, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			panic("blockpalette: " + err.Error())
		}
		sc := bufio.NewScanner(gz)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		seen := make(map[string]struct{}, 17000)
		old := &oldPalette{byKey: make(map[string]uint32, 17000), first: make(map[string]uint32, 1200),
			props: make(map[string]map[string]bool, 1200)}
		for sc.Scan() {
			line := sc.Text()
			if _, dup := seen[line]; dup || line == "" {
				continue
			}
			seen[line] = struct{}{}
			name, props := parseStateLine(line)
			oldHash := world.NetworkBlockHash(name, props)
			old.byKey[untypedKey(name, props)] = oldHash
			if _, first := old.first[name]; !first {
				old.first[name] = oldHash
				old.props[name] = make(map[string]bool, len(props))
				for k := range props {
					old.props[name][k] = true
				}
			}
			upName, upProps := name, cloneProps(props)
			for _, s := range schemas {
				if p.version < s.to {
					upName, upProps = s.apply(upName, upProps)
				}
			}
			rid, ok := latestKeys[untypedKey(upName, upProps)]
			if !ok {
				p.unmatchedOld++
				continue
			}
			if !mapped[rid] {
				mapped[rid] = true
				p.toOld[rid] = oldHash
			}
			if _, exists := p.toLatest[oldHash]; !exists {
				p.toLatest[oldHash] = rid
			}
		}

		infoUpdate := world.NetworkBlockHash("minecraft:info_update", nil)
		_, hasPlaceholder := p.toLatest[infoUpdate]
		p.unmatchedLatestNames = make(map[string]int)
		p.substituted = make(map[string]string)
		for rid := uint32(0); rid < n; rid++ {
			if mapped[rid] {
				continue
			}
			p.unmatchedLatest++
			name, props, ok := reg.RuntimeIDToState(rid)
			if !ok {
				continue
			}
			if h := p.substituteHash(name, props, mapped, old); h != 0 {
				p.toOld[rid] = h
				p.substituted[name] = substituteName(name)
				continue
			}
			p.unmatchedLatestNames[name]++
			if hasPlaceholder && !p.incomplete {
				p.toOld[rid] = infoUpdate
			}
		}
	})
}

func parseStateLine(line string) (string, map[string]any) {
	parts := strings.Split(line, "\t")
	props := make(map[string]any, len(parts)-1)
	for _, kv := range parts[1:] {
		eq := strings.IndexByte(kv, '=')
		if eq < 0 || eq+1 >= len(kv) {
			continue
		}
		k, tv := kv[:eq], kv[eq+1:]
		switch tv[0] {
		case 's':
			props[k] = tv[1:]
		case 'b':
			v, _ := strconv.Atoi(tv[1:])
			props[k] = uint8(v)
		case 'i':
			v, _ := strconv.Atoi(tv[1:])
			props[k] = int32(v)
		}
	}
	return parts[0], props
}

func cloneProps(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func untypedKey(name string, props map[string]any) string {
	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString(name)
	for _, k := range keys {
		sb.WriteByte('|')
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(valueString(props[k]))
	}
	return sb.String()
}

type schema struct {
	to int

	RenamedIds                  map[string]string                         `json:"renamedIds"`
	AddedProperties             map[string]map[string]propValue           `json:"addedProperties"`
	RemovedProperties           map[string][]string                       `json:"removedProperties"`
	RenamedProperties           map[string]map[string]string              `json:"renamedProperties"`
	RemappedPropertyValues      map[string]map[string]string              `json:"remappedPropertyValues"`
	RemappedPropertyValuesIndex map[string][]struct{ Old, New propValue } `json:"remappedPropertyValuesIndex"`
	FlattenedProperties         map[string]*flattenRule                   `json:"flattenedProperties"`
	RemappedStates              map[string][]remapRule                    `json:"remappedStates"`
}

type propValue struct {
	Int    *int32  `json:"int"`
	String *string `json:"string"`
	Byte   *int32  `json:"byte"`
}

func (v propValue) value() any {
	switch {
	case v.Int != nil:
		return *v.Int
	case v.String != nil:
		return *v.String
	case v.Byte != nil:
		return uint8(*v.Byte)
	}
	return nil
}

type flattenRule struct {
	Prefix                string            `json:"prefix"`
	FlattenedProperty     string            `json:"flattenedProperty"`
	FlattenedPropertyType string            `json:"flattenedPropertyType"`
	Suffix                string            `json:"suffix"`
	FlattenedValueRemaps  map[string]string `json:"flattenedValueRemaps"`
}

type remapRule struct {
	OldState         map[string]propValue `json:"oldState"`
	NewName          string               `json:"newName"`
	NewFlattenedName *flattenRule         `json:"newFlattenedName"`
	NewState         map[string]propValue `json:"newState"`
	CopiedState      []string             `json:"copiedState"`
}

func (s *schema) apply(name string, props map[string]any) (string, map[string]any) {
	if rules, ok := s.RemappedStates[name]; ok {
		for _, r := range rules {
			if len(r.OldState) > len(props) {
				continue
			}
			match := true
			for k, v := range r.OldState {
				if cur, ok := props[k]; !ok || cur != v.value() {
					match = false
					break
				}
			}
			if !match {
				continue
			}
			newName := r.NewName
			if r.NewFlattenedName != nil {
				newName, _ = r.NewFlattenedName.apply(name, cloneProps(props))
			}
			newProps := make(map[string]any, len(r.NewState)+len(r.CopiedState))
			for k, v := range r.NewState {
				newProps[k] = v.value()
			}
			for _, k := range r.CopiedState {
				if v, ok := props[k]; ok {
					newProps[k] = v
				}
			}
			return newName, newProps
		}
	}
	oldName := name
	if renamed, ok := s.RenamedIds[oldName]; ok {
		name = renamed
	} else if fl, ok := s.FlattenedProperties[oldName]; ok {
		name, props = fl.apply(oldName, props)
	}
	for k, v := range s.AddedProperties[oldName] {
		if _, exists := props[k]; !exists {
			props[k] = v.value()
		}
	}
	for _, k := range s.RemovedProperties[oldName] {
		delete(props, k)
	}
	for oldProp, newProp := range s.RenamedProperties[oldName] {
		if v, ok := props[oldProp]; ok {
			delete(props, oldProp)
			props[newProp] = s.remapValue(oldName, oldProp, v)
		}
	}
	for prop := range s.RemappedPropertyValues[oldName] {
		if v, ok := props[prop]; ok {
			props[prop] = s.remapValue(oldName, prop, v)
		}
	}
	return name, props
}

func (s *schema) remapValue(block, prop string, old any) any {
	idx, ok := s.RemappedPropertyValues[block][prop]
	if !ok {
		return old
	}
	for _, pair := range s.RemappedPropertyValuesIndex[idx] {
		if pair.Old.value() == old {
			return pair.New.value()
		}
	}
	return old
}

func (f *flattenRule) apply(name string, props map[string]any) (string, map[string]any) {
	v, ok := props[f.FlattenedProperty]
	if !ok {
		return name, props
	}
	var key string
	switch t := v.(type) {
	case string:
		if f.FlattenedPropertyType != "" && f.FlattenedPropertyType != "string" {
			return name, props
		}
		key = t
	case uint8:
		if f.FlattenedPropertyType != "byte" {
			return name, props
		}
		key = strconv.Itoa(int(t))
	case int32:
		if f.FlattenedPropertyType != "int" {
			return name, props
		}
		key = strconv.Itoa(int(t))
	default:
		return name, props
	}
	embed := key
	if r, ok := f.FlattenedValueRemaps[key]; ok {
		embed = r
	}
	delete(props, f.FlattenedProperty)
	return f.Prefix + embed + f.Suffix, props
}
