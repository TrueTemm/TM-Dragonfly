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
)

const (
	BossEventShow = iota
	BossEventRegisterPlayer
	BossEventHide
	BossEventUnregisterPlayer
	BossEventHealthPercentage
	BossEventTitle
	BossEventAppearanceProperties
	BossEventTexture
	BossEventRequest
)

const (
	BossEventColourGrey = iota
	BossEventColourBlue
	BossEventColourRed
	BossEventColourGreen
	BossEventColourYellow
	BossEventColourPurple
	BossEventColourWhite
)

type BossEvent struct {
	BossEntityUniqueID int64

	EventType uint32

	PlayerUniqueID int64

	BossBarTitle string

	FilteredBossBarTitle string

	HealthPercentage float32

	ScreenDarkening uint16

	Colour uint32

	Overlay uint32
}

func (*BossEvent) ID() uint32 {
	return IDBossEvent
}

func (pk *BossEvent) Marshal(io protocol.IO) {
	io.Varint64(&pk.BossEntityUniqueID)
	io.Varuint32(&pk.EventType)
	switch pk.EventType {
	case BossEventShow:
		io.String(&pk.BossBarTitle)
		io.String(&pk.FilteredBossBarTitle)
		io.Float32(&pk.HealthPercentage)
		io.Uint16(&pk.ScreenDarkening)
		io.Varuint32(&pk.Colour)
		io.Varuint32(&pk.Overlay)
	case BossEventRegisterPlayer, BossEventUnregisterPlayer, BossEventRequest:
		io.Varint64(&pk.PlayerUniqueID)
	case BossEventHide:

	case BossEventHealthPercentage:
		io.Float32(&pk.HealthPercentage)
	case BossEventTitle:
		io.String(&pk.BossBarTitle)
		io.String(&pk.FilteredBossBarTitle)
	case BossEventAppearanceProperties:
		io.Uint16(&pk.ScreenDarkening)
		io.Varuint32(&pk.Colour)
		io.Varuint32(&pk.Overlay)
	case BossEventTexture:
		io.Varuint32(&pk.Colour)
		io.Varuint32(&pk.Overlay)
	default:
		io.UnknownEnumOption(pk.EventType, "boss event type")
	}
}
