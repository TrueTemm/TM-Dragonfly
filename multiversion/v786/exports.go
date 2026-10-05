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
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func FromLatestPlayerList(pk *packet.PlayerList) *PlayerList { return fromLatestPlayerList(pk) }

func ToLatestPlayerList(pk *PlayerList) *packet.PlayerList { return toLatestPlayerList(pk) }

func FromLatestClientMovementPredictionSync(pk *packet.ClientMovementPredictionSync) *ClientMovementPredictionSync {
	return fromLatestClientMovementPredictionSync(pk)
}
func ToLatestClientMovementPredictionSync(pk *ClientMovementPredictionSync) *packet.ClientMovementPredictionSync {
	return toLatestClientMovementPredictionSync(pk)
}

func FromLatestBossEvent(pk *packet.BossEvent) *BossEvent { return fromLatestBossEvent(pk) }
func ToLatestBossEvent(pk *BossEvent) *packet.BossEvent   { return toLatestBossEvent(pk) }
func FromLatestCommandBlockUpdate(pk *packet.CommandBlockUpdate) *CommandBlockUpdate {
	return fromLatestCommandBlockUpdate(pk)
}
func ToLatestCommandBlockUpdate(pk *CommandBlockUpdate) *packet.CommandBlockUpdate {
	return toLatestCommandBlockUpdate(pk)
}
func FromLatestStructureBlockUpdate(pk *packet.StructureBlockUpdate) *StructureBlockUpdate {
	return fromLatestStructureBlockUpdate(pk)
}
func ToLatestStructureBlockUpdate(pk *StructureBlockUpdate) *packet.StructureBlockUpdate {
	return toLatestStructureBlockUpdate(pk)
}

func MarshalItemStackRequest786(io protocol.IO, pk *protocol.ItemStackRequest) {
	marshalItemStackRequest786(io, pk)
}
func PlayerBlockActions786(io protocol.IO, x *[]protocol.PlayerBlockAction) {
	playerBlockActions786(io, x)
}
func ToLatestPlayerAuthInput(pk *PlayerAuthInput) *packet.PlayerAuthInput {
	return toLatestPlayerAuthInput(pk)
}
func FromLatestPlayerAuthInput(pk *packet.PlayerAuthInput) *PlayerAuthInput {
	return fromLatestPlayerAuthInput(pk)
}

func FromLatestPlayerArmourDamage(pk *packet.PlayerArmourDamage) *PlayerArmourDamage {
	return fromLatestPlayerArmourDamage(pk)
}
func ToLatestPlayerArmourDamage(pk *PlayerArmourDamage) *packet.PlayerArmourDamage {
	return toLatestPlayerArmourDamage(pk)
}

func FromLatestMobEffect(pk *packet.MobEffect) *MobEffect { return fromLatestMobEffect(pk) }
func ToLatestMobEffect(pk *MobEffect) *packet.MobEffect   { return toLatestMobEffect(pk) }
func FromLatestInventorySlot(pk *packet.InventorySlot) *InventorySlot {
	return fromLatestInventorySlot(pk)
}
func ToLatestInventorySlot(pk *InventorySlot) *packet.InventorySlot { return toLatestInventorySlot(pk) }

func FromLatestTransfer(pk *packet.Transfer) *Transfer { return fromLatestTransfer(pk) }
func ToLatestTransfer(pk *Transfer) *packet.Transfer   { return toLatestTransfer(pk) }

func FromLatestPlaySound(pk *packet.PlaySound) *PlaySound { return fromLatestPlaySound(pk) }
func ToLatestPlaySound(pk *PlaySound) *packet.PlaySound   { return toLatestPlaySound(pk) }
func FromLatestChangeDimension(pk *packet.ChangeDimension) *ChangeDimension {
	return fromLatestChangeDimension(pk)
}
func FromLatestDisconnect(pk *packet.Disconnect) *Disconnect { return fromLatestDisconnect(pk) }
