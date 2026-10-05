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

package v786

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	ActorEventJump = iota + 1
	ActorEventHurt
	ActorEventDeath
	ActorEventStartAttacking
	ActorEventStopAttacking
	ActorEventTamingFailed
	ActorEventTamingSucceeded
	ActorEventShakeWetness
	ActorEventUseItem
	ActorEventEatGrass
	ActorEventFishhookBubble
	ActorEventFishhookFishPosition
	ActorEventFishhookHookTime
	ActorEventFishhookTease
	ActorEventSquidFleeing
	ActorEventZombieConverting
	ActorEventPlayAmbient
	ActorEventSpawnAlive
	ActorEventStartOfferFlower
	ActorEventStopOfferFlower
	ActorEventLoveHearts
	ActorEventVillagerAngry
	ActorEventVillagerHappy
	ActorEventWitchHatMagic
	ActorEventFireworksExplode
	ActorEventInLoveHearts
	ActorEventSilverfishMergeAnimation
	ActorEventGuardianAttackSound
	ActorEventDrinkPotion
	ActorEventThrowPotion
	ActorEventCartWithPrimeTNT
	ActorEventPrimeCreeper
	ActorEventAirSupply
	ActorEventAddPlayerLevels
	ActorEventGuardianMiningFatigue
	ActorEventAgentSwingArm
	ActorEventDragonStartDeathAnim
	ActorEventGroundDust
	ActorEventShake
)

const (
	ActorEventFeed = iota + 57
	_
	_
	ActorEventBabyEat
	ActorEventInstantDeath
	ActorEventNotifyTrade
	ActorEventLeashDestroyed
	ActorEventCaravanUpdated
	ActorEventTalismanActivate
	ActorEventUpdateStructureFeature
	ActorEventPlayerSpawnedMob
	ActorEventPuke
	ActorEventUpdateStackSize
	ActorEventStartSwimming
	ActorEventBalloonPop
	ActorEventTreasureHunt
	ActorEventSummonAgent
	ActorEventFinishedChargingItem
	ActorEventLandedOnGround
	ActorEventActorGrowUp
	ActorEventVibrationDetected
	ActorEventDrinkMilk
)

type ActorEvent struct {
	EntityRuntimeID uint64

	EventType byte

	EventData int32
}

func (*ActorEvent) ID() uint32 {
	return IDActorEvent
}

func (pk *ActorEvent) Marshal(io protocol.IO) {
	io.Varuint64(&pk.EntityRuntimeID)
	io.Uint8(&pk.EventType)
	io.Varint32(&pk.EventData)
}
