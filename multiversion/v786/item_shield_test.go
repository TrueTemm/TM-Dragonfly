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

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	_ "github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
)

func TestShieldIDIsTheClientVersions(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	const latestShield = 358
	boots712, shield712 := int32(358), itemdata.ShieldID(712)
	if shield712 == 0 || shield712 == boots712 {
		t.Fatalf("712 registry: shield=%d boots=%d — the collision this test relies on has moved", shield712, boots712)
	}
	for _, tc := range []struct {
		name string
		id   int32
	}{{"diamond boots (the ID that is latest's shield)", boots712}, {"the real 712 shield", shield712}} {

		buf := bytes.NewBuffer(nil)
		in := protocol.ItemInstance{StackNetworkID: 4, Stack: protocol.ItemStack{ItemType: protocol.ItemType{NetworkID: tc.id}, Count: 1,
			NBTData: map[string]any{"Damage": int32(3)}, BlockingTick: 12}}
		writeItemInstance786(protocol.NewWriter(buf, latestShield), shield712, &in)

		var out protocol.ItemInstance
		r := Reader786{Reader: protocol.NewReader(buf, latestShield, false), Proto: 712}
		func() {
			defer func() {
				if e := recover(); e != nil {
					t.Fatalf("%s: decode panicked: %v", tc.name, e)
				}
			}()
			r.ItemInstance(&out)
		}()
		if buf.Len() != 0 {
			t.Errorf("%s: %d bytes left unread", tc.name, buf.Len())
		}
		wantTick := int64(0)
		if tc.id == shield712 {
			wantTick = 12
		}
		if out.Stack.BlockingTick != wantTick {
			t.Errorf("%s: BlockingTick = %d, want %d", tc.name, out.Stack.BlockingTick, wantTick)
		}
		if out.Stack.NBTData["Damage"] != int32(3) {
			t.Errorf("%s: NBT lost: %v", tc.name, out.Stack.NBTData)
		}
	}
}
