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

package v748

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

const (
	inputFlag748PerformItemInteraction  = 1 << 34
	inputFlag748PerformBlockActions     = 1 << 35
	inputFlag748PerformItemStackRequest = 1 << 36
	inputFlag748ClientPredictedVehicle  = 1 << 45

	insertedAt = 53

	highestBit = 56
)

type PlayerAuthInput struct {
	Pitch, Yaw                 float32
	Position                   mgl32.Vec3
	MoveVector                 mgl32.Vec2
	HeadYaw                    float32
	InputData                  uint64
	InputMode                  uint32
	PlayMode                   uint32
	InteractionModel           uint32
	InteractPitch, InteractYaw float32
	Tick                       uint64
	Delta                      mgl32.Vec3
	ItemInteractionData        protocol.UseItemTransactionData
	ItemStackRequest           protocol.ItemStackRequest
	BlockActions               []protocol.PlayerBlockAction
	VehicleRotation            mgl32.Vec2
	ClientPredictedVehicle     int64
	AnalogueMoveVector         mgl32.Vec2
	CameraOrientation          mgl32.Vec3
}

func (*PlayerAuthInput) ID() uint32 { return v786.IDPlayerAuthInput }

func (pk *PlayerAuthInput) Marshal(io protocol.IO) {
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Vec3(&pk.Position)
	io.Vec2(&pk.MoveVector)
	io.Float32(&pk.HeadYaw)
	io.Varuint64(&pk.InputData)
	io.Varuint32(&pk.InputMode)
	io.Varuint32(&pk.PlayMode)
	io.Varuint32(&pk.InteractionModel)
	io.Float32(&pk.InteractPitch)
	io.Float32(&pk.InteractYaw)
	io.Varuint64(&pk.Tick)
	io.Vec3(&pk.Delta)
	if pk.InputData&inputFlag748PerformItemInteraction != 0 {
		io.PlayerInventoryAction(&pk.ItemInteractionData)
	}
	if pk.InputData&inputFlag748PerformItemStackRequest != 0 {
		v786.MarshalItemStackRequest786(io, &pk.ItemStackRequest)
	}
	if pk.InputData&inputFlag748PerformBlockActions != 0 {
		v786.PlayerBlockActions786(io, &pk.BlockActions)
	}
	if pk.InputData&inputFlag748ClientPredictedVehicle != 0 {
		io.Vec2(&pk.VehicleRotation)
		io.Varint64(&pk.ClientPredictedVehicle)
	}
	io.Vec2(&pk.AnalogueMoveVector)
	io.Vec3(&pk.CameraOrientation)
}

func bitsetFromMask(mask uint64) protocol.Bitset {
	b := protocol.NewBitset(v786.PlayerAuthInputBitsetSize)
	for bit := 0; bit <= highestBit; bit++ {
		if mask&(1<<uint(bit)) == 0 {
			continue
		}
		idx := bit
		if bit >= insertedAt {
			idx++
		}
		b.Set(idx)
	}
	return b
}

func maskFromBitset(b protocol.Bitset) uint64 {
	var mask uint64
	for idx := 0; idx < v786.PlayerAuthInputBitsetSize; idx++ {
		if !b.Load(idx) || idx == insertedAt {
			continue
		}
		bit := idx
		if idx > insertedAt {
			bit--
		}
		if bit > highestBit {
			continue
		}
		mask |= 1 << uint(bit)
	}
	return mask
}

func toLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	flags := bitsetFromMask(pk.InputData)
	if flags.Load(v786.InputFlagJumpDown) {

		flags.Set(v786.InputFlagJumpCurrentRaw)
	}
	return v786.ToLatestPlayerAuthInput(&v786.PlayerAuthInput{
		Pitch: pk.Pitch, Yaw: pk.Yaw, Position: pk.Position, MoveVector: pk.MoveVector, HeadYaw: pk.HeadYaw,
		InputData: flags, InputMode: pk.InputMode, PlayMode: pk.PlayMode,
		InteractionModel: pk.InteractionModel, InteractPitch: pk.InteractPitch, InteractYaw: pk.InteractYaw,
		Tick: pk.Tick, Delta: pk.Delta, ItemInteractionData: pk.ItemInteractionData,
		ItemStackRequest: pk.ItemStackRequest, BlockActions: pk.BlockActions, VehicleRotation: pk.VehicleRotation,
		ClientPredictedVehicle: pk.ClientPredictedVehicle, AnalogueMoveVector: pk.AnalogueMoveVector,
		CameraOrientation: pk.CameraOrientation,

		RawMoveVector: pk.AnalogueMoveVector,
	})
}

func fromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	x := v786.FromLatestPlayerAuthInput(pk)
	return &PlayerAuthInput{
		Pitch: x.Pitch, Yaw: x.Yaw, Position: x.Position, MoveVector: x.MoveVector, HeadYaw: x.HeadYaw,
		InputData: maskFromBitset(x.InputData), InputMode: x.InputMode, PlayMode: x.PlayMode,
		InteractionModel: x.InteractionModel, InteractPitch: x.InteractPitch, InteractYaw: x.InteractYaw,
		Tick: x.Tick, Delta: x.Delta, ItemInteractionData: x.ItemInteractionData, ItemStackRequest: x.ItemStackRequest,
		BlockActions: x.BlockActions, VehicleRotation: x.VehicleRotation, ClientPredictedVehicle: x.ClientPredictedVehicle,
		AnalogueMoveVector: x.AnalogueMoveVector, CameraOrientation: x.CameraOrientation,
	}
}

const (
	armourFlag748Helmet     = 1 << 1
	armourFlag748Chestplate = 1 << 2
	armourFlag748Leggings   = 1 << 3
	armourFlag748Boots      = 1 << 4
	armourFlag748Body       = 1 << 5
)

type PlayerArmourDamage struct {
	Bitset           uint8
	HelmetDamage     int32
	ChestplateDamage int32
	LeggingsDamage   int32
	BootsDamage      int32
	BodyDamage       int32
}

func (*PlayerArmourDamage) ID() uint32 { return v786.IDPlayerArmourDamage }

func (pk *PlayerArmourDamage) Marshal(io protocol.IO) {
	io.Uint8(&pk.Bitset)
	if pk.Bitset&armourFlag748Helmet != 0 {
		io.Varint32(&pk.HelmetDamage)
	}
	if pk.Bitset&armourFlag748Chestplate != 0 {
		io.Varint32(&pk.ChestplateDamage)
	}
	if pk.Bitset&armourFlag748Leggings != 0 {
		io.Varint32(&pk.LeggingsDamage)
	}
	if pk.Bitset&armourFlag748Boots != 0 {
		io.Varint32(&pk.BootsDamage)
	}
	if pk.Bitset&armourFlag748Body != 0 {
		io.Varint32(&pk.BodyDamage)
	}
}

func fromLatestPlayerArmourDamage(pk *packet.PlayerArmourDamage) *PlayerArmourDamage {

	x := v786.FromLatestPlayerArmourDamage(pk)
	return &PlayerArmourDamage{Bitset: x.Bitset << 1, HelmetDamage: x.HelmetDamage,
		ChestplateDamage: x.ChestplateDamage, LeggingsDamage: x.LeggingsDamage, BootsDamage: x.BootsDamage,
		BodyDamage: x.BodyDamage}
}

func toLatestPlayerArmourDamage(pk *PlayerArmourDamage) *packet.PlayerArmourDamage {
	return v786.ToLatestPlayerArmourDamage(&v786.PlayerArmourDamage{Bitset: pk.Bitset >> 1, HelmetDamage: pk.HelmetDamage,
		ChestplateDamage: pk.ChestplateDamage, LeggingsDamage: pk.LeggingsDamage, BootsDamage: pk.BootsDamage,
		BodyDamage: pk.BodyDamage})
}

type TexturePackInfo struct {
	UUID            string
	Version         string
	Size            uint64
	ContentKey      string
	SubPackName     string
	ContentIdentity string
	HasScripts      bool
	AddonPack       bool
	RTXEnabled      bool
	DownloadURL     string
}

func (x *TexturePackInfo) Marshal(r protocol.IO) {
	r.String(&x.UUID)
	r.String(&x.Version)
	r.Uint64(&x.Size)
	r.String(&x.ContentKey)
	r.String(&x.SubPackName)
	r.String(&x.ContentIdentity)
	r.Bool(&x.HasScripts)
	r.Bool(&x.AddonPack)
	r.Bool(&x.RTXEnabled)
	r.String(&x.DownloadURL)
}

type ResourcePacksInfo struct {
	TexturePackRequired bool
	HasAddons           bool
	HasScripts          bool
	TexturePacks        []TexturePackInfo
}

func (*ResourcePacksInfo) ID() uint32 { return v786.IDResourcePacksInfo }

func (pk *ResourcePacksInfo) Marshal(io protocol.IO) {
	io.Bool(&pk.TexturePackRequired)
	io.Bool(&pk.HasAddons)
	io.Bool(&pk.HasScripts)
	v786.SliceUint16Length786(io, &pk.TexturePacks)
}

func fromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	x := v786.FromLatestResourcePacksInfo786(pk)
	out := &ResourcePacksInfo{TexturePackRequired: x.TexturePackRequired, HasAddons: x.HasAddons,
		HasScripts: x.HasScripts, TexturePacks: make([]TexturePackInfo, len(x.TexturePacks))}
	for i, p := range x.TexturePacks {
		out.TexturePacks[i] = TexturePackInfo{UUID: p.UUID.String(), Version: p.Version, Size: p.Size,
			ContentKey: p.ContentKey, SubPackName: p.SubPackName, ContentIdentity: p.ContentIdentity,
			HasScripts: p.HasScripts, AddonPack: p.AddonPack, RTXEnabled: p.RTXEnabled, DownloadURL: p.DownloadURL}
	}
	return out
}
