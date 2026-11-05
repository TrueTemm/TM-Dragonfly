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
	StoreOfferTypeMarketplace = iota
	StoreOfferTypeDressingRoom
	StoreOfferTypeServerPage
)

type ShowStoreOffer struct {
	OfferID string

	Type byte
}

func (*ShowStoreOffer) ID() uint32 {
	return IDShowStoreOffer
}

func (pk *ShowStoreOffer) Marshal(io protocol.IO) {
	io.String(&pk.OfferID)
	if p := ProtoOf(io); p != 0 && p < 630 {
		showAll := pk.Type != 0
		io.Bool(&showAll)
		if showAll {
			pk.Type = 1
		} else {
			pk.Type = 0
		}
	} else {
		io.Uint8(&pk.Type)
	}
}
