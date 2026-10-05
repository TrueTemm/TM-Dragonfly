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
	"math"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Command786 struct {
	Name string

	Description string

	Flags uint16

	PermissionLevel byte

	AliasesOffset uint32

	ChainedSubcommandOffsets []uint16

	Overloads []protocol.CommandOverload
}

func (c *Command786) Marshal(r protocol.IO) {
	r.String(&c.Name)
	r.String(&c.Description)
	r.Uint16(&c.Flags)
	r.Uint8(&c.PermissionLevel)
	r.Uint32(&c.AliasesOffset)
	protocol.FuncSlice(r, &c.ChainedSubcommandOffsets, r.Uint16)
	protocol.Slice(r, &c.Overloads)
}

type CommandEnum786 struct {
	Type string

	ValueIndices []uint
}

type commandEnumContext786 struct {
	EnumValues []string
}

func (ctx commandEnumContext786) Marshal(r protocol.IO, x *CommandEnum786) {
	r.String(&x.Type)
	protocol.FuncIOSlice(r, &x.ValueIndices, ctx.enumOption)
}

func (ctx commandEnumContext786) enumOption(r protocol.IO, opt *uint) {
	n := len(ctx.EnumValues)
	switch {
	case n <= math.MaxUint8:
		val := byte(*opt)
		r.Uint8(&val)
		*opt = uint(val)
	case n <= math.MaxUint16:
		val := uint16(*opt)
		r.Uint16(&val)
		*opt = uint(val)
	default:
		val := uint32(*opt)
		r.Uint32(&val)
		*opt = uint(val)
	}
}

type ChainedSubcommandValue786 struct {
	Index uint16
	Value uint16
}

func (x *ChainedSubcommandValue786) Marshal(r protocol.IO) {
	r.Uint16(&x.Index)
	r.Uint16(&x.Value)
}

type ChainedSubcommand786 struct {
	Name   string
	Values []ChainedSubcommandValue786
}

func (x *ChainedSubcommand786) Marshal(r protocol.IO) {
	r.String(&x.Name)
	protocol.Slice(r, &x.Values)
}

type AvailableCommands struct {
	EnumValues              []string
	ChainedSubcommandValues []string
	Suffixes                []string
	Enums                   []CommandEnum786
	ChainedSubcommands      []ChainedSubcommand786
	Commands                []Command786

	DynamicEnums []protocol.DynamicEnum
	Constraints  []protocol.CommandEnumConstraint
}

func (*AvailableCommands) ID() uint32 {
	return IDAvailableCommands
}

func (pk *AvailableCommands) Marshal(io protocol.IO) {
	protocol.FuncSlice(io, &pk.EnumValues, io.String)
	protocol.FuncSlice(io, &pk.ChainedSubcommandValues, io.String)
	protocol.FuncSlice(io, &pk.Suffixes, io.String)
	protocol.FuncIOSlice(io, &pk.Enums, commandEnumContext786{EnumValues: pk.EnumValues}.Marshal)
	protocol.Slice(io, &pk.ChainedSubcommands)
	protocol.Slice(io, &pk.Commands)
	protocol.Slice(io, &pk.DynamicEnums)
	protocol.Slice(io, &pk.Constraints)
}
