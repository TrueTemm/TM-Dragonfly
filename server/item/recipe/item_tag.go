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

package recipe

import (
	_ "embed"
	"encoding/json"
)

var (
	//go:embed item_tags.json
	itemTagData []byte
	itemTags    = make(map[string][]string)
)

func init() {
	if err := json.Unmarshal(itemTagData, &itemTags); err != nil {
		panic(err)
	}
}

type ItemTag struct {
	tag   string
	count int

	items []string
}

func NewItemTag(tag string, count int) ItemTag {
	if count < 0 {
		count = 0
	}
	return ItemTag{tag: tag, count: count, items: itemTags[tag]}
}

func (i ItemTag) Count() int {
	return i.count
}

func (i ItemTag) Empty() bool {
	return i.count == 0 || i.tag == ""
}

func (i ItemTag) Tag() string {
	return i.tag
}

func (i ItemTag) Contains(name string) bool {
	for _, item := range i.items {
		if item == name {
			return true
		}
	}
	return false
}
