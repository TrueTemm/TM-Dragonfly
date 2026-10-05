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

type Attack struct {
	Damage bool

	sound
}

type Drowning struct{ sound }

type Burning struct{ sound }

type Fall struct {
	Distance float64

	sound
}

type Burp struct{ sound }

type Pop struct{ sound }

type Explosion struct{ sound }

type Thunder struct{ sound }

type LevelUp struct{ sound }

type Experience struct{ sound }

type GhastWarning struct{ sound }

type GhastShoot struct{ sound }

type FireworkLaunch struct{ sound }

type FireworkHugeBlast struct{ sound }

type FireworkBlast struct{ sound }

type FireworkTwinkle struct{ sound }
