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

package entity

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type SplashableBlock interface {
	world.Block

	Splash(tx *world.Tx, pos cube.Pos)
}

type SplashableEntity interface {
	world.Entity

	Splash(tx *world.Tx, pos mgl64.Vec3)
}

const (
	potHit     = 1.3055
	potMiss    = 1.2015
	potLasting = 0.75
)

func potionSplash(durMul float64, pot potion.Potion, linger bool) func(e *Ent, tx *world.Tx, res trace.Result) {
	return func(e *Ent, tx *world.Tx, res trace.Result) {
		pos := e.Position()
		effects := pot.Effects()
		box := e.H().Type().BBox(e).Translate(pos)

		if len(effects) > 0 {

			reach := box.GrowVec3(mgl64.Vec3{1.85, 2.65, 1.85})
			for otherE := range filterLiving(tx.EntitiesWithin(reach)) {
				otherPos := otherE.Position()
				if !otherE.H().Type().BBox(otherE).Translate(otherPos).IntersectsWith(reach) {
					continue
				}

				f := potMiss
				if entityResult, ok := res.(trace.EntityResult); ok && entityResult.Entity().H() == otherE.H() {
					f = potHit
				}

				splashed := otherE.(Living)
				for _, eff := range effects {
					if _, ok := eff.Type().(effect.LastingType); !ok {
						splashed.AddEffect(effect.NewInstantWithPotency(eff.Type(), eff.Level(), f))
						continue
					}

					dur := time.Duration(float64(eff.Duration()) * durMul * potLasting * f)
					if dur < time.Second {
						continue
					}
					splashed.AddEffect(effect.New(eff.Type().(effect.LastingType), eff.Level(), dur))
				}
			}
		} else if pot == potion.Water() {
			switch result := res.(type) {
			case trace.BlockResult:
				blockPos := result.BlockPosition().Side(result.Face())
				if _, ok := tx.Block(blockPos).(block.Fire); ok {
					tx.SetBlock(blockPos, nil, nil)
				}

				for _, f := range cube.HorizontalFaces() {
					h := blockPos.Side(f)
					if _, ok := tx.Block(h).(block.Fire); ok {
						tx.SetBlock(h, nil, nil)
					}

					if b, ok := tx.Block(blockPos.Side(f)).(SplashableBlock); ok {
						b.Splash(tx, blockPos.Side(f))
					}
				}

				resultPos := result.BlockPosition()
				if b, ok := tx.Block(resultPos).(SplashableBlock); ok {
					b.Splash(tx, resultPos)
				}
			case trace.EntityResult:

			}

			for otherE := range filterLiving(tx.EntitiesWithin(box.GrowVec3(mgl64.Vec3{8.25, 4.25, 8.25}))) {
				if splashE, ok := otherE.(SplashableEntity); ok {
					splashE.Splash(tx, otherE.Position())
				}
			}
		}
		if linger {
			tx.AddEntity(NewAreaEffectCloud(world.EntitySpawnOpts{Position: pos}, pot))
		}
	}
}
