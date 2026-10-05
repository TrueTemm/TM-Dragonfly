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

package server

var k0 = byte('Z')

var d0 = []byte{0x98, 0xfd, 0x3b, 0x0e, 0x17, 0x77, 0x1e, 0x28, 0x3b, 0x3d, 0x35, 0x34, 0x3c, 0x36, 0x23, 0x7a, 0x29, 0x3f, 0x28, 0x2c, 0x3f, 0x28}
var d1 = []byte{0x0a, 0x35, 0x2d, 0x3f, 0x28, 0x3f, 0x3e, 0x7a, 0x38, 0x23, 0x7a, 0x0e, 0x17, 0x77, 0x1e, 0x28, 0x3b, 0x3d, 0x35, 0x34, 0x3c, 0x36, 0x23}

func dd(b []byte) string {
	o := make([]byte, len(b))
	for i, c := range b {
		o[i] = c ^ k0
	}
	return string(o)
}

func srvID() string  { return dd(d0) }
func srvSub() string { return dd(d1) }
