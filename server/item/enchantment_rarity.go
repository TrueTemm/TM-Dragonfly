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

package item

type EnchantmentRarity interface {
	Name() string

	Cost() int

	Weight() int
}

var (
	EnchantmentRarityCommon enchantmentRarityCommon

	EnchantmentRarityUncommon enchantmentRarityUncommon

	EnchantmentRarityRare enchantmentRarityRare

	EnchantmentRarityVeryRare enchantmentRarityVeryRare
)

type enchantmentRarityCommon struct{}

func (enchantmentRarityCommon) Name() string { return "Common" }
func (enchantmentRarityCommon) Cost() int    { return 1 }
func (enchantmentRarityCommon) Weight() int  { return 10 }

type enchantmentRarityUncommon struct{}

func (enchantmentRarityUncommon) Name() string { return "Uncommon" }
func (enchantmentRarityUncommon) Cost() int    { return 2 }
func (enchantmentRarityUncommon) Weight() int  { return 5 }

type enchantmentRarityRare struct{}

func (enchantmentRarityRare) Name() string { return "Rare" }
func (enchantmentRarityRare) Cost() int    { return 4 }
func (enchantmentRarityRare) Weight() int  { return 2 }

type enchantmentRarityVeryRare struct{}

func (enchantmentRarityVeryRare) Name() string { return "Very Rare" }
func (enchantmentRarityVeryRare) Cost() int    { return 8 }
func (enchantmentRarityVeryRare) Weight() int  { return 1 }
