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
