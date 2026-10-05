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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func toLatestActorEvent(pk *ActorEvent) *packet.ActorEvent {
	return &packet.ActorEvent{EntityRuntimeID: pk.EntityRuntimeID, EventType: pk.EventType, EventData: pk.EventData}
}
func fromLatestActorEvent(pk *packet.ActorEvent) *ActorEvent {
	return &ActorEvent{EntityRuntimeID: pk.EntityRuntimeID, EventType: pk.EventType, EventData: pk.EventData}
}

func toLatestAddActor(pk *AddActor) *packet.AddActor {
	return &packet.AddActor{
		EntityUniqueID: pk.EntityUniqueID, EntityRuntimeID: pk.EntityRuntimeID, EntityType: pk.EntityType,
		Position: pk.Position, Velocity: pk.Velocity, Pitch: pk.Pitch, Yaw: pk.Yaw, HeadYaw: pk.HeadYaw,
		BodyYaw: pk.BodyYaw, Attributes: pk.Attributes, EntityMetadata: pk.EntityMetadata,
		EntityProperties: pk.EntityProperties, EntityLinks: pk.EntityLinks,
	}
}
func fromLatestAddActor(pk *packet.AddActor) *AddActor {
	return &AddActor{
		EntityUniqueID: pk.EntityUniqueID, EntityRuntimeID: pk.EntityRuntimeID, EntityType: pk.EntityType,
		Position: pk.Position, Velocity: pk.Velocity, Pitch: pk.Pitch, Yaw: pk.Yaw, HeadYaw: pk.HeadYaw,
		BodyYaw: pk.BodyYaw, Attributes: pk.Attributes, EntityMetadata: pk.EntityMetadata,
		EntityProperties: pk.EntityProperties, EntityLinks: pk.EntityLinks,
	}
}

func toLatestAddItemActor(pk *AddItemActor) *packet.AddItemActor {
	return &packet.AddItemActor{
		EntityUniqueID: pk.EntityUniqueID, EntityRuntimeID: pk.EntityRuntimeID, Item: pk.Item,
		Position: pk.Position, Velocity: pk.Velocity, EntityMetadata: pk.EntityMetadata, FromFishing: pk.FromFishing,
	}
}
func fromLatestAddItemActor(pk *packet.AddItemActor) *AddItemActor {
	return &AddItemActor{
		EntityUniqueID: pk.EntityUniqueID, EntityRuntimeID: pk.EntityRuntimeID, Item: pk.Item,
		Position: pk.Position, Velocity: pk.Velocity, EntityMetadata: pk.EntityMetadata, FromFishing: pk.FromFishing,
	}
}

func toLatestAddPlayer(pk *AddPlayer) *packet.AddPlayer {
	return &packet.AddPlayer{
		UUID: pk.UUID, Username: pk.Username, EntityRuntimeID: pk.EntityRuntimeID, PlatformChatID: pk.PlatformChatID,
		Position: pk.Position, Velocity: pk.Velocity, Pitch: pk.Pitch, Yaw: pk.Yaw, HeadYaw: pk.HeadYaw,
		HeldItem: pk.HeldItem, GameType: pk.GameType, EntityMetadata: pk.EntityMetadata,
		EntityProperties: pk.EntityProperties, AbilityData: pk.AbilityData, EntityLinks: pk.EntityLinks,
		DeviceID: pk.DeviceID, BuildPlatform: pk.BuildPlatform,
	}
}
func fromLatestAddPlayer(pk *packet.AddPlayer) *AddPlayer {
	return &AddPlayer{
		UUID: pk.UUID, Username: pk.Username, EntityRuntimeID: pk.EntityRuntimeID, PlatformChatID: pk.PlatformChatID,
		Position: pk.Position, Velocity: pk.Velocity, Pitch: pk.Pitch, Yaw: pk.Yaw, HeadYaw: pk.HeadYaw,
		HeldItem: pk.HeldItem, GameType: pk.GameType, EntityMetadata: pk.EntityMetadata,
		EntityProperties: pk.EntityProperties, AbilityData: pk.AbilityData, EntityLinks: pk.EntityLinks,
		DeviceID: pk.DeviceID, BuildPlatform: pk.BuildPlatform,
	}
}

func toLatestAddVolumeEntity(pk *AddVolumeEntity) *packet.AddVolumeEntity {
	return &packet.AddVolumeEntity{
		EntityRuntimeID: uint32(pk.EntityRuntimeID), EntityMetadata: pk.EntityMetadata,
		EncodingIdentifier: pk.EncodingIdentifier, InstanceIdentifier: pk.InstanceIdentifier,
		Bounds: pk.Bounds, Dimension: pk.Dimension, EngineVersion: pk.EngineVersion,
	}
}
func fromLatestAddVolumeEntity(pk *packet.AddVolumeEntity) *AddVolumeEntity {
	return &AddVolumeEntity{
		EntityRuntimeID: uint64(pk.EntityRuntimeID), EntityMetadata: pk.EntityMetadata,
		EncodingIdentifier: pk.EncodingIdentifier, InstanceIdentifier: pk.InstanceIdentifier,
		Bounds: pk.Bounds, Dimension: pk.Dimension, EngineVersion: pk.EngineVersion,
	}
}

func toLatestRemoveVolumeEntity(pk *RemoveVolumeEntity) *packet.RemoveVolumeEntity {
	return &packet.RemoveVolumeEntity{EntityRuntimeID: uint32(pk.EntityRuntimeID), Dimension: pk.Dimension}
}
func fromLatestRemoveVolumeEntity(pk *packet.RemoveVolumeEntity) *RemoveVolumeEntity {
	return &RemoveVolumeEntity{EntityRuntimeID: uint64(pk.EntityRuntimeID), Dimension: pk.Dimension}
}

func toLatestSetActorData(pk *SetActorData) *packet.SetActorData {
	return &packet.SetActorData{
		EntityRuntimeID: pk.EntityRuntimeID, EntityMetadata: pk.EntityMetadata,
		EntityProperties: pk.EntityProperties, Tick: pk.Tick,
	}
}
func fromLatestSetActorData(pk *packet.SetActorData) *SetActorData {
	return &SetActorData{
		EntityRuntimeID: pk.EntityRuntimeID, EntityMetadata: pk.EntityMetadata,
		EntityProperties: pk.EntityProperties, Tick: pk.Tick,
	}
}

func toLatestHurtArmour(pk *HurtArmour) *packet.HurtArmour {
	return &packet.HurtArmour{Cause: pk.Cause, Damage: pk.Damage, ArmourSlots: uint64(pk.ArmourSlots)}
}
func fromLatestHurtArmour(pk *packet.HurtArmour) *HurtArmour {
	return &HurtArmour{Cause: pk.Cause, Damage: pk.Damage, ArmourSlots: int64(pk.ArmourSlots)}
}

func toLatestMobEffect(pk *MobEffect) *packet.MobEffect {
	return &packet.MobEffect{
		EntityRuntimeID: pk.EntityRuntimeID, Operation: pk.Operation, EffectType: pk.EffectType,
		Amplifier: pk.Amplifier, Particles: pk.Particles, Duration: pk.Duration, Tick: pk.Tick,
	}
}
func fromLatestMobEffect(pk *packet.MobEffect) *MobEffect {
	return &MobEffect{
		EntityRuntimeID: pk.EntityRuntimeID, Operation: pk.Operation, EffectType: pk.EffectType,
		Amplifier: pk.Amplifier, Particles: pk.Particles, Duration: pk.Duration, Tick: pk.Tick,
	}
}

func toLatestChangeMobProperty(pk *ChangeMobProperty) *packet.ChangeMobProperty {
	return &packet.ChangeMobProperty{
		EntityUniqueID: int64(pk.EntityUniqueID), Property: pk.Property, BoolValue: pk.BoolValue,
		StringValue: pk.StringValue, IntValue: pk.IntValue, FloatValue: pk.FloatValue,
	}
}
func fromLatestChangeMobProperty(pk *packet.ChangeMobProperty) *ChangeMobProperty {
	return &ChangeMobProperty{
		EntityUniqueID: uint64(pk.EntityUniqueID), Property: pk.Property, BoolValue: pk.BoolValue,
		StringValue: pk.StringValue, IntValue: pk.IntValue, FloatValue: pk.FloatValue,
	}
}

func toLatestPlayerUpdateEntityOverrides(pk *PlayerUpdateEntityOverrides) *packet.PlayerUpdateEntityOverrides {
	return &packet.PlayerUpdateEntityOverrides{

		EntityUniqueID: int64(pk.EntityRuntimeID), PropertyIndex: pk.PropertyIndex, Type: uint32(pk.Type),
		IntValue: pk.IntValue, FloatValue: pk.FloatValue,
	}
}
func fromLatestPlayerUpdateEntityOverrides(pk *packet.PlayerUpdateEntityOverrides) *PlayerUpdateEntityOverrides {
	return &PlayerUpdateEntityOverrides{
		EntityRuntimeID: uint64(pk.EntityUniqueID), PropertyIndex: pk.PropertyIndex, Type: byte(pk.Type),
		IntValue: pk.IntValue, FloatValue: pk.FloatValue,
	}
}

func toLatestPlayerArmourDamage(pk *PlayerArmourDamage) *packet.PlayerArmourDamage {
	return &packet.PlayerArmourDamage{}
}
func fromLatestPlayerArmourDamage(pk *packet.PlayerArmourDamage) *PlayerArmourDamage {
	return &PlayerArmourDamage{}
}
