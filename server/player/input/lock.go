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

package input

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type Lock struct {
	lock
}

type lock uint32

func Camera() Lock {
	return Lock{lock(packet.ClientInputLockCamera)}
}

func Movement() Lock {
	return Lock{lock(packet.ClientInputLockMovement)}
}

func LateralMovement() Lock {
	return Lock{lock(packet.ClientInputLockLateralMovement)}
}

func Sneak() Lock {
	return Lock{lock(packet.ClientInputLockSneak)}
}

func Jump() Lock {
	return Lock{lock(packet.ClientInputLockJump)}
}

func Mount() Lock {
	return Lock{lock(packet.ClientInputLockMount)}
}

func Dismount() Lock {
	return Lock{lock(packet.ClientInputLockDismount)}
}

func MoveForward() Lock {
	return Lock{lock(packet.ClientInputLockMoveForward)}
}

func MoveBackward() Lock {
	return Lock{lock(packet.ClientInputLockMoveBackward)}
}

func MoveLeft() Lock {
	return Lock{lock(packet.ClientInputLockMoveLeft)}
}

func MoveRight() Lock {
	return Lock{lock(packet.ClientInputLockMoveRight)}
}

func (l lock) Uint32() uint32 {
	return uint32(l)
}

func All() []Lock {
	return []Lock{
		Camera(), Movement(), LateralMovement(), Sneak(), Jump(), Mount(), Dismount(),
		MoveForward(), MoveBackward(), MoveLeft(), MoveRight(),
	}
}
