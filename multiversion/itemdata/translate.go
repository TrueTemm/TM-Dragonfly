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
	"sync"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type itemPalette struct {
	entries     func() []protocol.ItemEntry
	once        sync.Once
	dragonToVer map[int32]int32
	verToDragon map[int32]int32
	shieldID    int32
}

var itemPalettes = map[uint32]*itemPalette{
	685:  {entries: Items686},
	686:  {entries: Items686},
	712:  {entries: Items712},
	729:  {entries: Items729},
	748:  {entries: Items748},
	766:  {entries: Items766},
	776:  {entries: Items776},
	786:  {entries: Items786},
	800:  {entries: Items800},
	818:  {entries: Items818},
	819:  {entries: Items819},
	827:  {entries: Items827},
	844:  {entries: Items844},
	859:  {entries: Items859},
	898:  {entries: Items898},
	924:  {entries: Items924},
	944:  {entries: Items944},
	975:  {entries: Items975},
	1001: {entries: Items1001},
}

func (p *itemPalette) load() {
	p.once.Do(func() {
		list := p.entries()
		p.dragonToVer = make(map[int32]int32, len(list))
		p.verToDragon = make(map[int32]int32, len(list))
		for _, e := range list {
			verID := int32(e.RuntimeID)
			if e.Name == "minecraft:shield" {
				p.shieldID = verID
			}
			dragonID, ok := world.ItemRuntimeIDByName(e.Name)
			if !ok {

				continue
			}
			if _, exists := p.dragonToVer[dragonID]; !exists {
				p.dragonToVer[dragonID] = verID
			}
			if _, exists := p.verToDragon[verID]; !exists {
				p.verToDragon[verID] = dragonID
			}
		}
	})
}

func Supported(protocolID uint32) bool {
	_, ok := itemPalettes[protocolID]
	return ok
}

func Translate(protocolID uint32, dragonflyID int32) int32 {
	if dragonflyID == 0 {
		return 0
	}
	p, ok := itemPalettes[protocolID]
	if !ok {
		return dragonflyID
	}
	p.load()
	if verID, hit := p.dragonToVer[dragonflyID]; hit {
		return verID
	}
	return 0
}

func ReverseTranslate(protocolID uint32, versionID int32) int32 {
	if versionID == 0 {
		return 0
	}
	p, ok := itemPalettes[protocolID]
	if !ok {
		return versionID
	}
	p.load()
	if dragonID, hit := p.verToDragon[versionID]; hit {
		return dragonID
	}
	return 0
}

func ShieldID(protocolID uint32) int32 {
	p, ok := itemPalettes[protocolID]
	if !ok {
		return 0
	}
	p.load()
	return p.shieldID
}
