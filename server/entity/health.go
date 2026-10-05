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

type HealthManager struct {
	health float64
	max    float64
}

func NewHealthManager(health, max float64) *HealthManager {
	if health > max {
		health = max
	}
	return &HealthManager{health: health, max: max}
}

func (m *HealthManager) Health() float64 {
	return m.health
}

func (m *HealthManager) AddHealth(health float64) {
	m.health = max(min(m.health+health, m.max), 0)
}

func (m *HealthManager) MaxHealth() float64 {
	return m.max
}

func (m *HealthManager) SetMaxHealth(max float64) {
	if max <= 0 {
		max = 1
	}
	m.max = max
	m.health = min(m.health, max)
}
