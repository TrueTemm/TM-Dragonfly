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

package session

import (
	"fmt"
	"math"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerAuthInputHandler struct{}

func (h PlayerAuthInputHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.PlayerAuthInput)
	if err := h.handleMovement(pk, s, c); err != nil {
		return err
	}
	return h.handleActions(pk, s, tx, c)
}

func (h PlayerAuthInputHandler) handleMovement(pk *packet.PlayerAuthInput, s *Session, c Controllable) error {
	yaw, pitch := c.Rotation().Elem()
	pos := c.Position()

	reference := []float64{pitch, yaw, yaw, pos[0], pos[1], pos[2]}
	for i, v := range [...]*float32{&pk.Pitch, &pk.Yaw, &pk.HeadYaw, &pk.Position[0], &pk.Position[1], &pk.Position[2]} {
		f := float64(*v)
		if math.IsNaN(f) || math.IsInf(f, 1) || math.IsInf(f, 0) {

			s.conf.Log.Debug("process packet: PlayerAuthInput: found nan/inf values. assuming server-side values", "pos", fmt.Sprint(pk.Position), "yaw", pk.Yaw, "head-yaw", pk.HeadYaw, "pitch", pk.Pitch)
			*v = float32(reference[i])
		}
	}

	s.clientTick.Store(pk.Tick)
	pk.Position = pk.Position.Sub(mgl32.Vec3{0, 1.62})

	newPos := vec32To64(pk.Position)
	deltaPos, deltaYaw, deltaPitch := newPos.Sub(pos), float64(pk.Yaw)-yaw, float64(pk.Pitch)-pitch

	if !mgl64.FloatEqual(deltaPos.Len(), 0) || !mgl64.FloatEqual(deltaYaw, 0) || !mgl64.FloatEqual(deltaPitch, 0) {
		if expected := s.teleportPos.Load(); expected != nil {
			if newPos.Sub(*expected).Len() > 1 {

				s.moving = true
				c.Move(mgl64.Vec3{}, deltaYaw, deltaPitch)
				return nil
			}
			s.teleportPos.Store(nil)
		}
	}

	s.moving = true
	c.Move(deltaPos, deltaYaw, deltaPitch)
	return nil
}

func (h PlayerAuthInputHandler) handleActions(pk *packet.PlayerAuthInput, s *Session, tx *world.Tx, c Controllable) error {
	if pk.InputData.Load(packet.InputFlagPerformItemInteraction) {
		data, ok := pk.ItemInteractionData.Value()
		if !ok {
			return fmt.Errorf("item interaction flag set without item interaction data")
		}
		if err := h.handleUseItemData(data, s, c); err != nil {
			return err
		}
	}
	if pk.InputData.Load(packet.InputFlagPerformBlockActions) {
		actions, ok := pk.BlockActions.Value()
		if !ok {
			return fmt.Errorf("block actions flag set without block actions")
		}
		if err := h.handleBlockActions(actions, s, c); err != nil {
			return err
		}
	}
	h.handleInputFlags(pk.InputData, s, c)

	if pk.InputData.Load(packet.InputFlagPerformItemStackRequest) {
		request, ok := pk.ItemStackRequest.Value()
		if !ok {
			return fmt.Errorf("item stack request flag set without item stack request")
		}
		s.inTransaction.Store(true)
		defer s.inTransaction.Store(false)

		sh := s.handlers[packet.IDItemStackRequest].(*ItemStackRequestHandler)
		if err := sh.handleRequest(request, s, tx, c); err != nil {

			s.conf.Log.Debug("process packet: PlayerAuthInput: resolve item stack request: " + err.Error())
		}
	}
	return nil
}

func (h PlayerAuthInputHandler) handleInputFlags(flags protocol.InputFlags, s *Session, c Controllable) {
	if flags.Load(packet.InputFlagStartSprinting) {
		c.StartSprinting()
	}
	if flags.Load(packet.InputFlagStopSprinting) {
		c.StopSprinting()
	}
	if sneaking := flags.Load(packet.InputFlagSneaking); sneaking != c.Sneaking() {
		if sneaking {
			c.StartSneaking()
		} else {
			c.StopSneaking()
		}
	}
	if flags.Load(packet.InputFlagStartSwimming) {
		c.StartSwimming()
	}
	if flags.Load(packet.InputFlagStopSwimming) {
		c.StopSwimming()
	}
	if flags.Load(packet.InputFlagStartGliding) {
		c.StartGliding()
	}
	if flags.Load(packet.InputFlagStopGliding) {
		c.StopGliding()
	}
	if flags.Load(packet.InputFlagStartJumping) {
		c.Jump()
	}

	jumpDown := flags.Load(packet.InputFlagJumpCurrentRaw)
	switch {
	case flags.Load(packet.InputFlagJumpPressedRaw):
		s.rawJump = true
		c.PressJump()
	case !s.rawJump && jumpDown && !s.jumpDown:
		c.PressJump()
	}
	s.jumpDown = jumpDown
	if flags.Load(packet.InputFlagStartCrawling) {
		c.StartCrawling()
	}
	if flags.Load(packet.InputFlagStopCrawling) {
		c.StopCrawling()
	}
	if flags.Load(packet.InputFlagMissedSwing) {
		s.swingingArm.Store(true)
		defer s.swingingArm.Store(false)
		c.PunchAir()
	}
	if flags.Load(packet.InputFlagStartFlying) {
		if !c.GameMode().AllowsFlying() {
			s.conf.Log.Debug("process packet: PlayerAuthInput: flying flag enabled while unable to fly")
			s.SendAbilities(c)
		} else {
			c.StartFlying()
		}
	}
	if flags.Load(packet.InputFlagStopFlying) && c.GameMode().HasCollision() {

		c.StopFlying()
	}
}

func (h PlayerAuthInputHandler) handleUseItemData(data protocol.UseItemTransactionData, s *Session, c Controllable) error {
	s.swingingArm.Store(true)
	defer s.swingingArm.Store(false)

	held, _ := c.HeldItems()
	if !held.Equal(stackToItem(s.br, data.HeldItem.Stack)) {
		s.conf.Log.Debug("process packet: PlayerAuthInput: UseItemTransaction: mismatch between actual held item and client held item")
		return nil
	}
	pos := cube.Pos{int(data.BlockPosition[0]), int(data.BlockPosition[1]), int(data.BlockPosition[2])}

	switch data.ActionType {
	case protocol.UseItemActionBreakBlock:
		c.BreakBlock(pos)
	default:
		return fmt.Errorf("unhandled UseItem ActionType for PlayerAuthInput packet %v", data.ActionType)
	}
	return nil
}

func (h PlayerAuthInputHandler) handleBlockActions(a []protocol.PlayerBlockAction, s *Session, c Controllable) error {
	for _, action := range a {
		if err := handlePlayerAction(action.Action, action.Face, action.BlockPos, selfEntityRuntimeID, s, c); err != nil {
			return err
		}
	}
	return nil
}
