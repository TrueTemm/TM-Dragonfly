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

package v818

import (
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v800"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func applyDeltas818(p packet.Pool) {

	delete(p, v786.IDPassengerJump)
	delete(p, v786.IDTickSync)
	delete(p, v786.IDPlayerInput)
	delete(p, v786.IDCompressedBiomeDefinitionList)
	delete(p, v786.IDSetMovementAuthority)

	p[IDPlayerList] = func() packet.Packet { return &v800.PlayerList{} }
	p[IDCameraAimAssistPresets] = func() packet.Packet { return &packet.CameraAimAssistPresets{} }

	p[IDClientMovementPredictionSync] = func() packet.Packet { return &v800.ClientMovementPredictionSync{} }
	p[IDPlayerLocation] = func() packet.Packet { return &v800.PlayerLocation{} }

	p[IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }

	p[IDClientBoundControlSchemeSet] = func() packet.Packet { return &packet.ClientBoundControlSchemeSet{} }
}

func NewServerPool() packet.Pool {
	p := v786.NewServerPool()
	applyDeltas818(p)
	return p
}

func NewClientPool() packet.Pool {
	p := v786.NewClientPool()
	applyDeltas818(p)
	return p
}
