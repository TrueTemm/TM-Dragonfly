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

package v686

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v786"
)

const (
	inputFlag686PerformItemInteraction  = 1 << 34
	inputFlag686PerformBlockActions     = 1 << 35
	inputFlag686PerformItemStackRequest = 1 << 36
	inputFlag686ClientPredictedVehicle  = 1 << 45
)

type PlayerAuthInput struct {
	v729.PlayerAuthInput
}

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
	if pk.InputData&inputFlag686PerformItemInteraction != 0 {
		io.PlayerInventoryAction(&pk.ItemInteractionData)
	}
	if pk.InputData&inputFlag686PerformItemStackRequest != 0 {
		v786.MarshalItemStackRequest786(io, &pk.ItemStackRequest)
	}
	if pk.InputData&inputFlag686ClientPredictedVehicle != 0 {
		io.Vec2(&pk.VehicleRotation)
		io.Varint64(&pk.ClientPredictedVehicle)
	}
	if pk.InputData&inputFlag686PerformBlockActions != 0 {
		v786.PlayerBlockActions786(io, &pk.BlockActions)
	}
	io.Vec2(&pk.AnalogueMoveVector)
}

type PlaySound struct {
	SoundName string
	Position  mgl32.Vec3
	Volume    float32
	Pitch     float32
}

func (*PlaySound) ID() uint32 { return v786.IDPlaySound }

func (pk *PlaySound) Marshal(io protocol.IO) {
	io.String(&pk.SoundName)
	b := protocol.BlockPos{int32(pk.Position[0] * 8), int32(pk.Position[1] * 8), int32(pk.Position[2] * 8)}
	io.BlockPos(&b)
	pk.Position = mgl32.Vec3{float32(b[0]) / 8, float32(b[1]) / 8, float32(b[2]) / 8}
	io.Float32(&pk.Volume)
	io.Float32(&pk.Pitch)
}

func fromLatestPlaySound(pk *packet.PlaySound) *PlaySound {
	x := v786.FromLatestPlaySound(pk)
	return &PlaySound{SoundName: x.SoundName, Position: x.Position, Volume: x.Volume, Pitch: x.Pitch}
}

type ChangeDimension struct {
	Dimension int32
	Position  mgl32.Vec3
	Respawn   bool
}

func (*ChangeDimension) ID() uint32 { return v786.IDChangeDimension }

func (pk *ChangeDimension) Marshal(io protocol.IO) {
	io.Varint32(&pk.Dimension)
	io.Vec3(&pk.Position)
	io.Bool(&pk.Respawn)
}

type Disconnect struct {
	Reason                  int32
	HideDisconnectionScreen bool
	Message                 string
}

func (*Disconnect) ID() uint32 { return v786.IDDisconnect }

func (pk *Disconnect) Marshal(io protocol.IO) {
	io.Varint32(&pk.Reason)
	io.Bool(&pk.HideDisconnectionScreen)
	if !pk.HideDisconnectionScreen {
		io.String(&pk.Message)
	}
}

type SetTitle struct {
	ActionType       int32
	Text             string
	FadeInDuration   int32
	RemainDuration   int32
	FadeOutDuration  int32
	XUID             string
	PlatformOnlineID string
}

func (*SetTitle) ID() uint32 { return v786.IDSetTitle }

func (pk *SetTitle) Marshal(io protocol.IO) {
	io.Varint32(&pk.ActionType)
	io.String(&pk.Text)
	io.Varint32(&pk.FadeInDuration)
	io.Varint32(&pk.RemainDuration)
	io.Varint32(&pk.FadeOutDuration)
	io.String(&pk.XUID)
	io.String(&pk.PlatformOnlineID)
}

type StopSound struct {
	SoundName string
	StopAll   bool
}

func (*StopSound) ID() uint32 { return v786.IDStopSound }

func (pk *StopSound) Marshal(io protocol.IO) {
	io.String(&pk.SoundName)
	io.Bool(&pk.StopAll)
}

type MobArmourEquipment struct {
	EntityRuntimeID uint64
	Helmet          protocol.ItemInstance
	Chestplate      protocol.ItemInstance
	Leggings        protocol.ItemInstance
	Boots           protocol.ItemInstance
}

func (*MobArmourEquipment) ID() uint32 { return v786.IDMobArmourEquipment }

func (pk *MobArmourEquipment) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.ItemInstance(&pk.Helmet)
	io.ItemInstance(&pk.Chestplate)
	io.ItemInstance(&pk.Leggings)
	io.ItemInstance(&pk.Boots)
}

type PlayerArmourDamage struct {
	Bitset           uint8
	HelmetDamage     int32
	ChestplateDamage int32
	LeggingsDamage   int32
	BootsDamage      int32
}

func (*PlayerArmourDamage) ID() uint32 { return v786.IDPlayerArmourDamage }

func (pk *PlayerArmourDamage) Marshal(io protocol.IO) {
	io.Uint8(&pk.Bitset)
	if pk.Bitset&0b0001 != 0 {
		io.Varint32(&pk.HelmetDamage)
	}
	if pk.Bitset&0b0010 != 0 {
		io.Varint32(&pk.ChestplateDamage)
	}
	if pk.Bitset&0b0100 != 0 {
		io.Varint32(&pk.LeggingsDamage)
	}
	if pk.Bitset&0b1000 != 0 {
		io.Varint32(&pk.BootsDamage)
	}
}

func fromLatestPlayerArmourDamage(pk *packet.PlayerArmourDamage) *PlayerArmourDamage {

	x := v786.FromLatestPlayerArmourDamage(pk)
	return &PlayerArmourDamage{Bitset: x.Bitset & 0b1111, HelmetDamage: x.HelmetDamage,
		ChestplateDamage: x.ChestplateDamage, LeggingsDamage: x.LeggingsDamage, BootsDamage: x.BootsDamage}
}

type InventoryContent struct {
	WindowID uint32
	Content  []protocol.ItemInstance
}

func (*InventoryContent) ID() uint32 { return v786.IDInventoryContent }

func (pk *InventoryContent) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	protocol.FuncSlice(io, &pk.Content, io.ItemInstance)
}

type InventorySlot struct {
	WindowID uint32
	Slot     uint32
	NewItem  protocol.ItemInstance
}

func (*InventorySlot) ID() uint32 { return v786.IDInventorySlot }

func (pk *InventorySlot) Marshal(io protocol.IO) {
	io.Varuint32(&pk.WindowID)
	io.Varuint32(&pk.Slot)
	io.ItemInstance(&pk.NewItem)
}

type TexturePackInfo struct {
	UUID            string
	Version         string
	Size            uint64
	ContentKey      string
	SubPackName     string
	ContentIdentity string
	HasScripts      bool
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
	r.Bool(&x.RTXEnabled)
}

type BehaviourPackInfo struct {
	UUID            string
	Version         string
	Size            uint64
	ContentKey      string
	SubPackName     string
	ContentIdentity string
	HasScripts      bool
}

func (x *BehaviourPackInfo) Marshal(r protocol.IO) {
	r.String(&x.UUID)
	r.String(&x.Version)
	r.Uint64(&x.Size)
	r.String(&x.ContentKey)
	r.String(&x.SubPackName)
	r.String(&x.ContentIdentity)
	r.Bool(&x.HasScripts)
}

type ResourcePacksInfo struct {
	TexturePackRequired bool
	HasAddons           bool
	HasScripts          bool
	ForcingServerPacks  bool
	BehaviourPacks      []BehaviourPackInfo
	TexturePacks        []TexturePackInfo
	PackURLs            []v729.PackURL
}

func (*ResourcePacksInfo) ID() uint32 { return v786.IDResourcePacksInfo }

func (pk *ResourcePacksInfo) Marshal(io protocol.IO) {
	io.Bool(&pk.TexturePackRequired)
	io.Bool(&pk.HasAddons)
	io.Bool(&pk.HasScripts)
	io.Bool(&pk.ForcingServerPacks)
	v786.SliceUint16Length786(io, &pk.BehaviourPacks)
	v786.SliceUint16Length786(io, &pk.TexturePacks)
	protocol.Slice(io, &pk.PackURLs)
}

func fromLatestResourcePacksInfo(pk *packet.ResourcePacksInfo) *ResourcePacksInfo {
	x := v729.FromLatestResourcePacksInfo(pk)
	out := &ResourcePacksInfo{TexturePackRequired: x.TexturePackRequired, HasAddons: x.HasAddons, HasScripts: x.HasScripts,
		BehaviourPacks: []BehaviourPackInfo{}, TexturePacks: make([]TexturePackInfo, len(x.TexturePacks)), PackURLs: x.PackURLs}
	for i, p := range x.TexturePacks {
		out.TexturePacks[i] = TexturePackInfo{UUID: p.UUID, Version: p.Version, Size: p.Size, ContentKey: p.ContentKey,
			SubPackName: p.SubPackName, ContentIdentity: p.ContentIdentity, HasScripts: p.HasScripts, RTXEnabled: p.RTXEnabled}
	}
	return out
}
