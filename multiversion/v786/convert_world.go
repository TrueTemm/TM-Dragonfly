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

package v786

import (
	_ "embed"

	"github.com/df-mc/dragonfly/multiversion/blockpalette"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

//go:embed biome_definitions_786.dat
var legacyBiomeDefinitionsNBT786 []byte

func fromLatestBiomeDefinitionList(*packet.BiomeDefinitionList) *BiomeDefinitionList {
	return &BiomeDefinitionList{SerialisedBiomeDefinitions: legacyBiomeDefinitionsNBT786}
}
func toLatestBiomeDefinitionList(pk *BiomeDefinitionList) *packet.BiomeDefinitionList {
	return &packet.BiomeDefinitionList{}
}

func toLatestLevelChunk(pk *LevelChunk) *packet.LevelChunk {
	out := &packet.LevelChunk{
		Position: pk.Position, Dimension: pk.Dimension, SubChunkCount: pk.SubChunkCount,
		CacheEnabled: pk.CacheEnabled, BlobHashes: pk.BlobHashes, RawPayload: pk.RawPayload,
	}
	if pk.SubChunkCount == SubChunkRequestModeLimited786 {
		out.SubChunkLimit = protocol.Option(int32(pk.HighestSubChunk))
	}
	return out
}
func fromLatestLevelChunk(proto uint32, pk *packet.LevelChunk) *LevelChunk {
	payload := pk.RawPayload

	if !pk.CacheEnabled && pk.SubChunkCount > 0 && pk.SubChunkCount < 0xfffffffe {
		payload = blockpalette.RemapChunkBlockPalette(pk.RawPayload, pk.SubChunkCount, proto)
	}
	out := &LevelChunk{
		Position: pk.Position, Dimension: pk.Dimension, SubChunkCount: pk.SubChunkCount,
		CacheEnabled: pk.CacheEnabled, BlobHashes: pk.BlobHashes, RawPayload: payload,
	}

	if v, ok := pk.SubChunkLimit.Value(); ok {
		out.SubChunkCount = SubChunkRequestModeLimited786
		out.HighestSubChunk = uint16(v)
	}
	return out
}

func toLatestSubChunk(pk *SubChunk) *packet.SubChunk {
	entries := make([]protocol.SubChunkEntry, len(pk.SubChunkEntries))
	for i, e := range pk.SubChunkEntries {
		entries[i] = toLatestSubChunkEntry786(e)
	}
	return &packet.SubChunk{CacheEnabled: pk.CacheEnabled, Dimension: pk.Dimension, Position: pk.Position, SubChunkEntries: entries}
}
func fromLatestSubChunk(proto uint32, pk *packet.SubChunk) *SubChunk {
	entries := make([]SubChunkEntry786, len(pk.SubChunkEntries))
	for i, e := range pk.SubChunkEntries {
		entries[i] = fromLatestSubChunkEntry786(proto, e)
	}
	return &SubChunk{CacheEnabled: pk.CacheEnabled, Dimension: pk.Dimension, Position: pk.Position, SubChunkEntries: entries}
}

func toLatestSubChunkEntry786(x SubChunkEntry786) protocol.SubChunkEntry {
	out := protocol.SubChunkEntry{Offset: x.Offset, Result: x.Result, HeightMapType: x.HeightMapType}
	if len(x.RawPayload) > 0 {
		out.RawPayload = protocol.Option(x.RawPayload)
	}
	if x.HeightMapType == protocol.HeightMapDataHasData {
		out.HeightMapData = protocol.Option(heightMapOf(x.HeightMapData))
	}
	return out
}

func heightMapOf(flat []int8) protocol.HeightMap {
	var out protocol.HeightMap
	for i, h := range flat {
		if i >= 256 {
			break
		}
		out[i>>4][i&15] = h
	}
	return out
}

func flatHeightMap(m protocol.HeightMap) []int8 {
	out := make([]int8, 0, 256)
	for z := range m {
		out = append(out, m[z][:]...)
	}
	return out
}
func fromLatestSubChunkEntry786(proto uint32, x protocol.SubChunkEntry) SubChunkEntry786 {
	out := SubChunkEntry786{Offset: x.Offset, Result: x.Result, HeightMapType: x.HeightMapType}
	if v, ok := x.RawPayload.Value(); ok {

		out.RawPayload = blockpalette.RemapChunkBlockPalette(v, 1, proto)
	}
	if v, ok := x.HeightMapData.Value(); ok {
		out.HeightMapData = flatHeightMap(v)
	}
	if v, ok := x.BlobHash.Value(); ok {
		out.BlobHash = v
	}
	return out
}

func ToLatestLevelSoundEvent786(proto uint32, pk *LevelSoundEvent) *packet.LevelSoundEvent {
	name := soundIDToName786[pk.SoundType]
	extra := pk.ExtraData
	if blockSoundExtraData[name] {
		extra = int32(blockpalette.RuntimeIDFor(proto, uint32(pk.ExtraData)))
	}
	return &packet.LevelSoundEvent{
		SoundType: name, Position: pk.Position, ExtraData: extra, EntityType: pk.EntityType,
		BabyMob: pk.BabyMob, DisableRelativeVolume: pk.DisableRelativeVolume, EntityUniqueID: pk.EntityUniqueID,
	}
}
func FromLatestLevelSoundEvent786(proto uint32, pk *packet.LevelSoundEvent) *LevelSoundEvent {
	extra := pk.ExtraData
	if blockSoundExtraData[pk.SoundType] {
		extra = int32(blockpalette.HashFor(proto, uint32(pk.ExtraData)))
	}
	return &LevelSoundEvent{
		SoundType: soundNameToID786[pk.SoundType], Position: pk.Position, ExtraData: extra, EntityType: pk.EntityType,
		BabyMob: pk.BabyMob, DisableRelativeVolume: pk.DisableRelativeVolume, EntityUniqueID: pk.EntityUniqueID,
	}
}

func toLatestPlaySound(pk *PlaySound) *packet.PlaySound {
	return &packet.PlaySound{SoundName: pk.SoundName, Position: pk.Position, Volume: pk.Volume, Pitch: pk.Pitch}
}
func fromLatestPlaySound(pk *packet.PlaySound) *PlaySound {
	return &PlaySound{SoundName: pk.SoundName, Position: pk.Position, Volume: pk.Volume, Pitch: pk.Pitch}
}

func toLatestJigsawStructureData(pk *JigsawStructureData) *packet.JigsawStructureData {
	return &packet.JigsawStructureData{}
}

func fromLatestJigsawStructureData(pk *packet.JigsawStructureData) *JigsawStructureData {
	data, err := nbt.MarshalEncoding(pk.StructureData, nbt.NetworkLittleEndian)
	if err != nil {

		data, _ = nbt.MarshalEncoding(map[string]any{}, nbt.NetworkLittleEndian)
	}
	return &JigsawStructureData{StructureData: data}
}

func toLatestOnScreenTextureAnimation(pk *OnScreenTextureAnimation) *packet.OnScreenTextureAnimation {
	return &packet.OnScreenTextureAnimation{AnimationType: uint32(pk.AnimationType)}
}
func fromLatestOnScreenTextureAnimation(pk *packet.OnScreenTextureAnimation) *OnScreenTextureAnimation {
	return &OnScreenTextureAnimation{AnimationType: int32(pk.AnimationType)}
}

func toLatestClientBoundMapItemData(pk *ClientBoundMapItemData) *packet.ClientBoundMapItemData {
	out := &packet.ClientBoundMapItemData{
		MapID: pk.MapID, Dimension: pk.Dimension, LockedMap: pk.LockedMap, Origin: pk.Origin,
	}
	if pk.UpdateFlags&(MapUpdateFlagInitialisation|MapUpdateFlagDecoration|MapUpdateFlagTexture) != 0 {
		out.Scale = protocol.Option(pk.Scale)
	}
	if pk.UpdateFlags&MapUpdateFlagInitialisation != 0 {
		out.MapsIncludedIn = protocol.Option(pk.MapsIncludedIn)
	}
	if pk.UpdateFlags&MapUpdateFlagDecoration != 0 {
		out.TrackedObjects = protocol.Option(pk.TrackedObjects)
		out.Decorations = protocol.Option(pk.Decorations)
	}
	if pk.UpdateFlags&MapUpdateFlagTexture != 0 {
		out.Width, out.Height = protocol.Option(pk.Width), protocol.Option(pk.Height)
		out.XOffset, out.YOffset = protocol.Option(pk.XOffset), protocol.Option(pk.YOffset)
		out.Pixels = protocol.Option(pk.Pixels)
	}
	return out
}
func fromLatestClientBoundMapItemData(pk *packet.ClientBoundMapItemData) *ClientBoundMapItemData {
	out := &ClientBoundMapItemData{MapID: pk.MapID, Dimension: pk.Dimension, LockedMap: pk.LockedMap, Origin: pk.Origin}
	if v, ok := pk.Scale.Value(); ok {
		out.Scale = v
	}
	if v, ok := pk.MapsIncludedIn.Value(); ok {
		out.UpdateFlags |= MapUpdateFlagInitialisation
		out.MapsIncludedIn = v
	}
	if v, ok := pk.TrackedObjects.Value(); ok {
		out.UpdateFlags |= MapUpdateFlagDecoration
		out.TrackedObjects = v
	}
	if v, ok := pk.Decorations.Value(); ok {
		out.Decorations = v
	}
	if v, ok := pk.Pixels.Value(); ok {
		out.UpdateFlags |= MapUpdateFlagTexture
		out.Pixels = v
	}
	if v, ok := pk.Width.Value(); ok {
		out.Width = v
	}
	if v, ok := pk.Height.Value(); ok {
		out.Height = v
	}
	if v, ok := pk.XOffset.Value(); ok {
		out.XOffset = v
	}
	if v, ok := pk.YOffset.Value(); ok {
		out.YOffset = v
	}
	return out
}

func fromLatestCreativeContent(pk *packet.CreativeContent) *CreativeContent {
	out := &CreativeContent{Items: pk.Items}
	out.Groups = make([]CreativeGroup786, len(pk.Groups))
	for i, g := range pk.Groups {
		out.Groups[i] = CreativeGroup786{Category: int32(g.Category), Name: g.Name, Icon: g.Icon}
	}
	return out
}
func toLatestCreativeContent(pk *CreativeContent) *packet.CreativeContent {
	out := &packet.CreativeContent{Items: pk.Items}
	out.Groups = make([]protocol.CreativeGroup, len(pk.Groups))
	for i, g := range pk.Groups {
		out.Groups[i] = protocol.CreativeGroup{Category: byte(g.Category), Name: g.Name, Icon: g.Icon}
	}
	return out
}

func ToLatestCreativeContent786(pk *CreativeContent) *packet.CreativeContent {
	return toLatestCreativeContent(pk)
}
func FromLatestCreativeContent786(pk *packet.CreativeContent) *CreativeContent {
	return fromLatestCreativeContent(pk)
}

func BlockSoundCarriesRuntimeID(sound string) bool { return blockSoundExtraData[sound] }
