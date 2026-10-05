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

package world

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type ChunkPos [2]int32

func (p ChunkPos) String() string {
	return fmt.Sprintf("(%v, %v)", p[0], p[1])
}

func (p ChunkPos) X() int32 {
	return p[0]
}

func (p ChunkPos) Z() int32 {
	return p[1]
}

type SubChunkPos [3]int32

func (p SubChunkPos) String() string {
	return fmt.Sprintf("(%v, %v, %v)", p[0], p[1], p[2])
}

func (p SubChunkPos) X() int32 {
	return p[0]
}

func (p SubChunkPos) Y() int32 {
	return p[1]
}

func (p SubChunkPos) Z() int32 {
	return p[2]
}

func chunkPosFromVec3(vec3 mgl64.Vec3) ChunkPos {
	return ChunkPos{int32(math.Floor(vec3[0])) >> 4, int32(math.Floor(vec3[2])) >> 4}
}

func chunkPosFromBlockPos(p cube.Pos) ChunkPos {
	return ChunkPos{int32(p[0] >> 4), int32(p[2] >> 4)}
}
