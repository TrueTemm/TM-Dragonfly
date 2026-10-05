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

type CreativeContent struct {
	Groups []CreativeGroup786

	Items []protocol.CreativeItem
}

type CreativeGroup786 struct {
	Category int32

	Name string

	Icon protocol.ItemStack
}

func (x *CreativeGroup786) Marshal(r protocol.IO) {
	r.Int32(&x.Category)
	r.String(&x.Name)
	r.Item(&x.Icon)
}

func (*CreativeContent) ID() uint32 {
	return IDCreativeContent
}

func (pk *CreativeContent) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Groups)
	protocol.Slice(io, &pk.Items)
}
