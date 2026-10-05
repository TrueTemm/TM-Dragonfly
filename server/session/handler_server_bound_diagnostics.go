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

package session

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ServerBoundDiagnosticsHandler struct{}

func (h *ServerBoundDiagnosticsHandler) Handle(p packet.Packet, _ *Session, _ *world.Tx, c Controllable) error {
	pk := p.(*packet.ServerBoundDiagnostics)
	c.UpdateDiagnostics(Diagnostics{
		AverageFramesPerSecond:        float64(pk.AverageFramesPerSecond),
		AverageServerSimTickTime:      float64(pk.AverageServerSimTickTime),
		AverageClientSimTickTime:      float64(pk.AverageClientSimTickTime),
		AverageBeginFrameTime:         float64(pk.AverageBeginFrameTime),
		AverageInputTime:              float64(pk.AverageInputTime),
		AverageRenderTime:             float64(pk.AverageRenderTime),
		AverageEndFrameTime:           float64(pk.AverageEndFrameTime),
		AverageRemainderTimePercent:   float64(pk.AverageRemainderTimePercent),
		AverageUnaccountedTimePercent: float64(pk.AverageUnaccountedTimePercent),
	})
	return nil
}
