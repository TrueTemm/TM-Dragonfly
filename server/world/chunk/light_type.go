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

var (
	SkyLight skyLight

	BlockLight blockLight
)

type (
	light interface {
		light(sub *SubChunk, x, y, z uint8) uint8
		setLight(sub *SubChunk, x, y, z, v uint8)
	}
	skyLight   struct{}
	blockLight struct{}
)

func (skyLight) light(sub *SubChunk, x, y, z uint8) uint8   { return sub.SkyLight(x, y, z) }
func (skyLight) setLight(sub *SubChunk, x, y, z, v uint8)   { sub.SetSkyLight(x, y, z, v) }
func (blockLight) light(sub *SubChunk, x, y, z uint8) uint8 { return sub.BlockLight(x, y, z) }
func (blockLight) setLight(sub *SubChunk, x, y, z, v uint8) { sub.SetBlockLight(x, y, z, v) }
