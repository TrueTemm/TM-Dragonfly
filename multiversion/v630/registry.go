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

package v630

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/df-mc/dragonfly/multiversion/v776"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

//go:embed biome_definitions.nbt
var biomeDefinitions []byte

//go:embed entity_identifiers.nbt
var entityIdentifiers []byte

//go:embed level_sound_id_map.json
var soundData []byte

var sounds map[string]uint32
var soundNames map[uint32]string

func init() {
	if err := json.Unmarshal(soundData, &sounds); err != nil {
		panic(err)
	}
	soundNames = make(map[uint32]string, len(sounds))
	for name, id := range sounds {
		soundNames[id] = name
	}
}

var particles = map[int32]int32{
	0:  0,
	1:  1,
	2:  2,
	3:  3,
	4:  4,
	5:  5,
	6:  6,
	7:  7,
	8:  8,
	9:  10,
	10: 11,
	11: 12,
	12: 13,
	13: 14,
	14: 15,
	15: 16,
	16: 17,
	17: 18,
	18: 19,
	19: 20,
	20: 21,
	21: 22,
	23: 24,
	24: 25,
	25: 26,
	26: 27,
	27: 28,
	28: 29,
	29: 30,
	30: 31,
	31: 32,
	32: 33,
	33: 34,
	34: 35,
	35: 36,
	36: 37,
	37: 38,
	38: 39,
	39: 40,
	40: 41,
	41: 42,
	42: 43,
	43: 44,
	44: 45,
	45: 46,
	46: 47,
	47: 48,
	48: 49,
	49: 50,
	50: 51,
	51: 52,
	52: 53,
	53: 54,
	54: 55,
	55: 56,
	56: 57,
	57: 58,
	58: 59,
	59: 60,
	60: 61,
	61: 62,
	62: 63,
	63: 64,
	64: 65,
	65: 66,
	66: 67,
	67: 68,
	68: 69,
	69: 70,
	71: 72,
	72: 73,
	73: 74,
	74: 75,
	79: 80,
	80: 81,
	81: 9,
	82: 82,
	83: 83,
	84: 84,
	87: 87,
	88: 88,
}

func fromLatestSound(proto uint32, pk *packet.LevelSoundEvent) []packet.Packet {
	id, ok := registryFor(proto).sounds[pk.SoundType]
	if !ok {
		return nil
	}
	extra := pk.ExtraData
	if v786.BlockSoundCarriesRuntimeID(pk.SoundType) && extra != -1 {
		extra = int32(blockpalette.HashFor(proto, uint32(extra)))
	}
	return []packet.Packet{&v776.LevelSoundEvent{SoundType: id, Position: pk.Position, ExtraData: extra, EntityType: pk.EntityType, BabyMob: pk.BabyMob, DisableRelativeVolume: pk.DisableRelativeVolume}}
}

func fromLatestEvent(proto uint32, pk *packet.LevelEvent) []packet.Packet {
	cp := *pk
	if pk.EventType&0x4000 != 0 {
		id, ok := registryFor(proto).particles[pk.EventType&^0x4000]
		if !ok {
			return nil
		}
		cp.EventType = 0x4000 | id
	}
	switch pk.EventType {
	case packet.LevelEventParticlesDestroyBlock, packet.LevelEventParticlesDestroyBlockNoSound:
		cp.EventData = int32(blockpalette.HashFor(proto, uint32(pk.EventData)))
	case packet.LevelEventParticlesCrackBlock:
		cp.EventData = int32(blockpalette.HashFor(proto, uint32(pk.EventData)&0xffffff)) | pk.EventData & ^int32(0xffffff)
	}
	return []packet.Packet{&cp}
}

var versionedVanillaPacks = map[string]bool{
	"d34cfa4b-2ad1-453d-a0db-668b429a3ea0": true,
	"b41c2785-c512-4a49-af56-3a87afd47c57": true,
	"a4df0cb3-17be-4163-88d7-fcf7002b935d": true,
	"d19adffe-a2e1-4b02-8436-ca4583368c89": true,
	"85d5603d-2824-4b21-8044-34f441f4fce1": true,
	"e977cd13-0a11-4618-96fb-03dfe9c43608": true,
	"0674721c-a0aa-41a1-9ba8-1ed33ea3e7ed": true,
}

//go:embed data/*
var oldRegistries embed.FS

var versions = map[uint32]string{589: "1.20.0", 594: "1.20.10", 598: "1.20.20", 618: "1.20.30", 622: "1.20.40", 630: "1.20.50"}

type registry struct {
	biomes, entities []byte
	sounds           map[string]uint32
	soundNames       map[uint32]string
	particles        map[int32]int32
}

var registries = map[uint32]*registry{}

func registryFor(proto uint32) *registry { return registries[proto] }

func init() {
	registries[630] = &registry{biomes: biomeDefinitions, entities: entityIdentifiers, sounds: sounds, soundNames: soundNames, particles: particles}
	for proto := range versions {
		if proto == 630 {
			continue
		}
		read := func(name string) []byte {
			data, err := oldRegistries.ReadFile(fmt.Sprintf("data/%d/%s", proto, name))
			if err != nil {
				panic(err)
			}
			return data
		}
		r := &registry{biomes: read("biome_definitions.nbt"), entities: read("entity_identifiers.nbt"), soundNames: map[uint32]string{}, particles: map[int32]int32{}}
		if err := json.Unmarshal(read("level_sound_id_map.json"), &r.sounds); err != nil {
			panic(err)
		}
		for name, id := range r.sounds {
			r.soundNames[id] = name
		}
		var available map[string]int32
		if err := json.Unmarshal(read("particle_id_map.json"), &available); err != nil {
			panic(err)
		}
		ids := map[int32]bool{}
		for _, id := range available {
			ids[id] = true
		}
		for native, id := range particles {
			if ids[id] {
				r.particles[native] = id
			}
		}
		registries[proto] = r
	}
}
