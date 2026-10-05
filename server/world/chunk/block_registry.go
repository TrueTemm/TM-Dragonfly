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

package chunk

type BlockRegistry interface {
	BlockCount() int

	AirRuntimeID() uint32

	RuntimeIDToState(runtimeID uint32) (name string, properties map[string]any, found bool)

	StateToRuntimeID(name string, properties map[string]any) (runtimeID uint32, found bool)

	FilteringBlock(rid uint32) uint8

	LightBlock(rid uint32) uint8

	RandomTickBlock(rid uint32) bool

	NBTBlock(rid uint32) bool

	LiquidDisplacingBlock(rid uint32) bool

	LiquidBlock(rid uint32) bool

	HashToRuntimeID(hash uint32) (rid uint32, ok bool)

	RuntimeIDToHash(runtimeID uint32) (hash uint32, ok bool)
}
