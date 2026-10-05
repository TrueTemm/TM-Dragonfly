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
	_ "embed"
	"os"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"golang.org/x/mod/sumdb/dirhash"
)

//go:embed pack_icon.png
var packIcon []byte

func BuildResourcePack(reg world.BlockRegistry) (*resource.Pack, bool) {
	dir, err := os.MkdirTemp("", "tm-dragonfly_resource_pack-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	var assets int
	var lang []string

	itemCount, itemLang := buildItems(dir)
	assets += itemCount
	lang = append(lang, itemLang...)

	blockCount, blockLang := buildBlocks(reg, dir)
	assets += blockCount
	lang = append(lang, blockLang...)

	if assets > 0 {
		buildLanguageFile(dir, lang)
		if err := os.WriteFile(dir+"/pack_icon.png", packIcon, 0666); err != nil {
			panic(err)
		}
		hash, err := dirhash.HashDir(dir, "", dirhash.Hash1)
		if err != nil {
			panic(err)
		}
		var header, module [16]byte
		copy(header[:], hash)
		copy(module[:], hash[16:])
		buildManifest(dir, header, module)
		return resource.MustReadPath(dir), true
	}
	return nil, false
}
