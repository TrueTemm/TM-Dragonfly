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

package v2169

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type EntityDiagnosticTimingInfo struct {
	protocol.EntityDiagnosticTimingInfo
}

func (x *EntityDiagnosticTimingInfo) Marshal(r protocol.IO) {
	r.String(&x.DisplayName)
	r.String(&x.Entity)
	r.Uint64(&x.DurationNanos)
	r.Uint8(&x.PercentOfTotal)
}

type ServerBoundDiagnostics struct {
	packet.ServerBoundDiagnostics
}

func (pk *ServerBoundDiagnostics) Marshal(io protocol.IO) {
	io.Float32(&pk.AverageFramesPerSecond)
	io.Float32(&pk.AverageServerSimTickTime)
	io.Float32(&pk.AverageClientSimTickTime)
	io.Float32(&pk.AverageBeginFrameTime)
	io.Float32(&pk.AverageInputTime)
	io.Float32(&pk.AverageRenderTime)
	io.Float32(&pk.AverageEndFrameTime)
	io.Float32(&pk.AverageRemainderTimePercent)
	io.Float32(&pk.AverageUnaccountedTimePercent)
	protocol.Slice(io, &pk.MemoryCategoryValues)
	entities := make([]EntityDiagnosticTimingInfo, len(pk.EntityDiagnostics))
	for i, e := range pk.EntityDiagnostics {
		entities[i] = EntityDiagnosticTimingInfo{EntityDiagnosticTimingInfo: e}
	}
	protocol.Slice(io, &entities)
	pk.EntityDiagnostics = make([]protocol.EntityDiagnosticTimingInfo, len(entities))
	for i, e := range entities {
		pk.EntityDiagnostics[i] = e.EntityDiagnosticTimingInfo
	}
	protocol.Slice(io, &pk.SystemDiagnostics)
	protocol.Slice(io, &pk.SystemCategories)
	protocol.Slice(io, &pk.WhiskerScopes)
}

func toLatestServerBoundDiagnostics(pk *ServerBoundDiagnostics) *packet.ServerBoundDiagnostics {
	out := pk.ServerBoundDiagnostics
	return &out
}

func fromLatestServerBoundDiagnostics(pk *packet.ServerBoundDiagnostics) *ServerBoundDiagnostics {
	return &ServerBoundDiagnostics{ServerBoundDiagnostics: *pk}
}
