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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"math"
	"time"
)

type Beacon struct {
	solid
	transparent
	clicksAndSticks
	sourceWaterDisplacer

	Primary, Secondary effect.LastingType

	level int
}

type BeaconSource interface {
	PowersBeacon() bool
}

func (b Beacon) BreakInfo() BreakInfo {
	return newBreakInfo(3, alwaysHarvestable, nothingEffective, oneOf(Beacon{}))
}

func (b Beacon) Activate(pos cube.Pos, _ cube.Face, tx *world.Tx, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos, tx)
		return true
	}
	return true
}

func (b Beacon) DecodeNBT(data map[string]any) any {
	b.level = int(nbtconv.Int32(data, "Levels"))
	if primary, ok := effect.ByID(int(nbtconv.Int32(data, "Primary"))); ok {
		b.Primary = primary.(effect.LastingType)
	}
	if secondary, ok := effect.ByID(int(nbtconv.Int32(data, "Secondary"))); ok {
		b.Secondary = secondary.(effect.LastingType)
	}
	return b
}

func (b Beacon) EncodeNBT() map[string]any {
	m := map[string]any{
		"id":     "Beacon",
		"Levels": int32(b.level),
	}
	if primary, ok := effect.ID(b.Primary); ok {
		m["Primary"] = int32(primary)
	}
	if secondary, ok := effect.ID(b.Secondary); ok {
		m["Secondary"] = int32(secondary)
	}
	return m
}

func (b Beacon) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

func (Beacon) LightEmissionLevel() uint8 {
	return 15
}

func (b Beacon) Level() int {
	return b.level
}

func (b Beacon) Tick(currentTick int64, pos cube.Pos, tx *world.Tx) {
	if currentTick%80 == 0 {
		before := b.level

		b.level = b.recalculateLevel(pos, tx)
		if before != b.level {
			tx.SetBlock(pos, b, nil)
		}
		if b.level == 0 {
			return
		}
		if !b.obstructed(pos, tx) {
			b.broadcastBeaconEffects(pos, tx)
		}
	}
}

func (b Beacon) recalculateLevel(pos cube.Pos, tx *world.Tx) int {
	var lvl int
	iter := 1

	for y := pos.Y() - 1; y >= pos.Y()-4; y-- {
		for x := pos.X() - iter; x <= pos.X()+iter; x++ {
			for z := pos.Z() - iter; z <= pos.Z()+iter; z++ {
				if s, ok := tx.Block(cube.Pos{x, y, z}).(BeaconSource); !ok || !s.PowersBeacon() {
					return lvl
				}
			}
		}
		iter++
		lvl++
	}
	return lvl
}

func (b Beacon) obstructed(pos cube.Pos, tx *world.Tx) bool {

	if tx.SkyLight(pos.Side(cube.FaceUp)) == 15 {
		return false
	}

	return tx.HighestLightBlocker(pos.X(), pos.Z()) > pos[1]
}

func (b Beacon) broadcastBeaconEffects(pos cube.Pos, tx *world.Tx) {
	seconds := 9 + b.level*2
	if b.level == 4 {
		seconds--
	}
	dur := time.Duration(seconds) * time.Second

	primary, secondary := b.Primary, effect.LastingType(nil)
	switch b.level {
	case 0:
		primary = nil
	case 1:
		switch primary {
		case effect.Resistance, effect.JumpBoost, effect.Strength:
			primary = nil
		}
	case 2:
		if primary == effect.Strength {
			primary = nil
		}
	case 3:

	default:
		secondary = b.Secondary
	}
	var primaryEff, secondaryEff effect.Effect

	if primary != nil {
		primaryEff = effect.NewAmbient(primary, 1, dur)

		if secondary != nil {

			if primary == secondary {
				primaryEff = effect.NewAmbient(primary, 2, dur)
			} else {
				secondaryEff = effect.NewAmbient(secondary, 1, dur)
			}
		}
	}

	r := 10 + (b.level * 10)
	entitiesInRange := tx.EntitiesWithin(cube.Box(
		float64(pos.X()-r), -math.MaxFloat64, float64(pos.Z()-r),
		float64(pos.X()+r), math.MaxFloat64, float64(pos.Z()+r),
	))
	for e := range entitiesInRange {
		if p, ok := e.(beaconAffected); ok {
			if primaryEff.Type() != nil {
				p.AddEffect(primaryEff)
			}
			if secondaryEff.Type() != nil {
				p.AddEffect(secondaryEff)
			}
		}
	}
}

type beaconAffected interface {
	AddEffect(e effect.Effect)

	BeaconAffected() bool
}

func (Beacon) EncodeItem() (name string, meta int16) {
	return "minecraft:beacon", 0
}

func (Beacon) EncodeBlock() (string, map[string]any) {
	return "minecraft:beacon", nil
}
