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
	"github.com/df-mc/dragonfly/server/world"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	_ "unsafe"
)

func buildItems(dir string) (count int, lang []string) {
	if err := os.Mkdir(filepath.Join(dir, "items"), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "textures/items"), os.ModePerm); err != nil {
		panic(err)
	}

	textureData := make(map[string]any)
	for _, item := range world.CustomItems() {
		identifier, _ := item.EncodeItem()
		lang = append(lang, fmt.Sprintf("item.%s.name=%s", identifier, item.Name()))

		name := strings.Split(identifier, ":")[1]
		textureData[identifier] = map[string]string{"textures": fmt.Sprintf("textures/items/%s.png", name)}

		buildItemTexture(dir, name, item.Texture())

		count++
	}

	buildItemAtlas(dir, map[string]any{
		"resource_pack_name": "vanilla",
		"texture_name":       "atlas.items",
		"texture_data":       textureData,
	})
	return
}

func buildItemTexture(dir, name string, img image.Image) {
	texture, err := os.Create(filepath.Join(dir, "textures/items", name+".png"))
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

func buildItemAtlas(dir string, atlas map[string]any) {
	b, err := json.Marshal(atlas)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "textures/item_texture.json"), b, 0666); err != nil {
		panic(err)
	}
}
