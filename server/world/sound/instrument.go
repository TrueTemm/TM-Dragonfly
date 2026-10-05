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

type Instrument struct {
	instrument
}

type instrument int32

func (i instrument) Int32() int32 {
	return int32(i)
}

func Piano() Instrument {
	return Instrument{0}
}

func BassDrum() Instrument {
	return Instrument{1}
}

func Snare() Instrument {
	return Instrument{2}
}

func ClicksAndSticks() Instrument {
	return Instrument{3}
}

func Bass() Instrument {
	return Instrument{4}
}

func Flute() Instrument {
	return Instrument{5}
}

func Bell() Instrument {
	return Instrument{6}
}

func Guitar() Instrument {
	return Instrument{7}
}

func Chimes() Instrument {
	return Instrument{8}
}

func Xylophone() Instrument {
	return Instrument{9}
}

func IronXylophone() Instrument {
	return Instrument{10}
}

func CowBell() Instrument {
	return Instrument{11}
}

func Didgeridoo() Instrument {
	return Instrument{12}
}

func Bit() Instrument {
	return Instrument{13}
}

func Banjo() Instrument {
	return Instrument{14}
}

func Pling() Instrument {
	return Instrument{15}
}
