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

package v776

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

type LevelSoundEvent struct {
	SoundType             uint32
	Position              mgl32.Vec3
	ExtraData             int32
	EntityType            string
	BabyMob               bool
	DisableRelativeVolume bool
}

func (*LevelSoundEvent) ID() uint32 { return v786.IDLevelSoundEvent }

func (pk *LevelSoundEvent) Marshal(io protocol.IO) {
	io.Varuint32(&pk.SoundType)
	io.Vec3(&pk.Position)
	io.Varint32(&pk.ExtraData)
	io.String(&pk.EntityType)
	io.Bool(&pk.BabyMob)
	io.Bool(&pk.DisableRelativeVolume)
}

func fromLatestLevelSoundEvent(proto uint32, pk *packet.LevelSoundEvent) *LevelSoundEvent {
	x := v786.FromLatestLevelSoundEvent786(proto, pk)
	return &LevelSoundEvent{SoundType: x.SoundType, Position: x.Position, ExtraData: x.ExtraData,
		EntityType: x.EntityType, BabyMob: x.BabyMob, DisableRelativeVolume: x.DisableRelativeVolume}
}

func toLatestLevelSoundEvent(proto uint32, pk *LevelSoundEvent) *packet.LevelSoundEvent {
	return v786.ToLatestLevelSoundEvent786(proto, &v786.LevelSoundEvent{SoundType: pk.SoundType, Position: pk.Position,
		ExtraData: pk.ExtraData, EntityType: pk.EntityType, BabyMob: pk.BabyMob,
		DisableRelativeVolume: pk.DisableRelativeVolume})
}

type ClientMovementPredictionSync struct {
	ActorFlags              protocol.Bitset
	BoundingBoxScale        float32
	BoundingBoxWidth        float32
	BoundingBoxHeight       float32
	MovementSpeed           float32
	UnderwaterMovementSpeed float32
	LavaMovementSpeed       float32
	JumpStrength            float32
	Health                  float32
	Hunger                  float32
}

func (*ClientMovementPredictionSync) ID() uint32 { return v786.IDClientMovementPredictionSync }

func (pk *ClientMovementPredictionSync) Marshal(io protocol.IO) {
	io.Bitset(&pk.ActorFlags, 120)
	io.Float32(&pk.BoundingBoxScale)
	io.Float32(&pk.BoundingBoxWidth)
	io.Float32(&pk.BoundingBoxHeight)
	io.Float32(&pk.MovementSpeed)
	io.Float32(&pk.UnderwaterMovementSpeed)
	io.Float32(&pk.LavaMovementSpeed)
	io.Float32(&pk.JumpStrength)
	io.Float32(&pk.Health)
	io.Float32(&pk.Hunger)
}

func fromLatestClientMovementPredictionSync(pk *packet.ClientMovementPredictionSync) *ClientMovementPredictionSync {
	x := v786.FromLatestClientMovementPredictionSync(pk)
	return &ClientMovementPredictionSync{ActorFlags: x.ActorFlags, BoundingBoxScale: x.BoundingBoxScale,
		BoundingBoxWidth: x.BoundingBoxWidth, BoundingBoxHeight: x.BoundingBoxHeight, MovementSpeed: x.MovementSpeed,
		UnderwaterMovementSpeed: x.UnderwaterMovementSpeed, LavaMovementSpeed: x.LavaMovementSpeed,
		JumpStrength: x.JumpStrength, Health: x.Health, Hunger: x.Hunger}
}

func toLatestClientMovementPredictionSync(pk *ClientMovementPredictionSync) *packet.ClientMovementPredictionSync {
	return v786.ToLatestClientMovementPredictionSync(&v786.ClientMovementPredictionSync{ActorFlags: pk.ActorFlags,
		BoundingBoxScale: pk.BoundingBoxScale, BoundingBoxWidth: pk.BoundingBoxWidth,
		BoundingBoxHeight: pk.BoundingBoxHeight, MovementSpeed: pk.MovementSpeed,
		UnderwaterMovementSpeed: pk.UnderwaterMovementSpeed, LavaMovementSpeed: pk.LavaMovementSpeed,
		JumpStrength: pk.JumpStrength, Health: pk.Health, Hunger: pk.Hunger})
}

type SetHud struct {
	Elements   []byte
	Visibility byte
}

func (*SetHud) ID() uint32 { return v786.IDSetHud }

func (pk *SetHud) Marshal(io protocol.IO) {
	protocol.FuncSlice(io, &pk.Elements, io.Uint8)
	io.Uint8(&pk.Visibility)
}

func fromLatestSetHud(pk *packet.SetHud) *SetHud {
	out := &SetHud{Elements: make([]byte, len(pk.Elements)), Visibility: byte(pk.Visibility)}
	for i, e := range pk.Elements {
		out.Elements[i] = byte(e)
	}
	return out
}

func toLatestSetHud(pk *SetHud) *packet.SetHud {
	out := &packet.SetHud{Elements: make([]int32, len(pk.Elements)), Visibility: int32(pk.Visibility)}
	for i, e := range pk.Elements {
		out.Elements[i] = int32(e)
	}
	return out
}
