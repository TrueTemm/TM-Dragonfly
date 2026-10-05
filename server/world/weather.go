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

package world

import (
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)

type weather struct{ w *World }

func (w weather) StopWeatherCycle() {
	w.enableWeatherCycle(false)
}

func (w weather) StartWeatherCycle() {
	w.enableWeatherCycle(true)
}

func (tx *Tx) snowingAt(pos cube.Pos) bool {
	w := tx.World()
	if w == nil || !w.Dimension().WeatherCycle() {
		return false
	}
	if b := tx.biome(pos); b.Rainfall() == 0 || tx.temperature(pos) > 0.15 {
		return false
	}
	w.set.Lock()
	raining := w.set.Raining
	w.set.Unlock()
	return raining && tx.highestObstructingBlock(pos[0], pos[2]) < pos[1]
}

func (tx *Tx) rainingAt(pos cube.Pos) bool {
	w := tx.World()
	if w == nil || !w.Dimension().WeatherCycle() {
		return false
	}
	if b := tx.biome(pos); b.Rainfall() == 0 || tx.temperature(pos) <= 0.15 {
		return false
	}
	w.set.Lock()
	a := w.set.Raining
	w.set.Unlock()
	return a && tx.highestObstructingBlock(pos[0], pos[2]) < pos[1]
}

func (tx *Tx) thunderingAt(pos cube.Pos) bool {
	w := tx.World()
	raining := tx.rainingAt(pos)
	w.set.Lock()
	a := w.set.Thundering && raining
	w.set.Unlock()
	return a && tx.highestObstructingBlock(pos[0], pos[2]) < pos[1]
}

func (w weather) raining() bool {
	if w.w == nil || !w.w.Dimension().WeatherCycle() {
		return false
	}
	w.w.set.Lock()
	defer w.w.set.Unlock()
	return w.w.set.Raining
}

func (w weather) thundering() bool {
	if w.w == nil || !w.w.Dimension().WeatherCycle() {
		return false
	}
	w.w.set.Lock()
	defer w.w.set.Unlock()
	return w.w.set.Thundering
}

func (w weather) StartRaining(dur time.Duration) {
	w.w.set.Lock()
	defer w.w.set.Unlock()
	w.setRaining(true, dur)
}

func (w weather) StopRaining() {
	w.w.set.Lock()
	defer w.w.set.Unlock()

	if w.w.set.Raining {
		w.setRaining(false, time.Second*(time.Duration(w.w.r.IntN(8400)+600)))
		if w.w.set.Thundering {

			w.setThunder(false, time.Second*(time.Duration(w.w.r.IntN(8400)+600)))
		}
	}
}

func (w weather) StartThundering(dur time.Duration) {
	w.w.set.Lock()
	defer w.w.set.Unlock()

	w.setThunder(true, dur)
	w.setRaining(true, dur)
}

func (w weather) StopThundering() {
	w.w.set.Lock()
	defer w.w.set.Unlock()
	if w.w.set.Thundering && w.w.set.Raining {
		w.setThunder(false, time.Second*(time.Duration(w.w.r.IntN(8400)+600)))
	}
}

func (w weather) advanceWeather() {
	w.w.set.RainTime--
	w.w.set.ThunderTime--

	if w.w.set.RainTime <= 0 {

		if w.w.set.Raining {
			w.w.setRaining(false, time.Second*(time.Duration(w.w.r.IntN(8400)+600)))
		} else {
			w.w.setRaining(true, time.Second*time.Duration(w.w.r.IntN(600)+600))
		}
	}
	if w.w.set.ThunderTime <= 0 {

		if w.w.set.Thundering {
			w.w.setThunder(false, time.Second*(time.Duration(w.w.r.IntN(8400)+600)))
		} else {
			w.w.setThunder(true, time.Second*time.Duration(w.w.r.IntN(620)+180))
		}
	}
}

func (w weather) setRaining(raining bool, x time.Duration) {
	w.w.set.Raining = raining
	w.w.set.RainTime = int64(x.Seconds() * 20)
}

func (w weather) setThunder(thundering bool, x time.Duration) {
	w.w.set.Thundering = thundering
	w.w.set.ThunderTime = int64(x.Seconds() * 20)
}

func (w weather) enableWeatherCycle(v bool) {
	if w.w == nil {
		return
	}
	w.w.set.Lock()
	defer w.w.set.Unlock()
	w.w.set.WeatherCycle = v
}

func (w weather) tickLightning(tx *Tx) {
	positions := make([]ChunkPos, 0, len(w.w.chunks)/100000)
	for pos := range w.w.chunks {

		if w.w.r.IntN(100000) == 0 {
			positions = append(positions, pos)
		}
	}

	for _, pos := range positions {
		w.w.strikeLightning(tx, pos)
	}
}

func (w weather) strikeLightning(tx *Tx, c ChunkPos) {
	if w.w.conf.Entities.conf.Lightning == nil {
		return
	}
	if pos := w.lightningPosition(tx, c); tx.ThunderingAt(cube.PosFromVec3(pos)) {
		tx.AddEntity(w.w.conf.Entities.conf.Lightning(EntitySpawnOpts{Position: pos}))
	}
}

func (w weather) lightningPosition(tx *Tx, c ChunkPos) mgl64.Vec3 {
	v := w.w.r.Int32()
	x, z := float64(c[0]<<4+(v&0xf)), float64(c[1]<<4+((v>>8)&0xf))

	vec := w.adjustPositionToEntities(tx, mgl64.Vec3{x, float64(tx.HighestBlock(int(x), int(z)) + 1), z})
	if pos := cube.PosFromVec3(vec); len(tx.Block(pos).Model().BBox(pos, tx)) != 0 {

		return vec.Add(mgl64.Vec3{0, 1})
	}
	return vec
}

func (w weather) adjustPositionToEntities(tx *Tx, vec mgl64.Vec3) mgl64.Vec3 {
	max := vec.Add(mgl64.Vec3{0, float64(w.w.Range().Max())})

	list := make([]mgl64.Vec3, 0, 16)
	for e := range tx.EntitiesWithin(cube.Box(vec[0], vec[1], vec[2], max[0], max[1], max[2]).GrowVec3(mgl64.Vec3{3, 3, 3})) {
		if h, ok := e.(interface{ Health() float64 }); ok && h.Health() > 0 {

			pos := cube.PosFromVec3(e.Position())
			if tx.HighestBlock(pos[0], pos[2]) < pos[1] {
				list = append(list, e.Position())
			}
		}
	}

	if len(list) > 0 {
		vec = list[w.w.r.IntN(len(list))]
	}
	return vec
}
