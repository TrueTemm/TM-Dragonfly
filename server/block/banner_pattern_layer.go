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

package block

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item"
)

type BannerPatternLayer struct {
	Type BannerPatternType

	Colour item.Colour
}

func (b BannerPatternLayer) EncodeNBT() map[string]any {
	return map[string]any{
		"Pattern": bannerPatternID(b.Type),
		"Color":   int32(invertColour(b.Colour)),
	}
}

func (b BannerPatternLayer) DecodeNBT(data map[string]any) any {
	id := nbtconv.String(data, "Pattern")
	pattern, exists := BannerPatternByID(id)
	if !exists {
		panic(fmt.Errorf("unknown banner pattern id %q", id))
	}
	b.Type = pattern
	b.Colour = invertColourID(int16(nbtconv.Int32(data, "Color")))
	return b
}
