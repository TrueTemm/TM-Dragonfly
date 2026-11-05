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
	"github.com/df-mc/worldupgrader/itemupgrader"
)

type legacyItem630 struct {
	id    int32
	meta  uint32
	fixed bool
}

type legacyPalette struct {
	once  sync.Once
	names map[int32]string
	toOld map[int32]legacyItem630
}

var legacyPalettes = map[uint32]*legacyPalette{
	589: {}, 594: {}, 598: {}, 618: {}, 622: {}, 630: {},
}

func (p *legacyPalette) load(proto uint32) {
	p.once.Do(func() {
		p.names = make(map[int32]string)
		p.toOld = make(map[int32]legacyItem630)
		entries := ItemsFor(proto)
		for _, e := range entries {
			id := int32(e.RuntimeID)
			p.names[id] = e.Name
			if native, ok := world.ItemRuntimeIDByName(e.Name); ok {
				p.toOld[native] = legacyItem630{id: id}
			}
		}
		for _, e := range entries {
			for meta := int16(0); meta < 16; meta++ {
				up := itemupgrader.Upgrade(itemupgrader.ItemMeta{Name: e.Name, Meta: meta})
				if up.Name == e.Name {
					continue
				}
				native, ok := world.ItemRuntimeIDByName(up.Name)
				if !ok {
					continue
				}
				if _, exists := p.toOld[native]; !exists {
					p.toOld[native] = legacyItem630{id: int32(e.RuntimeID), meta: uint32(meta), fixed: up.Meta != meta}
				}
			}
		}
	})
}

func TranslateStack(protocolID uint32, id int32, meta uint32) (int32, uint32) {
	p, legacy := legacyPalettes[protocolID]
	if !legacy {
		return Translate(protocolID, id), meta
	}
	if id == 0 {
		return 0, 0
	}
	p.load(protocolID)
	if old, ok := p.toOld[id]; ok {
		if old.fixed {
			meta = old.meta
		}
		return old.id, meta
	}
	return 0, 0
}

func ReverseTranslateStack(protocolID uint32, id int32, meta uint32) (int32, uint32) {
	p, legacy := legacyPalettes[protocolID]
	if !legacy {
		return ReverseTranslate(protocolID, id), meta
	}
	if id == 0 {
		return 0, 0
	}
	p.load(protocolID)
	name, ok := p.names[id]
	if !ok {
		return 0, 0
	}
	up := itemupgrader.Upgrade(itemupgrader.ItemMeta{Name: name, Meta: int16(meta)})
	if native, ok := world.ItemRuntimeIDByName(up.Name); ok {
		if _, supported := p.toOld[native]; supported { // reject future item variants
			return native, uint32(uint16(up.Meta))
		}
	}
	return 0, 0
}
