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

package native

import (
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestClientCacheIsRefused(t *testing.T) {
	out := Protocol{}.ConvertToLatest(&packet.ClientCacheStatus{Enabled: true}, nil)
	if len(out) != 1 {
		t.Fatalf("ClientCacheStatus converted into %v packets, want exactly 1", len(out))
	}
	status, ok := out[0].(*packet.ClientCacheStatus)
	if !ok {
		t.Fatalf("ClientCacheStatus converted into a %T", out[0])
	}
	if status.Enabled {
		t.Fatal("the client's blob cache was left on: the server will send chunks as blob hashes, the " +
			"blocks inside them are never translated, and the client is handed an empty world")
	}
}
