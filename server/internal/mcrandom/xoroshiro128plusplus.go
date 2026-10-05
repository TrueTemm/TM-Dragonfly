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

package mcrandom

import "math/bits"

type Xoroshiro128PlusPlus struct {
	seed0, seed1 uint64
}

func NewXoroshiro128PlusPlus(seed0, seed1 uint64) *Xoroshiro128PlusPlus {
	return &Xoroshiro128PlusPlus{seed0, seed1}
}

func (x *Xoroshiro128PlusPlus) Next() uint64 {
	s0 := x.seed0
	s1 := x.seed1
	result := bits.RotateLeft64(s0+s1, 17) + s0
	s1 ^= s0
	x.seed0 = bits.RotateLeft64(s0, 49) ^ s1 ^ (s1 << 21)
	x.seed1 = bits.RotateLeft64(s1, 28)
	return result
}
