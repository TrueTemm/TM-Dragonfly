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
	BookActionReplacePage = iota
	BookActionAddPage
	BookActionDeletePage
	BookActionSwapPages
	BookActionSign
)

type BookEdit struct {
	ActionType byte

	InventorySlot byte

	PageNumber byte

	SecondaryPageNumber byte

	Text string

	PhotoName string

	Title string

	Author string

	XUID string
}

func (*BookEdit) ID() uint32 {
	return IDBookEdit
}

func (pk *BookEdit) Marshal(io protocol.IO) {
	io.Uint8(&pk.ActionType)
	io.Uint8(&pk.InventorySlot)
	switch pk.ActionType {
	case BookActionReplacePage, BookActionAddPage:
		io.Uint8(&pk.PageNumber)
		io.String(&pk.Text)
		io.String(&pk.PhotoName)
	case BookActionDeletePage:
		io.Uint8(&pk.PageNumber)
	case BookActionSwapPages:
		io.Uint8(&pk.PageNumber)
		io.Uint8(&pk.SecondaryPageNumber)
	case BookActionSign:
		io.String(&pk.Title)
		io.String(&pk.Author)
		io.String(&pk.XUID)
	default:
		io.UnknownEnumOption(pk.ActionType, "book edit action type")
	}
}
