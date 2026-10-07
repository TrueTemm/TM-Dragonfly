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

package itemdata

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

//go:embed runtime_item_states_686.json
var runtimeStates686 []byte

//go:embed runtime_item_states_671.json
var runtimeStates671 []byte

var (
	items671Once sync.Once
	items671     []protocol.ItemEntry
)

func Items671() []protocol.ItemEntry {
	items671Once.Do(func() { items671 = loadStatesOnly(runtimeStates671) })
	return items671
}

//go:embed runtime_item_states_712.json
var runtimeStates712 []byte

//go:embed runtime_item_states_729.json
var runtimeStates729 []byte

//go:embed runtime_item_states_748.json
var runtimeStates748 []byte

//go:embed runtime_item_states_766.json
var runtimeStates766 []byte

//go:embed runtime_item_states_776.json
var runtimeStates776 []byte

//go:embed item_components_776.nbt
var itemComponents776 []byte

//go:embed runtime_item_states_944.json
var runtimeStates944 []byte

//go:embed item_components_944.nbt
var itemComponents944 []byte

//go:embed runtime_item_states_975.json
var runtimeStates975 []byte

//go:embed item_components_975.nbt
var itemComponents975 []byte

//go:embed runtime_item_states_1001.json
var runtimeStates1001 []byte

//go:embed item_components_1001.nbt
var itemComponents1001 []byte

//go:embed runtime_item_states_924.json
var runtimeStates924 []byte

//go:embed item_components_924.nbt
var itemComponents924 []byte

//go:embed runtime_item_states_786.json
var runtimeStates786 []byte

//go:embed item_components_786.nbt
var itemComponents786 []byte

//go:embed runtime_item_states_800.json
var runtimeStates800 []byte

//go:embed item_components_800.nbt
var itemComponents800 []byte

//go:embed runtime_item_states_818.json
var runtimeStates818 []byte

//go:embed item_components_818.nbt
var itemComponents818 []byte

//go:embed runtime_item_states_819.json
var runtimeStates819 []byte

//go:embed item_components_819.nbt
var itemComponents819 []byte

//go:embed runtime_item_states_827.json
var runtimeStates827 []byte

//go:embed item_components_827.nbt
var itemComponents827 []byte

//go:embed runtime_item_states_844.json
var runtimeStates844 []byte

//go:embed item_components_844.nbt
var itemComponents844 []byte

//go:embed runtime_item_states_859.json
var runtimeStates859 []byte

//go:embed item_components_859.nbt
var itemComponents859 []byte

//go:embed runtime_item_states_898.json
var runtimeStates898 []byte

//go:embed item_components_898.nbt
var itemComponents898 []byte

type runtimeItemState struct {
	Name           string `json:"name"`
	ID             int32  `json:"id"`
	Version        int32  `json:"version"`
	ComponentBased bool   `json:"componentBased"`
}

func load(jsonData, nbtGzData []byte) []protocol.ItemEntry {
	var states []runtimeItemState
	if err := json.Unmarshal(jsonData, &states); err != nil {
		panic(fmt.Errorf("itemdata: decode runtime item states: %w", err))
	}

	gz, err := gzip.NewReader(bytes.NewReader(nbtGzData))
	if err != nil {
		panic(fmt.Errorf("itemdata: gunzip item components: %w", err))
	}
	var components map[string]map[string]any
	if err := nbt.NewDecoderWithEncoding(gz, nbt.BigEndian).Decode(&components); err != nil {
		panic(fmt.Errorf("itemdata: decode item components NBT: %w", err))
	}

	entries := make([]protocol.ItemEntry, 0, len(states))
	for _, s := range states {
		data := components[s.Name]
		if data == nil {
			data = map[string]any{}
		}
		entries = append(entries, protocol.ItemEntry{
			Name:           s.Name,
			RuntimeID:      int16(s.ID),
			ComponentBased: s.ComponentBased,
			Version:        s.Version,
			Data:           data,
		})
	}
	return entries
}

func loadStatesOnly(jsonData []byte) []protocol.ItemEntry {
	var states []runtimeItemState
	if err := json.Unmarshal(jsonData, &states); err != nil {
		panic(fmt.Errorf("itemdata: decode runtime item states: %w", err))
	}
	entries := make([]protocol.ItemEntry, 0, len(states))
	for _, s := range states {
		entries = append(entries, protocol.ItemEntry{Name: s.Name, RuntimeID: int16(s.ID), ComponentBased: s.ComponentBased, Version: s.Version, Data: map[string]any{}})
	}
	return entries
}

var (
	items786Once, items800Once, items818Once, items819Once, items827Once sync.Once
	items786, items800, items818, items819, items827                     []protocol.ItemEntry
	items844Once, items859Once, items898Once                             sync.Once
	items844, items859, items898                                         []protocol.ItemEntry
	items776Once                                                         sync.Once
	items776                                                             []protocol.ItemEntry
	items924Once                                                         sync.Once
	items924                                                             []protocol.ItemEntry
	items944Once                                                         sync.Once
	items944                                                             []protocol.ItemEntry
	items975Once, items1001Once                                          sync.Once
	items975, items1001                                                  []protocol.ItemEntry
	items766Once                                                         sync.Once
	items766                                                             []protocol.ItemEntry
	items748Once                                                         sync.Once
	items748                                                             []protocol.ItemEntry
	items729Once                                                         sync.Once
	items729                                                             []protocol.ItemEntry
	items712Once                                                         sync.Once
	items712                                                             []protocol.ItemEntry
	items686Once                                                         sync.Once
	items686                                                             []protocol.ItemEntry
)

func Items686() []protocol.ItemEntry {
	items686Once.Do(func() { items686 = loadStatesOnly(runtimeStates686) })
	return items686
}

func Items712() []protocol.ItemEntry {
	items712Once.Do(func() { items712 = loadStatesOnly(runtimeStates712) })
	return items712
}

func Items729() []protocol.ItemEntry {
	items729Once.Do(func() { items729 = loadStatesOnly(runtimeStates729) })
	return items729
}

func Items748() []protocol.ItemEntry {
	items748Once.Do(func() { items748 = loadStatesOnly(runtimeStates748) })
	return items748
}

func Items766() []protocol.ItemEntry {
	items766Once.Do(func() { items766 = loadStatesOnly(runtimeStates766) })
	return items766
}

func Items944() []protocol.ItemEntry {
	items944Once.Do(func() { items944 = load(runtimeStates944, itemComponents944) })
	return items944
}

func Items975() []protocol.ItemEntry {
	items975Once.Do(func() { items975 = load(runtimeStates975, itemComponents975) })
	return items975
}

func Items1001() []protocol.ItemEntry {
	items1001Once.Do(func() { items1001 = load(runtimeStates1001, itemComponents1001) })
	return items1001
}

func Items924() []protocol.ItemEntry {
	items924Once.Do(func() { items924 = load(runtimeStates924, itemComponents924) })
	return items924
}

func Items776() []protocol.ItemEntry {
	items776Once.Do(func() { items776 = load(runtimeStates776, itemComponents776) })
	return items776
}

func Items786() []protocol.ItemEntry {
	items786Once.Do(func() { items786 = load(runtimeStates786, itemComponents786) })
	return items786
}

func Items800() []protocol.ItemEntry {
	items800Once.Do(func() { items800 = load(runtimeStates800, itemComponents800) })
	return items800
}

func Items818() []protocol.ItemEntry {
	items818Once.Do(func() { items818 = load(runtimeStates818, itemComponents818) })
	return items818
}

func Items819() []protocol.ItemEntry {
	items819Once.Do(func() { items819 = load(runtimeStates819, itemComponents819) })
	return items819
}

func Items827() []protocol.ItemEntry {
	items827Once.Do(func() { items827 = load(runtimeStates827, itemComponents827) })
	return items827
}

func Items844() []protocol.ItemEntry {
	items844Once.Do(func() { items844 = load(runtimeStates844, itemComponents844) })
	return items844
}

func Items859() []protocol.ItemEntry {
	items859Once.Do(func() { items859 = load(runtimeStates859, itemComponents859) })
	return items859
}

func Items898() []protocol.ItemEntry {
	items898Once.Do(func() { items898 = load(runtimeStates898, itemComponents898) })
	return items898
}
