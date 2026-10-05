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

package item

import (
	"time"
)

type Smeltable interface {
	SmeltInfo() SmeltInfo
}

type Fuel interface {
	FuelInfo() FuelInfo
}

type SmeltInfo struct {
	Product Stack

	Experience float64

	Food bool

	Ores bool
}

func newSmeltInfo(product Stack, experience float64) SmeltInfo {
	return SmeltInfo{
		Product:    product,
		Experience: experience,
	}
}

func newFoodSmeltInfo(product Stack, experience float64) SmeltInfo {
	return SmeltInfo{
		Product:    product,
		Experience: experience,
		Food:       true,
	}
}

func newOreSmeltInfo(product Stack, experience float64) SmeltInfo {
	return SmeltInfo{
		Product:    product,
		Experience: experience,
		Ores:       true,
	}
}

type FuelInfo struct {
	Duration time.Duration

	Residue Stack
}

func (f FuelInfo) WithResidue(residue Stack) FuelInfo {
	f.Residue = residue
	return f
}

func newFuelInfo(duration time.Duration) FuelInfo {
	return FuelInfo{Duration: duration}
}
