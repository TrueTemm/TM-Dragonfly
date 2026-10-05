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

package leveldat

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

type LevelDat struct {
	hdr  header
	data []byte
}

type header struct {
	StorageVersion int32
	FileLength     int32
}

func ReadFile(name string) (*LevelDat, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("level.dat: open file: %w", err)
	}
	defer f.Close()
	return Read(bufio.NewReader(f))
}

func Read(r io.Reader) (*LevelDat, error) {
	var ldat LevelDat
	if err := binary.Read(r, binary.LittleEndian, &ldat.hdr); err != nil {
		return nil, fmt.Errorf("level.dat: read header: %w", err)
	}
	ldat.data = make([]byte, ldat.hdr.FileLength)
	if n, err := io.ReadFull(r, ldat.data); err != nil || int32(n) != ldat.hdr.FileLength {
		return nil, fmt.Errorf("level.dat: read data: %w", err)
	}
	return &ldat, nil
}

func (ld *LevelDat) Unmarshal(dst any) error {
	err := nbt.UnmarshalEncoding(ld.data, dst, nbt.LittleEndian)
	if err == nil {
		return nil
	}
	var unexpected nbt.UnexpectedNamedTagError
	if !errors.As(err, &unexpected) {
		return fmt.Errorf("level.dat: decode nbt: %w", err)
	}

	if err = unmarshalLenient(ld.data, dst); err != nil {
		return fmt.Errorf("level.dat: decode nbt: %w", err)
	}
	return nil
}

func (ld *LevelDat) Ver() int {
	return int(ld.hdr.StorageVersion)
}

func (ld *LevelDat) Marshal(src any) error {
	var err error
	ld.data, err = nbt.MarshalEncoding(src, nbt.LittleEndian)
	if err != nil {
		return fmt.Errorf("level.dat: encode nbt: %w", err)
	}
	ld.hdr = header{
		StorageVersion: Version,
		FileLength:     int32(len(ld.data)),
	}
	return nil
}

func (ld *LevelDat) Write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, ld.hdr); err != nil {
		return fmt.Errorf("level.dat: write header: %w", err)
	}
	if _, err := w.Write(ld.data); err != nil {
		return fmt.Errorf("level.dat: write data: %w", err)
	}
	return nil
}

func (ld *LevelDat) WriteFile(name string) error {
	f, err := os.OpenFile(name, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("level.dat: open file: %w", err)
	}
	w := bufio.NewWriter(f)
	defer func() {
		_ = w.Flush()
		_ = f.Close()
	}()
	return ld.Write(w)
}
