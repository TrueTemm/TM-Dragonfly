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

package chunk

import (
	"bytes"
	_ "embed"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type legacyBlockEntry struct {
	Name string `nbt:"name"`
	Meta int16  `nbt:"meta"`
}

var (
	//go:embed legacy_states.nbt
	legacyMappingsData []byte

	legacyMappings = make(map[legacyBlockEntry]blockEntry)
)

func upgradeLegacyEntry(name string, meta int16) (blockEntry, bool) {
	entry, ok := legacyMappings[legacyBlockEntry{Name: name, Meta: meta}]
	if !ok {

		entry, ok = legacyMappings[legacyBlockEntry{Name: name}]
	}
	return entry, ok
}

func init() {
	var entry struct {
		Legacy  legacyBlockEntry `nbt:"legacy"`
		Updated blockEntry       `nbt:"updated"`
	}
	dec := nbt.NewDecoder(bytes.NewBuffer(legacyMappingsData))
	for {
		if err := dec.Decode(&entry); err != nil {
			break
		}
		legacyMappings[entry.Legacy] = entry.Updated
	}
}
