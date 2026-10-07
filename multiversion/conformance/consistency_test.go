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
	"reflect"
	"strings"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/world"
)

func latestPools() map[uint32]func() packet.Packet {
	all := map[uint32]func() packet.Packet{}
	for id, fn := range packet.NewClientPool() {
		all[id] = fn
	}
	for id, fn := range packet.NewServerPool() {
		all[id] = fn
	}
	return all
}

func isLatest(pk packet.Packet) bool {
	return strings.HasSuffix(reflect.TypeOf(pk).Elem().PkgPath(), "/protocol/packet")
}

var knownDisagreements = map[int32]map[uint32]bool{
	671:  {130: true, 131: true, 135: true, 147: true}, // 130/131 swapped at 671
	685:  {135: true, 147: true},
	686:  {135: true, 147: true},
	712:  {135: true, 147: true},
	729:  {135: true, 147: true},
	748:  {135: true, 147: true},
	766:  {135: true, 147: true},
	776:  {135: true, 147: true},
	786:  {135: true, 147: true},
	800:  {122: true, 135: true, 147: true, 300: true, 316: true},
	818:  {135: true, 147: true},
	819:  {135: true, 147: true},
	827:  {135: true, 147: true, 307: true},
	844:  {135: true, 147: true, 307: true},
	859:  {69: true, 135: true, 147: true, 185: true, 300: true, 307: true},
	898:  {69: true, 135: true, 147: true, 185: true, 300: true, 307: true},
	924:  {69: true, 135: true, 147: true, 185: true, 300: true, 307: true},
	944:  {69: true, 135: true, 147: true, 185: true, 300: true, 307: true},
	975:  {69: true, 135: true, 147: true, 185: true, 300: true, 307: true},
	1001: {69: true, 135: true, 147: true, 185: true, 300: true, 307: true},
}

func TestPoolAndConverterAgree(t *testing.T) {

	finaliseOnce.Do(world.DefaultBlockRegistry.Finalize)

	latest := latestPools()
	for _, proto := range protocols() {
		pool := proto.Packets(true)
		for id, newOld := range pool {
			newLatest, ok := latest[id]
			if !ok {
				continue
			}
			decoded := newOld()
			out := proto.ConvertFromLatest(newLatest(), nil)
			if len(out) != 1 {
				continue
			}
			if got, want := isLatest(out[0]), isLatest(decoded); got != want {
				if knownDisagreements[proto.ID()][id] {
					continue
				}
				t.Errorf("protocol %v, packet id %v: pool decodes it as %T but ConvertFromLatest encodes it "+
					"as %T — pool.go and convert.go disagree about whether this packet still uses an older "+
					"wire at this protocol; exactly one of them is wrong", proto.ID(), id, decoded, out[0])
			}
		}
	}
}
