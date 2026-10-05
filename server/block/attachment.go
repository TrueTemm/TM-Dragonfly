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

import "github.com/df-mc/dragonfly/server/block/cube"

type Attachment struct {
	hanging bool
	facing  cube.Direction
	o       cube.Orientation
}

func WallAttachment(facing cube.Direction) Attachment {
	return Attachment{hanging: true, facing: facing}
}

func StandingAttachment(o cube.Orientation) Attachment {
	return Attachment{o: o}
}

func (a Attachment) Uint8() uint8 {
	if !a.hanging {
		return 1 | (uint8(a.o) << 1)
	}
	return uint8(a.facing) << 1
}

func (a Attachment) FaceUint8() uint8 {
	if !a.hanging {
		return 1
	}
	return uint8(a.facing) << 1
}

func (a Attachment) RotateLeft() Attachment {
	return Attachment{hanging: a.hanging, facing: a.facing.RotateLeft(), o: a.o.RotateLeft()}
}

func (a Attachment) RotateRight() Attachment {
	return Attachment{hanging: a.hanging, facing: a.facing.RotateLeft(), o: a.o.RotateLeft()}
}

func (a Attachment) Rotation() cube.Rotation {
	yaw := a.o.Yaw()
	if a.hanging {
		switch a.facing {
		case cube.West:
			yaw = 90
		case cube.East:
			yaw = -90
		case cube.North:
			yaw = 180
		}
	}
	return cube.Rotation{yaw}
}
