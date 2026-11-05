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

import "github.com/sandertv/gophertunnel/minecraft/protocol"

func GameRuleLegacy(io protocol.IO, x *protocol.GameRule) {
	io.String(&x.Name)
	io.Bool(&x.CanBeModifiedByPlayer)
	var kind uint32
	switch x.Value.(type) {
	case bool:
		kind = 1
	case uint32:
		kind = 2
	case float32:
		kind = 3
	case nil:
	default:
		io.UnknownEnumOption(x.Value, "game rule type")
	}
	io.Varuint32(&kind)
	switch kind {
	case 0:
		x.Value = nil
	case 1:
		v, _ := x.Value.(bool)
		io.Bool(&v)
		x.Value = v
	case 2:
		v, _ := x.Value.(uint32)
		io.Varuint32(&v) // legacy integer rule
		x.Value = v
	case 3:
		v, _ := x.Value.(float32)
		io.Float32(&v)
		x.Value = v
	default:
		io.UnknownEnumOption(kind, "game rule type")
	}
}

func (r Reader786) GameRule(x *protocol.GameRule) {
	if r.Proto != 0 && r.Proto < 844 {
		GameRuleLegacy(r, x)
		return
	}
	r.Reader.GameRule(x)
}

func (w Writer786) GameRule(x *protocol.GameRule) {
	if w.Proto != 0 && w.Proto < 844 {
		GameRuleLegacy(w, x)
		return
	}
	w.Writer.GameRule(x)
}
