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
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"net"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player/bossbar"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/debug"
	"github.com/df-mc/dragonfly/server/player/dialogue"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/player/hud"
	"github.com/df-mc/dragonfly/server/player/input"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type playerData struct {
	xuid              string
	locale            language.Tag
	nameTag, scoreTag string
	alwaysShowNameTag bool
	absorptionHealth  float64
	scale             float64

	gameMode world.GameMode
	skin     skin.Skin
	s        *session.Session
	h        Handler

	closed atomic.Bool

	inv, offHand, enderChest, ui *inventory.Inventory
	armour                       *inventory.Armour
	heldSlot                     *uint32

	sneaking, sprinting, swimming, gliding, crawling, flying,
	invisible, immobile, onGround, usingItem bool

	sleeping bool
	sleepPos cube.Pos

	usingSince time.Time

	glideTicks   int64
	fireTicks    int64
	fallDistance float64

	breathing         bool
	airSupplyTicks    int
	maxAirSupplyTicks int

	cooldowns map[string]time.Time

	speed               float64
	flightSpeed         float64
	verticalFlightSpeed float64

	health     *entity.HealthManager
	experience *entity.ExperienceManager
	effects    *entity.EffectManager

	lastXPPickup *time.Time

	lastDamage  float64
	immuneUntil time.Time

	deathPos       *mgl64.Vec3
	deathDimension world.Dimension

	enchantSeed int64

	mc           *entity.MovementComputer
	portalTravel *entity.PortalTravelComputer

	collidedVertically, collidedHorizontally bool

	collisionBoxes []cube.BBox

	breaking          bool
	breakingPos       cube.Pos
	breakingFace      cube.Face
	lastBreakDuration time.Duration

	breakCounter uint32

	hunger *hungerManager

	once sync.Once

	prevWorld *world.World

	lastJumpTick uint64

	turned bool

	pendingTeleport *pendingTeleport
}

type pendingTeleport struct {
	pos   mgl64.Vec3
	since time.Time
}

const pendingTeleportPatience = time.Millisecond * 150

type Player struct {
	tx     *world.Tx
	handle *world.EntityHandle
	data   *world.EntityData
	*playerData
}

func (p *Player) H() *world.EntityHandle {
	return p.handle
}

func (p *Player) Tx() *world.Tx {
	return p.tx
}

func (p *Player) Name() string {

	return p.data.Name
}

func (p *Player) UUID() uuid.UUID {
	return p.handle.UUID()
}

func (p *Player) XUID() string {
	return p.xuid
}

func (p *Player) DeviceID() string {
	if p.session() == session.Nop {
		return ""
	}
	return string(p.session().ClientData().DeviceID)
}

func (p *Player) DeviceModel() string {
	if p.session() == session.Nop {
		return ""
	}
	return p.session().ClientData().DeviceModel
}

func (p *Player) DeviceOS() int {
	if p.session() == session.Nop {
		return 0
	}
	return int(p.session().ClientData().DeviceOS)
}

func (p *Player) InputMode() (current, standard int) {
	if p.session() == session.Nop {
		return 0, 0
	}
	data := p.session().ClientData()
	return data.CurrentInputMode, data.DefaultInputMode
}

func (p *Player) GameVersion() string {
	if p.session() == session.Nop {
		return ""
	}
	return p.session().ClientData().GameVersion
}

func (p *Player) MemoryTier() int {
	if p.session() == session.Nop {
		return 0
	}
	return int(p.session().ClientData().MemoryTier)
}

func (p *Player) SelfSignedID() string {
	if p.session() == session.Nop {
		return ""
	}
	return p.session().ClientData().SelfSignedID
}

func (p *Player) Addr() net.Addr {
	if p.session() == session.Nop {
		return nil
	}
	return p.session().Addr()
}

func (p *Player) Skin() skin.Skin {
	return p.skin
}

func (p *Player) SetSkin(skin skin.Skin) {
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleSkinChange(ctx, &skin); ctx.Cancelled() {
		p.session().ViewSkin(p)
		return
	}
	p.skin = skin
	for _, v := range p.viewers() {
		v.ViewSkin(p)
	}
}

func (p *Player) Locale() language.Tag {
	return p.locale
}

func (p *Player) Handle(h Handler) {
	if h == nil {
		h = NopHandler{}
	}
	p.h = h
}

func (p *Player) Message(a ...any) {
	p.session().SendMessage(format(a))
}

func (p *Player) Messagef(f string, a ...any) {
	p.session().SendMessage(fmt.Sprintf(f, a...))
}

func (p *Player) Messaget(t chat.Translation, a ...any) {
	p.session().SendTranslation(t, p.locale, a)
}

func (p *Player) SendPopup(a ...any) {
	p.session().SendPopup(format(a))
}

func (p *Player) SendTip(a ...any) {
	p.session().SendTip(format(a))
}

func (p *Player) SendJukeboxPopup(a ...any) {
	p.session().SendJukeboxPopup(format(a))
}

func (p *Player) SendToast(title, message string) {
	p.session().SendToast(title, message)
}

func (p *Player) ResetFallDistance() {
	p.fallDistance = 0
}

func (p *Player) FallDistance() float64 {
	return p.fallDistance
}

func (p *Player) SendTitle(t title.Title) {
	p.session().SetTitleDurations(t.FadeInDuration(), t.Duration(), t.FadeOutDuration())
	if t.Text() != "" || t.Subtitle() != "" {
		p.session().SendTitle(t.Text())
		if t.Subtitle() != "" {
			p.session().SendSubtitle(t.Subtitle())
		}
	}
	if t.ActionText() != "" {
		p.session().SendActionBarMessage(t.ActionText())
	}
}

func (p *Player) SendScoreboard(scoreboard *scoreboard.Scoreboard) {
	p.session().SendScoreboard(scoreboard)
}

func (p *Player) RemoveScoreboard() {
	p.session().RemoveScoreboard()
}

func (p *Player) SendBossBar(bar bossbar.BossBar) {
	p.session().SendBossBar(bar.Text(), bar.Colour().Uint8(), bar.HealthPercentage())
}

func (p *Player) RemoveBossBar() {
	p.session().RemoveBossBar()
}

func (p *Player) Chat(msg ...any) {
	message := format(msg)
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleChat(ctx, &message); ctx.Cancelled() {
		return
	}
	_, _ = fmt.Fprintf(chat.Global, "<%v> %v\n", p.Name(), message)
}

func (p *Player) ExecuteCommand(commandLine string) {
	if p.Dead() {
		return
	}
	args := strings.Split(commandLine, " ")

	name, ok := strings.CutPrefix(args[0], "/")
	if !ok {
		return
	}

	command, ok := cmd.ByAlias(name)
	if !ok {
		o := &cmd.Output{}
		o.Errort(cmd.MessageUnknown, name)
		p.SendCommandOutput(o)
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleCommandExecution(ctx, command, args[1:]); ctx.Cancelled() {
		return
	}
	command.Execute(strings.Join(args[1:], " "), p, p.tx)
}

func (p *Player) Transfer(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleTransfer(ctx, addr); ctx.Cancelled() {
		return nil
	}
	p.session().Transfer(addr.IP, addr.Port)
	return nil
}

func (p *Player) SendCommandOutput(output *cmd.Output) {
	p.session().SendCommandOutput(output, p.locale)
}

func (p *Player) SendDialogue(d dialogue.Dialogue, e world.Entity) {
	p.session().SendDialogue(d, e)
}

func (p *Player) CloseDialogue() {
	p.session().CloseDialogue()
}

func (p *Player) SendForm(f form.Form) {
	p.session().SendForm(f)
}

func (p *Player) CloseForm() {
	p.session().CloseForm()
}

func (p *Player) ShowCoordinates() {
	p.session().EnableCoordinates(true)
}

func (p *Player) HideCoordinates() {
	p.session().EnableCoordinates(false)
}

func (p *Player) EnableInstantRespawn() {
	p.session().EnableInstantRespawn(true)
}

func (p *Player) DisableInstantRespawn() {
	p.session().EnableInstantRespawn(false)
}

func (p *Player) SetNameTag(name string) {
	p.nameTag = name
	p.updateState()
}

func (p *Player) NameTag() string {
	return p.nameTag
}

func (p *Player) SetAlwaysShowNameTag(alwaysShow bool) {
	p.alwaysShowNameTag = alwaysShow
	p.updateState()
}

func (p *Player) AlwaysShowNameTag() bool {
	return p.alwaysShowNameTag
}

func (p *Player) SetScoreTag(a ...any) {
	tag := format(a)
	p.scoreTag = tag
	p.updateState()
}

func (p *Player) ScoreTag() string {
	return p.scoreTag
}

func (p *Player) SetSpeed(speed float64) {
	p.speed = speed
	p.session().SendSpeed(speed)
}

func (p *Player) Speed() float64 {
	return p.speed
}

func (p *Player) SetFlightSpeed(flightSpeed float64) {
	p.flightSpeed = flightSpeed
	p.session().SendAbilities(p)
}

func (p *Player) FlightSpeed() float64 {
	return p.flightSpeed
}

func (p *Player) SetVerticalFlightSpeed(flightSpeed float64) {
	p.verticalFlightSpeed = flightSpeed
	p.session().SendAbilities(p)
}

func (p *Player) VerticalFlightSpeed() float64 {
	return p.verticalFlightSpeed
}

func (p *Player) Health() float64 {
	return p.health.Health()
}

func (p *Player) MaxHealth() float64 {
	return p.health.MaxHealth()
}

func (p *Player) SetMaxHealth(health float64) {
	p.health.SetMaxHealth(health)
	p.session().SendHealth(p.Health(), p.MaxHealth(), p.absorptionHealth)
}

func (p *Player) addHealth(health float64) {
	p.health.AddHealth(health)
	p.session().SendHealth(p.Health(), p.MaxHealth(), p.absorptionHealth)
}

func (p *Player) Heal(health float64, source world.HealingSource) float64 {
	if p.Dead() || health < 0 || !p.GameMode().AllowsTakingDamage() {
		return 0
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleHeal(ctx, &health, source); ctx.Cancelled() {
		return 0
	}
	oldHealth := p.Health()
	p.addHealth(health)
	return p.Health() - oldHealth
}

func (p *Player) updateFallState(distanceThisTick float64) {
	switch {
	case p.OnGround():
		p.fallDistance -= distanceThisTick
		if p.fallDistance > 3 {
			p.fall(p.fallDistance)
		}
		p.ResetFallDistance()
	case distanceThisTick < 0 && distanceThisTick < p.fallDistance:
		p.fallDistance -= distanceThisTick
	default:
		p.ResetFallDistance()
	}
}

func (p *Player) fall(distance float64) {
	if pos, lander, ok := p.landedOn(); ok {
		lander.EntityLand(pos, p.tx, p, &distance)
	}
	if !p.tx.World().FallDamage() {

		return
	}
	dmg := distance - 3
	if boost, ok := p.Effect(effect.JumpBoost); ok {
		dmg -= float64(boost.Level())
	}
	if dmg < 0.5 {
		return
	}
	p.Hurt(math.Ceil(dmg), entity.FallDamageSource{})
}

func (p *Player) landedOn() (cube.Pos, block.EntityLander, bool) {
	low, high := p.blocksUnder()
	for x := low[0]; x <= high[0]; x++ {
		for z := low[2]; z <= high[2]; z++ {
			pos := cube.Pos{x, low[1], z}
			if lander, ok := p.tx.Block(pos).(block.EntityLander); ok {
				return pos, lander, true
			}
		}
	}
	return cube.Pos{}, nil, false
}

func (p *Player) blocksUnder() (low, high cube.Pos) {
	box := Type.BBox(p).Translate(p.Position())

	y := int(math.Floor(box.Min()[1] - 0.0001))
	horizontal := box.Grow(-0.0001)
	low, high = cube.PosFromVec3(horizontal.Min()), cube.PosFromVec3(horizontal.Max())
	low[1], high[1] = y, y
	return low, high
}

func (p *Player) Hurt(dmg float64, src world.DamageSource) (float64, bool) {
	if _, ok := p.Effect(effect.FireResistance); (ok && src.Fire()) || p.Dead() || !p.GameMode().AllowsTakingDamage() || dmg < 0 {
		return 0, false
	}
	totalDamage := p.FinalDamageFrom(dmg, src)
	damageLeft := totalDamage

	immune := p.AttackImmune()
	if immune {

		if _, melee := src.(entity.AttackDamageSource); melee {
			return 0, false
		}
		if damageLeft -= p.lastDamage; damageLeft <= 0 {
			return 0, false
		}
	}

	immunity := attackCooldownTicks * time.Second / 20
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleHurt(ctx, &damageLeft, immune, &immunity, src); ctx.Cancelled() {
		return 0, false
	}
	p.setAttackImmunity(immunity, totalDamage)

	if a := p.Absorption(); a > 0 {
		remaining := a - damageLeft
		p.SetAbsorption(remaining)
		damageLeft = max(0, damageLeft-a)
		if _, exists := p.Effect(effect.Absorption); exists && remaining <= 0 {
			p.RemoveEffect(effect.Absorption)
		}
	}

	if p.Health()-damageLeft <= mgl64.Epsilon && !src.IgnoreTotem() {
		hand, offHand := p.HeldItems()
		if _, ok := offHand.Item().(item.Totem); ok {
			p.applyTotemEffects()
			p.SetHeldItems(hand, offHand.Grow(-1))
			return 0, false
		} else if _, ok := hand.Item().(item.Totem); ok {
			p.applyTotemEffects()
			p.SetHeldItems(hand.Grow(-1), offHand)
			return 0, false
		}
	}

	p.addHealth(-damageLeft)

	if src.ReducedByArmour() {
		p.Exhaust(0.1)
		p.Armour().Damage(dmg, p.damageItem)

		var origin world.Entity
		if s, ok := src.(entity.AttackDamageSource); ok {
			origin = s.Attacker
		} else if s, ok := src.(entity.ProjectileDamageSource); ok {
			origin = s.Owner
		}
		if l, ok := origin.(entity.Living); ok {
			if thornsDmg := p.Armour().ThornsDamage(p.damageItem); thornsDmg > 0 {
				l.Hurt(thornsDmg, enchantment.ThornsDamageSource{Owner: p})
			}
		}
	}

	pos := p.Position()
	for _, viewer := range p.viewers() {
		viewer.ViewEntityAction(p, entity.HurtAction{})
	}
	if src.Fire() {
		p.tx.PlaySound(pos, sound.Burning{})
	} else if _, ok := src.(entity.DrowningDamageSource); ok {
		p.tx.PlaySound(pos, sound.Drowning{})
	}

	p.Wake()

	if p.Dead() {
		p.kill(src)
	}
	return totalDamage, true
}

func (p *Player) applyTotemEffects() {
	p.addHealth(2 - p.Health())

	for _, e := range p.Effects() {
		p.RemoveEffect(e.Type())
	}

	p.AddEffect(effect.New(effect.Regeneration, 2, time.Second*40))
	p.AddEffect(effect.New(effect.FireResistance, 1, time.Second*40))
	p.AddEffect(effect.New(effect.Absorption, 2, time.Second*5))

	p.tx.PlaySound(p.Position(), sound.Totem{})

	for _, viewer := range p.viewers() {
		viewer.ViewEntityAction(p, entity.TotemUseAction{})
	}
}

func (p *Player) FinalDamageFrom(dmg float64, src world.DamageSource) float64 {
	dmg = max(dmg, 0)

	dmg -= p.Armour().DamageReduction(dmg, src)
	if res, ok := p.Effect(effect.Resistance); ok {
		dmg *= effect.Resistance.Multiplier(src, res.Level())
	}
	return dmg
}

func (p *Player) Explode(src world.ExplosionSource, impact float64) {
	explosionPos := src.Position()
	diff := p.Position().Sub(explosionPos)
	p.Hurt(math.Floor((impact*impact+impact)*3.5*src.Size()*2+1), entity.ExplosionDamageSource{Source: src})
	p.knockBack(explosionPos, impact, diff[1]/diff.Len()*impact)
}

func (p *Player) SetAbsorption(health float64) {
	p.absorptionHealth = max(health, 0)
	p.session().SendHealth(p.Health(), p.MaxHealth(), p.absorptionHealth)
}

func (p *Player) Absorption() float64 {
	return p.absorptionHealth
}

func (p *Player) KnockBack(src mgl64.Vec3, force, height float64) {
	if p.Dead() || !p.GameMode().AllowsTakingDamage() {
		return
	}
	p.knockBack(src, force, height)
}

var (
	attackCooldownTicks = time.Duration(8)

	attackImmunitySlack = 20 * time.Millisecond
)

func (p *Player) AttackImmune() bool {
	return time.Now().Add(attackImmunitySlack).Before(p.immuneUntil)
}

func (p *Player) knockBack(src mgl64.Vec3, force, height float64) {
	dir := p.Position().Sub(src)
	dir[1] = 0
	if dir.Len() > 0 {
		dir = dir.Normalize()
	}
	resist := p.Armour().KnockBackResistance()
	if resist >= 1 {
		return
	}

	motion := mgl64.Vec3{dir[0] * force, height, dir[2] * force}
	if resist > 0 {
		motion = motion.Mul(1 - resist)
	}
	p.SetVelocity(motion)
}

func (p *Player) setAttackImmunity(d time.Duration, dmg float64) {
	p.immuneUntil = time.Now().Add(d)
	p.lastDamage = dmg
}

func (p *Player) Food() int {
	return p.hunger.Food()
}

func (p *Player) SetFood(level int) {
	p.hunger.SetFood(level)
	p.sendFood()
}

func (p *Player) AddFood(points int) {
	p.hunger.AddFood(points)
	p.sendFood()
}

func (p *Player) Saturate(food int, saturation float64) {
	p.hunger.saturate(food, saturation)
	p.sendFood()
}

func (p *Player) sendFood() {
	p.session().SendFood(p.hunger.foodLevel, p.hunger.saturationLevel, p.hunger.exhaustionLevel)
}

func (p *Player) AddEffect(e effect.Effect) {
	p.session().SendEffect(p.effects.Add(e, p))
	p.updateState()
}

func (p *Player) RemoveEffect(e effect.Type) {
	p.effects.Remove(e, p)
	p.session().SendEffectRemoval(e)
	p.updateState()
}

func (p *Player) Effect(e effect.Type) (effect.Effect, bool) {
	return p.effects.Effect(e)
}

func (p *Player) Effects() []effect.Effect {
	return p.effects.Effects()
}

func (*Player) BeaconAffected() bool {
	return true
}

var HungerDepletes = false

var NaturalRegeneration = false

func (p *Player) Exhaust(points float64) {

	if !HungerDepletes || !p.GameMode().AllowsTakingDamage() || p.tx.World().Difficulty().FoodRegenerates() {
		return
	}
	before := p.hunger.Food()
	p.hunger.exhaust(points)
	if after := p.hunger.Food(); before != after {

		p.hunger.SetFood(before)

		ctx := NewEventContext(p.tx, p)
		if p.Handler().HandleFoodLoss(ctx, before, &after); ctx.Cancelled() {

			p.hunger.resetExhaustion()
			return
		}
		p.hunger.SetFood(after)
		if before >= 7 && after <= 6 {

			p.StopSprinting()
		}
	}
	p.sendFood()
}

func (p *Player) Dead() bool {
	return p.Health() <= mgl64.Epsilon
}

func (p *Player) DeathPosition() (mgl64.Vec3, world.Dimension, bool) {
	if p.deathPos == nil {
		return mgl64.Vec3{}, nil, false
	}
	return *p.deathPos, p.deathDimension, true
}

func (p *Player) kill(src world.DamageSource) {
	for _, viewer := range p.viewers() {
		viewer.ViewEntityAction(p, entity.DeathAction{})
	}

	p.addHealth(-p.MaxHealth())

	keepInv := false
	p.Handler().HandleDeath(p, src, &keepInv)
	p.StopSneaking()
	p.StopSprinting()

	pos := p.Position()
	if !keepInv {
		p.dropItems()
	}
	for _, e := range p.Effects() {
		p.RemoveEffect(e.Type())
	}

	p.deathPos, p.deathDimension = &pos, p.tx.World().Dimension()

	DoAfter(p.handle, time.Millisecond*1100, func(_ *world.Tx, p *Player) {
		finishDying(p)
	})
}

func finishDying(p *Player) {
	if p.session() == session.Nop {
		_ = p.Close()
		return
	}
	if p.Dead() {
		p.SetInvisible()

		pos, _, _, _ := p.spawnLocation()

		p.data.Pos = pos.Vec3()
	}
}

func (p *Player) dropItems() {
	pos := p.Position()
	for _, orb := range entity.NewExperienceOrbs(pos, int(math.Min(float64(p.experience.Level()*7), 100))) {
		p.tx.AddEntity(orb)
	}
	p.experience.Reset()
	p.session().SendExperience(p.ExperienceLevel(), p.ExperienceProgress())

	p.MoveItemsToInventory()
	for _, it := range append(p.inv.Clear(), append(p.armour.Clear(), p.offHand.Clear()...)...) {
		if _, ok := it.Enchantment(enchantment.CurseOfVanishing); ok {
			continue
		}
		opts := world.EntitySpawnOpts{Position: pos, Velocity: mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1}}
		p.tx.AddEntity(entity.NewItem(opts, it))
	}
}

func (p *Player) MoveItemsToInventory() {
	for _, i := range p.ui.Clear() {
		if n, err := p.inv.AddItem(i); err != nil {

			p.Drop(i.Grow(-n))
		}
	}
}

func (p *Player) Respawn() *world.EntityHandle {
	p.respawn(nil)
	return p.handle
}

func (p *Player) respawn(f func(p *Player)) {
	if !p.Dead() || p.session() == session.Nop {
		return
	}

	blockPos, w, spawnObstructed, _ := p.spawnLocation()
	pos := blockPos.Vec3Middle()

	if spawnObstructed {
		p.Messaget(chat.MessageBedNotValid)
	}

	p.addHealth(p.MaxHealth())
	p.hunger.Reset()
	p.sendFood()
	p.Extinguish()
	p.ResetFallDistance()

	p.Handler().HandleRespawn(p, &pos, &w)

	sess := p.session()
	src := p.tx.World()
	handle := p.tx.RemoveEntity(p)

	restore := func(tx *world.Tx) {
		np := tx.AddEntity(handle).(*Player)
		if f != nil {
			f(np)
			return
		}
		np.quit("respawn failed")
	}
	task := w.Do(func(tx *world.Tx) {
		np := tx.AddEntity(handle).(*Player)
		np.Teleport(pos)
		np.session().SendRespawn(pos, p)
		np.SetVisible()
		if f != nil {
			f(np)
		}
	})
	if errors.Is(task.Err(), world.ErrWorldClosed) {

		restore(p.tx)
		return
	}
	task.OnDone(func(err error) {

		if !errors.Is(err, world.ErrWorldClosed) {
			return
		}

		src.Do(restore).OnDone(func(err error) {
			if err == nil || errors.Is(err, world.ErrTaskPanicked) {
				return
			}

			_ = handle.Close()
			sess.Disconnect("respawn failed")
			sess.Close(nil, p)
			sess.CloseConnection()
		})
	})
}

func (p *Player) spawnLocation() (playerSpawn cube.Pos, w *world.World, spawnBlockBroken bool, previousDimension world.Dimension) {
	tx := p.tx
	w = tx.World()
	previousDimension = w.Dimension()
	playerSpawn = w.PlayerSpawn(p.UUID())
	if b, ok := tx.Block(playerSpawn).(block.Bed); ok && b.CanRespawnOn() {
		pos, ok := b.SafeSpawn(playerSpawn, tx)
		if ok {
			return pos, w, false, previousDimension
		}
	}

	w = w.PortalDestination(w.Dimension())
	worldSpawn := w.Spawn()
	return worldSpawn, w, playerSpawn != worldSpawn, previousDimension
}

func (p *Player) StartSprinting() {
	if !p.hunger.canSprint() && p.GameMode().AllowsTakingDamage() || p.crawling || p.sprinting {
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleToggleSprint(ctx, true); ctx.Cancelled() {
		return
	}
	p.StopSneaking()
	p.sprinting = true
	p.SetSpeed(p.speed * 1.3)
	p.updateState()
}

func (p *Player) Sprinting() bool {
	return p.sprinting
}

func (p *Player) StopSprinting() {
	if !p.sprinting {
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleToggleSprint(ctx, false); ctx.Cancelled() {
		return
	}
	p.sprinting = false
	p.SetSpeed(p.speed / 1.3)
	p.updateState()
}

func (p *Player) StartSneaking() {
	if p.sneaking {
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleToggleSneak(ctx, true); ctx.Cancelled() {
		return
	}
	if !p.Flying() {
		p.StopSprinting()
	}
	p.sneaking = true
	p.updateState()
}

func (p *Player) Sneaking() bool {
	return p.sneaking
}

func (p *Player) StopSneaking() {
	if !p.sneaking {
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleToggleSneak(ctx, false); ctx.Cancelled() {
		return
	}
	p.sneaking = false
	p.updateState()
}

func (p *Player) StartSwimming() {
	if p.swimming {
		return
	}
	p.StopSneaking()
	p.swimming = true
	p.updateState()
}

func (p *Player) Swimming() bool {
	return p.swimming
}

func (p *Player) StopSwimming() {
	if !p.swimming {
		return
	}
	p.swimming = false
	p.updateState()
}

func (p *Player) Splash(*world.Tx, mgl64.Vec3) {
	if d := p.OnFireDuration(); d.Seconds() <= 0 {
		return
	}
	p.Extinguish()
}

func (p *Player) StartCrawling() {
	if p.crawling {
		return
	}
	for _, corner := range p.H().Type().BBox(p).Translate(p.Position()).Corners() {
		if _, isAir := p.tx.Block(cube.PosFromVec3(corner).Add(cube.Pos{0, 1, 0})).(block.Air); !isAir {
			p.crawling = true
			break
		}
	}
	p.StopSneaking()
	p.updateState()
}

func (p *Player) StopCrawling() {
	if !p.crawling {
		return
	}
	p.crawling = false
	p.updateState()
}

func (p *Player) Crawling() bool {
	return p.crawling
}

func (p *Player) StartGliding() {
	if p.gliding {
		return
	}
	chest := p.Armour().Chestplate()
	if _, ok := chest.Item().(item.Elytra); !ok || chest.Durability() < 2 {
		return
	}
	p.gliding = true
	p.updateState()
}

func (p *Player) Gliding() bool {
	return p.gliding
}

func (p *Player) StopGliding() {
	if !p.gliding {
		return
	}
	p.gliding = false
	p.glideTicks = 0
	p.updateState()
}

func (p *Player) StartFlying() {
	if !p.GameMode().AllowsFlying() || p.Flying() {
		return
	}
	p.flying = true
	p.session().SendAbilities(p)
}

func (p *Player) Flying() bool {
	return p.flying
}

func (p *Player) StopFlying() {
	if !p.flying {
		return
	}
	p.flying = false
	p.session().SendAbilities(p)
}

func (p *Player) Jump() {
	if p.Dead() {
		return
	}

	p.Handler().HandleJump(p)
	if p.OnGround() {
		jumpVel := 0.42
		if e, ok := p.Effect(effect.JumpBoost); ok {
			jumpVel = float64(e.Level()) / 10
		}
		p.data.Vel = mgl64.Vec3{0, jumpVel}
	}
	if p.Sprinting() {
		p.Exhaust(0.2)
	} else {
		p.Exhaust(0.05)
	}
}

func (p *Player) Sleep(pos cube.Pos) {
	if p.sleeping {

		return
	}

	tx := p.tx
	b, ok := tx.Block(pos).(block.Bed)
	if !ok || b.Sleeper != nil {

		return
	}

	ctx, sendReminder := NewEventContext(p.tx, p), true
	if p.Handler().HandleSleep(ctx, &sendReminder); ctx.Cancelled() {
		return
	}

	b.Sleeper = p.H()
	tx.SetBlock(pos, b, nil)

	tx.World().SetRequiredSleepDuration(time.Millisecond * 5050)

	p.data.Pos = pos.Vec3Middle().Add(mgl64.Vec3{0, 0.5625})
	p.sleeping = true
	p.sleepPos = pos

	p.session().SendPlayerSpawn(pos.Vec3())

	if sendReminder {
		tx.BroadcastSleepingReminder(p)
	}

	tx.BroadcastSleepingIndicator()
	p.updateState()
}

func (p *Player) Wake() {
	if !p.sleeping {
		return
	}
	p.sleeping = false

	tx := p.tx
	tx.BroadcastSleepingIndicator()

	for _, v := range p.viewers() {
		v.ViewEntityWake(p)
	}
	p.updateState()

	pos := p.sleepPos
	if b, ok := tx.Block(pos).(block.Bed); ok {
		b.Sleeper = nil
		tx.SetBlock(pos, b, nil)
	}
}

func (p *Player) Sleeping() (cube.Pos, bool) {
	if !p.sleeping {
		return cube.Pos{}, false
	}
	return p.sleepPos, true
}

func (p *Player) SendSleepingIndicator(sleeping, max int) {
	p.session().ViewSleepingPlayers(sleeping, max)
}

func (p *Player) SetInvisible() {
	if p.Invisible() {
		return
	}
	p.invisible = true
	p.updateState()
}

func (p *Player) SetVisible() {
	if _, ok := p.Effect(effect.Invisibility); ok || !p.GameMode().Visible() || !p.invisible {
		return
	}
	p.invisible = false
	p.updateState()
}

func (p *Player) Invisible() bool {
	return p.invisible
}

func (p *Player) SetImmobile() {
	if p.Immobile() {
		return
	}
	p.immobile = true
	p.updateState()
}

func (p *Player) SetMobile() {
	if !p.Immobile() {
		return
	}
	p.immobile = false
	p.updateState()
}

func (p *Player) Immobile() bool {
	return p.immobile
}

func (p *Player) FireProof() bool {
	if _, ok := p.Effect(effect.FireResistance); ok {
		return true
	}
	return !p.GameMode().AllowsTakingDamage()
}

func (p *Player) OnFireDuration() time.Duration {
	return time.Duration(p.fireTicks) * time.Second / 20
}

func (p *Player) SetOnFire(duration time.Duration) {
	ticks := int64(duration.Seconds() * 20)
	if level := p.Armour().HighestEnchantmentLevel(enchantment.FireProtection); level > 0 {
		ticks -= int64(math.Floor(float64(ticks) * float64(level) * 0.15))
	}

	duration = time.Duration(ticks) * time.Second / 20
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleSetOnFire(ctx, &duration); ctx.Cancelled() {
		return
	}

	ticks = int64(duration.Seconds() * 20)
	p.fireTicks = ticks
	p.updateState()
}

func (p *Player) Extinguish() {
	p.SetOnFire(0)
}

func (p *Player) Inventory() *inventory.Inventory {
	return p.inv
}

func (p *Player) Armour() *inventory.Armour {
	return p.armour
}

func (p *Player) HeldItems() (mainHand, offHand item.Stack) {
	offHand, _ = p.offHand.Item(0)
	mainHand, _ = p.inv.Item(int(*p.heldSlot))
	return mainHand, offHand
}

func (p *Player) OffHand() *inventory.Inventory {
	return p.offHand
}

func (p *Player) SetHeldItems(mainHand, offHand item.Stack) {
	_ = p.inv.SetItem(int(*p.heldSlot), mainHand)
	_ = p.offHand.SetItem(0, offHand)
}

func (p *Player) HeldSlot() int {
	return int(*p.heldSlot)
}

func (p *Player) SetHeldSlot(to int) error {

	if to < 0 || to > 8 {
		return fmt.Errorf("held slot exceeds hotbar range 0-8: slot is %v", to)
	}
	from := int(*p.heldSlot)
	if from == to {

		return nil
	}

	ctx := NewEventContext(p.tx, p)
	p.Handler().HandleHeldSlotChange(ctx, from, to)
	if ctx.Cancelled() {

		p.session().SendHeldSlot(from, p, true)
		return nil
	}
	*p.heldSlot = uint32(to)
	if p.usingItem {

		p.usingItem = false
		p.updateState()
	}

	for _, viewer := range p.viewers() {
		viewer.ViewEntityItems(p)
	}
	p.session().SendHeldSlot(to, p, false)
	return nil
}

func (p *Player) EnderChestInventory() *inventory.Inventory {
	return p.enderChest
}

func (p *Player) SetGameMode(mode world.GameMode) {
	previous := p.GameMode()
	p.gameMode = mode

	if !mode.AllowsFlying() {
		p.StopFlying()
	}
	if !mode.Visible() {
		p.SetInvisible()
	} else if !previous.Visible() {
		p.SetVisible()
	}

	p.session().SendGameMode(p)
	for _, v := range p.viewers() {
		v.ViewEntityGameMode(p)
	}
	if mode.AllowsTakingDamage() {
		p.session().SendHealth(p.Health(), p.MaxHealth(), p.absorptionHealth)
	}
}

func (p *Player) GameMode() world.GameMode {
	return p.gameMode
}

func (p *Player) HasCooldown(item world.Item) bool {
	if item == nil {
		return false
	}
	name, _ := item.EncodeItem()
	otherTime, ok := p.cooldowns[name]
	if !ok {
		return false
	}
	if time.Now().After(otherTime) {
		delete(p.cooldowns, name)
		return false
	}
	return true
}

func (p *Player) SetCooldown(item world.Item, cooldown time.Duration) {
	if item == nil {
		return
	}
	name, _ := item.EncodeItem()
	p.cooldowns[name] = time.Now().Add(cooldown)
	p.session().ViewItemCooldown(item, cooldown)
}

func (p *Player) ResetCooldowns() {
	for name := range p.cooldowns {
		delete(p.cooldowns, name)
		if it, ok := world.ItemByName(name, 0); ok {
			p.session().ViewItemCooldown(it, 0)
		}
	}
}

func (p *Player) UseItem() {
	i, _ := p.HeldItems()
	ctx := NewEventContext(p.tx, p)
	if p.HasCooldown(i.Item()) {
		return
	}
	if p.Handler().HandleItemUse(ctx); ctx.Cancelled() {
		return
	}
	i, left := p.HeldItems()
	it := i.Item()

	if cd, ok := it.(item.Cooldown); ok {
		p.SetCooldown(it, cd.Cooldown())
	}

	if _, ok := it.(item.Releasable); ok {
		if !p.canRelease() {
			return
		}
		p.usingSince, p.usingItem = time.Now(), true
		p.updateState()
	}
	switch usable := it.(type) {
	case item.Chargeable:
		useCtx := p.useContext()
		if !p.usingItem {
			if !usable.ReleaseCharge(p, p.tx, useCtx) && usable.CanCharge(p, p.tx, useCtx) {

				p.usingSince, p.usingItem = time.Now(), true
			}
			p.handleUseContext(useCtx)
			p.updateState()
			return
		}

		p.usingItem = false
		dur := p.useDuration()
		if usable.Charge(p, p.tx, useCtx, dur) {
			p.session().SendChargeItemComplete()
		}
		p.handleUseContext(useCtx)
		p.updateState()
	case item.Usable:
		useCtx := p.useContext()
		if !usable.Use(p.tx, p, useCtx) {
			return
		}

		p.SwingArm()
		p.SetHeldItems(p.subtractItem(p.damageItem(i, useCtx.Damage), useCtx.CountSub), left)
		p.addNewItem(useCtx)
	case item.Consumable:
		if c, ok := usable.(interface{ CanConsume() bool }); ok && !c.CanConsume() {
			p.ReleaseItem()
			return
		}
		if !usable.AlwaysConsumable() && p.GameMode().AllowsTakingDamage() && p.Food() >= 20 {

			p.ReleaseItem()
			return
		}
		if !p.usingItem {

			p.usingItem, p.usingSince = true, time.Now()
			p.updateState()
			return
		}

		useCtx, dur := p.useContext(), p.useDuration()
		if dur < usable.ConsumeDuration() {

			return
		}

		p.usingSince = time.Now()
		ctx := NewEventContext(p.tx, p)
		if p.Handler().HandleItemConsume(ctx, i); ctx.Cancelled() {
			return
		}
		useCtx.CountSub, useCtx.NewItem = 1, usable.Consume(p.tx, p)
		p.handleUseContext(useCtx)
		p.tx.PlaySound(p.Position().Add(mgl64.Vec3{0, 1.5}), sound.Burp{})
	}
}

func (p *Player) ReleaseItem() {
	if !p.usingItem || !p.canRelease() || !p.GameMode().AllowsInteraction() {
		wasUsing := p.usingItem
		p.usingItem = false
		if wasUsing {

			p.updateState()
		}
		return
	}
	p.usingItem = false

	useCtx, dur := p.useContext(), p.useDuration()
	i, _ := p.HeldItems()
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleItemRelease(ctx, i, dur); ctx.Cancelled() {
		return
	}
	i.Item().(item.Releasable).Release(p, p.tx, useCtx, dur)
	p.handleUseContext(useCtx)
	p.updateState()
}

func (p *Player) canRelease() bool {
	held, left := p.HeldItems()
	releasable, ok := held.Item().(item.Releasable)
	if !ok {
		return false
	}
	if p.GameMode().CreativeInventory() {
		return true
	}
	for _, req := range releasable.Requirements() {
		reqName, _ := req.Item().EncodeItem()

		if !left.Empty() {
			leftName, _ := left.Item().EncodeItem()
			if leftName == reqName {
				continue
			}
		}

		_, found := p.Inventory().FirstFunc(func(stack item.Stack) bool {
			name, _ := stack.Item().EncodeItem()
			return name == reqName
		})
		if !found {
			return false
		}
	}
	return true
}

func (p *Player) handleUseContext(ctx *item.UseContext) {
	i, left := p.HeldItems()

	p.SetHeldItems(p.subtractItem(p.damageItem(i, ctx.Damage), ctx.CountSub), left)
	p.addNewItem(ctx)
	for _, it := range ctx.ConsumedItems {
		_, offHand := p.HeldItems()
		if offHand.Comparable(it) {
			if err := p.offHand.RemoveItem(it); err == nil {
				continue
			}

			it = it.Grow(-offHand.Count())
		}

		_ = p.Inventory().RemoveItem(it)
	}
}

func (p *Player) useDuration() time.Duration {
	return time.Since(p.usingSince) + time.Second/20
}

func (p *Player) UsingItem() bool {
	return p.usingItem
}

func (p *Player) UseItemOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	if _, ok := p.tx.Block(pos).(block.Air); ok || !p.canReach(pos.Vec3Centre()) {

		p.resendNearbyBlocks(pos, face)
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleItemUseOnBlock(ctx, pos, face, clickPos); ctx.Cancelled() {
		p.resendNearbyBlocks(pos, face)
		return
	}
	i, left := p.HeldItems()
	b := p.tx.Block(pos)
	if act, ok := b.(block.Activatable); ok {

		if !p.Sneaking() || i.Empty() {

			if useCtx := p.useContext(); act.Activate(pos, face, p.tx, p, useCtx) {
				p.SwingArm()
				p.SetHeldItems(p.subtractItem(p.damageItem(i, useCtx.Damage), useCtx.CountSub), left)
				p.addNewItem(useCtx)
				return
			}
		}
	}
	if i.Empty() {
		return
	}
	switch ib := i.Item().(type) {
	case item.UsableOnBlock:

		useCtx := p.useContext()
		if !ib.UseOnBlock(pos, face, clickPos, p.tx, p, useCtx) {

			p.resendNearbyBlocks(pos, face)
			return
		}
		p.SwingArm()
		p.SetHeldItems(p.subtractItem(p.damageItem(i, useCtx.Damage), useCtx.CountSub), left)
		p.addNewItem(useCtx)
	case world.Block:

		replacedPos := pos
		if replaceable, ok := b.(block.Replaceable); !ok || !replaceable.ReplaceableBy(ib) {

			replacedPos = pos.Side(face)
		}
		if replaceable, ok := p.tx.Block(replacedPos).(block.Replaceable); !ok || !replaceable.ReplaceableBy(ib) || replacedPos.OutOfBounds(p.tx.Range()) {
			return
		}
		if !p.placeBlock(replacedPos, ib, false) || p.GameMode().CreativeInventory() {
			return
		}
		p.SetHeldItems(p.subtractItem(i, 1), left)
	}
}

func (p *Player) UseItemOnEntity(e world.Entity) bool {
	if !p.canReach(e.Position()) {
		return false
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleItemUseOnEntity(ctx, e); ctx.Cancelled() {
		return false
	}
	i, left := p.HeldItems()
	usable, ok := i.Item().(item.UsableOnEntity)
	if !ok {
		return true
	}
	useCtx := p.useContext()
	if !usable.UseOnEntity(e, p.tx, p, useCtx) {
		return true
	}
	p.SwingArm()
	p.SetHeldItems(p.subtractItem(p.damageItem(i, useCtx.Damage), useCtx.CountSub), left)
	p.addNewItem(useCtx)
	return true
}

func (p *Player) AttackEntity(e world.Entity) bool {
	if !p.canReach(e.Position()) {
		return false
	}

	living, isLiving := e.(entity.Living)
	if isLiving && living.Dead() {
		return false
	}

	var (
		force, height  = 0.44, 0.385
		_, slowFalling = p.Effect(effect.SlowFalling)
		_, blind       = p.Effect(effect.Blindness)

		critical = !p.Sprinting() && !p.Flying() && p.FallDistance() > 0 && !slowFalling && !blind
	)

	i, _ := p.HeldItems()
	if k, ok := i.Enchantment(enchantment.Knockback); ok {
		inc := enchantment.Knockback.Force(k.Level())
		force += inc
		height += inc
	}

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleAttackEntity(ctx, e, &force, &height, &critical); ctx.Cancelled() {
		return false
	}
	p.SwingArm()

	if !isLiving {
		if !entity.DamageableEntity(e) {
			return false
		}
		i, left := p.HeldItems()
		if durable, ok := i.Item().(item.Durable); ok {
			p.SetHeldItems(p.damageItem(i, durable.DurabilityInfo().AttackDurability), left)
		}
		n, vulnerable, _ := entity.HurtEntity(e, i.AttackDamage(), entity.AttackDamageSource{Attacker: p})
		p.tx.PlaySound(entity.EyePosition(e), sound.Attack{Damage: !mgl64.FloatEqual(n, 0)})
		if vulnerable {
			p.Exhaust(0.1)
		}
		return true
	}

	dmg := i.AttackDamage()
	if strength, ok := p.Effect(effect.Strength); ok {
		dmg += dmg * effect.Strength.Multiplier(strength.Level())
	}
	if weakness, ok := p.Effect(effect.Weakness); ok {
		dmg -= dmg * effect.Weakness.Multiplier(weakness.Level())
	}
	if s, ok := i.Enchantment(enchantment.Sharpness); ok {
		dmg += enchantment.Sharpness.Addend(s.Level())
		for _, v := range p.tx.Viewers(living.Position()) {
			v.ViewEntityAction(living, entity.EnchantedHitAction{})
		}
	}
	if critical {
		dmg *= 1.5
	}

	n, vulnerable := living.Hurt(dmg, entity.AttackDamageSource{Attacker: p})
	i, left := p.HeldItems()

	if durable, ok := i.Item().(item.Durable); ok {
		p.SetHeldItems(p.damageItem(i, durable.DurabilityInfo().AttackDurability), left)
	}

	p.tx.PlaySound(entity.EyePosition(e), sound.Attack{Damage: !mgl64.FloatEqual(n, 0)})
	if !vulnerable {

		return true
	}
	if critical {
		for _, v := range p.tx.Viewers(living.Position()) {
			v.ViewEntityAction(living, entity.CriticalHitAction{})
		}
	}

	p.Exhaust(0.1)

	living.KnockBack(p.Position(), force, height)

	if f, ok := i.Enchantment(enchantment.FireAspect); ok {
		if flammable, ok := living.(entity.Flammable); ok {
			flammable.SetOnFire(enchantment.FireAspect.Duration(f.Level()))
		}
	}
	return true
}

func (p *Player) StartBreaking(pos cube.Pos, face cube.Face) {
	p.AbortBreaking()
	if _, air := p.tx.Block(pos).(block.Air); air || !p.canReach(pos.Vec3Centre()) {

		return
	}
	if _, ok := p.tx.Block(pos.Side(face)).(block.Fire); ok {
		ctx := NewEventContext(p.tx, p)
		if p.Handler().HandleFireExtinguish(ctx, pos); ctx.Cancelled() {

			p.resendNearbyBlocks(pos, face)
			return
		}

		p.tx.SetBlock(pos.Side(face), nil, nil)
		p.tx.PlaySound(pos.Vec3(), sound.FireExtinguish{})
		return
	}

	held, _ := p.HeldItems()
	if _, ok := held.Item().(item.Sword); ok && p.GameMode().CreativeInventory() {

		return
	}

	p.breakingPos = pos

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleStartBreak(ctx, pos); ctx.Cancelled() {
		return
	}
	if punchable, ok := p.tx.Block(pos).(block.Punchable); ok {
		punchable.Punch(pos, face, p.tx, p)
	}

	p.breaking, p.breakingFace = true, face
	p.SwingArm()

	if p.GameMode().CreativeInventory() {
		return
	}
	p.lastBreakDuration = p.breakTime(pos)
	for _, viewer := range p.viewers() {
		viewer.ViewBlockAction(pos, block.StartCrackAction{BreakTime: p.lastBreakDuration})
	}
}

func (p *Player) breakTime(pos cube.Pos) time.Duration {
	held, _ := p.HeldItems()
	return block.BreakDuration(p.tx.Block(pos), held, p.breakContext())
}

func (p *Player) breakContext() block.BreakContext {
	_, aquaAffinity := p.Armour().Helmet().Enchantment(enchantment.AquaAffinity)
	ctx := block.BreakContext{
		Underwater:   p.insideOfWater(),
		AquaAffinity: aquaAffinity,
		Airborne:     !p.OnGround(),
		Flying:       p.Flying(),
	}
	if e, ok := p.Effect(effect.Haste); ok {
		ctx.HasteLevel = e.Level()
	}
	if e, ok := p.Effect(effect.ConduitPower); ok {
		ctx.ConduitPowerLevel = e.Level()
	}
	if e, ok := p.Effect(effect.MiningFatigue); ok {
		ctx.MiningFatigueLevel = e.Level()
	}
	return ctx
}

func (p *Player) FinishBreaking() {
	if !p.breaking {
		p.resendNearbyBlock(p.breakingPos)
		return
	}
	p.AbortBreaking()
	p.BreakBlock(p.breakingPos)
}

func (p *Player) AbortBreaking() {
	if !p.breaking {
		return
	}
	p.breaking, p.breakCounter = false, 0
	for _, viewer := range p.viewers() {
		viewer.ViewBlockAction(p.breakingPos, block.StopCrackAction{})
	}
}

func (p *Player) ContinueBreaking(face cube.Face) {
	if !p.breaking {
		return
	}
	pos := p.breakingPos
	b := p.tx.Block(pos)
	p.tx.AddParticle(pos.Vec3(), particle.PunchBlock{Block: b, Face: face})

	if p.breakCounter++; p.breakCounter%5 == 0 {
		p.SwingArm()

		p.tx.PlaySound(pos.Vec3(), sound.BlockBreaking{Block: b})
	}
	if breakTime := p.breakTime(pos); breakTime != p.lastBreakDuration {
		for _, viewer := range p.viewers() {
			viewer.ViewBlockAction(pos, block.ContinueCrackAction{BreakTime: breakTime})
		}
		p.lastBreakDuration = breakTime
	}
}

func (p *Player) PlaceBlock(pos cube.Pos, b world.Block, ctx *item.UseContext) {
	ignoreBBox := ctx != nil && ctx.IgnoreBBox
	if !p.placeBlock(pos, b, ignoreBBox) {
		return
	}
	if ctx != nil {
		ctx.CountSub++
	}
}

func (p *Player) placeBlock(pos cube.Pos, b world.Block, ignoreBBox bool) bool {
	if !p.canReach(pos.Vec3Centre()) || !p.GameMode().AllowsEditing() {
		p.resendNearbyBlocks(pos, cube.Faces()...)
		return false
	}
	if obstructed, selfOnly := p.obstructedPos(pos, b); obstructed && !ignoreBBox {
		if !selfOnly {

			p.resendNearbyBlocks(pos, cube.Faces()...)
		}
		return false
	}

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleBlockPlace(ctx, pos, b); ctx.Cancelled() {
		p.resendNearbyBlocks(pos, cube.Faces()...)
		return false
	}
	p.tx.SetBlock(pos, b, nil)
	p.tx.PlaySound(pos.Vec3(), sound.BlockPlace{Block: b})
	p.SwingArm()
	return true
}

func (p *Player) obstructedPos(pos cube.Pos, b world.Block) (obstructed, selfOnly bool) {

	blockBoxes := slices.Clone(b.Model().BBox(pos, p.tx))
	for i, box := range blockBoxes {
		blockBoxes[i] = box.Translate(pos.Vec3())
	}

	for e := range p.tx.EntitiesWithin(cube.Box(-3, -3, -3, 3, 3, 3).Translate(pos.Vec3())) {
		t := e.H().Type()
		switch t {
		case entity.ItemType, entity.ArrowType, entity.ExperienceOrbType:
			continue
		default:
			if cube.AnyIntersections(blockBoxes, t.BBox(e).Translate(e.Position()).Grow(-1e-4)) {
				obstructed = true
				if e.H() == p.handle {
					continue
				}
				return true, false
			}
		}
	}
	return obstructed, true
}

func (p *Player) BreakBlock(pos cube.Pos) {
	b := p.tx.Block(pos)
	if _, air := b.(block.Air); air {

		return
	}
	if !p.canReach(pos.Vec3Centre()) || !p.GameMode().AllowsEditing() {
		p.resendNearbyBlocks(pos)
		return
	}
	if _, breakable := b.(block.Breakable); !breakable && !p.GameMode().CreativeInventory() {
		p.resendNearbyBlocks(pos)
		return
	}
	held, _ := p.HeldItems()
	drops := p.drops(held, b)

	xp := 0
	if breakable, ok := b.(block.Breakable); ok && !p.GameMode().CreativeInventory() {
		if _, hasSilkTouch := held.Enchantment(enchantment.SilkTouch); !hasSilkTouch {
			xp = breakable.BreakInfo().XPDrops.RandomValue()
		}
	}

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleBlockBreak(ctx, pos, &drops, &xp); ctx.Cancelled() {
		p.resendNearbyBlocks(pos)
		return
	}
	held, left := p.HeldItems()

	p.SwingArm()
	p.tx.SetBlock(pos, nil, nil)
	p.tx.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: b})

	if breakable, ok := b.(block.Breakable); ok {
		info := breakable.BreakInfo()
		if info.BreakHandler != nil {
			info.BreakHandler(pos, p.tx, p)
		}
		for _, orb := range entity.NewExperienceOrbs(pos.Vec3Centre(), xp) {
			p.tx.AddEntity(orb)
		}
	}
	for _, drop := range drops {
		opts := world.EntitySpawnOpts{Position: pos.Vec3Centre(), Velocity: mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1}}
		p.tx.AddEntity(entity.NewItem(opts, drop))
	}

	p.Exhaust(0.005)

	if block.BreaksInstantly(b) {
		return
	}
	if durable, ok := held.Item().(item.Durable); ok {
		p.SetHeldItems(p.damageItem(held, durable.DurabilityInfo().BreakDurability), left)
	}
}

func (p *Player) drops(held item.Stack, b world.Block) []item.Stack {
	t, ok := held.Item().(item.Tool)
	if !ok {
		t = item.ToolNone{}
	}
	var drops []item.Stack
	if breakable, ok := b.(block.Breakable); ok && !p.GameMode().CreativeInventory() {
		if breakable.BreakInfo().Harvestable(t) {
			drops = breakable.BreakInfo().Drops(t, held.Enchantments())
		}
	} else if it, ok := b.(world.Item); ok && !p.GameMode().CreativeInventory() {
		drops = []item.Stack{item.NewStack(it, 1)}
	}
	return drops
}

func (p *Player) PickBlock(pos cube.Pos) {
	if !p.canReach(pos.Vec3()) {
		return
	}

	b := p.tx.Block(pos)

	var pickedItem item.Stack
	if pi, ok := b.(block.Pickable); ok {
		pickedItem = pi.Pick()
	} else if i, ok := b.(world.Item); ok {
		it, _ := world.ItemByName(i.EncodeItem())
		pickedItem = item.NewStack(it, 1)
	} else {
		return
	}

	slot, found := p.Inventory().First(pickedItem)
	if !found && !p.GameMode().CreativeInventory() {
		return
	}

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleBlockPick(ctx, pos, b); ctx.Cancelled() {
		return
	}
	_, offhand := p.HeldItems()

	if found {
		if slot < 9 {
			_ = p.SetHeldSlot(slot)
			return
		}
		_ = p.Inventory().Swap(slot, int(*p.heldSlot))
		return
	}

	firstEmpty, emptyFound := p.Inventory().FirstEmpty()
	if !emptyFound {
		p.SetHeldItems(pickedItem, offhand)
		return
	}
	if firstEmpty < 9 {
		_ = p.SetHeldSlot(firstEmpty)
		_ = p.Inventory().SetItem(firstEmpty, pickedItem)
		return
	}
	_ = p.Inventory().Swap(firstEmpty, int(*p.heldSlot))
	p.SetHeldItems(pickedItem, offhand)
}

func (p *Player) Teleport(pos mgl64.Vec3) {
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleTeleport(ctx, pos); ctx.Cancelled() {
		return
	}
	p.forceTeleport(pos)
}

func (p *Player) TransferTo(w *world.World, pos mgl64.Vec3) {
	if w == nil || w == p.tx.World() {
		p.Teleport(pos)
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleTeleport(ctx, pos); ctx.Cancelled() {
		return
	}
	p.Wake()
	p.StopSneaking()
	p.StopSwimming()
	p.StopGliding()
	h, source, origin := p.handle, p.tx.World(), p.Position()
	p.tx.Defer(func(tx *world.Tx) {
		e, ok := h.Entity(tx)
		if !ok {
			return
		}
		handle := tx.RemoveEntity(e)
		if handle == nil {
			return
		}
		go func() {
			_, err := world.Call(context.Background(), w, func(tx *world.Tx) (struct{}, error) {
				if p, ok := tx.AddEntityAt(handle, pos).(*Player); ok {
					if p.closed.Load() {

						tx.RemoveEntity(p)
						_ = handle.Close()
						return struct{}{}, nil
					}

					p.session().EnterWorld(tx, p)
				}
				return struct{}{}, nil
			})
			if err == nil {
				return
			}

			if _, err = world.Call(context.Background(), source, func(tx *world.Tx) (struct{}, error) {
				tx.AddEntityAt(handle, origin)
				return struct{}{}, nil
			}); err != nil {
				_ = handle.Close()
			}
		}()
	})
}

func (p *Player) TeleportFromProjectile(pos mgl64.Vec3) {
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleTeleport(ctx, pos); ctx.Cancelled() {
		return
	}
	p.Wake()
	p.pendingTeleport = &pendingTeleport{pos: pos, since: time.Now()}
}

func (p *Player) settleTeleport() {
	pending := p.pendingTeleport
	if pending == nil {
		return
	}
	p.pendingTeleport = nil
	p.teleportBecause(pending.pos, session.TeleportCauseProjectile)
}

func (p *Player) forceTeleport(pos mgl64.Vec3) {
	p.Wake()
	p.teleport(pos)
}

func (p *Player) teleport(pos mgl64.Vec3) {
	p.teleportBecause(pos, session.TeleportCauseUnknown)
}

func (p *Player) teleportBecause(pos mgl64.Vec3, cause session.TeleportCause) {
	for _, v := range p.viewers() {
		if s, ok := v.(*session.Session); ok {
			s.ViewEntityTeleportCause(p, pos, cause)
			continue
		}
		v.ViewEntityTeleport(p, pos)
	}
	p.data.Pos = pos
	p.data.Vel = mgl64.Vec3{}
	p.ResetFallDistance()
}

func (p *Player) Move(deltaPos mgl64.Vec3, deltaYaw, deltaPitch float64) {

	defer p.settleTeleport()
	if p.Dead() || (deltaPos.ApproxEqual(mgl64.Vec3{}) && mgl64.FloatEqual(deltaYaw, 0) && mgl64.FloatEqual(deltaPitch, 0)) {
		p.onGround = true
		p.updateFallState(deltaPos.Y())
		return
	}
	if p.immobile {
		if mgl64.FloatEqual(deltaYaw, 0) && mgl64.FloatEqual(deltaPitch, 0) {

			return
		}

		deltaPos = mgl64.Vec3{}
	}
	var (
		pos         = p.Position()
		res, resRot = pos.Add(deltaPos), p.Rotation().Add(cube.Rotation{deltaYaw, deltaPitch})
	)
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleMove(ctx, res, resRot); ctx.Cancelled() {
		if p.session() != session.Nop && pos.ApproxEqual(p.Position()) {

			p.teleport(pos)
		}
		return
	}
	for _, v := range p.viewers() {
		v.ViewEntityMovement(p, res, resRot, p.OnGround())
	}

	p.data.Pos = res
	p.data.Rot = resRot
	if deltaPos.Len() <= 3 {

		p.data.Vel = deltaPos
		p.checkBlockCollisions(deltaPos)
	}

	horizontalVel := deltaPos
	horizontalVel[1] = 0
	if p.Gliding() {
		if deltaPos.Y() >= -0.5 {
			p.fallDistance = 1.0
		}
		if p.collidedHorizontally {
			if force := horizontalVel.Len()*10.0 - 3.0; force > 0.0 {
				p.tx.PlaySound(p.Position(), sound.Fall{Distance: force})
				p.Hurt(force, entity.GlideDamageSource{})
			}
		}
	}

	_, submergedBefore := p.tx.Liquid(cube.PosFromVec3(pos.Add(mgl64.Vec3{0, p.EyeHeight()})))
	_, submergedAfter := p.tx.Liquid(cube.PosFromVec3(res.Add(mgl64.Vec3{0, p.EyeHeight()})))
	if submergedBefore != submergedAfter {

		p.session().ViewEntityState(p)
	}

	p.onGround = p.checkOnGround(deltaPos)
	p.updateFallState(deltaPos.Y())

	if p.Swimming() {
		p.Exhaust(0.01 * horizontalVel.Len())
	} else if p.Sprinting() {
		p.Exhaust(0.1 * horizontalVel.Len())
	}
}

func (p *Player) Displace(deltaPos mgl64.Vec3) {
	if p.Dead() || deltaPos.ApproxEqual(mgl64.Vec3{}) {
		return
	}
	pos := p.Position()
	deltaPos, velocity := p.mc.CheckCollision(p.tx, p, pos, deltaPos)
	if deltaPos.ApproxEqual(mgl64.Vec3{}) {
		return
	}
	res := pos.Add(deltaPos)
	for _, v := range p.viewers() {
		v.ViewEntityDisplacement(p, res, p.Rotation(), p.OnGround())
	}
	p.data.Pos, p.data.Vel = res, velocity
	p.checkBlockCollisions(deltaPos)
	p.onGround = p.checkOnGround(deltaPos)
	p.updateFallState(deltaPos[1])
}

func (p *Player) Position() mgl64.Vec3 {
	return p.data.Pos
}

func (p *Player) Velocity() mgl64.Vec3 {
	return p.data.Vel
}

func (p *Player) SetVelocity(velocity mgl64.Vec3) {
	if p.session() == session.Nop {

		p.data.Vel = velocity
		if velocity[1] > 0 {
			p.mc.Airborne()
		}
		return
	}
	for _, v := range p.viewers() {
		v.ViewEntityVelocity(p, velocity)
	}
}

func (p *Player) Steer(delta mgl64.Vec3) {
	if p.session() != session.Nop {
		return
	}
	p.data.Vel = p.data.Vel.Add(delta)
}

func (p *Player) Look(rot cube.Rotation) {
	if p.session() != session.Nop || p.Dead() {
		return
	}
	p.data.Rot = rot
	p.turned = true
}

func (p *Player) Rotation() cube.Rotation {
	return p.data.Rot
}

func (p *Player) Collect(s item.Stack) (int, bool) {
	if p.Dead() || !p.GameMode().AllowsInteraction() {
		return 0, false
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleItemPickup(ctx, &s); ctx.Cancelled() {
		return 0, false
	}
	var added int
	if _, offHand := p.HeldItems(); !offHand.Empty() && offHand.Comparable(s) {
		added, _ = p.offHand.AddItem(s)
	}
	if s.Count() != added {
		n, _ := p.Inventory().AddItem(s.Grow(-added))
		added += n
	}
	return added, true
}

func (p *Player) Experience() int {
	return p.experience.Experience()
}

func (p *Player) EnchantmentSeed() int64 {
	return p.enchantSeed
}

func (p *Player) ResetEnchantmentSeed() {
	p.enchantSeed = rand.Int64()
}

func (p *Player) AddExperience(amount int) int {
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleExperienceGain(ctx, &amount); ctx.Cancelled() {
		return 0
	}
	before := p.experience.Level()
	level, _ := p.experience.Add(amount)
	if level/5 > before/5 {
		p.PlaySound(sound.LevelUp{})
	} else if amount > 0 {
		p.PlaySound(sound.Experience{})
	}
	p.session().SendExperience(p.ExperienceLevel(), p.ExperienceProgress())
	return amount
}

func (p *Player) RemoveExperience(amount int) {
	p.experience.Add(-amount)
	p.session().SendExperience(p.ExperienceLevel(), p.ExperienceProgress())
}

func (p *Player) ExperienceLevel() int {
	return p.experience.Level()
}

func (p *Player) SetExperienceLevel(level int) {
	p.experience.SetLevel(level)
	p.session().SendExperience(p.ExperienceLevel(), p.ExperienceProgress())
}

func (p *Player) ExperienceProgress() float64 {
	return p.experience.Progress()
}

func (p *Player) SetExperienceProgress(progress float64) {
	p.experience.SetProgress(progress)
	p.session().SendExperience(p.ExperienceLevel(), p.ExperienceProgress())
}

func (p *Player) CanCollectExperience() bool {
	if p.Dead() || !p.GameMode().AllowsInteraction() {
		return false
	}
	if last := p.lastXPPickup; last != nil && time.Since(*last) < time.Millisecond*100 {
		return false
	}
	return true
}

func (p *Player) CollectExperience(value int) bool {
	if !p.CanCollectExperience() {
		return false
	}
	value = p.mendItems(value)
	now := time.Now()
	p.lastXPPickup = &now
	if value > 0 {
		return p.AddExperience(value) > 0
	}

	p.PlaySound(sound.Experience{})
	return true
}

func (p *Player) mendItems(xp int) int {
	mendingItems := make([]item.Stack, 0, 6)
	held, offHand := p.HeldItems()
	if _, ok := offHand.Enchantment(enchantment.Mending); ok && offHand.Durability() < offHand.MaxDurability() {
		mendingItems = append(mendingItems, offHand)
	}
	if _, ok := held.Enchantment(enchantment.Mending); ok && held.Durability() < held.MaxDurability() {
		mendingItems = append(mendingItems, held)
	}
	for _, i := range p.Armour().Items() {
		if i.Durability() == i.MaxDurability() {
			continue
		}
		if _, ok := i.Enchantment(enchantment.Mending); ok {
			mendingItems = append(mendingItems, i)
		}
	}
	length := len(mendingItems)
	if length == 0 {
		return xp
	}
	foundItem := mendingItems[rand.IntN(length)]
	repairAmount := math.Min(float64(foundItem.MaxDurability()-foundItem.Durability()), float64(xp*2))
	repairedItem := foundItem.WithDurability(foundItem.Durability() + int(repairAmount))
	if repairAmount >= 2 {

		xp -= int(math.Ceil(repairAmount / 2))
	}
	if offHand.Equal(foundItem) {
		p.SetHeldItems(held, repairedItem)
	} else if held.Equal(foundItem) {
		p.SetHeldItems(repairedItem, offHand)
	} else if slot, ok := p.Armour().Inventory().First(foundItem); ok {
		_ = p.Armour().Inventory().SetItem(slot, repairedItem)
	}
	return xp
}

func (p *Player) Drop(s item.Stack) int {
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleItemDrop(ctx, s); ctx.Cancelled() {
		return 0
	}
	opts := world.EntitySpawnOpts{Position: p.Position().Add(mgl64.Vec3{0, 1.4}), Velocity: p.Rotation().Vec3().Mul(0.4)}
	p.tx.AddEntity(entity.NewItemPickupDelay(opts, s, time.Second*2))
	return s.Count()
}

func (p *Player) OpenBlockContainer(pos cube.Pos, tx *world.Tx) {
	if p.session() != session.Nop {
		p.session().OpenBlockContainer(pos, tx)
	}
}

func (p *Player) HideEntity(e world.Entity) {
	if p.session() != session.Nop && p.H() != e.H() {
		p.session().StopShowingEntity(e)
	}
}

func (p *Player) ForgetHiddenEntity(id uuid.UUID) {
	if p.session() != session.Nop {
		p.session().ForgetHiddenEntity(id)
	}
}

func (p *Player) ShowEntity(e world.Entity) {
	if p.session() != session.Nop {
		p.session().StartShowingEntity(e)
	}
}

func (p *Player) Latency() time.Duration {
	if p.session() == session.Nop {
		return 0
	}
	return p.session().Latency()
}

func (p *Player) Tick(tx *world.Tx, current int64) {

	p.session().Tick(tx, p)
	if p.Dead() {
		return
	}
	if pending := p.pendingTeleport; pending != nil && time.Since(pending.since) >= pendingTeleportPatience {

		p.settleTeleport()
	}
	if _, ok := p.tx.Liquid(cube.PosFromVec3(p.Position())); !ok {
		p.StopSwimming()
		if _, ok := p.Armour().Helmet().Item().(item.TurtleShell); ok {
			p.AddEffect(effect.New(effect.WaterBreathing, 1, time.Second*10).WithoutParticles())
		}
	}

	if _, ok := p.Armour().Chestplate().Item().(item.Elytra); ok && p.Gliding() {
		if p.glideTicks += 1; p.glideTicks%20 == 0 {
			d := p.damageItem(p.Armour().Chestplate(), 1)
			p.armour.SetChestplate(d)
			if d.Durability() < 2 {
				p.StopGliding()
			}
		}
	}

	if p.GameMode().HasCollision() {

		p.checkBlockCollisions(p.data.Vel)
		p.onGround = p.checkOnGround(mgl64.Vec3{})
		p.checkEntitySteppers()
	}

	p.effects.Tick(p, p.tx)

	p.tickFood()
	if p.GameMode().HasCollision() {
		p.tickAirSupply()
	}

	if p.Position()[1] < float64(p.tx.Range()[0]) {
		p.Hurt(4, entity.VoidDamageSource{})
	}
	if p.GameMode().HasCollision() && p.insideOfSolid() {
		p.Hurt(1, entity.SuffocationDamageSource{})
	}

	if p.OnFireDuration() > 0 {
		p.fireTicks -= 1
		if !p.GameMode().AllowsTakingDamage() || p.OnFireDuration() <= 0 || p.tx.RainingAt(cube.PosFromVec3(p.Position())) {
			p.Extinguish()
		}
		if p.OnFireDuration()%time.Second == 0 {
			p.Hurt(1, block.FireDamageSource{})
		}
	}

	held, _ := p.HeldItems()
	if p.usingItem {
		switch held.Item().(type) {
		case item.Consumable, item.Releasable, item.Chargeable:
		default:

			p.usingItem = false
			p.updateState()
		}
	}
	if current%4 == 0 && p.usingItem {
		if _, ok := held.Item().(item.Consumable); ok {

			for _, v := range p.viewers() {
				v.ViewEntityAction(p, entity.EatAction{})
			}
		}
	}

	if p.usingItem {
		if c, ok := held.Item().(item.Chargeable); ok {
			c.ContinueCharge(p, tx, p.useContext(), p.useDuration())
		}
	}
	if p.breaking {
		p.ContinueBreaking(p.breakingFace)
	}

	for it, ti := range p.cooldowns {
		if time.Now().After(ti) {
			delete(p.cooldowns, it)
		}
	}

	p.session().SendDebugShapes(tx.World().Dimension())
	p.session().SendHudUpdates()

	if p.prevWorld != tx.World() && p.prevWorld != nil {
		p.Handler().HandleChangeWorld(p, p.prevWorld, tx.World())
	}
	p.prevWorld = tx.World()

	if p.session() == session.Nop && !p.Immobile() {

		m := p.mc.TickMovementVanilla(p, p.Position(), p.Velocity(), p.Rotation(), p.tx)
		if delta := m.Position().Sub(p.Position()); !delta.ApproxEqualThreshold(mgl64.Vec3{}, 0.001) {
			p.Move(delta, 0, 0)
			p.turned = false
		} else if p.turned {
			p.turned = false
			for _, v := range p.viewers() {
				v.ViewEntityMovement(p, p.Position(), p.Rotation(), p.OnGround())
			}
		}
		p.data.Vel = m.Velocity()
	} else {
		p.data.Vel = mgl64.Vec3{}
	}

	p.portalTravel.StopPortalContact()
}

func (p *Player) SelfViewing() bool {
	return p.session() != session.Nop
}

func (p *Player) TravelThroughPortal(tx *world.Tx, target world.Dimension) {
	if !p.GameMode().HasCollision() {

		return
	}
	p.portalTravel.EnterPortal(p, tx, target)
}

func (p *Player) ViewLayer() *world.ViewLayer {
	return p.session().ViewLayer()
}

func (p *Player) ViewNameTag(entity world.Entity, nameTag string) {
	p.session().ViewNameTag(entity, nameTag)
}

func (p *Player) ViewPublicNameTag(entity world.Entity) {
	p.session().ViewPublicNameTag(entity)
}

func (p *Player) ViewAlwaysShowNameTag(entity world.Entity, alwaysShow bool) {
	p.session().ViewAlwaysShowNameTag(entity, alwaysShow)
}

func (p *Player) ViewPublicAlwaysShowNameTag(entity world.Entity) {
	p.session().ViewPublicAlwaysShowNameTag(entity)
}

func (p *Player) ViewScoreTag(entity world.Entity, scoreTag string) {
	p.session().ViewScoreTag(entity, scoreTag)
}

func (p *Player) ViewPublicScoreTag(entity world.Entity) {
	p.session().ViewPublicScoreTag(entity)
}

func (p *Player) ViewVisibility(entity world.Entity, level world.VisibilityLevel) {
	p.session().ViewVisibility(entity, level)
}

func (p *Player) RemoveViewLayer(entity world.Entity) {
	p.session().RemoveViewLayer(entity)
}

func (p *Player) tickAirSupply() {
	if !p.canBreathe() {
		if r, ok := p.Armour().Helmet().Enchantment(enchantment.Respiration); ok && rand.Float64() < enchantment.Respiration.Chance(r.Level()) {

			return
		}
		if p.airSupplyTicks -= 1; p.airSupplyTicks <= -20 {
			p.airSupplyTicks = 0
			p.Hurt(2, entity.DrowningDamageSource{})
		}
		if p.breathing {

			p.breathing = false
			p.updateState()
		}
	} else if !p.breathing && p.airSupplyTicks < p.maxAirSupplyTicks {
		p.airSupplyTicks = min(p.airSupplyTicks+5, p.maxAirSupplyTicks)
		p.breathing = p.airSupplyTicks == p.maxAirSupplyTicks
		p.updateState()
	}
}

func (p *Player) tickFood() {
	if p.hunger.foodTick%10 == 0 && (p.hunger.canQuicklyRegenerate() || p.tx.World().Difficulty().FoodRegenerates()) {
		if p.tx.World().Difficulty().FoodRegenerates() {
			p.AddFood(1)
		}
		if p.hunger.foodTick%20 == 0 {
			p.regenerate(true)
		}
	}
	if p.hunger.foodTick == 1 {
		if p.hunger.canRegenerate() {
			p.regenerate(!p.tx.World().Difficulty().FoodRegenerates())
		} else if p.hunger.starving() {
			p.starve()
		}
	}

	if !p.hunger.canSprint() {
		p.StopSprinting()
	}

	p.hunger.foodTick++
	if p.hunger.foodTick > 80 {
		p.hunger.foodTick = 1
	}
}

func (p *Player) regenerate(exhaust bool) {

	if !NaturalRegeneration || p.Health() == p.MaxHealth() {
		return
	}
	regenerated := p.Heal(1, entity.FoodHealingSource{QuickRegeneration: exhaust})
	if exhaust && regenerated > 0 {
		p.Exhaust(6)
	}
}

func (p *Player) starve() {
	if p.Health() > p.tx.World().Difficulty().StarvationHealthLimit() {
		p.Hurt(1, StarvationDamageSource{})
	}
}

func (p *Player) AirSupply() time.Duration {
	return time.Duration(p.airSupplyTicks) * time.Second / 20
}

func (p *Player) SetAirSupply(duration time.Duration) {
	p.airSupplyTicks = int(duration.Milliseconds() / 50)
	p.updateState()
}

func (p *Player) MaxAirSupply() time.Duration {
	return time.Duration(p.maxAirSupplyTicks) * time.Second / 20
}

func (p *Player) SetMaxAirSupply(duration time.Duration) {
	p.maxAirSupplyTicks = int(duration.Milliseconds() / 50)
	p.updateState()
}

func (p *Player) canBreathe() bool {
	canTakeDamage := p.GameMode().AllowsTakingDamage()
	_, waterBreathing := p.effects.Effect(effect.WaterBreathing)
	_, conduitPower := p.effects.Effect(effect.ConduitPower)
	return !canTakeDamage || waterBreathing || conduitPower || (!p.insideOfWater() && !p.insideOfSolid())
}

const breathingDistanceBelowEyes = 0.11111111

func (p *Player) insideOfWater() bool {
	pos := cube.PosFromVec3(entity.EyePosition(p))
	if l, ok := p.tx.Liquid(pos); ok {
		if _, ok := l.(block.Water); ok {
			d := float64(l.SpreadDecay()) + 1
			if l.LiquidFalling() {
				d = 1
			}
			return p.Position().Y() < (pos.Side(cube.FaceUp).Vec3().Y())-(d/9-breathingDistanceBelowEyes)
		}
	}
	return false
}

func (p *Player) insideOfSolid() bool {
	pos := cube.PosFromVec3(entity.EyePosition(p))
	b, box := p.tx.Block(pos), p.handle.Type().BBox(p).Translate(p.Position())

	_, solid := b.Model().(model.Solid)
	if !solid {

		return false
	}
	d, diffuses := b.(block.LightDiffuser)
	if diffuses && d.LightDiffusionLevel() == 0 {

		return false
	}
	if immune, ok := b.(block.NonSuffocating); ok && immune.PreventsSuffocation() {
		return false
	}
	for _, blockBox := range b.Model().BBox(pos, p.tx) {
		if blockBox.Translate(pos.Vec3()).IntersectsWith(box) {
			return true
		}
	}
	return false
}

func (p *Player) checkBlockCollisions(vel mgl64.Vec3) {
	entityBBox := Type.BBox(p).Translate(p.Position())
	deltaX, deltaY, deltaZ := vel[0], vel[1], vel[2]

	p.checkEntityInsiders(entityBBox)

	grown := entityBBox.Extend(vel).Grow(0.25)
	low, high := grown.Min(), grown.Max()
	minX, minY, minZ := int(math.Floor(low[0])), int(math.Floor(low[1])), int(math.Floor(low[2]))
	maxX, maxY, maxZ := int(math.Ceil(high[0])), int(math.Ceil(high[1])), int(math.Ceil(high[2]))

	blocks := p.collisionBoxes[:0]
	defer func() { p.collisionBoxes = blocks[:0] }()
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			for z := minZ; z <= maxZ; z++ {
				pos := cube.Pos{x, y, z}
				boxes := p.tx.Block(pos).Model().BBox(pos, p.tx)
				for _, box := range boxes {
					blocks = append(blocks, box.Translate(pos.Vec3()))
				}
			}
		}
	}

	const epsilon = 0.001

	if !mgl64.FloatEqualThreshold(deltaY, 0, epsilon) {

		for _, blockBBox := range blocks {
			deltaY = entityBBox.YOffset(blockBBox, deltaY)
		}
		entityBBox = entityBBox.Translate(mgl64.Vec3{0, deltaY})
	}
	if !mgl64.FloatEqualThreshold(deltaX, 0, epsilon) {

		for _, blockBBox := range blocks {
			deltaX = entityBBox.XOffset(blockBBox, deltaX)
		}
		entityBBox = entityBBox.Translate(mgl64.Vec3{deltaX})
	}
	if !mgl64.FloatEqualThreshold(deltaZ, 0, epsilon) {

		for _, blockBBox := range blocks {
			deltaZ = entityBBox.ZOffset(blockBBox, deltaZ)
		}
	}

	p.collidedHorizontally = !mgl64.FloatEqual(deltaX, vel[0]) || !mgl64.FloatEqual(deltaZ, vel[2])
	p.collidedVertically = !mgl64.FloatEqual(deltaY, vel[1])
}

func (p *Player) checkEntityInsiders(entityBBox cube.BBox) {
	box := entityBBox.Grow(-0.0001)
	low, high := cube.PosFromVec3(box.Min()), cube.PosFromVec3(box.Max())

	for y := low[1]; y <= high[1]; y++ {
		for x := low[0]; x <= high[0]; x++ {
			for z := low[2]; z <= high[2]; z++ {
				blockPos := cube.Pos{x, y, z}
				b := p.tx.Block(blockPos)
				if collide, ok := b.(block.EntityInsider); ok {
					collide.EntityInside(blockPos, p.tx, p)
					if _, liquid := b.(world.Liquid); liquid {
						continue
					}
				}

				if l, ok := p.tx.Liquid(blockPos); ok {
					if collide, ok := l.(block.EntityInsider); ok {
						collide.EntityInside(blockPos, p.tx, p)
					}
				}
			}
		}
	}
}

func (p *Player) checkEntitySteppers() {
	if !p.OnGround() {
		return
	}
	low, high := p.blocksUnder()
	for x := low[0]; x <= high[0]; x++ {
		for z := low[2]; z <= high[2]; z++ {
			pos := cube.Pos{x, low[1], z}
			if stepper, ok := p.tx.Block(pos).(block.EntityStepper); ok {
				stepper.EntityStepOn(pos, p.tx, p)
				return
			}
		}
	}
}

func (p *Player) checkOnGround(deltaPos mgl64.Vec3) bool {
	box := Type.BBox(p).Translate(p.Position()).Extend(mgl64.Vec3{0, -0.05}).Extend(deltaPos.Mul(-1.0))
	b := box.Grow(1)

	epsilon := mgl64.Vec3{mgl64.Epsilon, mgl64.Epsilon, mgl64.Epsilon}
	low, high := cube.PosFromVec3(b.Min().Add(epsilon)), cube.PosFromVec3(b.Max().Sub(epsilon))
	for x := low[0]; x <= high[0]; x++ {
		for z := low[2]; z <= high[2]; z++ {
			for y := low[1]; y < high[1]; y++ {
				pos := cube.Pos{x, y, z}
				for _, bb := range p.tx.Block(pos).Model().BBox(pos, p.tx) {
					if bb.Translate(pos.Vec3()).IntersectsWith(box) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (p *Player) Scale() float64 {
	return p.scale
}

func (p *Player) SetScale(s float64) {
	p.scale = s
	p.updateState()
}

func (p *Player) OnGround() bool {
	if p.session() == session.Nop {
		return p.mc.OnGround()
	}
	return p.onGround
}

func (p *Player) EyeHeight() float64 {
	switch {
	case p.swimming || p.crawling || p.gliding:
		return 0.4
	case p.sneaking:
		return 1.27
	default:
		return 1.62
	}
}

func (p *Player) TorsoHeight() float64 {
	return p.EyeHeight() - 0.1
}

func (p *Player) PlaySound(sound world.Sound) {
	p.session().PlaySound(sound, entity.EyePosition(p))
}

func (p *Player) ShowParticle(pos mgl64.Vec3, particle world.Particle) {
	p.session().ViewParticle(pos, particle)
}

func (p *Player) OpenSign(pos cube.Pos, frontSide bool) {
	p.session().OpenSign(pos, frontSide)
}

func (p *Player) EditSign(pos cube.Pos, frontText, backText string) error {
	sign, ok := p.tx.Block(pos).(block.Sign)
	if !ok {
		return fmt.Errorf("edit sign: no sign at position %v", pos)
	}

	if sign.Waxed {
		return nil
	} else if frontText == sign.Front.Text && backText == sign.Back.Text {
		return nil
	}

	ctx := NewEventContext(p.tx, p)
	if frontText != sign.Front.Text {
		if p.Handler().HandleSignEdit(ctx, pos, true, sign.Front.Text, frontText); ctx.Cancelled() {
			p.resendNearbyBlock(pos)
			return nil
		}
		sign.Front.Text = frontText
		sign.Front.Owner = p.XUID()
	} else {
		if p.Handler().HandleSignEdit(ctx, pos, false, sign.Back.Text, backText); ctx.Cancelled() {
			p.resendNearbyBlock(pos)
			return nil
		}
		sign.Back.Text = backText
		sign.Back.Owner = p.XUID()
	}
	p.tx.SetBlock(pos, sign, nil)
	return nil
}

func (p *Player) TurnLecternPage(pos cube.Pos, page int) error {
	lectern, ok := p.tx.Block(pos).(block.Lectern)
	if !ok {
		return fmt.Errorf("edit lectern: no lectern at position %v", pos)
	}

	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleLecternPageTurn(ctx, pos, lectern.Page, &page); ctx.Cancelled() {
		return nil
	}

	lectern.Page = page
	p.tx.SetBlock(pos, lectern, nil)
	return nil
}

func (p *Player) updateState() {
	for _, v := range p.viewers() {
		v.ViewEntityState(p)
	}
}

func (p *Player) Breathing() bool {
	_, breathing := p.Effect(effect.WaterBreathing)
	_, conduitPower := p.Effect(effect.ConduitPower)
	_, submerged := p.tx.Liquid(cube.PosFromVec3(entity.EyePosition(p)))
	return !p.GameMode().AllowsTakingDamage() || !submerged || breathing || conduitPower
}

func (p *Player) SwingArm() {
	if p.Dead() {
		return
	}
	for _, v := range p.viewers() {
		v.ViewEntityAction(p, entity.SwingArmAction{})
	}
}

const doubleJump = 6

func (p *Player) PressJump() {
	tick := p.session().ClientTick()
	mode := p.GameMode()
	if !mode.AllowsFlying() || mode.CreativeInventory() || tick == 0 {
		return
	}
	if p.lastJumpTick != 0 && tick-p.lastJumpTick <= doubleJump {
		p.lastJumpTick = 0
		if p.Flying() {
			p.StopFlying()
		} else {
			p.StartFlying()
		}
		return
	}
	p.lastJumpTick = tick
}

func (p *Player) ReportSwing() {
	if p.Dead() {
		return
	}
	p.Handler().HandleSwing(NewEventContext(p.tx, p))
}

func (p *Player) PunchAir() {
	if p.Dead() {
		return
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandlePunchAir(ctx); ctx.Cancelled() {
		return
	}
	p.SwingArm()
	p.tx.PlaySound(p.Position(), sound.Attack{})
}

func (p *Player) UpdateDiagnostics(d session.Diagnostics) {
	p.Handler().HandleDiagnostics(p, d)
}

func (p *Player) ShowHudElement(e hud.Element) {
	p.session().ShowHudElement(e)
}

func (p *Player) HideHudElement(e hud.Element) {
	p.session().HideHudElement(e)
}

func (p *Player) HudElementHidden(e hud.Element) bool {
	return p.session().HudElementHidden(e)
}

func (p *Player) AddDebugShape(shape debug.Shape) {
	p.session().AddDebugShape(shape)
}

func (p *Player) RemoveDebugShape(shape debug.Shape) {
	p.session().RemoveDebugShape(shape)
}

func (p *Player) VisibleDebugShapes() []debug.Shape {
	return p.session().VisibleDebugShapes()
}

func (p *Player) RemoveAllDebugShapes() {
	p.session().RemoveAllDebugShapes()
}

func (p *Player) LockInput(l input.Lock) {
	p.session().LockInput(l)
	p.session().SendInputLocks()
}

func (p *Player) UnlockInput(l input.Lock) {
	p.session().UnlockInput(l)
	p.session().SendInputLocks()
}

func (p *Player) ClearInputLocks() {
	p.session().ClearInputLocks()
	p.session().SendInputLocks()
}

func (p *Player) InputLocked(l input.Lock) bool {
	return p.session().InputLocked(l)
}

func (p *Player) damageItem(s item.Stack, d int) item.Stack {
	if p.GameMode().CreativeInventory() || d == 0 || s.MaxDurability() == -1 {
		return s
	}
	ctx := NewEventContext(p.tx, p)
	if p.Handler().HandleItemDamage(ctx, s, &d); ctx.Cancelled() || d <= 0 {
		return s
	}
	if e, ok := s.Enchantment(enchantment.Unbreaking); ok {
		d = enchantment.Unbreaking.Reduce(s.Item(), e.Level(), d)
	}
	if s = s.Damage(d); s.Empty() {
		p.tx.PlaySound(p.Position(), sound.ItemBreak{})
	}
	return s
}

func (p *Player) subtractItem(s item.Stack, d int) item.Stack {
	if !p.GameMode().CreativeInventory() && d != 0 {
		return s.Grow(-d)
	}
	return s
}

func (p *Player) addNewItem(ctx *item.UseContext) {
	if (ctx.NewItemSurvivalOnly && p.GameMode().CreativeInventory()) || ctx.NewItem.Empty() {
		return
	}
	held, left := p.HeldItems()
	if held.Empty() {
		p.SetHeldItems(ctx.NewItem, left)
		return
	}
	n, err := p.Inventory().AddItem(ctx.NewItem)
	if err != nil {

		p.Drop(ctx.NewItem.Grow(-n))
	}
	if p.Dead() {
		p.dropItems()
	}
}

func (p *Player) canReach(pos mgl64.Vec3) bool {
	dist := entity.EyePosition(p).Sub(pos).Len()
	return !p.Dead() && p.GameMode().AllowsInteraction() &&
		(dist <= 8.0 || (dist <= 14.0 && p.GameMode().CreativeInventory()))
}

func (p *Player) Disconnect(msg ...any) {
	p.once.Do(func() {
		p.close(format(msg))
	})
}

func (p *Player) Close() error {
	p.once.Do(func() {
		p.close("Connection closed.")
	})
	return nil
}

func (p *Player) close(msg string) {

	p.closed.Store(true)

	if p.Dead() && p.session() != nil {
		p.respawn(func(np *Player) {
			np.quit(msg)
		})
		return
	}
	p.quit(msg)
}

func (p *Player) quit(msg string) {
	p.h.HandleQuit(p)
	p.h = NopHandler{}

	if s := p.s; s != nil {
		s.Disconnect(msg)

		s.Close(p.tx, p)
		s.CloseConnection()
		return
	}

	p.tx.RemoveEntity(p)
	_ = p.handle.Close()
}

func (p *Player) Data() Config {
	p.hunger.mu.RLock()
	defer p.hunger.mu.RUnlock()
	return Config{
		Session:             p.s,
		Skin:                p.skin,
		XUID:                p.xuid,
		UUID:                p.UUID(),
		Name:                p.nameTag,
		Locale:              p.locale,
		GameMode:            p.gameMode,
		Position:            p.Position(),
		Rotation:            p.Rotation(),
		Velocity:            p.Velocity(),
		Health:              p.Health(),
		MaxHealth:           p.MaxHealth(),
		FoodTick:            p.hunger.foodTick,
		Food:                p.hunger.foodLevel,
		Exhaustion:          p.hunger.exhaustionLevel,
		Saturation:          p.hunger.saturationLevel,
		AirSupply:           p.airSupplyTicks,
		MaxAirSupply:        p.maxAirSupplyTicks,
		EnchantmentSeed:     p.enchantSeed,
		Experience:          p.experience.Experience(),
		HeldSlot:            int(*p.heldSlot),
		Inventory:           p.inv,
		OffHand:             p.offHand,
		Armour:              p.armour,
		EnderChestInventory: p.enderChest,
		FireTicks:           p.fireTicks,
		FallDistance:        p.fallDistance,
		Effects:             p.Effects(),
	}
}

func (p *Player) session() *session.Session {
	if s := p.s; s != nil {
		return s
	}
	return session.Nop
}

func (p *Player) useContext() *item.UseContext {
	call := func(ctx *inventory.Context, slot int, it item.Stack, f func(ctx *inventory.Context, slot int, it item.Stack)) error {
		if ctx.Cancelled() {
			return fmt.Errorf("action was cancelled")
		}
		f(ctx, slot, it)
		if ctx.Cancelled() {
			return fmt.Errorf("action was cancelled")
		}
		return nil
	}
	return &item.UseContext{
		SwapHeldWithArmour: func(i int) {
			src, dst, srcInv, dstInv := int(*p.heldSlot), i, p.inv, p.armour.Inventory()
			srcIt, _ := srcInv.Item(src)
			dstIt, _ := dstInv.Item(dst)

			ctx := event.C(inventory.Holder(p))
			_ = call(ctx, src, srcIt, srcInv.Handler().HandleTake)
			_ = call(ctx, src, dstIt, srcInv.Handler().HandlePlace)
			_ = call(ctx, dst, dstIt, dstInv.Handler().HandleTake)
			if err := call(ctx, dst, srcIt, dstInv.Handler().HandlePlace); err == nil {
				_ = srcInv.SetItem(src, dstIt)
				_ = dstInv.SetItem(dst, srcIt)
				p.PlaySound(sound.EquipItem{Item: srcIt.Item()})
			}
		},
		FirstFunc: func(comparable func(item.Stack) bool) (item.Stack, bool) {
			_, left := p.HeldItems()
			if !left.Empty() && comparable(left) {
				return left, true
			}
			inv := p.Inventory()
			s, ok := inv.FirstFunc(comparable)
			if !ok {
				return item.Stack{}, false
			}
			it, _ := inv.Item(s)
			return it, ok
		},
	}
}

func (p *Player) Handler() Handler {
	return p.h
}

func (p *Player) broadcastItems(int, item.Stack, item.Stack) {
	for _, viewer := range p.viewers() {
		viewer.ViewEntityItems(p)
	}
}

func (p *Player) broadcastArmour(_ int, before, after item.Stack) {
	if before.Comparable(after) && before.Empty() == after.Empty() {

		return
	}
	for _, viewer := range p.viewers() {
		viewer.ViewEntityArmour(p)
	}
}

func (p *Player) viewers() []world.Viewer {
	viewers := p.tx.Viewers(p.Position())
	var s world.Viewer = p.session()
	if slices.Index(viewers, s) == -1 && p.s != nil {
		return append(viewers, p.s)
	}
	return viewers
}

func (p *Player) withinChunkRadius(pos mgl64.Vec3) bool {
	playerChunkX, playerChunkZ := int(p.Position().X())>>4, int(p.Position().Z())>>4
	posChunkX, posChunkZ := int(pos.X())>>4, int(pos.Z())>>4
	dx, dz := playerChunkX-posChunkX, playerChunkZ-posChunkZ
	if dx < 0 {
		dx = -dx
	}
	if dz < 0 {
		dz = -dz
	}
	r := int(p.session().ChunkRadius())
	return dx <= r && dz <= r
}

func (p *Player) resendNearbyBlocks(pos cube.Pos, faces ...cube.Face) {
	if p.session() == session.Nop {
		return
	}
	p.resendNearbyBlock(pos)
	for _, f := range faces {
		p.resendNearbyBlock(pos.Side(f))
	}
}

func (p *Player) resendNearbyBlock(pos cube.Pos) {
	if p.session() == session.Nop {
		return
	}
	if !p.withinChunkRadius(pos.Vec3()) {

		return
	}
	b := p.tx.Block(pos)
	p.session().ViewBlockUpdate(pos, b, 0)
	if _, ok := b.(world.LiquidDisplacer); ok {
		liq, _ := p.tx.Liquid(pos)
		p.session().ViewBlockUpdate(pos, liq, 1)
	}
}

func format(a []any) string {
	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(a...), "\n"), "\n")
}
