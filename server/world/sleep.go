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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/google/uuid"
)

type Sleeper interface {
	Entity

	Name() string
	UUID() uuid.UUID

	Messaget(t chat.Translation, a ...any)
	SendSleepingIndicator(sleeping, max int)

	Sleep(pos cube.Pos)
	Sleeping() (cube.Pos, bool)
	Wake()
}

const (
	TimeSleep         = 12542
	TimeWake          = 23459
	TimeSleepWithRain = 12010
	TimeWakeWithRain  = 23991
	TimeFull          = 24000
)

func (ticker) tryAdvanceDay(tx *Tx, timeCycle bool) {
	sleepers := tx.Sleepers()
	time := tx.w.Time() % TimeFull

	for s := range sleepers {
		if !tx.Thundering() {
			if !tx.Raining() && (time <= TimeSleep || time >= TimeWake) {
				return
			}
			if time <= TimeSleepWithRain || time >= TimeWakeWithRain {
				return
			}
		}

		if _, ok := s.Sleeping(); !ok {

			return
		}
	}

	for s := range sleepers {
		s.Wake()
	}

	totalTime := tx.w.Time()
	if timeCycle {
		tx.w.SetTime(totalTime + TimeFull - time)
	}
	tx.w.StopRaining()
}
