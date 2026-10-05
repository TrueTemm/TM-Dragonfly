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
	"github.com/df-mc/dragonfly/server/world"
	"math"
	"time"
)

type StationaryBehaviourConfig struct {
	ExistenceDuration time.Duration

	SpawnSounds []world.Sound

	Tick func(e *Ent, tx *world.Tx)
}

func (conf StationaryBehaviourConfig) Apply(data *world.EntityData) {
	data.Data = conf.New()
}

func (conf StationaryBehaviourConfig) New() *StationaryBehaviour {
	if conf.ExistenceDuration == 0 {
		conf.ExistenceDuration = math.MaxInt64
	}
	return &StationaryBehaviour{BaseBehaviour: NewBaseBehaviour(), conf: conf}
}

type StationaryBehaviour struct {
	BaseBehaviour

	conf  StationaryBehaviourConfig
	close bool
}

func (s *StationaryBehaviour) Tick(e *Ent, tx *world.Tx) *Movement {
	if s.close {
		_ = e.Close()
		return nil
	}

	if e.Age() == 0 {
		for _, ss := range s.conf.SpawnSounds {
			tx.PlaySound(e.Position(), ss)
		}
	}
	if s.conf.Tick != nil {
		s.conf.Tick(e, tx)
	}

	if e.Age() > s.conf.ExistenceDuration {
		s.close = true
	}

	return nil
}

func (s *StationaryBehaviour) Immobile() bool {
	return true
}
