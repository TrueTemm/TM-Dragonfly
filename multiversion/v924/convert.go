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

package v924

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v898"
)

func applyDeltas924(p packet.Pool) {
	p[v898.IDStartGame] = func() packet.Packet { return &StartGame{} }
	p[v898.IDText] = func() packet.Packet { return &packet.Text{} }
	p[packet.IDBookEdit] = func() packet.Packet { return &packet.BookEdit{} }
	p[packet.IDServerBoundDiagnostics] = func() packet.Packet { return &ServerBoundDiagnostics{} }

	p[packet.IDVoxelShapes] = func() packet.Packet { return &VoxelShapes{} }
	p[packet.IDClientBoundDataDrivenUIReload] = func() packet.Packet { return &packet.ClientBoundDataDrivenUIReload{} }
	p[packet.IDClientBoundTextureShift] = func() packet.Packet { return &packet.ClientBoundTextureShift{} }
}

func NewClientPool() packet.Pool {
	p := v898.NewClientPool()
	applyDeltas924(p)
	return p
}

func NewServerPool() packet.Pool {
	p := v898.NewServerPool()
	applyDeltas924(p)
	return p
}

func FromLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertFromLatest(proto, pk); ok {
		return out, true
	}
	return v898.FromLatestShared(proto, pk)
}

func ToLatestShared(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	if out, ok := convertToLatest(proto, pk); ok {
		return out, true
	}
	return v898.ToLatestShared(proto, pk)
}

func convertFromLatest(proto uint32, pk packet.Packet) ([]packet.Packet, bool) {
	switch pk := pk.(type) {
	case *packet.StartGame:
		return []packet.Packet{FromLatestStartGame924(pk)}, true
	case *packet.Text, *packet.BookEdit:
		return []packet.Packet{pk}, true
	case *packet.VoxelShapes:
		return []packet.Packet{&VoxelShapes{Shapes: pk.Shapes, NameMap: pk.NameMap}}, true
	case *packet.ItemRegistry:
		return []packet.Packet{&packet.ItemRegistry{Items: itemdata.Items924()}}, true
	case *packet.ClientBoundDataDrivenUIReload, *packet.ClientBoundTextureShift:
		return []packet.Packet{pk}, true
	case *packet.CameraSpline, *packet.CameraAimAssistActorPriority, *packet.ClientBoundDataDrivenUIShowScreen,
		*packet.ClientBoundDataDrivenUICloseScreen, *packet.GraphicsOverrideParameter:

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
		return []packet.Packet{ToLatestStartGame924(pk)}, true
	case *packet.Text, *packet.BookEdit:
		return []packet.Packet{pk}, true
	case *ServerBoundDiagnostics:
		return []packet.Packet{&packet.ServerBoundDiagnostics{AverageFramesPerSecond: pk.AverageFramesPerSecond,
			AverageServerSimTickTime: pk.AverageServerSimTickTime, AverageClientSimTickTime: pk.AverageClientSimTickTime,
			AverageBeginFrameTime: pk.AverageBeginFrameTime, AverageInputTime: pk.AverageInputTime,
			AverageRenderTime: pk.AverageRenderTime, AverageEndFrameTime: pk.AverageEndFrameTime,
			AverageRemainderTimePercent: pk.AverageRemainderTimePercent, AverageUnaccountedTimePercent: pk.AverageUnaccountedTimePercent,
			MemoryCategoryValues: pk.MemoryCategoryValues}}, true
	case *VoxelShapes:
		return []packet.Packet{&packet.VoxelShapes{Shapes: pk.Shapes, NameMap: pk.NameMap}}, true
	}
	return nil, false
}
