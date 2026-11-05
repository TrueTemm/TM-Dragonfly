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
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/worldupgrader/itemupgrader"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func (w Writer786) ItemDescriptorCount(x *protocol.ItemDescriptorCount) {
	if w.Proto > 630 || w.Proto == 0 {
		w.Writer.ItemDescriptorCount(x)
		return
	}
	var kind byte
	switch x.Descriptor.(type) {
	case nil, *protocol.InvalidItemDescriptor:
		kind = 0
	case *protocol.DefaultItemDescriptor:
		kind = 1
	case *protocol.MoLangItemDescriptor:
		kind = 2
	case *protocol.ItemTagItemDescriptor:
		kind = 3
	default:
		w.UnknownEnumOption(x.Descriptor, "legacy ingredient")
		return
	}
	w.Uint8(&kind)
	switch d := x.Descriptor.(type) {
	case *protocol.DefaultItemDescriptor:
		native, _ := world.ItemRuntimeIDByName(d.Name)
		id, meta := itemdata.TranslateStack(w.Proto, native, uint32(d.MetadataValue))
		networkID := int16(id)
		w.Int16(&networkID)
		if networkID != 0 {
			aux := int16(meta)
			w.Int16(&aux)
		}
	case *protocol.MoLangItemDescriptor:
		value := d.Expression
		w.String(&value)
		version := byte(d.Version)
		w.Uint8(&version)
	case *protocol.ItemTagItemDescriptor:
		tag := d.Tag
		w.String(&tag)
	}
	count := x.Count
	w.Varint32(&count)
}

func (r Reader786) ItemDescriptorCount(x *protocol.ItemDescriptorCount) {
	if r.Proto > 630 || r.Proto == 0 {
		r.Reader.ItemDescriptorCount(x)
		return
	}
	var kind byte
	r.Uint8(&kind)
	switch kind {
	case 0:
		x.Descriptor = &protocol.InvalidItemDescriptor{}
	case 1:
		var id, meta int16
		r.Int16(&id)
		if id != 0 {
			r.Int16(&meta)
		}
		native, aux := itemdata.ReverseTranslateStack(r.Proto, int32(id), uint32(uint16(meta)))
		name, _ := world.ItemNameByRuntimeID(native)
		x.Descriptor = &protocol.DefaultItemDescriptor{Name: name, MetadataValue: int32(aux)}
	case 2:
		var value string
		var version byte
		r.String(&value)
		r.Uint8(&version)
		x.Descriptor = &protocol.MoLangItemDescriptor{Expression: value, Version: int16(version)}
	case 3:
		var tag string
		r.String(&tag)
		x.Descriptor = &protocol.ItemTagItemDescriptor{Tag: tag}
	case 4, 5:
		var name string
		var meta int16
		r.String(&name)
		if kind == 4 {
			r.Int16(&meta)
		}
		up := itemupgrader.Upgrade(itemupgrader.ItemMeta{Name: name, Meta: meta})
		x.Descriptor = &protocol.DefaultItemDescriptor{Name: up.Name, MetadataValue: int32(up.Meta)}
	default:
		r.UnknownEnumOption(kind, "legacy ingredient")
		return
	}
	r.Varint32(&x.Count)
}
