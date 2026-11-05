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

package v786

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

//go:embed sound_ids.json.gz
var soundRegistryData []byte

type soundRegistry struct {
	Revision string            `json:"revision"`
	IDs      map[string]uint32 `json:"ids"`
	names    map[uint32]string
}

var soundRegistries map[uint32]*soundRegistry

func init() {
	r, err := gzip.NewReader(bytes.NewReader(soundRegistryData))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	if err := json.NewDecoder(r).Decode(&soundRegistries); err != nil {
		panic(err)
	}
	for _, registry := range soundRegistries {
		registry.names = make(map[uint32]string, len(registry.IDs))
		for name, id := range registry.IDs {
			if _, exists := registry.names[id]; exists {
				panic("duplicate sound ID")
			}
			registry.names[id] = name
		}
	}
}

func SoundIDFor(proto uint32, name string) (uint32, bool) {
	if proto == 0 {
		proto = 786
	}
	r := soundRegistries[proto]
	if r == nil {
		return 0, false
	}
	id, ok := r.IDs[name]
	return id, ok
}

func SoundNameFor(proto uint32, id uint32) (string, bool) {
	if proto == 0 {
		proto = 786
	}
	r := soundRegistries[proto]
	if r == nil {
		return "", false
	}
	name, ok := r.names[id]
	return name, ok
}

func LevelSoundPacketsFromLatest(proto uint32, pk *packet.LevelSoundEvent) []packet.Packet {
	out := FromLatestLevelSoundEvent786(proto, pk)
	if out == nil {
		return nil
	}
	return []packet.Packet{out}
}

func LevelSoundPacketsToLatest(proto uint32, pk *LevelSoundEvent) []packet.Packet {
	out := ToLatestLevelSoundEvent786(proto, pk)
	if out == nil {
		return nil
	}
	return []packet.Packet{out}
}

var blockSoundExtraData = map[string]bool{
	"place":            true,
	"hit":              true,
	"item.use.on":      true,
	"door.open":        true,
	"door.close":       true,
	"trapdoor.open":    true,
	"trapdoor.close":   true,
	"fence_gate.open":  true,
	"fence_gate.close": true,
}
