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
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestAvailableCommandsRoundTrips(t *testing.T) {
	in := &AvailableCommands{
		EnumValues:              []string{"true", "false", "yes", "no"},
		ChainedSubcommandValues: []string{"as", "at"},
		Suffixes:                []string{"m", "s"},
		Enums: []CommandEnum786{
			{Type: "boolArg", ValueIndices: []uint{0, 1}},
		},
		ChainedSubcommands: []ChainedSubcommand786{
			{Name: "as", Values: []ChainedSubcommandValue786{{Index: 0, Value: 8}}},
		},
		Commands: []Command786{
			{
				Name: "gamerule", Description: "commands.gamerule.description", Flags: 0,
				PermissionLevel: 1, AliasesOffset: 0xFFFFFFFF,
				ChainedSubcommandOffsets: []uint16{0},
				Overloads: []protocol.CommandOverload{
					{Chaining: false, Parameters: []protocol.CommandParameter{
						{Name: "rule", Type: protocol.CommandArgTypeString | protocol.CommandArgValid, Optional: false},
					}},
				},
			},
		},
		DynamicEnums: []protocol.DynamicEnum{{Type: "dyn", Values: []string{"a", "b"}}},
		Constraints:  []protocol.CommandEnumConstraint{{EnumValueIndex: 0, EnumIndex: 0, Constraints: []byte{0}}},
	}

	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	in.Marshal(w)

	if bytes.Contains(buf.Bytes(), []byte("gamedirectors")) {
		t.Fatalf("encoded bytes contain latest's string-based permission encoding — wrong Marshal was used")
	}

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	out := &AvailableCommands{}
	out.Marshal(r)

	if len(out.Commands) != 1 || out.Commands[0].Name != "gamerule" {
		t.Fatalf("expected 1 command named gamerule, got %+v", out.Commands)
	}
	if out.Commands[0].PermissionLevel != 1 {
		t.Fatalf("expected PermissionLevel=1, got %d", out.Commands[0].PermissionLevel)
	}
	if len(out.Commands[0].Overloads) != 1 || len(out.Commands[0].Overloads[0].Parameters) != 1 {
		t.Fatalf("overload/parameter round-trip failed: %+v", out.Commands[0])
	}
	if len(out.Enums) != 1 || out.Enums[0].Type != "boolArg" || len(out.Enums[0].ValueIndices) != 2 {
		t.Fatalf("enum round-trip failed: %+v", out.Enums)
	}
	if len(out.ChainedSubcommands) != 1 || out.ChainedSubcommands[0].Values[0].Value != 8 {
		t.Fatalf("chained subcommand round-trip failed: %+v", out.ChainedSubcommands)
	}
}

func TestCommandEnumWidthEscalates(t *testing.T) {
	small := make([]string, 5)
	for i := range small {
		small[i] = "x"
	}
	large := make([]string, 300)
	for i := range large {
		large[i] = "x"
	}

	for _, tc := range []struct {
		name        string
		enumValues  []string
		wantBodyLen int
	}{
		{"byteWidth", small, 1},
		{"uint16Width", large, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := commandEnumContext786{EnumValues: tc.enumValues}
			buf := new(bytes.Buffer)
			w := protocol.NewWriter(buf, 0)
			enum := &CommandEnum786{Type: "t", ValueIndices: []uint{1}}
			ctx.Marshal(w, enum)

			r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
			out := &CommandEnum786{}
			ctx.Marshal(r, out)
			if len(out.ValueIndices) != 1 || out.ValueIndices[0] != 1 {
				t.Fatalf("round-trip failed for %s: %+v", tc.name, out)
			}
		})
	}
}
