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

package playerdb

import "github.com/df-mc/dragonfly/server/entity/effect"

func effectsToData(effects []effect.Effect) []jsonEffect {
	data := make([]jsonEffect, len(effects))
	for key, eff := range effects {
		id, ok := effect.ID(eff.Type())
		if !ok {
			continue
		}
		data[key] = jsonEffect{
			ID:              id,
			Duration:        eff.Duration(),
			Level:           eff.Level(),
			Ambient:         eff.Ambient(),
			ParticlesHidden: eff.ParticlesHidden(),
			Infinite:        eff.Infinite(),
		}
	}
	return data
}

func dataToEffects(data []jsonEffect) []effect.Effect {
	effects := make([]effect.Effect, len(data))
	for i, d := range data {
		e, ok := effect.ByID(d.ID)
		if !ok {
			continue
		}
		switch eff := e.(type) {
		case effect.LastingType:
			switch {
			case d.Ambient:
				effects[i] = effect.NewAmbient(eff, d.Level, d.Duration)
			case d.Infinite:
				effects[i] = effect.NewInfinite(eff, d.Level)
			default:
				effects[i] = effect.New(eff, d.Level, d.Duration)
			}

			if d.ParticlesHidden {
				effects[i] = effects[i].WithoutParticles()
			}
		default:
			effects[i] = effect.NewInstant(eff, d.Level)
		}
	}
	return effects
}
