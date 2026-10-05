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

package v800

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func FromLatestBiomeDefinitionList(pk *packet.BiomeDefinitionList) *BiomeDefinitionList {
	return fromLatestBiomeDefinitionList800(pk)
}

func FromLatestClientMovementPredictionSync(pk *packet.ClientMovementPredictionSync) *ClientMovementPredictionSync {
	return fromLatestClientMovementPredictionSync(pk)
}
func ToLatestClientMovementPredictionSync(pk *ClientMovementPredictionSync) *packet.ClientMovementPredictionSync {
	return toLatestClientMovementPredictionSync(pk)
}

func FromLatestPlayerLocation(pk *packet.PlayerLocation) *PlayerLocation {
	return fromLatestPlayerLocation(pk)
}
func ToLatestPlayerLocation(pk *PlayerLocation) *packet.PlayerLocation {
	return toLatestPlayerLocation(pk)
}
