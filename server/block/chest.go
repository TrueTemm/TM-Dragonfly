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

type Chest struct {
	chest
	transparent
	bass
	sourceWaterDisplacer

	Facing cube.Direction

	CustomName string

	paired       bool
	pairX, pairZ int
	pairInv      *inventory.Inventory

	inventory *inventory.Inventory
	viewerMu  *sync.RWMutex
	viewers   map[ContainerViewer]struct{}
}

func NewChest() Chest {
	c := Chest{
		viewerMu: new(sync.RWMutex),
		viewers:  make(map[ContainerViewer]struct{}, 1),
	}

	c.inventory = inventory.New(27, func(slot int, _, after item.Stack) {
		c.viewerMu.RLock()
		defer c.viewerMu.RUnlock()
		for viewer := range c.viewers {
			viewer.ViewSlotChange(slot, after)
		}
	})
	return c
}

func (c Chest) ContainerSize() int {
	if c.paired {
		return 54
	}
	return 27
}

func (c Chest) Inventory(tx *world.Tx, pos cube.Pos) *inventory.Inventory {
	inv, _ := c.tryPair(tx, pos)
	return inv
}

func (c Chest) tryPair(tx *world.Tx, pos cube.Pos) (*inventory.Inventory, bool) {
	if c.paired {
		if c.pairInv == nil {
			if ch, pair, ok := c.pair(tx, pos, c.pairPos(pos)); ok {
				tx.SetBlock(pos, ch, nil)
				tx.SetBlock(c.pairPos(pos), pair, nil)
				return ch.pairInv, true
			}
			c.paired = false
			tx.SetBlock(pos, c, nil)
			return c.inventory, true
		}
		return c.pairInv, false
	}
	return c.inventory, false
}

func (c Chest) WithName(a ...any) world.Item {
	c.CustomName = strings.TrimSuffix(fmt.Sprintln(a...), "\n")
	return c
}

func (Chest) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (c Chest) open(tx *world.Tx, pos cube.Pos) {
	for _, v := range tx.Viewers(pos.Vec3()) {
		if c.paired {
			v.ViewBlockAction(c.pairPos(pos), OpenAction{})
		}
		v.ViewBlockAction(pos, OpenAction{})
	}
	tx.PlaySound(pos.Vec3Centre(), sound.ChestOpen{})
}

func (c Chest) close(tx *world.Tx, pos cube.Pos) {
	for _, v := range tx.Viewers(pos.Vec3()) {
		if c.paired {
			v.ViewBlockAction(c.pairPos(pos), CloseAction{})
		}
		v.ViewBlockAction(pos, CloseAction{})
	}
	tx.PlaySound(pos.Vec3Centre(), sound.ChestClose{})
}

func (c Chest) AddViewer(v ContainerViewer, tx *world.Tx, pos cube.Pos) {
	if _, changed := c.tryPair(tx, pos); changed {
		c = tx.Block(pos).(Chest)
	}
	c.viewerMu.Lock()
	defer c.viewerMu.Unlock()
	if len(c.viewers) == 0 {
		c.open(tx, pos)
	}
	c.viewers[v] = struct{}{}
}

func (c Chest) RemoveViewer(v ContainerViewer, tx *world.Tx, pos cube.Pos) {
	if _, changed := c.tryPair(tx, pos); changed {
		c = tx.Block(pos).(Chest)
	}
	c.viewerMu.Lock()
	defer c.viewerMu.Unlock()
	if len(c.viewers) == 0 {
		return
	}
	delete(c.viewers, v)
	if len(c.viewers) == 0 {
		c.close(tx, pos)
	}
}

func (c Chest) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		if c.paired {
			if d, ok := tx.Block(c.pairPos(pos).Side(cube.FaceUp)).(LightDiffuser); !ok || d.LightDiffusionLevel() > 2 {
				return false
			}
		}
		if d, ok := tx.Block(pos.Side(cube.FaceUp)).(LightDiffuser); ok && d.LightDiffusionLevel() <= 2 {
			opener.OpenBlockContainer(pos, tx)
		}
		return true
	}
	return false
}

func (c Chest) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, c)
	if !used {
		return
	}

	c = NewChest()
	c.Facing = user.Rotation().Direction().Opposite()

	for _, dir := range []cube.Direction{c.Facing.RotateLeft(), c.Facing.RotateRight()} {
		if ch, pair, ok := c.pair(tx, pos, pos.Side(dir.Face())); ok {
			place(tx, pos, ch, user, ctx)
			tx.SetBlock(ch.pairPos(pos), pair, nil)
			return placed(ctx)
		}
	}

	place(tx, pos, c, user, ctx)
	return placed(ctx)
}

func (c Chest) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, oneOf(c)).withBreakHandler(func(pos cube.Pos, tx *world.Tx, u item.User) {
		if c.paired {
			pairPos := c.pairPos(pos)
			if _, pair, ok := c.unpair(tx, pos); ok {
				c.paired = false
				tx.SetBlock(pairPos, pair, nil)
			}
		}

		for _, i := range c.Inventory(tx, pos).Clear() {
			dropItem(tx, i, pos.Vec3Centre())
		}
	})
}

func (Chest) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

func (c Chest) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, true)
}

func (c Chest) Paired() bool {
	return c.paired
}

func (c Chest) pair(tx *world.Tx, pos, pairPos cube.Pos) (ch, pair Chest, ok bool) {
	pair, ok = tx.Block(pairPos).(Chest)
	if !ok || c.Facing != pair.Facing || pair.paired && (pair.pairX != pos[0] || pair.pairZ != pos[2]) {
		return c, pair, false
	}
	m := new(sync.RWMutex)
	v := make(map[ContainerViewer]struct{})
	left, right := c.inventory.Clone(nil), pair.inventory.Clone(nil)
	if pos.Side(c.Facing.RotateRight().Face()) == pairPos {
		left, right = right, left
	}
	double := left.Merge(right, func(slot int, _, item item.Stack) {
		if slot < 27 {
			_ = left.SetItem(slot, item)
		} else {
			_ = right.SetItem(slot-27, item)
		}
		m.RLock()
		defer m.RUnlock()
		for viewer := range v {
			viewer.ViewSlotChange(slot, item)
		}
	})

	c.inventory, pair.inventory = left, right
	if pos.Side(c.Facing.RotateRight().Face()) == pairPos {
		c.inventory, pair.inventory = right, left
	}
	c.pairX, c.pairZ, c.paired = pairPos[0], pairPos[2], true
	pair.pairX, pair.pairZ, pair.paired = pos[0], pos[2], true
	c.viewerMu, pair.viewerMu = m, m
	c.viewers, pair.viewers = v, v
	c.pairInv, pair.pairInv = double, double
	return c, pair, true
}

func (c Chest) unpair(tx *world.Tx, pos cube.Pos) (ch, pair Chest, ok bool) {
	if !c.paired {
		return c, Chest{}, false
	}

	pair, ok = tx.Block(c.pairPos(pos)).(Chest)
	if !ok || c.Facing != pair.Facing || pair.paired && (pair.pairX != pos[0] || pair.pairZ != pos[2]) {
		return c, pair, false
	}

	if len(c.viewers) != 0 {
		c.close(tx, pos)
	}

	c.inventory = c.inventory.Clone(func(slot int, _, after item.Stack) {
		c.viewerMu.RLock()
		defer c.viewerMu.RUnlock()
		for viewer := range c.viewers {
			viewer.ViewSlotChange(slot, after)
		}
	})
	pair.inventory = pair.inventory.Clone(func(slot int, _, after item.Stack) {
		pair.viewerMu.RLock()
		defer pair.viewerMu.RUnlock()
		for viewer := range pair.viewers {
			viewer.ViewSlotChange(slot, after)
		}
	})
	c.paired, pair.paired = false, false
	c.viewerMu, pair.viewerMu = new(sync.RWMutex), new(sync.RWMutex)
	c.viewers, pair.viewers = make(map[ContainerViewer]struct{}, 1), make(map[ContainerViewer]struct{}, 1)
	c.pairInv, pair.pairInv = nil, nil
	return c, pair, true
}

func (c Chest) pairPos(pos cube.Pos) cube.Pos {
	return cube.Pos{c.pairX, pos[1], c.pairZ}
}

func (c Chest) DecodeNBT(data map[string]any) any {
	facing := c.Facing

	c = NewChest()
	c.Facing = facing
	c.CustomName = nbtconv.String(data, "CustomName")

	pairX, ok := data["pairx"]
	pairZ, ok2 := data["pairz"]
	if ok && ok2 {
		pairX, ok := pairX.(int32)
		pairZ, ok2 := pairZ.(int32)
		if ok && ok2 {
			c.paired = true
			c.pairX, c.pairZ = int(pairX), int(pairZ)
		}
	}

	nbtconv.InvFromNBT(c.inventory, nbtconv.Slice(data, "Items"))
	return c
}

func (c Chest) EncodeNBT() map[string]any {
	if c.inventory == nil {
		facing, customName := c.Facing, c.CustomName

		c = NewChest()
		c.Facing, c.CustomName = facing, customName
	}
	m := map[string]any{
		"Items": nbtconv.InvToNBT(c.inventory),
		"id":    "Chest",
	}
	if c.CustomName != "" {
		m["CustomName"] = c.CustomName
	}

	if c.paired {
		m["pairx"] = int32(c.pairX)
		m["pairz"] = int32(c.pairZ)
	}
	return m
}

func (Chest) EncodeItem() (name string, meta int16) {
	return "minecraft:chest", 0
}

func (c Chest) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:chest", map[string]any{"minecraft:cardinal_direction": c.Facing.String()}
}

func allChests() (chests []world.Block) {
	for _, direction := range cube.Directions() {
		chests = append(chests, Chest{Facing: direction})
	}
	return
}
