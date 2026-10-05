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

package sound

import "github.com/df-mc/dragonfly/server/world"

type ItemBreak struct{ sound }

type ItemThrow struct{ sound }

type ItemUseOn struct {
	Block world.Block

	sound
}

type EquipItem struct {
	Item world.Item

	sound
}

type BucketFill struct {
	Liquid world.Liquid

	sound
}

type BucketEmpty struct {
	Liquid world.Liquid

	sound
}

type BowShoot struct{ sound }

type CrossbowLoad struct {
	Stage int

	QuickCharge bool

	sound
}

type CrossbowShoot struct{ sound }

const (
	CrossbowLoadingStart = iota

	CrossbowLoadingMiddle

	CrossbowLoadingEnd
)

type ArrowHit struct{ sound }

type Teleport struct{ sound }

type UseSpyglass struct{ sound }

type StopUsingSpyglass struct{ sound }

type GoatHorn struct {
	Horn Horn

	sound
}

type FireCharge struct{ sound }

type Totem struct{ sound }
