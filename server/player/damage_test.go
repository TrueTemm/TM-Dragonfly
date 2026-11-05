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

package player

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func withDamagePlayer(t *testing.T, f func(*world.Tx, *Player)) {
	t.Helper()
	w := (world.Config{Synchronous: true}).New()
	defer w.Close()
	w.SetFallDamage(true)
	h := (world.EntitySpawnOpts{}).New(Type, Config{Position: mgl64.Vec3{0.5, 65, 0.5}})
	task := w.Do(func(tx *world.Tx) { f(tx, tx.AddEntity(h).(*Player)) })
	if err := task.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestDamageSources(t *testing.T) {
	sources := map[string]world.DamageSource{
		"attack": entity.AttackDamageSource{}, "projectile": entity.ProjectileDamageSource{},
		"explosion": entity.ExplosionDamageSource{}, "fall": entity.FallDamageSource{},
		"glide": entity.GlideDamageSource{}, "lightning": entity.LightningDamageSource{},
		"void": entity.VoidDamageSource{}, "suffocation": entity.SuffocationDamageSource{},
		"drowning": entity.DrowningDamageSource{}, "fire": block.FireDamageSource{},
		"lava": block.LavaDamageSource{}, "magma": block.MagmaDamageSource{},
		"cactus": block.DamageSource{Block: block.Cactus{}}, "poison": effect.PoisonDamageSource{},
		"fatal_poison": effect.PoisonDamageSource{Fatal: true}, "wither": effect.WitherDamageSource{},
		"instant_damage": effect.InstantDamageSource{}, "starvation": StarvationDamageSource{},
		"thorns": enchantment.ThornsDamageSource{},
	}
	for name, src := range sources {
		t.Run(name, func(t *testing.T) {
			withDamagePlayer(t, func(_ *world.Tx, p *Player) {
				damage, vulnerable := p.Hurt(4, src)
				if damage != 4 || !vulnerable || p.Health() != 16 {
					t.Fatalf("damage=%v vulnerable=%v health=%v", damage, vulnerable, p.Health())
				}
			})
		})
	}
}

func TestDamageTriggers(t *testing.T) {
	cases := []struct {
		name   string
		damage float64
		apply  func(*world.Tx, *Player)
	}{
		{"fall", 3, func(_ *world.Tx, p *Player) { p.fall(6) }},
		{"fire", 1, func(tx *world.Tx, p *Player) { (block.Fire{}).EntityInside(cube.Pos{}, tx, p) }},
		{"soul_fire", 2, func(tx *world.Tx, p *Player) { (block.Fire{Type: block.SoulFire()}).EntityInside(cube.Pos{}, tx, p) }},
		{"lava", 4, func(tx *world.Tx, p *Player) { (block.Lava{}).EntityInside(cube.Pos{}, tx, p) }},
		{"magma", 1, func(tx *world.Tx, p *Player) { (block.Magma{}).EntityStepOn(cube.Pos{}, tx, p) }},
		{"cactus", 0.5, func(tx *world.Tx, p *Player) { (block.Cactus{}).EntityInside(cube.Pos{}, tx, p) }},
		{"campfire", 1, func(tx *world.Tx, p *Player) { (block.Campfire{}).EntityInside(cube.Pos{}, tx, p) }},
		{"poison", 1, func(_ *world.Tx, p *Player) { effect.Poison.Apply(p, effect.New(effect.Poison, 1, time.Second)) }},
		{"fatal_poison", 1, func(_ *world.Tx, p *Player) {
			effect.FatalPoison.Apply(p, effect.New(effect.FatalPoison, 1, time.Second))
		}},
		{"wither", 1, func(_ *world.Tx, p *Player) { effect.Wither.Apply(p, effect.New(effect.Wither, 1, time.Second)) }},
		{"instant_damage", 6, func(_ *world.Tx, p *Player) { p.AddEffect(effect.NewInstant(effect.InstantDamage, 1)) }},
		{"drowning", 2, func(tx *world.Tx, p *Player) {
			tx.SetLiquid(cube.Pos{0, 66, 0}, block.Water{Depth: 8})
			p.airSupplyTicks = -19
			p.tickAirSupply()
		}},
		{"suffocation", 1, func(tx *world.Tx, p *Player) { tx.SetBlock(cube.Pos{0, 66, 0}, block.Stone{}, nil); p.Tick(tx, 1) }},
		{"void", 4, func(tx *world.Tx, p *Player) { p.data.Pos[1] = float64(tx.Range()[0] - 1); p.Tick(tx, 1) }},
		{"starvation", 1, func(_ *world.Tx, p *Player) { p.starve() }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDamagePlayer(t, func(tx *world.Tx, p *Player) {
				c.apply(tx, p)
				if p.Health() != 20-c.damage {
					t.Fatalf("health=%v want=%v", p.Health(), 20-c.damage)
				}
			})
		})
	}
}

func TestFallDamageRules(t *testing.T) {
	cases := []struct {
		name             string
		distance, health float64
		setup            func(*world.Tx, *Player)
	}{
		{"safe_height", 3, 20, nil},
		{"normal_landing", 9, 14, func(tx *world.Tx, _ *Player) { tx.SetBlock(cube.Pos{0, 64, 0}, block.Stone{}, nil) }},
		{"disabled", 9, 20, func(tx *world.Tx, _ *Player) { tx.World().SetFallDamage(false) }},
		{"creative", 9, 20, func(_ *world.Tx, p *Player) { p.SetGameMode(world.GameModeCreative) }},
		{"spectator", 9, 20, func(_ *world.Tx, p *Player) { p.SetGameMode(world.GameModeSpectator) }},
		{"jump_boost", 9, 16, func(_ *world.Tx, p *Player) { p.AddEffect(effect.New(effect.JumpBoost, 2, time.Minute)) }},
		{"resistance", 8, 16, func(_ *world.Tx, p *Player) { p.AddEffect(effect.New(effect.Resistance, 1, time.Minute)) }},
		{"slime", 9, 20, func(tx *world.Tx, _ *Player) { tx.SetBlock(cube.Pos{0, 64, 0}, block.Slime{}, nil) }},
		{"hay_bale", 10, 20, func(tx *world.Tx, _ *Player) { tx.SetBlock(cube.Pos{0, 64, 0}, block.HayBale{}, nil) }},
		{"bed", 10, 18, func(tx *world.Tx, _ *Player) { tx.SetBlock(cube.Pos{0, 64, 0}, block.Bed{}, nil) }},
		{"armour_does_not_reduce_fall", 9, 14, func(_ *world.Tx, p *Player) {
			p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1))
		}},
		{"feather_falling", 9, 16.88, func(_ *world.Tx, p *Player) {
			p.Armour().SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.FeatherFalling, 4)))
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withDamagePlayer(t, func(tx *world.Tx, p *Player) {
				if c.setup != nil {
					c.setup(tx, p)
				}
				p.fall(c.distance)
				if math.Abs(p.Health()-c.health) > 1e-6 {
					t.Fatalf("health=%v want=%v", p.Health(), c.health)
				}
			})
		})
	}
}

func TestDamageProtectionAndImmunity(t *testing.T) {
	t.Run("fire_resistance", func(t *testing.T) {
		withDamagePlayer(t, func(_ *world.Tx, p *Player) {
			p.AddEffect(effect.New(effect.FireResistance, 1, time.Minute))
			for _, src := range []world.DamageSource{block.FireDamageSource{}, block.LavaDamageSource{}, block.MagmaDamageSource{}} {
				p.Hurt(4, src)
			}
			if p.Health() != 20 {
				t.Fatal("fire resistance allowed fire damage")
			}
			p.Hurt(4, entity.DrowningDamageSource{})
			if p.Health() != 16 {
				t.Fatal("fire resistance prevented drowning")
			}
		})
	})
	t.Run("absorption", func(t *testing.T) {
		withDamagePlayer(t, func(_ *world.Tx, p *Player) {
			p.SetAbsorption(3)
			p.Hurt(5, entity.AttackDamageSource{})
			if p.Absorption() != 0 || p.Health() != 18 {
				t.Fatalf("absorption=%v health=%v", p.Absorption(), p.Health())
			}
		})
	})
	t.Run("cooldown", func(t *testing.T) {
		withDamagePlayer(t, func(_ *world.Tx, p *Player) {
			p.Hurt(4, entity.AttackDamageSource{})
			p.Hurt(4, entity.AttackDamageSource{})
			if p.Health() != 16 {
				t.Fatal("repeated melee bypassed cooldown")
			}
			p.Hurt(6, entity.FallDamageSource{})
			if p.Health() != 14 {
				t.Fatal("larger damage during cooldown did not apply the difference")
			}
			p.immuneUntil = time.Time{}
			p.Hurt(2, entity.AttackDamageSource{})
			if p.Health() != 12 {
				t.Fatal("damage remained blocked after cooldown")
			}
		})
	})
}

type cancelledDamageHandler struct{ NopHandler }

func (cancelledDamageHandler) HandleHurt(ctx *Context, _ *float64, _ bool, _ *time.Duration, _ world.DamageSource) {
	ctx.Cancel()
}

func TestDamageCancellationAndDeath(t *testing.T) {
	t.Run("cancelled", func(t *testing.T) {
		withDamagePlayer(t, func(_ *world.Tx, p *Player) {
			p.Handle(cancelledDamageHandler{})
			damage, vulnerable := p.Hurt(4, entity.AttackDamageSource{})
			if damage != 0 || vulnerable || p.Health() != 20 {
				t.Fatal("cancelled damage changed health")
			}
		})
	})
	t.Run("totem", func(t *testing.T) {
		withDamagePlayer(t, func(_ *world.Tx, p *Player) {
			p.SetHeldItems(item.NewStack(item.Totem{}, 1), item.Stack{})
			p.Hurt(100, entity.FallDamageSource{})
			hand, _ := p.HeldItems()
			if p.Dead() || !hand.Empty() {
				t.Fatal("totem did not prevent death and consume itself")
			}
		})
	})
	t.Run("void", func(t *testing.T) {
		withDamagePlayer(t, func(_ *world.Tx, p *Player) {
			p.Hurt(100, entity.VoidDamageSource{})
			if !p.Dead() || p.Health() != 0 {
				t.Fatal("lethal damage did not kill player")
			}
			p.Hurt(4, entity.AttackDamageSource{})
			if p.Health() != 0 {
				t.Fatal("damage changed dead player health")
			}
		})
	})
}
