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

package effect

import (
	"image/color"
	"time"

	"github.com/df-mc/dragonfly/server/world"
)

type LastingType interface {
	Type

	Start(e world.Entity, lvl int)

	End(e world.Entity, lvl int)
}

type Type interface {
	RGBA() color.RGBA

	Apply(e world.Entity, eff Effect)
}

type Effect struct {
	t                        Type
	d                        time.Duration
	lvl                      int
	potency                  float64
	ambient, particlesHidden bool
	infinite                 bool
	tick                     int
}

func NewInstant(t Type, lvl int) Effect {
	return NewInstantWithPotency(t, lvl, 1)
}

func NewInstantWithPotency(t Type, lvl int, potency float64) Effect {
	return Effect{t: t, lvl: lvl, potency: potency}
}

func New(t LastingType, lvl int, d time.Duration) Effect {
	return Effect{t: t, lvl: lvl, d: d}
}

func NewAmbient(t LastingType, lvl int, d time.Duration) Effect {
	return Effect{t: t, lvl: lvl, d: d, ambient: true}
}

func NewInfinite(t LastingType, lvl int) Effect {
	return Effect{t: t, lvl: lvl, infinite: true}
}

func (e Effect) WithoutParticles() Effect {
	e.particlesHidden = true
	return e
}

func (e Effect) ParticlesHidden() bool {
	return e.particlesHidden
}

func (e Effect) Level() int {
	return e.lvl
}

func (e Effect) Duration() time.Duration {
	return e.d
}

func (e Effect) Ambient() bool {
	return e.ambient
}

func (e Effect) Infinite() bool {
	return e.infinite
}

func (e Effect) Type() Type {
	return e.t
}

func (e Effect) TickDuration() Effect {
	if _, ok := e.t.(LastingType); ok {
		if !e.Infinite() {
			e.d -= time.Second / 20
		}
		e.tick++
	}
	return e
}

func (e Effect) Tick() int {
	return e.tick
}

type nopLasting struct{}

func (nopLasting) Apply(world.Entity, Effect) {}
func (nopLasting) End(world.Entity, int)      {}
func (nopLasting) Start(world.Entity, int)    {}

func ResultingColour(effects []Effect) (color.RGBA, bool) {
	r, g, b, a, n := 0, 0, 0, 0, 0
	ambient := true
	for _, e := range effects {
		if e.particlesHidden {

			continue
		}
		c := e.Type().RGBA()
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
		a += int(c.A)
		n++
		if !e.Ambient() {
			ambient = false
		}
	}
	if n == 0 {
		return color.RGBA{R: 0x38, G: 0x5d, B: 0xc6, A: 0xff}, false
	}
	return color.RGBA{R: uint8(r / n), G: uint8(g / n), B: uint8(b / n), A: uint8(a / n)}, ambient
}

type living interface {
	world.Entity

	Health() float64

	MaxHealth() float64

	SetMaxHealth(v float64)

	Hurt(damage float64, source world.DamageSource) (n float64, vulnerable bool)

	Heal(health float64, source world.HealingSource) float64

	Speed() float64

	SetSpeed(float64)
}
