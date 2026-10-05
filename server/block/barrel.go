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

package block

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"strings"
	"sync"
	"time"
)

type Barrel struct {
	solid
	bass

	Facing cube.Face

	Open bool

	CustomName string

	inventory *inventory.Inventory
	viewerMu  *sync.RWMutex
	viewers   map[ContainerViewer]struct{}
}

func NewBarrel() Barrel {
	m := new(sync.RWMutex)
	v := make(map[ContainerViewer]struct{}, 1)
	return Barrel{
		inventory: inventory.New(27, func(slot int, _, item item.Stack) {
			m.RLock()
			defer m.RUnlock()
			for viewer := range v {
				viewer.ViewSlotChange(slot, item)
			}
		}),
		viewerMu: m,
		viewers:  v,
	}
}

func (Barrel) ContainerSize() int { return 27 }

func (b Barrel) Inventory(*world.Tx, cube.Pos) *inventory.Inventory {
	return b.inventory
}

func (b Barrel) WithName(a ...any) world.Item {
	b.CustomName = strings.TrimSuffix(fmt.Sprintln(a...), "\n")
	return b
}

func (b Barrel) open(tx *world.Tx, pos cube.Pos) {
	b.Open = true
	tx.PlaySound(pos.Vec3Centre(), sound.BarrelOpen{})
	tx.SetBlock(pos, b, nil)
}

func (b Barrel) close(tx *world.Tx, pos cube.Pos) {
	b.Open = false
	tx.PlaySound(pos.Vec3Centre(), sound.BarrelClose{})
	tx.SetBlock(pos, b, nil)
}

func (b Barrel) AddViewer(v ContainerViewer, tx *world.Tx, pos cube.Pos) {
	b.viewerMu.Lock()
	defer b.viewerMu.Unlock()
	if len(b.viewers) == 0 {
		b.open(tx, pos)
	}
	b.viewers[v] = struct{}{}
}

func (b Barrel) RemoveViewer(v ContainerViewer, tx *world.Tx, pos cube.Pos) {
	b.viewerMu.Lock()
	defer b.viewerMu.Unlock()
	if len(b.viewers) == 0 {
		return
	}
	delete(b.viewers, v)
	if len(b.viewers) == 0 {
		b.close(tx, pos)
	}
}

func (b Barrel) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return false
}

func (b Barrel) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, b)
	if !used {
		return
	}

	b = NewBarrel()
	b.Facing = calculateFace(user, pos)

	place(tx, pos, b, user, ctx)
	return placed(ctx)
}

func (b Barrel) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, oneOf(b)).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		for _, i := range b.Inventory(tx, pos).Clear() {
			dropItem(tx, i, pos.Vec3())
		}
	})
}

func (b Barrel) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, true)
}

func (Barrel) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (b Barrel) DecodeNBT(data map[string]any) any {
	facing := b.Facing

	b = NewBarrel()
	b.Facing = facing
	b.CustomName = nbtconv.String(data, "CustomName")
	nbtconv.InvFromNBT(b.inventory, nbtconv.Slice(data, "Items"))
	return b
}

func (b Barrel) EncodeNBT() map[string]any {
	if b.inventory == nil {
		facing, customName := b.Facing, b.CustomName

		b = NewBarrel()
		b.Facing, b.CustomName = facing, customName
	}
	m := map[string]any{
		"Items": nbtconv.InvToNBT(b.inventory),
		"id":    "Barrel",
	}
	if b.CustomName != "" {
		m["CustomName"] = b.CustomName
	}
	return m
}

func (b Barrel) EncodeBlock() (string, map[string]any) {
	return "minecraft:barrel", map[string]any{"open_bit": boolByte(b.Open), "facing_direction": int32(b.Facing)}
}

func (b Barrel) EncodeItem() (name string, meta int16) {
	return "minecraft:barrel", 0
}

func allBarrels() (b []world.Block) {
	for i := cube.Face(0); i < 6; i++ {
		b = append(b, Barrel{Facing: i})
		b = append(b, Barrel{Facing: i, Open: true})
	}
	return
}
