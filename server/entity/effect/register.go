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

func Register(id int, e Type) {
	effects[id] = e
	effectIds[e] = id
}

func init() {
	Register(1, Speed)
	Register(2, Slowness)
	Register(3, Haste)
	Register(4, MiningFatigue)
	Register(5, Strength)
	Register(6, InstantHealth)
	Register(7, InstantDamage)
	Register(8, JumpBoost)
	Register(9, Nausea)
	Register(10, Regeneration)
	Register(11, Resistance)
	Register(12, FireResistance)
	Register(13, WaterBreathing)
	Register(14, Invisibility)
	Register(15, Blindness)
	Register(16, NightVision)
	Register(17, Hunger)
	Register(18, Weakness)
	Register(19, Poison)
	Register(20, Wither)
	Register(21, HealthBoost)
	Register(22, Absorption)
	Register(23, Saturation)
	Register(24, Levitation)
	Register(25, FatalPoison)
	Register(26, ConduitPower)
	Register(27, SlowFalling)

	Register(30, Darkness)
}

var (
	effects   = map[int]Type{}
	effectIds = map[Type]int{}
)

func ByID(id int) (Type, bool) {
	effect, ok := effects[id]
	return effect, ok
}

func ID(e Type) (int, bool) {
	id, ok := effectIds[e]
	return id, ok
}
