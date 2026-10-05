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
	"math/rand/v2"
	"sync"
	"sync/atomic"

	"github.com/df-mc/dragonfly/server/block/cube"
)

type Settings struct {
	sync.Mutex
	ref atomic.Int32

	Name string

	Spawn cube.Pos

	Time int64

	TimeCycle bool

	RainTime int64

	Raining bool

	ThunderTime int64

	Thundering bool

	WeatherCycle bool

	RequiredSleepTicks int64

	CurrentTick int64

	DefaultGameMode GameMode

	Difficulty Difficulty

	TickRange int32

	FallDamage bool
}

func defaultSettings() *Settings {
	return &Settings{
		Name:            "World",
		DefaultGameMode: GameModeSurvival,
		Difficulty:      DifficultyNormal,
		TimeCycle:       true,

		WeatherCycle: false,
		RainTime:     int64(rand.IntN(8400)+600) * 20,
		ThunderTime:  int64(rand.IntN(8400)+600) * 20,
		TickRange:    6,
		FallDamage:   true,
	}
}
