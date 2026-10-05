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

package block

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

type solid struct{}

func (solid) Model() world.BlockModel {
	return model.Solid{}
}

type empty struct{}

func (empty) Model() world.BlockModel {
	return model.Empty{}
}

type chest struct{}

func (chest) Model() world.BlockModel {
	return model.Chest{}
}

type carpet struct{}

func (carpet) Model() world.BlockModel {
	return model.Carpet{}
}

type tilledGrass struct{}

func (tilledGrass) Model() world.BlockModel {
	return model.TilledGrass{}
}

type leaves struct{}

func (leaves) Model() world.BlockModel {
	return model.Leaves{}
}

type thin struct{}

func (thin) Model() world.BlockModel {
	return model.Thin{}
}
