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
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

func TestUnmarshalForeignProperties(t *testing.T) {
	raw, err := nbt.MarshalEncoding(map[string]any{
		"LevelName":     "Lobby",
		"SpawnX":        int32(-7),
		"SpawnY":        int32(-37),
		"SpawnZ":        int32(-1),
		"GameType":      int32(0),
		"generatorName": "flat",
		"hardcore":      byte(0),
		"abilities": map[string]any{
			"flying":       byte(0),
			"worldbuilder": byte(1),
		},
		"experiments": map[string]any{
			"experiments_ever_used": byte(1),
		},
	}, nbt.LittleEndian)
	if err != nil {
		t.Fatal(err)
	}
	ld := &LevelDat{data: raw}

	var d Data
	if err = ld.Unmarshal(&d); err != nil {
		t.Fatalf("a level.dat with foreign properties must still be read: %v", err)
	}
	if d.LevelName != "Lobby" {
		t.Errorf("LevelName = %q, want Lobby", d.LevelName)
	}
	if d.SpawnX != -7 || d.SpawnY != -37 || d.SpawnZ != -1 {
		t.Errorf("spawn = %d %d %d, want -7 -37 -1", d.SpawnX, d.SpawnY, d.SpawnZ)
	}
	if len(d.Experiments) != 1 {
		t.Errorf("experiments = %v, want the compound to be kept as it is", d.Experiments)
	}
}
