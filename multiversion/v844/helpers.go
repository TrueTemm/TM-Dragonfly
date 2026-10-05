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

package v844

import "github.com/sandertv/gophertunnel/minecraft/protocol"

func GameRuleLegacy(io protocol.IO, x *protocol.GameRule) {
	io.String(&x.Name)
	io.Bool(&x.CanBeModifiedByPlayer)

	var t uint32
	switch x.Value.(type) {
	case bool:
		t = 1
	case uint32:
		t = 2
	case float32:
		t = 3
	}
	io.Varuint32(&t)

	switch t {
	case 1:
		v, _ := x.Value.(bool)
		io.Bool(&v)
		x.Value = v
	case 2:
		v, _ := x.Value.(uint32)
		io.Varuint32(&v)
		x.Value = v
	case 3:
		v, _ := x.Value.(float32)
		io.Float32(&v)
		x.Value = v
	}
}
