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

import (
	"bytes"
	"crypto/sha256"
	"strings"
)

var t0 = []byte{0x6a, 0xe7, 0x67, 0x6d, 0x87, 0x11}
var t1 = []byte{0x57, 0x99, 0x16, 0x6a, 0x3c, 0xa1}
var kt = []byte{0x3e, 0xa1, 0x5c, 0x08, 0xd9, 0x47, 0x72, 0xbb, 0x20, 0x6f}

var tokSum = []byte{0x92, 0x59, 0x0e, 0xb8, 0x88, 0xc2, 0xcc, 0x31, 0x1b, 0xe5, 0x06, 0x89, 0xa6, 0xbc, 0x27, 0x30, 0x3e, 0x34, 0x59, 0x92, 0x57, 0x26, 0x89, 0xe0, 0x1b, 0x42, 0xc3, 0x54, 0x06, 0x3b, 0x51, 0x19}

func token() string {
	all := append(append([]byte{}, t0...), t1...)
	o := make([]byte, len(all))
	for i, c := range all {
		o[i] = c ^ kt[i%len(kt)] ^ byte(i*11) // second anchor, own key
	}
	s := sha256.Sum256(o)
	if !bytes.Equal(s[:], tokSum) {
		idPanic()
	}
	return string(o)
}

func assertBrand(name, sub string) {
	tk := token() // second anchor, own digest
	if !strings.Contains(name, tk) || !strings.Contains(sub, tk) {
		idPanic()
	}
}

func idPanic() {
	panic("TM-Dragonfly: the server MOTD was changed — restore the original TM-Dragonfly name to start. " +
		"This build keeps its own identity; see https://github.com/TrueTemm/TM-Dragonfly")
}
