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
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Skin786 struct {
	SkinID                          string
	PlayFabID                       string
	SkinResourcePatch               []byte
	SkinImageWidth, SkinImageHeight uint32
	SkinData                        []byte
	Animations                      []SkinAnimation786
	CapeImageWidth, CapeImageHeight uint32
	CapeData                        []byte
	SkinGeometry                    []byte
	GeometryDataEngineVersion       []byte
	AnimationData                   []byte
	CapeID                          string
	FullID                          string

	ArmSize string

	SkinColour               string
	PersonaPieces            []PersonaPiece786
	PieceTintColours         []PersonaPieceTintColour786
	PremiumSkin              bool
	PersonaSkin              bool
	PersonaCapeOnClassicSkin bool
	PrimaryUser              bool
	OverrideAppearance       bool

	Trusted bool
}

func (x *Skin786) Marshal(io protocol.IO) {
	io.String(&x.SkinID)
	io.String(&x.PlayFabID)
	io.ByteSlice(&x.SkinResourcePatch)
	io.Uint32(&x.SkinImageWidth)
	io.Uint32(&x.SkinImageHeight)
	io.ByteSlice(&x.SkinData)
	protocol.SliceUint32Length(io, &x.Animations)
	io.Uint32(&x.CapeImageWidth)
	io.Uint32(&x.CapeImageHeight)
	io.ByteSlice(&x.CapeData)
	io.ByteSlice(&x.SkinGeometry)
	io.ByteSlice(&x.GeometryDataEngineVersion)
	io.ByteSlice(&x.AnimationData)
	io.String(&x.CapeID)
	io.String(&x.FullID)
	io.String(&x.ArmSize)
	io.String(&x.SkinColour)
	protocol.SliceUint32Length(io, &x.PersonaPieces)
	protocol.SliceUint32Length(io, &x.PieceTintColours)
	io.Bool(&x.PremiumSkin)
	io.Bool(&x.PersonaSkin)
	io.Bool(&x.PersonaCapeOnClassicSkin)
	io.Bool(&x.PrimaryUser)
	io.Bool(&x.OverrideAppearance)
}

type SkinAnimation786 struct {
	ImageWidth, ImageHeight uint32
	ImageData               []byte
	AnimationType           uint32
	FrameCount              float32
	ExpressionType          uint32
}

func (x *SkinAnimation786) Marshal(io protocol.IO) {
	io.Uint32(&x.ImageWidth)
	io.Uint32(&x.ImageHeight)
	io.ByteSlice(&x.ImageData)
	io.Uint32(&x.AnimationType)
	io.Float32(&x.FrameCount)
	io.Uint32(&x.ExpressionType)
}

type PersonaPiece786 struct {
	PieceID   string
	PieceType string
	PackID    string
	Default   bool
	ProductID string
}

func (x *PersonaPiece786) Marshal(io protocol.IO) {
	io.String(&x.PieceID)
	io.String(&x.PieceType)
	io.String(&x.PackID)
	io.Bool(&x.Default)
	io.String(&x.ProductID)
}

type PersonaPieceTintColour786 struct {
	PieceType string
	Colours   []string
}

func (x *PersonaPieceTintColour786) Marshal(io protocol.IO) {
	io.String(&x.PieceType)
	protocol.FuncSliceUint32Length(io, &x.Colours, ioString(io))
}

func ioString(io protocol.IO) func(*string) { return io.String }

var personaPieceTypeNames = map[uint32]string{
	protocol.PieceTypeSkeleton:      "persona_skeleton",
	protocol.PieceTypeBody:          "persona_body",
	protocol.PieceTypeSkin:          "persona_skin",
	protocol.PieceTypeBottom:        "persona_bottom",
	protocol.PieceTypeFeet:          "persona_feet",
	protocol.PieceTypeDress:         "persona_dress",
	protocol.PieceTypeTop:           "persona_top",
	protocol.PieceTypeHighPants:     "persona_high_pants",
	protocol.PieceTypeHands:         "persona_hand",
	protocol.PieceTypeOuterwear:     "persona_outerwear",
	protocol.PieceTypeFacialHair:    "persona_facial_hair",
	protocol.PieceTypeMouth:         "persona_mouth",
	protocol.PieceTypeEyes:          "persona_eyes",
	protocol.PieceTypeHair:          "persona_hair",
	protocol.PieceTypeHood:          "persona_hood",
	protocol.PieceTypeBack:          "persona_back",
	protocol.PieceTypeFaceAccessory: "persona_face_accessory",
	protocol.PieceTypeHead:          "persona_head",
	protocol.PieceTypeLegs:          "persona_legs",
	protocol.PieceTypeLeftLeg:       "persona_left_leg",
	protocol.PieceTypeRightLeg:      "persona_right_leg",
	protocol.PieceTypeArms:          "persona_arms",
	protocol.PieceTypeLeftArm:       "persona_left_arm",
	protocol.PieceTypeRightArm:      "persona_right_arm",
	protocol.PieceTypeCapes:         "persona_cape",
	protocol.PieceTypeClassicSkin:   "persona_classic_skin",
	protocol.PieceTypeEmote:         "persona_emote",
	protocol.PieceTypeUnsupported:   "persona_unsupported",
}

func personaPieceTypeName(t uint32) string {
	if s, ok := personaPieceTypeNames[t]; ok {
		return s
	}
	return "persona_unsupported"
}

func personaPieceTypeFromName(s string) uint32 {
	for t, name := range personaPieceTypeNames {
		if name == s {
			return t
		}
	}
	return protocol.PieceTypeUnsupported
}

func uuidStringOrNil(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}
	}
	return id
}
