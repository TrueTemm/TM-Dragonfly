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

package v898

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Text struct {
	packet.Text
}

func (*Text) ID() uint32 {
	return IDText
}

func stringConst(io protocol.IO, x string) {
	io.String(&x)
}

func (pk *Text) Marshal(io protocol.IO) {
	io.Bool(&pk.NeedsTranslation)

	var categoryType uint8
	switch pk.TextType {
	case packet.TextTypeRaw, packet.TextTypeTip, packet.TextTypeSystem, packet.TextTypeObjectWhisper,
		packet.TextTypeObjectAnnouncement, packet.TextTypeObject:
		categoryType = 0
	case packet.TextTypeChat, packet.TextTypeWhisper, packet.TextTypeAnnouncement:
		categoryType = 1
	default:
		categoryType = 2
	}
	io.Uint8(&categoryType)

	switch categoryType {
	case 0:
		stringConst(io, "raw")
		stringConst(io, "tip")
		stringConst(io, "systemMessage")
		stringConst(io, "textObjectWhisper")
		stringConst(io, "textObjectAnnouncement")
		stringConst(io, "textObject")
	case 1:
		stringConst(io, "chat")
		stringConst(io, "whisper")
		stringConst(io, "announcement")
	default:
		stringConst(io, "translate")
		stringConst(io, "popup")
		stringConst(io, "jukeboxPopup")
	}

	io.Uint8(&pk.TextType)
	switch pk.TextType {
	case packet.TextTypeChat, packet.TextTypeWhisper, packet.TextTypeAnnouncement:
		io.String(&pk.SourceName)
		io.String(&pk.Message)
	case packet.TextTypeRaw, packet.TextTypeTip, packet.TextTypeSystem, packet.TextTypeObject,
		packet.TextTypeObjectWhisper, packet.TextTypeObjectAnnouncement:
		io.String(&pk.Message)
	case packet.TextTypeTranslation, packet.TextTypePopup, packet.TextTypeJukeboxPopup:
		io.String(&pk.Message)
		protocol.FuncSlice(io, &pk.Parameters, io.String)
	}

	io.String(&pk.XUID)
	io.String(&pk.PlatformChatID)
	protocol.OptionalFunc(io, &pk.FilteredMessage, io.String)
}

func ToLatestText898(pk *Text) *packet.Text {
	return &pk.Text
}

func FromLatestText898(pk *packet.Text) *Text {
	return &Text{Text: *pk}
}
