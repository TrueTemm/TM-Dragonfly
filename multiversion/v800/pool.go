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

package v800

import (
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func applyDeltas800(p packet.Pool) {

	delete(p, v786.IDPassengerJump)
	delete(p, v786.IDTickSync)
	delete(p, v786.IDPlayerInput)
	delete(p, v786.IDCompressedBiomeDefinitionList)

	p[IDBiomeDefinitionList] = func() packet.Packet { return &packet.BiomeDefinitionList{} }
	p[IDPlayerList] = func() packet.Packet { return &PlayerList{} }
	p[IDCameraAimAssist] = func() packet.Packet { return &packet.CameraAimAssist{} }
	p[IDCameraAimAssistPresets] = func() packet.Packet { return &packet.CameraAimAssistPresets{} }
	p[IDCameraInstruction] = func() packet.Packet { return &packet.CameraInstruction{} }

	p[IDClientMovementPredictionSync] = func() packet.Packet { return &ClientMovementPredictionSync{} }
	p[IDPlayerLocation] = func() packet.Packet { return &PlayerLocation{} }
	p[IDClientBoundControlSchemeSet] = func() packet.Packet { return &packet.ClientBoundControlSchemeSet{} }
}

func NewServerPool() packet.Pool {
	p := v786.NewServerPool()
	applyDeltas800(p)
	return p
}

func NewClientPool() packet.Pool {
	p := v786.NewClientPool()
	applyDeltas800(p)
	return p
}
