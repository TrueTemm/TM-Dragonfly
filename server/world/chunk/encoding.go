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

package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type (
	Encoding interface {
		encodePalette(buf *bytes.Buffer, p *Palette, e paletteEncoding)
		decodePalette(buf *bytes.Buffer, blockSize paletteSize, e paletteEncoding) (*Palette, error)
		network() byte
	}

	paletteEncoding interface {
		encode(buf *bytes.Buffer, v uint32)
		decode(buf *bytes.Buffer) (uint32, error)
	}
)

var (
	DiskEncoding diskEncoding

	NetworkEncoding networkEncoding

	BiomePaletteEncoding biomePaletteEncoding
)

type biomePaletteEncoding struct{}

func (biomePaletteEncoding) encode(buf *bytes.Buffer, v uint32) {
	_ = binary.Write(buf, binary.LittleEndian, v)
}
func (biomePaletteEncoding) decode(buf *bytes.Buffer) (uint32, error) {
	var v uint32
	return v, binary.Read(buf, binary.LittleEndian, &v)
}

type BlockPaletteEncoding struct {
	Blocks BlockRegistry
}

func (bpe BlockPaletteEncoding) encode(buf *bytes.Buffer, v uint32) {
	_ = nbt.NewEncoderWithEncoding(buf, nbt.LittleEndian).Encode(bpe.EncodeBlockState(v))
}
func (bpe BlockPaletteEncoding) decode(buf *bytes.Buffer) (uint32, error) {
	var m map[string]any
	if err := nbt.NewDecoderWithEncoding(buf, nbt.LittleEndian).Decode(&m); err != nil {
		return 0, fmt.Errorf("error decoding block palette entry: %w", err)
	}
	return bpe.DecodeBlockState(m)
}

func (bpe BlockPaletteEncoding) EncodeBlockState(v uint32) blockEntry {

	name, props, _ := bpe.Blocks.RuntimeIDToState(v)
	return blockEntry{Name: name, State: props, Version: CurrentBlockVersion}
}

func (bpe BlockPaletteEncoding) DecodeBlockState(m map[string]any) (uint32, error) {

	name, _ := m["name"].(string)
	version, _ := m["version"].(int32)

	stateI, ok := m["states"]
	if version < 17694723 {

		meta, _ := m["val"].(int16)

		state, ok := upgradeLegacyEntry(name, meta)
		if !ok {
			return 0, fmt.Errorf("cannot find mapping for legacy block entry: %v, %v", name, meta)
		}

		name = state.Name
		stateI = state.State
		version = state.Version
	} else if !ok {

		stateI = make(map[string]any)
	}
	state, ok := stateI.(map[string]any)
	if !ok {
		return 0, fmt.Errorf("invalid state in block entry")
	}

	upgraded := blockupgrader.Upgrade(blockupgrader.BlockState{
		Name:       name,
		Properties: state,
		Version:    version,
	})

	v, ok := bpe.Blocks.StateToRuntimeID(upgraded.Name, upgraded.Properties)
	if !ok {
		return 0, fmt.Errorf("cannot get runtime ID of block state %v{%+v} %v", upgraded.Name, upgraded.Properties, upgraded.Version)
	}
	return v, nil
}

type diskEncoding struct{}

func (diskEncoding) network() byte { return 0 }
func (diskEncoding) encodePalette(buf *bytes.Buffer, p *Palette, e paletteEncoding) {
	if p.size != 0 {
		_ = binary.Write(buf, binary.LittleEndian, uint32(p.Len()))
	}
	for _, v := range p.values {
		e.encode(buf, v)
	}
}
func (diskEncoding) decodePalette(buf *bytes.Buffer, blockSize paletteSize, e paletteEncoding) (*Palette, error) {
	paletteCount := uint32(1)
	if blockSize != 0 {
		if err := binary.Read(buf, binary.LittleEndian, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
	}

	var err error
	palette := newPalette(blockSize, make([]uint32, paletteCount))
	for i := uint32(0); i < paletteCount; i++ {
		palette.values[i], err = e.decode(buf)
		if err != nil {
			return nil, err
		}
	}
	if paletteCount == 0 {
		return palette, fmt.Errorf("invalid palette entry count: found 0, but palette with %v bits per block must have at least 1 value", blockSize)
	}
	return palette, nil
}

type networkEncoding struct{}

func (networkEncoding) network() byte { return 1 }
func (networkEncoding) encodePalette(buf *bytes.Buffer, p *Palette, _ paletteEncoding) {
	if p.size != 0 {
		_ = protocol.WriteVarint32(buf, int32(p.Len()))
	}
	for _, val := range p.values {
		_ = protocol.WriteVarint32(buf, int32(val))
	}
}
func (networkEncoding) decodePalette(buf *bytes.Buffer, blockSize paletteSize, _ paletteEncoding) (*Palette, error) {
	var paletteCount int32 = 1
	if blockSize != 0 {
		if err := protocol.Varint32(buf, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
		if paletteCount <= 0 {
			return nil, fmt.Errorf("invalid palette entry count %v", paletteCount)
		}
	}

	blocks, temp := make([]uint32, paletteCount), int32(0)
	for i := int32(0); i < paletteCount; i++ {
		if err := protocol.Varint32(buf, &temp); err != nil {
			return nil, fmt.Errorf("error decoding palette entry: %w", err)
		}
		blocks[i] = uint32(temp)
	}
	return &Palette{values: blocks, size: blockSize}, nil
}
