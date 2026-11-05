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
	"fmt"
	"log/slog"
	"net"
)

func logNetherNetListenFailure(log *slog.Logger, signalingAddress, iceAddress string, err error) {
	_, tcpPort, tcpErr := net.SplitHostPort(signalingAddress)
	_, udpPort, udpErr := net.SplitHostPort(iceAddress)
	if tcpErr != nil || udpErr != nil {
		log.Warn("NetherNet listener failed", "error", err, "tcp", signalingAddress, "udp", iceAddress)
		return
	}
	log.Warn("NetherNet listener failed", "error", err, "tcp", signalingAddress, "udp", iceAddress)
	log.Warn("Check listener ports in UFW and provider firewall",
		"commands", fmt.Sprintf("sudo ufw allow %s/tcp && sudo ufw allow %s/udp", tcpPort, udpPort))
}
