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

package packbuilder

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func buildManifest(dir string, headerUUID, moduleUUID uuid.UUID) {
	m, err := json.Marshal(resource.Manifest{
		FormatVersion: 2,
		Header: resource.Header{
			Name:               "TM-Dragonfly",
			Description:        "TM-Dragonfly resources",
			UUID:               headerUUID,
			Version:            [3]int{0, 0, 1},
			MinimumGameVersion: parseVersion(protocol.CurrentVersion),
		},
		Modules: []resource.Module{
			{
				UUID:        moduleUUID.String(),
				Description: "TM-Dragonfly resources",
				Type:        "resources",
				Version:     [3]int{0, 0, 1},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), m, 0666); err != nil {
		panic(err)
	}
}

func parseVersion(ver string) [3]int {
	frag := strings.Split(ver, ".")
	if len(frag) != 3 {
		panic("invalid version number " + ver)
	}
	a, _ := strconv.ParseInt(frag[0], 10, 64)
	b, _ := strconv.ParseInt(frag[1], 10, 64)
	c, _ := strconv.ParseInt(frag[2], 10, 64)
	return [3]int{int(a), int(b), int(c)}
}
