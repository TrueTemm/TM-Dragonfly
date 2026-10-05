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

package bossbar

type Colour struct{ colour }

func Pink() Colour {
	return Colour{colour(0)}
}

func Blue() Colour {
	return Colour{colour(1)}
}

func Red() Colour {
	return Colour{colour(2)}
}

func Green() Colour {
	return Colour{colour(3)}
}

func Yellow() Colour {
	return Colour{colour(4)}
}

func Purple() Colour {
	return Colour{colour(5)}
}

func RebeccaPurple() Colour {
	return Colour{colour(6)}
}

func White() Colour {
	return Colour{colour(7)}
}

type colour uint8

func (c colour) Uint8() uint8 {
	return uint8(c)
}
