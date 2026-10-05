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
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestPlayerListEntry786RoundTrips(t *testing.T) {
	in := &PlayerListEntry786{
		UUID:           uuid.New(),
		EntityUniqueID: 42,
		Username:       "Steve",
		XUID:           "1234567890",
		PlatformChatID: "",
		BuildPlatform:  7,
		Skin: Skin786{
			SkinID:            "geometry.humanoid.custom",
			PlayFabID:         "playfab",
			SkinResourcePatch: []byte(`{"geometry":{"default":"geometry.humanoid.custom"}}`),
			SkinImageWidth:    64, SkinImageHeight: 64,
			SkinData:   make([]byte, 64*64*4),
			CapeData:   nil,
			ArmSize:    "wide",
			SkinColour: "#b37b62",
		},
		Teacher:   false,
		Host:      true,
		SubClient: false,
	}

	buf := new(bytes.Buffer)
	w := protocol.NewWriter(buf, 0)
	protocol.Single[PlayerListEntry786](w, in)

	var latestBuf bytes.Buffer
	latestW := protocol.NewWriter(&latestBuf, 0)
	latestEntry := toLatestPlayerListEntryForTest(in)
	protocol.Single[protocol.PlayerListEntry](latestW, &latestEntry)
	if bytes.Equal(buf.Bytes(), latestBuf.Bytes()) {
		t.Fatalf("expected 786 encoding to differ from latest's, both were %d bytes", buf.Len())
	}

	r := protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, false)
	out := &PlayerListEntry786{}
	protocol.Single[PlayerListEntry786](r, out)

	if out.UUID != in.UUID || out.EntityUniqueID != in.EntityUniqueID || out.Username != in.Username ||
		out.XUID != in.XUID || out.BuildPlatform != in.BuildPlatform || out.Host != in.Host {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", out, in)
	}
	if out.Skin.ArmSize != in.Skin.ArmSize || out.Skin.SkinColour != in.Skin.SkinColour {
		t.Fatalf("skin round-trip mismatch: got %+v, want %+v", out.Skin, in.Skin)
	}
	if len(out.Skin.SkinData) != len(in.Skin.SkinData) {
		t.Fatalf("skin data length mismatch: got %d, want %d", len(out.Skin.SkinData), len(in.Skin.SkinData))
	}
}

func toLatestPlayerListEntryForTest(x *PlayerListEntry786) protocol.PlayerListEntry {
	return protocol.PlayerListEntry{
		ActionType: PlayerListActionAdd, UUID: x.UUID, EntityUniqueID: x.EntityUniqueID, Username: x.Username,
		XUID: x.XUID, PlatformChatID: x.PlatformChatID, BuildPlatform: x.BuildPlatform,
		Skin: toLatestSkin786(x.Skin), Teacher: x.Teacher, Host: x.Host, SubClient: x.SubClient,
	}
}

func TestSkinColourHexRoundTrips(t *testing.T) {
	c := hexStringRGB("#b37b62")
	if got := rgbHexString(c); got != "#b37b62" {
		t.Fatalf("rgb hex round-trip mismatch: got %s, want #b37b62", got)
	}
}

func TestPersonaPieceTypeRoundTrips(t *testing.T) {
	name := personaPieceTypeName(protocol.PieceTypeHair)
	if name != "persona_hair" {
		t.Fatalf("expected persona_hair, got %s", name)
	}
	if got := personaPieceTypeFromName(name); got != protocol.PieceTypeHair {
		t.Fatalf("expected PieceTypeHair, got %d", got)
	}
}
