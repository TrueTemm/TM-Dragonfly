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

package v729

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v748"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

const (
	inputFlag729PerformItemInteraction  = 1 << 34
	inputFlag729PerformBlockActions     = 1 << 35
	inputFlag729PerformItemStackRequest = 1 << 36
	inputFlag729ClientPredictedVehicle  = 1 << 45

	mask729 = 1<<53 - 1
)

type PlayerAuthInput struct {
	Pitch, Yaw             float32
	Position               mgl32.Vec3
	MoveVector             mgl32.Vec2
	HeadYaw                float32
	InputData              uint64
	InputMode              uint32
	PlayMode               uint32
	InteractionModel       uint32
	GazeDirection          mgl32.Vec3
	Tick                   uint64
	Delta                  mgl32.Vec3
	ItemInteractionData    protocol.UseItemTransactionData
	ItemStackRequest       protocol.ItemStackRequest
	BlockActions           []protocol.PlayerBlockAction
	VehicleRotation        mgl32.Vec2
	ClientPredictedVehicle int64
	AnalogueMoveVector     mgl32.Vec2
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
	if pk.PlayMode == packet.PlayModeReality {
		io.Vec3(&pk.GazeDirection)
	}
	io.Varuint64(&pk.Tick)
	io.Vec3(&pk.Delta)
	if pk.InputData&inputFlag729PerformItemInteraction != 0 {
		io.PlayerInventoryAction(&pk.ItemInteractionData)
	}
	if pk.InputData&inputFlag729PerformItemStackRequest != 0 {
		v786.MarshalItemStackRequest786(io, &pk.ItemStackRequest)
	}
	if pk.InputData&inputFlag729PerformBlockActions != 0 {
		v786.PlayerBlockActions786(io, &pk.BlockActions)
	}
	if pk.InputData&inputFlag729ClientPredictedVehicle != 0 {
		io.Vec2(&pk.VehicleRotation)
		io.Varint64(&pk.ClientPredictedVehicle)
	}
	io.Vec2(&pk.AnalogueMoveVector)
}

func toLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	return v748.ToLatestPlayerAuthInput(&v748.PlayerAuthInput{
		Pitch: pk.Pitch, Yaw: pk.Yaw, Position: pk.Position, MoveVector: pk.MoveVector, HeadYaw: pk.HeadYaw,
		InputData: pk.InputData & mask729, InputMode: pk.InputMode, PlayMode: pk.PlayMode,
		InteractionModel: pk.InteractionModel,

		InteractPitch: pk.Pitch, InteractYaw: pk.HeadYaw,
		Tick: pk.Tick, Delta: pk.Delta, ItemInteractionData: pk.ItemInteractionData,
		ItemStackRequest: pk.ItemStackRequest, BlockActions: pk.BlockActions, VehicleRotation: pk.VehicleRotation,
		ClientPredictedVehicle: pk.ClientPredictedVehicle, AnalogueMoveVector: pk.AnalogueMoveVector,
	})
}

func fromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	x := v748.FromLatestPlayerAuthInput(pk)
	return &PlayerAuthInput{
		Pitch: x.Pitch, Yaw: x.Yaw, Position: x.Position, MoveVector: x.MoveVector, HeadYaw: x.HeadYaw,
		InputData: x.InputData & mask729, InputMode: x.InputMode, PlayMode: x.PlayMode,
		InteractionModel: x.InteractionModel, Tick: x.Tick, Delta: x.Delta, ItemInteractionData: x.ItemInteractionData,
		ItemStackRequest: x.ItemStackRequest, BlockActions: x.BlockActions, VehicleRotation: x.VehicleRotation,
		ClientPredictedVehicle: x.ClientPredictedVehicle, AnalogueMoveVector: x.AnalogueMoveVector,
	}
}

type InventoryContent struct {
	WindowID             uint32
	Content              []protocol.ItemInstance
	Container            protocol.FullContainerName
	DynamicContainerSize uint32
}

func (*InventoryContent) ID() uint32 { return v786.IDInventoryContent }

func (pk *InventoryContent) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	protocol.FuncSlice(io, &pk.Content, io.ItemInstance)
	protocol.Single(io, &pk.Container)
	io.Varuint32(&pk.DynamicContainerSize)
}

func fromLatestInventoryContent(pk *packet.InventoryContent) *InventoryContent {

	return &InventoryContent{WindowID: pk.WindowID, Content: pk.Content, Container: pk.Container}
}

func toLatestInventoryContent(pk *InventoryContent) *packet.InventoryContent {
	return &packet.InventoryContent{WindowID: pk.WindowID, Content: pk.Content, Container: pk.Container}
}

type InventorySlot struct {
	WindowID             uint32
	Slot                 uint32
	Container            protocol.FullContainerName
	DynamicContainerSize uint32
	NewItem              protocol.ItemInstance
}

func (*InventorySlot) ID() uint32 { return v786.IDInventorySlot }

func (pk *InventorySlot) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	io.Varuint32(&pk.Slot)
	protocol.Single(io, &pk.Container)
	io.Varuint32(&pk.DynamicContainerSize)
	io.ItemInstance(&pk.NewItem)
}

func fromLatestInventorySlot(pk *packet.InventorySlot) *InventorySlot {
	x := v786.FromLatestInventorySlot(pk)
	return &InventorySlot{WindowID: x.WindowID, Slot: x.Slot, Container: x.Container, NewItem: x.NewItem}
}

func toLatestInventorySlot(pk *InventorySlot) *packet.InventorySlot {
	return v786.ToLatestInventorySlot(&v786.InventorySlot{WindowID: pk.WindowID, Slot: pk.Slot, Container: pk.Container, NewItem: pk.NewItem})
}

type MobEffect struct {
	EntityRuntimeID uint64
	Operation       byte
	EffectType      int32
	Amplifier       int32
	Particles       bool
	Duration        int32
	Tick            uint64
}

func (*MobEffect) ID() uint32 { return v786.IDMobEffect }

func (pk *MobEffect) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.Uint8(&pk.Operation)
	io.Varint32(&pk.EffectType)
	io.Varint32(&pk.Amplifier)
	io.Bool(&pk.Particles)
	io.Varint32(&pk.Duration)
	io.Uint64(&pk.Tick)
}

func fromLatestMobEffect(pk *packet.MobEffect) *MobEffect {
	x := v786.FromLatestMobEffect(pk)
	return &MobEffect{EntityRuntimeID: x.EntityRuntimeID, Operation: x.Operation, EffectType: x.EffectType,
		Amplifier: x.Amplifier, Particles: x.Particles, Duration: x.Duration, Tick: x.Tick}
}

func toLatestMobEffect(pk *MobEffect) *packet.MobEffect {
	return v786.ToLatestMobEffect(&v786.MobEffect{EntityRuntimeID: pk.EntityRuntimeID, Operation: pk.Operation,
		EffectType: pk.EffectType, Amplifier: pk.Amplifier, Particles: pk.Particles, Duration: pk.Duration, Tick: pk.Tick})
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
}

type PackURL struct {
	UUIDVersion string
	URL         string
}

func (x *PackURL) Marshal(r protocol.IO) {
	r.String(&x.UUIDVersion)
	r.String(&x.URL)
}

type ResourcePacksInfo struct {
	TexturePackRequired bool
	HasAddons           bool
	HasScripts          bool
	TexturePacks        []TexturePackInfo
	PackURLs            []PackURL
}

func (*ResourcePacksInfo) ID() uint32 { return v786.IDResourcePacksInfo }

func (pk *ResourcePacksInfo) Marshal(io protocol.IO) {
	io.Bool(&pk.TexturePackRequired)
	io.Bool(&pk.HasAddons)
	io.Bool(&pk.HasScripts)
	v786.SliceUint16Length786(io, &pk.TexturePacks)
	protocol.Slice(io, &pk.PackURLs)
}

func fromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	x := v748.FromLatestResourcePacksInfo(pk)
	out := &ResourcePacksInfo{TexturePackRequired: x.TexturePackRequired, HasAddons: x.HasAddons,
		HasScripts: x.HasScripts, TexturePacks: make([]TexturePackInfo, len(x.TexturePacks)), PackURLs: []PackURL{}}
	for i, p := range x.TexturePacks {
		out.TexturePacks[i] = TexturePackInfo{UUID: p.UUID, Version: p.Version, Size: p.Size, ContentKey: p.ContentKey,
			SubPackName: p.SubPackName, ContentIdentity: p.ContentIdentity, HasScripts: p.HasScripts,
			AddonPack: p.AddonPack, RTXEnabled: p.RTXEnabled}

		if p.DownloadURL != "" {
			out.PackURLs = append(out.PackURLs, PackURL{UUIDVersion: p.UUID + "_" + p.Version, URL: p.DownloadURL})
		}
	}
	return out
}
