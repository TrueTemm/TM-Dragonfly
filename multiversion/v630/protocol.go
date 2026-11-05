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

package v630

import (
	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/df-mc/dragonfly/multiversion/itemdata"
	"github.com/df-mc/dragonfly/multiversion/v649"
	"github.com/df-mc/dragonfly/multiversion/v686"
	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v776"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{ Version int32 }

func (p Protocol) ID() int32 {
	if p.Version == 0 {
		return 630
	}
	return p.Version
}
func (p Protocol) Ver() string { return versions[uint32(p.ID())] }

func (p Protocol) Packets(listener bool) packet.Pool {
	pool := (v649.Protocol{}).Packets(listener)
	delete(pool, v786.IDSetHud)
	pool[v786.IDCraftingData] = func() packet.Packet { return &CraftingData{} }
	pool[v786.IDLevelChunk] = func() packet.Packet { return &LevelChunk{} }
	pool[v786.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	pool[v786.IDPlayerList] = func() packet.Packet { return &PlayerList{} }
	if p.ID() < 630 {
		pool[53] = func() packet.Packet { return &CraftingEvent{} }
		delete(pool, v786.IDPlayerToggleCrafterSlotRequest)
		delete(pool, v786.IDSetPlayerInventoryOptions)
	}
	if p.ID() < 618 {
		delete(pool, 305)
	}
	if p.ID() < 594 {
		delete(pool, 304)
	}
	return pool
}

func (p Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return v786.Reader786{Reader: protocol.NewReader(r, shieldID, enableLimits), Proto: uint32(p.ID())}
}

func (p Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return v786.Writer786{Writer: protocol.NewWriter(w, shieldID), Proto: uint32(p.ID())}
}

func (p Protocol) ConvertToLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	switch pk := pk.(type) {
	case *packet.ClientCacheStatus:
		cp := *pk
		cp.Enabled = false
		return []packet.Packet{&cp}
	case *v776.LevelSoundEvent:
		name, ok := registryFor(uint32(p.ID())).soundNames[pk.SoundType]
		if !ok {
			return nil
		}
		extra := pk.ExtraData
		if v786.BlockSoundCarriesRuntimeID(name) && extra != -1 {
			extra = int32(blockpalette.RuntimeIDFor(uint32(p.ID()), uint32(extra)))
		}
		return []packet.Packet{&packet.LevelSoundEvent{SoundType: name, Position: pk.Position, ExtraData: extra, EntityType: pk.EntityType, BabyMob: pk.BabyMob, DisableRelativeVolume: pk.DisableRelativeVolume}}

	case *PlayerAuthInput:
		mask := uint64(1<<45 - 1)
		if p.ID() < 594 {
			mask = 1<<39 - 1
		} else if p.ID() < 618 {
			mask = 1<<42 - 1
		} else if p.ID() < 622 {
			mask = 1<<44 - 1
		}
		pk.InputData &= mask
		return []packet.Packet{v729.ToLatestPlayerAuthInput(&pk.PlayerAuthInput)}
	case *LevelChunk, *PlayerList, *CraftingEvent, *v786.TickSync:
		return nil
	}
	if pk.ID() == 71 {
		return nil
	}
	if out, ok := v686.ToLatestShared(uint32(p.ID()), pk); ok {
		return out
	}
	return []packet.Packet{pk}
}

func (p Protocol) ConvertFromLatest(pk packet.Packet, conn *minecraft.Conn) []packet.Packet {
	switch pk := pk.(type) {
	case *packet.PlayStatus:
		if pk.Status == packet.PlayStatusPlayerSpawn {
			data := registryFor(uint32(p.ID())).biomes
			return []packet.Packet{&v786.BiomeDefinitionList{SerialisedBiomeDefinitions: data}, pk} // required before spawn
		}
		return []packet.Packet{pk}
	case *packet.CraftingData:
		return []packet.Packet{fromLatestCrafting(uint32(p.ID()), pk)}
	case *packet.CreativeContent:
		cp := *pk
		cp.Items = make([]protocol.CreativeItem, 0, len(pk.Items))
		for _, item := range pk.Items {
			if id, _ := itemdata.TranslateStack(uint32(p.ID()), item.Item.NetworkID, item.Item.MetadataValue); id != 0 {
				cp.Items = append(cp.Items, item) // preserve creative network IDs
			}
		}
		out, _ := v686.FromLatestShared(uint32(p.ID()), &cp)
		return out
	case *packet.ResourcePackStack:
		out := v786.FromLatestResourcePackStack786(pk)
		out.BaseGameVersion = p.Ver()
		out.TexturePacks = make([]protocol.StackResourcePack, 0, len(pk.TexturePacks))
		for _, pack := range pk.TexturePacks {
			if !versionedVanillaPacks[pack.UUID] {
				out.TexturePacks = append(out.TexturePacks, pack)
			}
		}
		return []packet.Packet{out}
	case *packet.BiomeDefinitionList:
		return []packet.Packet{&v786.BiomeDefinitionList{SerialisedBiomeDefinitions: registryFor(uint32(p.ID())).biomes}}
	case *packet.AvailableActorIdentifiers:
		return []packet.Packet{&packet.AvailableActorIdentifiers{SerialisedEntityIdentifiers: registryFor(uint32(p.ID())).entities}}
	case *packet.LevelSoundEvent:
		return fromLatestSound(uint32(p.ID()), pk)
	case *packet.LevelEvent:
		return fromLatestEvent(uint32(p.ID()), pk)
	case *packet.AddActor:
		cp := *pk
		if variant, ok := pk.EntityMetadata[protocol.EntityDataKeyVariant].(int32); ok && pk.EntityType == "minecraft:falling_block" {
			cp.EntityMetadata = make(protocol.EntityMetadata, len(pk.EntityMetadata))
			for key, value := range pk.EntityMetadata {
				cp.EntityMetadata[key] = value
			}
			cp.EntityMetadata[protocol.EntityDataKeyVariant] = int32(blockpalette.HashFor(uint32(p.ID()), uint32(variant)))
		}
		out, _ := v686.FromLatestShared(uint32(p.ID()), &cp)
		return out

	case *packet.PlayerToggleCrafterSlotRequest, *packet.SetPlayerInventoryOptions:
		if p.ID() < 630 {
			return nil
		}
	case *packet.SetHud, *packet.ClientBoundCloseForm, *packet.AwardAchievement:
		return nil
	case *packet.MapCreateLockedCopy, *packet.OnScreenTextureAnimation:
		return (v649.Protocol{}).ConvertFromLatest(pk, conn)
	case *packet.AvailableCommands:
		out, _ := v686.FromLatestShared(uint32(p.ID()), v649.RemapCommands(pk))
		return out
	case *packet.StartGame:
		out, _ := v686.FromLatestShared(uint32(p.ID()), pk)
		sg := out[0].(*v766.StartGame)
		sg.BaseGameVersion, sg.GameVersion = p.Ver(), p.Ver()
		sg.UseBlockNetworkIDHashes = false // legacy palette indices
		items := itemdata.ItemsFor(uint32(p.ID()))
		sg.Items = make([]v766.ItemEntry, len(items))
		for i, it := range items {
			sg.Items[i] = v766.ItemEntry{Name: it.Name, RuntimeID: it.RuntimeID}
		}
		return out
	case *packet.PlayerAuthInput:
		return []packet.Packet{&PlayerAuthInput{PlayerAuthInput: *v729.FromLatestPlayerAuthInput(pk)}}
	case *packet.PlayerList:
		return []packet.Packet{&PlayerList{PlayerList: *v786.FromLatestPlayerList(pk)}}
	}
	if out, ok := v686.FromLatestShared(uint32(p.ID()), pk); ok {
		if len(out) == 1 {
			if chunk, ok := out[0].(*v786.LevelChunk); ok {
				return []packet.Packet{&LevelChunk{LevelChunk: *chunk}}
			}
		}
		return out
	}
	return []packet.Packet{pk}
}
