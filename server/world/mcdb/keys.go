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

package mcdb

const (
	keySubChunkData = '/'
)

const (
	keyVersion = ','

	keyVersionOld = 'v'

	keyBlockEntities = '1'

	keyEntitiesOld = '2'

	keyPendingScheduledTicks = '3'

	keyFinalisation = '6'

	key3DData = '+'

	key2DData = '-'

	keyChecksums = ';'

	keyEntityIdentifiers = "digp"

	keyEntity = "actorprefix"
)

const (
	keyAutonomousEntities = "AutonomousEntities"
	keyOverworld          = "Overworld"
	keyMobEvents          = "mobevents"
	keyBiomeData          = "BiomeData"
	keyScoreboard         = "scoreboard"
	keyLocalPlayer        = "~local_player"
)

const (
	finalisationGenerated = iota + 1
	finalisationPopulated
)
