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
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/debug"
	"github.com/df-mc/dragonfly/server/player/dialogue"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/player/hud"
	"github.com/df-mc/dragonfly/server/player/input"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type Controllable interface {
	Name() string
	world.Entity
	item.User
	dialogue.Submitter
	form.Submitter
	cmd.Source
	chat.Subscriber
	hud.Renderer
	debug.Renderer
	input.Restricter

	Locale() language.Tag

	SetHeldItems(right, left item.Stack)
	SetHeldSlot(slot int) error

	Move(deltaPos mgl64.Vec3, deltaYaw, deltaPitch float64)

	Speed() float64
	FlightSpeed() float64
	VerticalFlightSpeed() float64

	Sleep(pos cube.Pos)
	Wake()

	Chat(msg ...any)
	ExecuteCommand(commandLine string)
	GameMode() world.GameMode
	SetGameMode(mode world.GameMode)
	Effects() []effect.Effect

	UseItem()
	ReleaseItem()
	UseItemOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
	UseItemOnEntity(e world.Entity) bool
	BreakBlock(pos cube.Pos)
	PickBlock(pos cube.Pos)
	AttackEntity(e world.Entity) bool
	Drop(s item.Stack) (n int)
	SwingArm()
	ReportSwing()
	PunchAir()
	PressJump()

	Health() float64
	MaxHealth() float64
	Absorption() float64
	Food() int

	ExperienceLevel() int
	ExperienceProgress() float64
	SetExperienceLevel(level int)

	EnchantmentSeed() int64
	ResetEnchantmentSeed()

	Respawn() *world.EntityHandle
	Dead() bool

	StartSneaking()
	Sneaking() bool
	StopSneaking()
	StartSprinting()
	Sprinting() bool
	StopSprinting()
	StartSwimming()
	Swimming() bool
	StopSwimming()
	StartCrawling()
	Crawling() bool
	StopCrawling()
	StartFlying()
	Flying() bool
	StopFlying()
	StartGliding()
	Gliding() bool
	StopGliding()
	Jump()

	StartBreaking(pos cube.Pos, face cube.Face)
	ContinueBreaking(face cube.Face)
	FinishBreaking()
	AbortBreaking()

	Exhaust(points float64)

	OpenSign(pos cube.Pos, frontSide bool)
	EditSign(pos cube.Pos, frontText, backText string) error
	TurnLecternPage(pos cube.Pos, page int) error

	EnderChestInventory() *inventory.Inventory
	MoveItemsToInventory()

	UUID() uuid.UUID

	XUID() string

	Skin() skin.Skin
	SetSkin(skin.Skin)

	UpdateDiagnostics(Diagnostics)
}
