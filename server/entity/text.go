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

package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func NewText(text string, pos mgl64.Vec3) *world.EntityHandle {
	return world.EntitySpawnOpts{Position: pos, NameTag: text}.New(TextType, textConf)
}

var textConf StationaryBehaviourConfig

var TextType textType

type textType struct{}

func (t textType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}
func (textType) EncodeEntity() string        { return "tm-dragonfly:text" }
func (textType) BBox(world.Entity) cube.BBox { return cube.BBox{} }
func (textType) NetworkEncodeEntity() string { return "minecraft:falling_block" }

func (textType) DecodeNBT(_ map[string]any, data *world.EntityData) { data.Data = textConf.New() }
func (textType) EncodeNBT(_ *world.EntityData) map[string]any       { return nil }
