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

package v944

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v898"
	"github.com/df-mc/dragonfly/multiversion/v924"
)

func applyDeltas944(p packet.Pool) {
	p[v898.IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[packet.IDVoxelShapes] = func() packet.Packet { return &packet.VoxelShapes{} }
	p[packet.IDUpdateClientInputLocks] = func() packet.Packet { return &packet.UpdateClientInputLocks{} }
	p[packet.IDServerBoundDiagnostics] = func() packet.Packet { return &ServerBoundDiagnostics{} }
	p[packet.IDClientBoundDataDrivenUIShowScreen] = func() packet.Packet { return &packet.ClientBoundDataDrivenUIShowScreen{} }
	p[packet.IDClientBoundDataDrivenUICloseScreen] = func() packet.Packet { return &packet.ClientBoundDataDrivenUICloseScreen{} }

	p[packet.IDResourcePacksReadyForValidation] = func() packet.Packet { return &packet.ResourcePacksReadyForValidation{} }
	p[packet.IDSyncWorldClocks] = func() packet.Packet { return &packet.SyncWorldClocks{} }
}

func NewClientPool() packet.Pool {
	p := v924.NewClientPool()
	applyDeltas944(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v924.NewServerPool()
	applyDeltas944(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v924.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v924.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		return []packet.Packet{FromLatestStartGame944(pk)}, true
	case *packet.VoxelShapes, *packet.UpdateClientInputLocks, *packet.ClientBoundDataDrivenUIShowScreen,
		*packet.ClientBoundDataDrivenUICloseScreen, *packet.SyncWorldClocks:
		return []packet.Packet{pk}, true
	case *packet.ItemRegistry:
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items944()}}, true
	case *packet.LocatorBar, *packet.ClientBoundDataStore, *packet.ClientBoundAttributeLayerSync, *packet.PartyChanged:

		return nil, true
	case *packet.ServerBoundDiagnostics:
		return []packet.Packet{&ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
			AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
			AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
			AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
			AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
			MemoryCategoryValues: pk.MemoryCategoryValues}}, true
	}
	return nil, false
}

func convertToLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *StartGame:
		return []packet.Packet{ToLatestStartGame944(pk)}, true
	case *packet.VoxelShapes, *packet.UpdateClientInputLocks, *packet.ClientBoundDataDrivenUIShowScreen,
		*packet.ClientBoundDataDrivenUICloseScreen, *packet.SyncWorldClocks, *packet.ResourcePacksReadyForValidation:
		return []packet.Packet{pk}, true
	case *ServerBoundDiagnostics:
		return []packet.Packet{&packet.ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
			AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
			AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
			AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
			AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
			MemoryCategoryValues: pk.MemoryCategoryValues}}, true
	}
	return nil, false
}
