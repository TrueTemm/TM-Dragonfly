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

package conformance

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestHistoricalIntegerGameRule(t *testing.T) {
	for _, p := range protocols() {
		t.Run(p.Ver(), func(t *testing.T) {
			want := []byte{3, 'r', 'u', 'l', 1, 2, 0xac, 2}
			if p.ID() >= 844 {
				want = []byte{3, 'r', 'u', 'l', 1, 2, 0x2c, 1, 0, 0}
			}
			input := protocol.GameRule{Name: "rul", CanBeModifiedByPlayer: true, Value: uint32(300)}
			encoded := new(bytes.Buffer)
			p.NewWriter(encoded, 0).GameRule(&input)
			if !bytes.Equal(encoded.Bytes(), want) {
				t.Fatalf("integer rule wire %x want %x", encoded.Bytes(), want)
			}
			raw := bytes.NewBuffer(want)
			var output protocol.GameRule
			p.NewReader(raw, 0, true).GameRule(&output)
			if output != input || raw.Len() != 0 {
				t.Fatalf("rule decoded as %+v with %d unread bytes", output, raw.Len())
			}
		})
	}
}
