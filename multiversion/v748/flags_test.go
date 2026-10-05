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

package v748

import (
	"testing"

	"github.com/df-mc/dragonfly/multiversion/v786"
)

func TestInputFlagMapping(t *testing.T) {
	cases := []struct {
		bit748 int
		idx786 int
	}{
		{0, v786.InputFlagAscend},
		{20, v786.InputFlagSprinting},
		{34, v786.InputFlagPerformItemInteraction},
		{36, v786.InputFlagPerformItemStackRequest},
		{45, v786.InputFlagClientPredictedVehicle},
		{52, v786.InputFlagDownRight},
		{53, v786.InputFlagCameraRelativeMovementEnabled},
		{56, v786.InputFlagStopSpinAttack},
	}
	for _, c := range cases {
		b := bitsetFromMask(1 << uint(c.bit748))
		if !b.Load(c.idx786) {
			t.Errorf("748 bit %d should map to 786 index %d", c.bit748, c.idx786)
		}
		for i := 0; i < v786.PlayerAuthInputBitsetSize; i++ {
			if i != c.idx786 && b.Load(i) {
				t.Errorf("748 bit %d also set 786 index %d", c.bit748, i)
			}
		}
		if got := maskFromBitset(b); got != 1<<uint(c.bit748) {
			t.Errorf("round trip of 748 bit %d gave mask %#x", c.bit748, got)
		}
	}

	b := bitsetFromMask(0)
	b.Set(v786.InputFlagStartUsingItem)
	if got := maskFromBitset(b); got != 0 {
		t.Errorf("StartUsingItem leaked into the 748 mask as %#x", got)
	}

	if inputFlag748PerformItemInteraction != 1<<34 || inputFlag748PerformItemStackRequest != 1<<36 ||
		inputFlag748PerformBlockActions != 1<<35 || inputFlag748ClientPredictedVehicle != 1<<45 {
		t.Error("branch flag constants drifted from the 748 table")
	}
}
