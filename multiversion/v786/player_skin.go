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
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerSkin struct {
	UUID        uuid.UUID
	Skin        Skin786
	NewSkinName string
	OldSkinName string
}

func (*PlayerSkin) ID() uint32 { return packet.IDPlayerSkin }

func (pk *PlayerSkin) Marshal(io protocol.IO) {
	io.UUID(&pk.UUID)
	protocol.Single(io, &pk.Skin)
	io.String(&pk.NewSkinName)
	io.String(&pk.OldSkinName)
	io.Bool(&pk.Skin.Trusted)
}

func ToLatestPlayerSkin(pk *PlayerSkin) *packet.PlayerSkin {
	return &packet.PlayerSkin{
		UUID:        pk.UUID,
		Skin:        toLatestSkin786(pk.Skin),
		NewSkinName: pk.NewSkinName,
		OldSkinName: pk.OldSkinName,
	}
}

func FromLatestPlayerSkin(pk *packet.PlayerSkin) *PlayerSkin {
	return &PlayerSkin{
		UUID:        pk.UUID,
		Skin:        fromLatestSkin786(pk.Skin),
		NewSkinName: pk.NewSkinName,
		OldSkinName: pk.OldSkinName,
	}
}
