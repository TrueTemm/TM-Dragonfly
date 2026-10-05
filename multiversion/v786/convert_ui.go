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
	"fmt"
	"image/color"
	"strings"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func toLatestText(pk *Text) *packet.Text {
	return &packet.Text{
		TextType: pk.TextType, NeedsTranslation: pk.NeedsTranslation, SourceName: pk.SourceName,
		Message: pk.Message, Parameters: pk.Parameters, XUID: pk.XUID, PlatformChatID: pk.PlatformChatID,
		FilteredMessage: protocol.Option(pk.FilteredMessage),
	}
}
func fromLatestText(pk *packet.Text) *Text {
	filtered, _ := pk.FilteredMessage.Value()
	return &Text{
		TextType: pk.TextType, NeedsTranslation: pk.NeedsTranslation, SourceName: pk.SourceName,
		Message: pk.Message, Parameters: pk.Parameters, XUID: pk.XUID, PlatformChatID: pk.PlatformChatID,
		FilteredMessage: filtered,
	}
}

func toLatestBossEvent(pk *BossEvent) *packet.BossEvent {
	return &packet.BossEvent{
		BossEntityUniqueID: pk.BossEntityUniqueID, EventType: uint8(pk.EventType),
		BossBarTitle: pk.BossBarTitle, FilteredBossBarTitle: pk.FilteredBossBarTitle,
		HealthPercentage: pk.HealthPercentage, Colour: uint8(pk.Colour), Overlay: uint8(pk.Overlay),
	}
}
func fromLatestBossEvent(pk *packet.BossEvent) *BossEvent {
	return &BossEvent{
		BossEntityUniqueID: pk.BossEntityUniqueID, EventType: uint32(pk.EventType),
		BossBarTitle: pk.BossBarTitle, FilteredBossBarTitle: pk.FilteredBossBarTitle,
		HealthPercentage: pk.HealthPercentage, Colour: uint32(pk.Colour), Overlay: uint32(pk.Overlay),
	}
}

func toLatestCameraAimAssist(pk *CameraAimAssist) *packet.CameraAimAssist {
	return &packet.CameraAimAssist{
		Preset: pk.Preset, Angle: pk.Angle, Distance: pk.Distance, TargetMode: pk.TargetMode, Action: pk.Action,
	}
}
func fromLatestCameraAimAssist(pk *packet.CameraAimAssist) *CameraAimAssist {
	return &CameraAimAssist{Preset: pk.Preset, Angle: pk.Angle, Distance: pk.Distance, TargetMode: pk.TargetMode, Action: pk.Action}
}

func toLatestCameraAimAssistPresets(pk *CameraAimAssistPresets) *packet.CameraAimAssistPresets {
	return &packet.CameraAimAssistPresets{Operation: pk.Operation}
}
func fromLatestCameraAimAssistPresets(pk *packet.CameraAimAssistPresets) *CameraAimAssistPresets {
	return &CameraAimAssistPresets{Operation: pk.Operation}
}

func toLatestCameraInstruction(pk *CameraInstruction) *packet.CameraInstruction {
	return &packet.CameraInstruction{Set: pk.Set, Clear: pk.Clear, Fade: pk.Fade, Target: pk.Target, RemoveTarget: pk.RemoveTarget}
}
func fromLatestCameraInstruction(pk *packet.CameraInstruction) *CameraInstruction {
	return &CameraInstruction{Set: pk.Set, Clear: pk.Clear, Fade: pk.Fade, Target: pk.Target, RemoveTarget: pk.RemoveTarget}
}

func toLatestCommandOutput(pk *CommandOutput) *packet.CommandOutput {
	out := &packet.CommandOutput{
		CommandOrigin: pk.CommandOrigin, OutputType: pk.OutputType, SuccessCount: pk.SuccessCount,
		OutputMessages: pk.OutputMessages,
	}
	if pk.OutputType == CommandOutputTypeDataSet {
		out.DataSet = protocol.Option(pk.DataSet)
	}
	return out
}
func fromLatestCommandOutput(pk *packet.CommandOutput) *CommandOutput {
	out := &CommandOutput{
		CommandOrigin: pk.CommandOrigin, OutputType: pk.OutputType, SuccessCount: pk.SuccessCount,
		OutputMessages: pk.OutputMessages,
	}
	if v, ok := pk.DataSet.Value(); ok {
		out.DataSet = v
	}
	return out
}

func toLatestCommandRequest(pk *CommandRequest) *packet.CommandRequest {
	return &packet.CommandRequest{CommandLine: pk.CommandLine, CommandOrigin: pk.CommandOrigin, Internal: pk.Internal}
}
func fromLatestCommandRequest(pk *packet.CommandRequest) *CommandRequest {
	return &CommandRequest{CommandLine: pk.CommandLine, CommandOrigin: pk.CommandOrigin, Internal: pk.Internal}
}

func ToLatestCommandRequest786(pk *CommandRequest) *packet.CommandRequest {
	return toLatestCommandRequest(pk)
}
func FromLatestCommandRequest786(pk *packet.CommandRequest) *CommandRequest {
	return fromLatestCommandRequest(pk)
}
func ToLatestCommandOutput786(pk *CommandOutput) *packet.CommandOutput {
	return toLatestCommandOutput(pk)
}
func FromLatestCommandOutput786(pk *packet.CommandOutput) *CommandOutput {
	return fromLatestCommandOutput(pk)
}
func ToLatestAvailableCommands786(pk *AvailableCommands) *packet.AvailableCommands {
	return toLatestAvailableCommands(pk)
}
func FromLatestAvailableCommands786(pk *packet.AvailableCommands) *AvailableCommands {
	return fromLatestAvailableCommands(pk)
}

func fromLatestAvailableCommands(pk *packet.AvailableCommands) *AvailableCommands {
	out := &AvailableCommands{
		EnumValues:              pk.EnumValues,
		ChainedSubcommandValues: pk.ChainedSubcommandValues,
		Suffixes:                pk.Suffixes,
		DynamicEnums:            pk.DynamicEnums,
		Constraints:             pk.Constraints,
	}
	out.Enums = make([]CommandEnum786, len(pk.Enums))
	for i, e := range pk.Enums {
		vi := make([]uint, len(e.ValueIndices))
		for j, v := range e.ValueIndices {
			vi[j] = uint(v)
		}
		out.Enums[i] = CommandEnum786{Type: e.Type, ValueIndices: vi}
	}
	out.ChainedSubcommands = make([]ChainedSubcommand786, len(pk.ChainedSubcommands))
	for i, cs := range pk.ChainedSubcommands {
		vals := make([]ChainedSubcommandValue786, len(cs.Values))
		for j, v := range cs.Values {
			vals[j] = ChainedSubcommandValue786{Index: uint16(v.Index), Value: uint16(v.Value)}
		}
		out.ChainedSubcommands[i] = ChainedSubcommand786{Name: cs.Name, Values: vals}
	}
	out.Commands = make([]Command786, len(pk.Commands))
	for i, c := range pk.Commands {
		offs := make([]uint16, len(c.ChainedSubcommandOffsets))
		for j, o := range c.ChainedSubcommandOffsets {
			offs[j] = uint16(o)
		}
		out.Commands[i] = Command786{
			Name: c.Name, Description: c.Description, Flags: c.Flags, PermissionLevel: c.PermissionLevel,
			AliasesOffset: c.AliasesOffset, ChainedSubcommandOffsets: offs, Overloads: c.Overloads,
		}
	}
	return out
}

func toLatestAvailableCommands(pk *AvailableCommands) *packet.AvailableCommands {
	out := &packet.AvailableCommands{
		EnumValues:              pk.EnumValues,
		ChainedSubcommandValues: pk.ChainedSubcommandValues,
		Suffixes:                pk.Suffixes,
		DynamicEnums:            pk.DynamicEnums,
		Constraints:             pk.Constraints,
	}
	out.Enums = make([]protocol.CommandEnum, len(pk.Enums))
	for i, e := range pk.Enums {
		vi := make([]uint32, len(e.ValueIndices))
		for j, v := range e.ValueIndices {
			vi[j] = uint32(v)
		}
		out.Enums[i] = protocol.CommandEnum{Type: e.Type, ValueIndices: vi}
	}
	out.ChainedSubcommands = make([]protocol.ChainedSubcommand, len(pk.ChainedSubcommands))
	for i, cs := range pk.ChainedSubcommands {
		vals := make([]protocol.ChainedSubcommandValue, len(cs.Values))
		for j, v := range cs.Values {
			vals[j] = protocol.ChainedSubcommandValue{Index: uint32(v.Index), Value: uint32(v.Value)}
		}
		out.ChainedSubcommands[i] = protocol.ChainedSubcommand{Name: cs.Name, Values: vals}
	}
	out.Commands = make([]protocol.Command, len(pk.Commands))
	for i, c := range pk.Commands {
		offs := make([]uint32, len(c.ChainedSubcommandOffsets))
		for j, o := range c.ChainedSubcommandOffsets {
			offs[j] = uint32(o)
		}
		out.Commands[i] = protocol.Command{
			Name: c.Name, Description: c.Description, Flags: c.Flags, PermissionLevel: c.PermissionLevel,
			AliasesOffset: c.AliasesOffset, ChainedSubcommandOffsets: offs, Overloads: c.Overloads,
		}
	}
	return out
}

func toLatestPlayerList(pk *PlayerList) *packet.PlayerList {
	entries := make([]protocol.PlayerListEntry, len(pk.Entries))
	for i, e := range pk.Entries {
		entries[i] = protocol.PlayerListEntry{
			ActionType: pk.ActionType, UUID: e.UUID, EntityUniqueID: e.EntityUniqueID, Username: e.Username,
			XUID: e.XUID, PlatformChatID: e.PlatformChatID, BuildPlatform: e.BuildPlatform,
			Skin: toLatestSkin786(e.Skin), Teacher: e.Teacher, Host: e.Host, SubClient: e.SubClient,
		}
	}
	return &packet.PlayerList{Entries: entries}
}

func fromLatestPlayerList(pk *packet.PlayerList) *PlayerList {
	action := byte(PlayerListActionAdd)
	if len(pk.Entries) > 0 {
		action = pk.Entries[0].ActionType
	}
	entries := make([]PlayerListEntry786, len(pk.Entries))
	for i, e := range pk.Entries {
		entries[i] = PlayerListEntry786{
			UUID: e.UUID, EntityUniqueID: e.EntityUniqueID, Username: e.Username, XUID: e.XUID,
			PlatformChatID: e.PlatformChatID, BuildPlatform: e.BuildPlatform, Skin: fromLatestSkin786(e.Skin),
			Teacher: e.Teacher, Host: e.Host, SubClient: e.SubClient,
		}
	}
	return &PlayerList{ActionType: action, Entries: entries}
}

func fromLatestSkin786(s protocol.Skin) Skin786 {
	animations := make([]SkinAnimation786, len(s.Animations))
	for i, a := range s.Animations {
		animations[i] = SkinAnimation786{
			ImageWidth: a.ImageWidth, ImageHeight: a.ImageHeight, ImageData: a.ImageData,
			AnimationType: a.AnimationType, FrameCount: a.FrameCount, ExpressionType: a.ExpressionType,
		}
	}
	pieces := make([]PersonaPiece786, len(s.PersonaPieces))
	for i, p := range s.PersonaPieces {
		pieces[i] = PersonaPiece786{
			PieceID: p.PieceID, PieceType: personaPieceTypeName(p.PieceType), PackID: p.PackID.String(),
			Default: p.Default, ProductID: p.ProductID,
		}
	}
	tints := make([]PersonaPieceTintColour786, len(s.PieceTintColours))
	for i, t := range s.PieceTintColours {

		n := len(t.Colours)
		for n > 0 && t.Colours[n-1] == (color.RGBA{}) {
			n--
		}
		colours := make([]string, n)
		for j := 0; j < n; j++ {
			colours[j] = argbHexString(t.Colours[j])
		}
		tints[i] = PersonaPieceTintColour786{PieceType: t.PieceType, Colours: colours}
	}
	armSize := "slim"
	if s.ArmSize == protocol.ArmSizeWide {
		armSize = "wide"
	}
	return Skin786{
		SkinID: s.SkinID, PlayFabID: s.PlayFabID, SkinResourcePatch: s.SkinResourcePatch,
		SkinImageWidth: s.SkinImageWidth, SkinImageHeight: s.SkinImageHeight, SkinData: s.SkinData,
		Animations: animations, CapeImageWidth: s.CapeImageWidth, CapeImageHeight: s.CapeImageHeight,
		CapeData: s.CapeData, SkinGeometry: s.SkinGeometry, GeometryDataEngineVersion: s.GeometryDataEngineVersion,
		AnimationData: s.AnimationData, CapeID: s.CapeID, FullID: s.FullID, ArmSize: armSize,
		SkinColour: rgbHexString(s.SkinColour), PersonaPieces: pieces, PieceTintColours: tints,
		PremiumSkin: s.PremiumSkin, PersonaSkin: s.PersonaSkin, PersonaCapeOnClassicSkin: s.PersonaCapeOnClassicSkin,
		PrimaryUser: s.PrimaryUser, OverrideAppearance: s.OverrideAppearance, Trusted: s.Trusted,
	}
}

func toLatestSkin786(s Skin786) protocol.Skin {
	animations := make([]protocol.SkinAnimation, len(s.Animations))
	for i, a := range s.Animations {
		animations[i] = protocol.SkinAnimation{
			ImageWidth: a.ImageWidth, ImageHeight: a.ImageHeight, ImageData: a.ImageData,
			AnimationType: a.AnimationType, FrameCount: a.FrameCount, ExpressionType: a.ExpressionType,
		}
	}
	pieces := make([]protocol.PersonaPiece, len(s.PersonaPieces))
	for i, p := range s.PersonaPieces {
		pieces[i] = protocol.PersonaPiece{
			PieceID: p.PieceID, PieceType: personaPieceTypeFromName(p.PieceType), PackID: uuidStringOrNil(p.PackID),
			Default: p.Default, ProductID: p.ProductID,
		}
	}
	tints := make([]protocol.PersonaPieceTintColour, len(s.PieceTintColours))
	for i, t := range s.PieceTintColours {
		var colours [4]color.RGBA
		for j := 0; j < 4 && j < len(t.Colours); j++ {
			colours[j] = hexStringARGB(t.Colours[j])
		}
		tints[i] = protocol.PersonaPieceTintColour{PieceType: t.PieceType, Colours: colours}
	}
	armSize := uint8(protocol.ArmSizeSlim)
	if s.ArmSize == "wide" {
		armSize = protocol.ArmSizeWide
	}
	return protocol.Skin{
		SkinID: s.SkinID, PlayFabID: s.PlayFabID, SkinResourcePatch: s.SkinResourcePatch,
		SkinImageWidth: s.SkinImageWidth, SkinImageHeight: s.SkinImageHeight, SkinData: s.SkinData,
		Animations: animations, CapeImageWidth: s.CapeImageWidth, CapeImageHeight: s.CapeImageHeight,
		CapeData: s.CapeData, SkinGeometry: s.SkinGeometry, GeometryDataEngineVersion: s.GeometryDataEngineVersion,
		AnimationData: s.AnimationData, CapeID: s.CapeID, FullID: s.FullID, ArmSize: armSize,
		SkinColour: hexStringRGB(s.SkinColour), PersonaPieces: pieces, PieceTintColours: tints,
		PremiumSkin: s.PremiumSkin, PersonaSkin: s.PersonaSkin, PersonaCapeOnClassicSkin: s.PersonaCapeOnClassicSkin,
		PrimaryUser: s.PrimaryUser, OverrideAppearance: s.OverrideAppearance, Trusted: s.Trusted,
	}
}

func rgbHexString(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
func hexStringRGB(s string) color.RGBA {
	s = strings.TrimPrefix(s, "#")
	var r, g, b uint8
	_, _ = fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 0xff}
}
func argbHexString(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.A, c.R, c.G, c.B)
}
func hexStringARGB(s string) color.RGBA {
	s = strings.TrimPrefix(s, "#")
	var a, r, g, b uint8
	if len(s) >= 8 {
		_, _ = fmt.Sscanf(s, "%02x%02x%02x%02x", &a, &r, &g, &b)
		return color.RGBA{R: r, G: g, B: b, A: a}
	}
	_, _ = fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 0xff}
}

func toLatestSimpleEvent(pk *SimpleEvent) *packet.SimpleEvent {
	return &packet.SimpleEvent{EventType: uint16(pk.EventType)}
}
func fromLatestSimpleEvent(pk *packet.SimpleEvent) *SimpleEvent {
	return &SimpleEvent{EventType: int16(pk.EventType)}
}

func toLatestShowStoreOffer(pk *ShowStoreOffer) *packet.ShowStoreOffer {
	id, _ := uuid.Parse(pk.OfferID)
	return &packet.ShowStoreOffer{OfferID: id, Type: pk.Type}
}
func fromLatestShowStoreOffer(pk *packet.ShowStoreOffer) *ShowStoreOffer {
	return &ShowStoreOffer{OfferID: pk.OfferID.String(), Type: pk.Type}
}

func toLatestUpdateClientInputLocks(pk *UpdateClientInputLocks) *packet.UpdateClientInputLocks {
	return &packet.UpdateClientInputLocks{Locks: pk.Locks}
}
func fromLatestUpdateClientInputLocks(pk *packet.UpdateClientInputLocks) *UpdateClientInputLocks {
	return &UpdateClientInputLocks{Locks: pk.Locks}
}

func toLatestUpdateClientOptions(pk *UpdateClientOptions) *packet.UpdateClientOptions {
	return &packet.UpdateClientOptions{GraphicsMode: pk.GraphicsMode}

}
func fromLatestUpdateClientOptions(pk *packet.UpdateClientOptions) *UpdateClientOptions {
	return &UpdateClientOptions{GraphicsMode: pk.GraphicsMode}
}

func toLatestCommandBlockUpdate(pk *CommandBlockUpdate) *packet.CommandBlockUpdate {
	return &packet.CommandBlockUpdate{
		Block: pk.Block, Position: pk.Position, Mode: pk.Mode, NeedsRedstone: pk.NeedsRedstone,
		Conditional: pk.Conditional, MinecartEntityRuntimeID: pk.MinecartEntityRuntimeID, Command: pk.Command,
		LastOutput: pk.LastOutput, Name: pk.Name, FilteredName: pk.FilteredName,
		ShouldTrackOutput: pk.ShouldTrackOutput, TickDelay: uint32(pk.TickDelay), ExecuteOnFirstTick: pk.ExecuteOnFirstTick,
	}
}
func fromLatestCommandBlockUpdate(pk *packet.CommandBlockUpdate) *CommandBlockUpdate {
	return &CommandBlockUpdate{
		Block: pk.Block, Position: pk.Position, Mode: pk.Mode, NeedsRedstone: pk.NeedsRedstone,
		Conditional: pk.Conditional, MinecartEntityRuntimeID: pk.MinecartEntityRuntimeID, Command: pk.Command,
		LastOutput: pk.LastOutput, Name: pk.Name, FilteredName: pk.FilteredName,
		ShouldTrackOutput: pk.ShouldTrackOutput, TickDelay: int32(pk.TickDelay), ExecuteOnFirstTick: pk.ExecuteOnFirstTick,
	}
}

func toLatestStructureBlockUpdate(pk *StructureBlockUpdate) *packet.StructureBlockUpdate {
	return &packet.StructureBlockUpdate{
		Position: pk.Position, StructureName: pk.StructureName, FilteredStructureName: pk.FilteredStructureName,
		DataField: pk.DataField, IncludePlayers: pk.IncludePlayers, ShowBoundingBox: pk.ShowBoundingBox,
		StructureBlockType: pk.StructureBlockType, Settings: pk.Settings, RedstoneSaveMode: uint8(pk.RedstoneSaveMode),
		ShouldTrigger: pk.ShouldTrigger, Waterlogged: pk.Waterlogged,
	}
}
func fromLatestStructureBlockUpdate(pk *packet.StructureBlockUpdate) *StructureBlockUpdate {
	return &StructureBlockUpdate{
		Position: pk.Position, StructureName: pk.StructureName, FilteredStructureName: pk.FilteredStructureName,
		DataField: pk.DataField, IncludePlayers: pk.IncludePlayers, ShowBoundingBox: pk.ShowBoundingBox,
		StructureBlockType: pk.StructureBlockType, Settings: pk.Settings, RedstoneSaveMode: int32(pk.RedstoneSaveMode),
		ShouldTrigger: pk.ShouldTrigger, Waterlogged: pk.Waterlogged,
	}
}

func toLatestBookEdit(pk *BookEdit) *packet.BookEdit {
	return &packet.BookEdit{
		ActionType: uint32(pk.ActionType), InventorySlot: int32(pk.InventorySlot), PageNumber: int32(pk.PageNumber),
		SecondaryPageNumber: int32(pk.SecondaryPageNumber), Text: pk.Text, PhotoName: pk.PhotoName,
		Title: pk.Title, Author: pk.Author, XUID: pk.XUID,
	}
}
func fromLatestBookEdit(pk *packet.BookEdit) *BookEdit {
	return &BookEdit{
		ActionType: byte(pk.ActionType), InventorySlot: byte(pk.InventorySlot), PageNumber: byte(pk.PageNumber),
		SecondaryPageNumber: byte(pk.SecondaryPageNumber), Text: pk.Text, PhotoName: pk.PhotoName,
		Title: pk.Title, Author: pk.Author, XUID: pk.XUID,
	}
}

func toLatestLessonProgress(pk *LessonProgress) *packet.LessonProgress {
	return &packet.LessonProgress{Identifier: pk.Identifier, Action: int32(pk.Action), Score: pk.Score}
}
func fromLatestLessonProgress(pk *packet.LessonProgress) *LessonProgress {
	return &LessonProgress{Identifier: pk.Identifier, Action: uint8(pk.Action), Score: pk.Score}
}

func toLatestAnvilDamage(pk *AnvilDamage) *packet.AnvilDamage {
	return &packet.AnvilDamage{AnvilPosition: pk.AnvilPosition}
}
func fromLatestAnvilDamage(pk *packet.AnvilDamage) *AnvilDamage {
	return &AnvilDamage{AnvilPosition: pk.AnvilPosition}
}
