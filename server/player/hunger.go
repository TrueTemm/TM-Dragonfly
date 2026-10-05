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
	"sync"
)

type hungerManager struct {
	mu              sync.RWMutex
	foodLevel       int
	saturationLevel float64
	exhaustionLevel float64
	foodTick        int
}

func newHungerManager() *hungerManager {
	return &hungerManager{foodLevel: 20, saturationLevel: 5, foodTick: 1}
}

func (m *hungerManager) Food() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.foodLevel
}

func (m *hungerManager) SetFood(level int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.foodLevel = max(min(level, 20), 0)
}

func (m *hungerManager) AddFood(points int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.foodLevel = max(min(m.foodLevel+points, 20), 0)
}

func (m *hungerManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.foodLevel = 20
	m.saturationLevel = 5
	m.exhaustionLevel = 0
	m.foodTick = 1
}

func (m *hungerManager) resetExhaustion() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exhaustionLevel = 0
	m.saturationLevel = 0
	m.foodTick = 1
}

func (m *hungerManager) exhaust(points float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.exhaustionLevel += points
	for m.exhaustionLevel >= 4 {

		m.exhaustionLevel -= 4
		m.desaturate()
	}
}

func (m *hungerManager) saturate(food int, saturation float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.foodLevel = max(min(m.foodLevel+food, 20), 0)
	m.saturationLevel = max(min(m.saturationLevel+saturation, float64(m.foodLevel)), 0)
}

func (m *hungerManager) desaturate() {
	if m.saturationLevel <= 0 && m.foodLevel != 0 {
		m.foodLevel--
	} else if m.saturationLevel > 0 {
		m.saturationLevel = max(m.saturationLevel-1, 0)
	}
}

func (m *hungerManager) canQuicklyRegenerate() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.foodLevel == 20 && m.saturationLevel > 0
}

func (m *hungerManager) canRegenerate() bool {
	return m.Food() >= 18
}

func (m *hungerManager) canSprint() bool {
	return m.Food() > 6
}

func (m *hungerManager) starving() bool {
	return m.Food() == 0
}

type StarvationDamageSource struct{}

func (StarvationDamageSource) ReducedByArmour() bool     { return false }
func (StarvationDamageSource) ReducedByResistance() bool { return false }
func (StarvationDamageSource) Fire() bool                { return false }
func (StarvationDamageSource) IgnoreTotem() bool         { return false }
