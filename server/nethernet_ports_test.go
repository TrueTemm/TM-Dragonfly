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

package server

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestNetherNetPortHints(t *testing.T) {
	for _, tc := range []struct {
		tcp, udp, expected string
	}{
		{":19132", ":19132", "sudo ufw allow 19132/tcp && sudo ufw allow 19132/udp"},
		{"0.0.0.0:19132", "[::]:19133", "sudo ufw allow 19132/tcp && sudo ufw allow 19133/udp"},
	} {
		var b bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&b, nil))
		logNetherNetListenFailure(logger, tc.tcp, tc.udp, fmt.Errorf("bind failed"))
		if !strings.Contains(b.String(), "NetherNet listener failed") || !strings.Contains(b.String(), tc.expected) {
			t.Fatalf("%s %s: missing port hint in %q", tc.tcp, tc.udp, b.String())
		}
	}
}

func TestNetherNetListenHintOnlyOnFailure(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, nil))
	if output.Len() != 0 {
		t.Fatal("port hints printed before listener error")
	}
	logNetherNetListenFailure(logger, ":19132", ":19132", fmt.Errorf("address already in use"))
	if strings.Count(output.String(), "Check listener ports") != 1 {
		t.Fatalf("expected exactly one port hint on failure: %q", output.String())
	}
}
