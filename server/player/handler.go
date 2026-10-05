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

package player

import (
	"net"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type Handler interface {
	HandleMove(ctx *Context, newPos mgl64.Vec3, newRot cube.Rotation)

	HandleJump(p *Player)

	HandleTeleport(ctx *Context, pos mgl64.Vec3)

	HandleChangeWorld(p *Player, before, after *world.World)

	HandleToggleSprint(ctx *Context, after bool)

	HandleToggleSneak(ctx *Context, after bool)

	HandleChat(ctx *Context, message *string)

	HandleFoodLoss(ctx *Context, from int, to *int)

	HandleHeal(ctx *Context, health *float64, src world.HealingSource)

	HandleHurt(ctx *Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource)

	HandleSetOnFire(ctx *Context, duration *time.Duration)

	HandleDeath(p *Player, src world.DamageSource, keepInv *bool)

	HandleRespawn(p *Player, pos *mgl64.Vec3, w **world.World)

	HandleSkinChange(ctx *Context, skin *skin.Skin)

	HandleFireExtinguish(ctx *Context, pos cube.Pos)

	HandleStartBreak(ctx *Context, pos cube.Pos)

	HandleBlockBreak(ctx *Context, pos cube.Pos, drops *[]item.Stack, xp *int)

	HandleBlockPlace(ctx *Context, pos cube.Pos, b world.Block)

	HandleBlockPick(ctx *Context, pos cube.Pos, b world.Block)

	HandleItemUse(ctx *Context)

	HandleItemUseOnBlock(ctx *Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)

	HandleItemUseOnEntity(ctx *Context, e world.Entity)

	HandleItemRelease(ctx *Context, item item.Stack, dur time.Duration)

	HandleItemConsume(ctx *Context, item item.Stack)

	HandleAttackEntity(ctx *Context, e world.Entity, force, height *float64, critical *bool)

	HandleExperienceGain(ctx *Context, amount *int)

	HandlePunchAir(ctx *Context)

	HandleSwing(ctx *Context)

	HandleSignEdit(ctx *Context, pos cube.Pos, frontSide bool, oldText, newText string)

	HandleSleep(ctx *Context, sendReminder *bool)

	HandleLecternPageTurn(ctx *Context, pos cube.Pos, oldPage int, newPage *int)

	HandleItemDamage(ctx *Context, i item.Stack, damage *int)

	HandleItemPickup(ctx *Context, i *item.Stack)

	HandleHeldSlotChange(ctx *Context, from, to int)

	HandleItemDrop(ctx *Context, s item.Stack)

	HandleTransfer(ctx *Context, addr *net.UDPAddr)

	HandleCommandExecution(ctx *Context, command cmd.Command, args []string)

	HandleQuit(p *Player)

	HandleDiagnostics(p *Player, d session.Diagnostics)
}

type NopHandler struct{}

var _ Handler = NopHandler{}

func (NopHandler) HandleItemDrop(*Context, item.Stack)                                     {}
func (NopHandler) HandleHeldSlotChange(*Context, int, int)                                 {}
func (NopHandler) HandleMove(*Context, mgl64.Vec3, cube.Rotation)                          {}
func (NopHandler) HandleJump(*Player)                                                      {}
func (NopHandler) HandleTeleport(*Context, mgl64.Vec3)                                     {}
func (NopHandler) HandleChangeWorld(*Player, *world.World, *world.World)                   {}
func (NopHandler) HandleToggleSprint(*Context, bool)                                       {}
func (NopHandler) HandleToggleSneak(*Context, bool)                                        {}
func (NopHandler) HandleCommandExecution(*Context, cmd.Command, []string)                  {}
func (NopHandler) HandleTransfer(*Context, *net.UDPAddr)                                   {}
func (NopHandler) HandleChat(*Context, *string)                                            {}
func (NopHandler) HandleSkinChange(*Context, *skin.Skin)                                   {}
func (NopHandler) HandleFireExtinguish(*Context, cube.Pos)                                 {}
func (NopHandler) HandleStartBreak(*Context, cube.Pos)                                     {}
func (NopHandler) HandleBlockBreak(*Context, cube.Pos, *[]item.Stack, *int)                {}
func (NopHandler) HandleBlockPlace(*Context, cube.Pos, world.Block)                        {}
func (NopHandler) HandleBlockPick(*Context, cube.Pos, world.Block)                         {}
func (NopHandler) HandleSignEdit(*Context, cube.Pos, bool, string, string)                 {}
func (NopHandler) HandleSleep(*Context, *bool)                                             {}
func (NopHandler) HandleLecternPageTurn(*Context, cube.Pos, int, *int)                     {}
func (NopHandler) HandleItemPickup(*Context, *item.Stack)                                  {}
func (NopHandler) HandleItemUse(*Context)                                                  {}
func (NopHandler) HandleItemUseOnBlock(*Context, cube.Pos, cube.Face, mgl64.Vec3)          {}
func (NopHandler) HandleItemUseOnEntity(*Context, world.Entity)                            {}
func (NopHandler) HandleItemRelease(ctx *Context, item item.Stack, dur time.Duration)      {}
func (NopHandler) HandleItemConsume(*Context, item.Stack)                                  {}
func (NopHandler) HandleItemDamage(*Context, item.Stack, *int)                             {}
func (NopHandler) HandleAttackEntity(*Context, world.Entity, *float64, *float64, *bool)    {}
func (NopHandler) HandleExperienceGain(*Context, *int)                                     {}
func (NopHandler) HandleSwing(*Context)                                                    {}
func (NopHandler) HandlePunchAir(*Context)                                                 {}
func (NopHandler) HandleHurt(*Context, *float64, bool, *time.Duration, world.DamageSource) {}
func (NopHandler) HandleSetOnFire(*Context, *time.Duration)                                {}
func (NopHandler) HandleHeal(*Context, *float64, world.HealingSource)                      {}
func (NopHandler) HandleFoodLoss(*Context, int, *int)                                      {}
func (NopHandler) HandleDeath(*Player, world.DamageSource, *bool)                          {}
func (NopHandler) HandleRespawn(*Player, *mgl64.Vec3, **world.World)                       {}
func (NopHandler) HandleQuit(*Player)                                                      {}
func (NopHandler) HandleDiagnostics(*Player, session.Diagnostics)                          {}
