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

type CopperGolemPose struct {
	pose
}

type pose uint8

func StandingPose() CopperGolemPose {
	return CopperGolemPose{0}
}

func SittingPose() CopperGolemPose {
	return CopperGolemPose{1}
}

func RunningPose() CopperGolemPose {
	return CopperGolemPose{2}
}

func StarPose() CopperGolemPose {
	return CopperGolemPose{3}
}

func (p pose) Uint8() uint8 {
	return uint8(p)
}

func (p pose) Name() string {
	switch p {
	case 0:
		return "Standing"
	case 1:
		return "Sitting"
	case 2:
		return "Running"
	case 3:
		return "Star"
	}
	panic("unknown copper golem pose")
}

func CopperGolemPoses() []CopperGolemPose {
	return []CopperGolemPose{StandingPose(), SittingPose(), RunningPose(), StarPose()}
}
