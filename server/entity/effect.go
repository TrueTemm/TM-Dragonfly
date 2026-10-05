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
	"fmt"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/world"
	"maps"
	"reflect"
	"slices"
)

type EffectManager struct {
	initialEffects []effect.Effect
	effects        map[reflect.Type]effect.Effect
}

func NewEffectManager(eff ...effect.Effect) *EffectManager {
	return &EffectManager{effects: make(map[reflect.Type]effect.Effect), initialEffects: eff}
}

func (m *EffectManager) Add(e effect.Effect, entity Living) effect.Effect {
	lvl, dur := e.Level(), e.Duration()
	if lvl <= 0 {
		panic(fmt.Sprintf("(*EffectManager).Add: effect cannot have level of 0 or below: %v", lvl))
	}
	if dur < 0 {
		panic(fmt.Sprintf("(*EffectManager).Add: effect cannot have negative duration: %v", dur))
	}

	m.flushInitialEffects(entity)

	t, ok := e.Type().(effect.LastingType)
	if !ok {
		e.Type().Apply(entity, e)
		return e
	}
	typ := reflect.TypeOf(e.Type())

	existing, ok := m.effects[typ]
	if !ok {
		m.effects[typ] = e

		t.Start(entity, lvl)
		return e
	}
	if existing.Level() > lvl || (existing.Level() == lvl && ((existing.Duration() > dur && !e.Infinite()) || existing.Infinite())) {
		return existing
	}
	m.effects[typ] = e

	existing.Type().(effect.LastingType).End(entity, existing.Level())
	t.Start(entity, lvl)
	return e
}

func (m *EffectManager) Remove(e effect.Type, entity Living) {
	m.flushInitialEffects(entity)

	t := reflect.TypeOf(e)
	if existing, ok := m.effects[t]; ok {
		delete(m.effects, t)
		existing.Type().(effect.LastingType).End(entity, existing.Level())
	}
}

func (m *EffectManager) Effect(e effect.Type) (effect.Effect, bool) {
	for _, eff := range m.initialEffects {
		if eff.Type() == e {
			return eff, true
		}
	}

	existing, ok := m.effects[reflect.TypeOf(e)]
	return existing, ok
}

func (m *EffectManager) Effects() []effect.Effect {
	return append(slices.Collect(maps.Values(m.effects)), m.initialEffects...)
}

func (m *EffectManager) Tick(entity Living, tx *world.Tx) {
	update := false

	m.flushInitialEffects(entity)

	for i, eff := range m.effects {
		if m.expired(eff) {
			delete(m.effects, i)
			eff.Type().(effect.LastingType).End(entity, eff.Level())
			update = true
			continue
		}
		eff.Type().Apply(entity, eff)
		m.effects[i] = eff.TickDuration()
	}

	if update {
		for _, v := range tx.Viewers(entity.Position()) {
			v.ViewEntityState(entity)
		}
	}
}

func (m *EffectManager) flushInitialEffects(entity Living) {
	initialEffects := m.initialEffects
	m.initialEffects = nil
	for _, e := range initialEffects {
		m.Add(e, entity)
	}
}

func (m *EffectManager) expired(e effect.Effect) bool {
	return e.Duration() <= 0 && !e.Infinite()
}
