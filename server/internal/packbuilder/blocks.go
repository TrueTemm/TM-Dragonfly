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
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/world"
)

func buildBlocks(reg world.BlockRegistry, dir string) (count int, lang []string) {
	if err := os.MkdirAll(filepath.Join(dir, "models/blocks"), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "textures/blocks"), os.ModePerm); err != nil {
		panic(err)
	}

	textureData := make(map[string]any)
	for identifier, blk := range reg.CustomBlocks() {
		b, ok := blk.(world.CustomBlockBuildable)
		if !ok {
			continue
		}

		name := strings.Split(identifier, ":")[1]
		lang = append(lang, fmt.Sprintf("tile.%s.name=%s", identifier, b.Name()))
		for name, texture := range b.Textures() {
			textureData[name] = map[string]string{"textures": "textures/blocks/" + name}
			buildBlockTexture(dir, name, texture)
		}
		if b.Geometry() != nil {
			if err := os.WriteFile(filepath.Join(dir, "models/blocks", fmt.Sprintf("%s.geo.json", name)), b.Geometry(), 0666); err != nil {
				panic(err)
			}
		}
		count++
	}

	buildBlockAtlas(dir, map[string]any{
		"resource_pack_name": "vanilla",
		"texture_name":       "atlas.terrain",
		"padding":            8,
		"num_mip_levels":     4,
		"texture_data":       textureData,
	})
	return
}

func buildBlockTexture(dir, name string, img image.Image) {
	texture, err := os.Create(filepath.Join(dir, fmt.Sprintf("textures/blocks/%s.png", name)))
	if err != nil {
		panic(err)
	}
	if err := png.Encode(texture, img); err != nil {
		_ = texture.Close()
		panic(err)
	}
	if err := texture.Close(); err != nil {
		panic(err)
	}
}

func buildBlockAtlas(dir string, atlas map[string]any) {
	b, err := json.Marshal(atlas)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "textures/terrain_texture.json"), b, 0666); err != nil {
		panic(err)
	}
}
