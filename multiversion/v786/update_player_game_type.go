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

type UpdatePlayerGameType struct {
	GameType       int32
	PlayerUniqueID int64
	Tick           uint64
}

func (*UpdatePlayerGameType) ID() uint32 { return IDUpdatePlayerGameType }
func (pk *UpdatePlayerGameType) Marshal(io protocol.IO) {
	io.Varint32(&pk.GameType)
	io.Varint64(&pk.PlayerUniqueID)
	if p := ProtoOf(io); p == 0 || p >= 671 { // added in 671
		io.Varuint64(&pk.Tick)
	}
}
func toLatestUpdatePlayerGameType(pk *UpdatePlayerGameType) *packet.UpdatePlayerGameType {
	return &packet.UpdatePlayerGameType{GameType: pk.GameType, PlayerUniqueID: pk.PlayerUniqueID, Tick: pk.Tick}
}
func fromLatestUpdatePlayerGameType(pk *packet.UpdatePlayerGameType) *UpdatePlayerGameType {
	return &UpdatePlayerGameType{GameType: pk.GameType, PlayerUniqueID: pk.PlayerUniqueID, Tick: pk.Tick}
}
