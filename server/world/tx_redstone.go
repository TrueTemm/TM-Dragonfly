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

package world

import "github.com/df-mc/dragonfly/server/block/cube"

func (tx *Tx) Redstone() RedstoneTransaction {
	return RedstoneTransaction{tx: tx}
}

type RedstoneTransaction struct {
	tx *Tx
}

func (r RedstoneTransaction) ScheduleUpdate(pos cube.Pos) {
	r.tx.World().redstone.invalidateAround(pos, pos, RedstoneUpdateCauseScheduledTick, r.tx.Range())
}

func (r RedstoneTransaction) Torch(pos cube.Pos) RedstoneTorchTransaction {
	return RedstoneTorchTransaction{tx: r.tx, pos: pos}
}

type RedstoneTorchTransaction struct {
	tx  *Tx
	pos cube.Pos
}

func (t RedstoneTorchTransaction) BurnoutStatus() (burnedOut, recoverable bool) {
	return t.tx.World().redstone.torchBurnoutStatus(t.pos, t.tx.CurrentTick())
}

func (t RedstoneTorchTransaction) RecordTurnOff() (burnsOut bool) {
	return t.tx.World().redstone.recordTorchTurnOff(t.pos, t.tx.CurrentTick())
}

func (t RedstoneTorchTransaction) MarkSelfTriggered() {
	t.tx.World().redstone.markTorchSelfTriggered(t.pos)
}

func (t RedstoneTorchTransaction) ConsumeSelfTriggered() bool {
	return t.tx.World().redstone.consumeTorchSelfTriggered(t.pos)
}

func (t RedstoneTorchTransaction) ClearBurnout() {
	t.tx.World().redstone.clearTorchBurnout(t.pos)
}

func (tx *Tx) RedstonePower(pos cube.Pos) int {
	return tx.World().redstone.powerTo(pos, tx)
}

func (tx *Tx) RedstoneDirectPower(pos cube.Pos) int {
	return tx.World().redstone.directPower(pos, tx)
}

func (tx *Tx) RedstoneStrongPower(pos cube.Pos) int {
	return tx.World().redstone.strongPower(pos, tx)
}

func (tx *Tx) RedstoneConductivePower(pos cube.Pos) int {
	return tx.World().redstone.conductivePowerTo(pos, tx)
}

func (tx *Tx) RedstonePowerFrom(pos cube.Pos, face cube.Face) int {
	return tx.World().redstone.powerFrom(pos, tx, face)
}

func (tx *Tx) RedstoneDirectPowerFrom(pos cube.Pos, face cube.Face) int {
	return tx.World().redstone.directPowerFrom(pos, tx, face)
}

func (tx *Tx) RedstoneStrongPowerFrom(pos cube.Pos, face cube.Face) int {
	return tx.World().redstone.strongPowerFrom(pos, tx, face)
}
