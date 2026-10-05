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

package biome

import (
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/world"
)

func init() {
	world.RegisterBiome(BadlandsPlateau{})
	world.RegisterBiome(Badlands{})
	world.RegisterBiome(BambooJungleHills{})
	world.RegisterBiome(BambooJungle{})
	world.RegisterBiome(BasaltDeltas{})
	world.RegisterBiome(Beach{})
	world.RegisterBiome(BirchForestHills{})
	world.RegisterBiome(BirchForest{})
	world.RegisterBiome(CherryGrove{})
	world.RegisterBiome(ColdOcean{})
	world.RegisterBiome(CrimsonForest{})
	world.RegisterBiome(DarkForestHills{})
	world.RegisterBiome(DarkForest{})
	world.RegisterBiome(DeepColdOcean{})
	world.RegisterBiome(DeepDark{})
	world.RegisterBiome(DeepFrozenOcean{})
	world.RegisterBiome(DeepLukewarmOcean{})
	world.RegisterBiome(DeepOcean{})
	world.RegisterBiome(DeepWarmOcean{})
	world.RegisterBiome(DesertHills{})
	world.RegisterBiome(DesertLakes{})
	world.RegisterBiome(Desert{})
	world.RegisterBiome(DripstoneCaves{})
	world.RegisterBiome(End{})
	world.RegisterBiome(ErodedBadlands{})
	world.RegisterBiome(FlowerForest{})
	world.RegisterBiome(Forest{})
	world.RegisterBiome(FrozenOcean{})
	world.RegisterBiome(FrozenPeaks{})
	world.RegisterBiome(FrozenRiver{})
	world.RegisterBiome(GiantSpruceTaigaHills{})
	world.RegisterBiome(GiantTreeTaigaHills{})
	world.RegisterBiome(GravellyMountainsPlus{})
	world.RegisterBiome(Grove{})
	world.RegisterBiome(IceSpikes{})
	world.RegisterBiome(JaggedPeaks{})
	world.RegisterBiome(JungleEdge{})
	world.RegisterBiome(JungleHills{})
	world.RegisterBiome(Jungle{})
	world.RegisterBiome(LegacyFrozenOcean{})
	world.RegisterBiome(LukewarmOcean{})
	world.RegisterBiome(LushCaves{})
	world.RegisterBiome(MangroveSwamp{})
	world.RegisterBiome(Meadow{})
	world.RegisterBiome(ModifiedBadlandsPlateau{})
	world.RegisterBiome(ModifiedJungleEdge{})
	world.RegisterBiome(ModifiedJungle{})
	world.RegisterBiome(ModifiedWoodedBadlandsPlateau{})
	world.RegisterBiome(MountainEdge{})
	world.RegisterBiome(MushroomFieldShore{})
	world.RegisterBiome(MushroomFields{})
	world.RegisterBiome(NetherWastes{})
	world.RegisterBiome(Ocean{})
	world.RegisterBiome(OldGrowthBirchForest{})
	world.RegisterBiome(OldGrowthPineTaiga{})
	world.RegisterBiome(OldGrowthSpruceTaiga{})
	world.RegisterBiome(PaleGarden{})
	world.RegisterBiome(Plains{})
	world.RegisterBiome(River{})
	world.RegisterBiome(SavannaPlateau{})
	world.RegisterBiome(Savanna{})
	world.RegisterBiome(ShatteredSavannaPlateau{})
	world.RegisterBiome(SnowyBeach{})
	world.RegisterBiome(SnowyMountains{})
	world.RegisterBiome(SnowyPlains{})
	world.RegisterBiome(SnowySlopes{})
	world.RegisterBiome(SnowyTaigaHills{})
	world.RegisterBiome(SnowyTaigaMountains{})
	world.RegisterBiome(SnowyTaiga{})
	world.RegisterBiome(SoulSandValley{})
	world.RegisterBiome(StonyPeaks{})
	world.RegisterBiome(StonyShore{})
	world.RegisterBiome(SulfurCaves{})
	world.RegisterBiome(SunflowerPlains{})
	world.RegisterBiome(SwampHills{})
	world.RegisterBiome(Swamp{})
	world.RegisterBiome(TaigaHills{})
	world.RegisterBiome(TaigaMountains{})
	world.RegisterBiome(Taiga{})
	world.RegisterBiome(TallBirchHills{})
	world.RegisterBiome(WarmOcean{})
	world.RegisterBiome(WarpedForest{})
	world.RegisterBiome(WindsweptForest{})
	world.RegisterBiome(WindsweptGravellyHills{})
	world.RegisterBiome(WindsweptHills{})
	world.RegisterBiome(WindsweptSavanna{})
	world.RegisterBiome(WoodedBadlandsPlateau{})
	world.RegisterBiome(WoodedHills{})

	world_finaliseBiomeRegistry()
}

//go:linkname world_finaliseBiomeRegistry github.com/df-mc/dragonfly/server/world.finaliseBiomeRegistry
func world_finaliseBiomeRegistry()
