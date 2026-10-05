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
	"encoding/json"
	"fmt"
	"image/color"
	"maps"
	"math"
	"net"
	"slices"
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/item/recipe"
	"github.com/df-mc/dragonfly/server/player/debug"
	"github.com/df-mc/dragonfly/server/player/dialogue"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/player/hud"
	"github.com/df-mc/dragonfly/server/player/input"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (s *Session) StopShowingEntity(e world.Entity) {
	s.entityMutex.Lock()
	_, ok := s.hiddenEntities[e.H().UUID()]
	if !ok {
		s.hiddenEntities[e.H().UUID()] = struct{}{}
	}
	s.entityMutex.Unlock()

	if !ok {
		s.HideEntity(e)
	}
}

func (s *Session) ForgetHiddenEntity(id uuid.UUID) {
	s.entityMutex.Lock()
	delete(s.hiddenEntities, id)
	s.entityMutex.Unlock()
}

func (s *Session) StartShowingEntity(e world.Entity) {
	s.entityMutex.Lock()
	_, ok := s.hiddenEntities[e.H().UUID()]
	if ok {
		delete(s.hiddenEntities, e.H().UUID())
	}
	s.entityMutex.Unlock()

	if ok {
		s.ViewEntity(e)
		s.ViewEntityState(e)
		s.ViewEntityItems(e)
		s.ViewEntityArmour(e)
	}
}

func (s *Session) closeCurrentContainer(tx *world.Tx, clientRequested bool) {
	if !s.closeWindow(clientRequested) {
		return
	}

	pos := *s.openedPos.Load()
	b := tx.Block(pos)
	if container, ok := b.(block.Container); ok {
		container.RemoveViewer(s, tx, pos)
	} else if enderChest, ok := b.(block.EnderChest); ok {
		enderChest.RemoveViewer(tx, pos)
	}
}

func (s *Session) SendRespawn(pos mgl64.Vec3, c Controllable) {
	s.writePacket(&packet.Respawn{
		Position:        vec64To32(pos.Add(entityOffset(c))),
		State:           packet.RespawnStateReadyToSpawn,
		EntityRuntimeID: selfEntityRuntimeID,
	})
}

func (s *Session) SendPlayerSpawn(pos mgl64.Vec3) {
	blockPos := protocol.BlockPos{int32(pos[0]), int32(pos[1]), int32(pos[2])}
	s.writePacket(&packet.SetSpawnPosition{
		SpawnType:     packet.SpawnTypePlayer,
		Position:      blockPos,
		Dimension:     packet.DimensionOverworld,
		SpawnPosition: blockPos,
	})
}

func (s *Session) sendBiomes() {
	s.writePacket(joinDataFor(s.br).biomes)
}

func buildBiomeDefinitionList() *packet.BiomeDefinitionList {
	definitions, stringList := world.BiomeDefinitions()
	return &packet.BiomeDefinitionList{
		BiomeDefinitions: definitions,
		StringList:       stringList,
	}
}

func (s *Session) sendRecipes() {
	data := joinDataFor(s.br)
	s.recipes = data.recipes
	s.writePacket(data.crafting)
}

func buildCraftingData(br world.BlockRegistry) (*packet.CraftingData, map[uint32]recipe.Recipe) {
	var (
		recipes = make(map[uint32]recipe.Recipe)

		shapedRecipes            []protocol.ShapedRecipe
		shapelessRecipes         []protocol.ShapelessRecipe
		userDataShapelessRecipes []protocol.UserDataShapelessRecipe
		multiRecipes             []protocol.MultiRecipe
		smithingTransformRecipes []protocol.SmithingTransformRecipe
		smithingTrimRecipes      []protocol.SmithingTrimRecipe
		potionRecipes            []protocol.PotionRecipe
		potionContainerChange    []protocol.PotionContainerChangeRecipe
	)

	for index, i := range recipe.Recipes() {
		networkID := uint32(index) + 1
		recipes[networkID] = i

		switch i := i.(type) {
		case recipe.Shapeless:
			shapelessRecipes = append(shapelessRecipes, protocol.ShapelessRecipe{
				RecipeID:        uuid.New().String(),
				Priority:        int32(i.Priority()),
				Input:           stacksToIngredientItems(br, i.Input()),
				Output:          stacksToRecipeStacks(br, i.Output()),
				Block:           i.Block(),
				RecipeNetworkID: networkID,
			})
		case recipe.UserDataShapeless:
			userDataShapelessRecipes = append(userDataShapelessRecipes, protocol.UserDataShapelessRecipe{ShapelessRecipe: protocol.ShapelessRecipe{
				RecipeID:        uuid.New().String(),
				Priority:        int32(i.Priority()),
				Input:           stacksToIngredientItems(br, i.Input()),
				Output:          stacksToRecipeStacks(br, i.Output()),
				Block:           i.Block(),
				RecipeNetworkID: networkID,
			}})
		case recipe.Multi:
			multiRecipes = append(multiRecipes, protocol.MultiRecipe{
				UUID:            i.UUID(),
				RecipeNetworkID: networkID,
			})
		case recipe.Shaped:
			shapedRecipes = append(shapedRecipes, protocol.ShapedRecipe{
				RecipeID:        uuid.New().String(),
				Priority:        int32(i.Priority()),
				Width:           int32(i.Shape().Width()),
				Height:          int32(i.Shape().Height()),
				Input:           stacksToIngredientItems(br, i.Input()),
				Output:          stacksToRecipeStacks(br, i.Output()),
				Block:           i.Block(),
				AssumeSymmetry:  true,
				RecipeNetworkID: networkID,
			})
		case recipe.SmithingTransform:
			input, output := stacksToIngredientItems(br, i.Input()), stacksToRecipeStacks(br, i.Output())
			smithingTransformRecipes = append(smithingTransformRecipes, protocol.SmithingTransformRecipe{
				RecipeID:        uuid.New().String(),
				Base:            input[0],
				Addition:        input[1],
				Template:        input[2],
				Result:          output[0],
				Block:           i.Block(),
				RecipeNetworkID: networkID,
			})
		case recipe.SmithingTrim:
			input := stacksToIngredientItems(br, i.Input())
			smithingTrimRecipes = append(smithingTrimRecipes, protocol.SmithingTrimRecipe{
				RecipeID:        uuid.New().String(),
				Base:            input[0],
				Addition:        input[1],
				Template:        input[2],
				Block:           i.Block(),
				RecipeNetworkID: networkID,
			})
		case recipe.Potion:
			inputRuntimeID, inputMeta, _ := world.ItemRuntimeID(i.Input()[0].(item.Stack).Item())
			reagentRuntimeID, reagentMeta, _ := world.ItemRuntimeID(i.Input()[1].(item.Stack).Item())
			outputRuntimeID, outputMeta, _ := world.ItemRuntimeID(i.Output()[0].Item())

			potionRecipes = append(potionRecipes, protocol.PotionRecipe{
				InputPotionID:        inputRuntimeID,
				InputPotionMetadata:  int32(inputMeta),
				ReagentItemID:        reagentRuntimeID,
				ReagentItemMetadata:  int32(reagentMeta),
				OutputPotionID:       outputRuntimeID,
				OutputPotionMetadata: int32(outputMeta),
			})

		case recipe.PotionContainerChange:
			inputRuntimeID, _, _ := world.ItemRuntimeID(i.Input()[0].(item.Stack).Item())
			reagentRuntimeID, _, _ := world.ItemRuntimeID(i.Input()[1].(item.Stack).Item())
			outputRuntimeID, _, _ := world.ItemRuntimeID(i.Output()[0].Item())

			potionContainerChange = append(potionContainerChange, protocol.PotionContainerChangeRecipe{
				InputItemID:   inputRuntimeID,
				ReagentItemID: reagentRuntimeID,
				OutputItemID:  outputRuntimeID,
			})
		}
	}
	return &packet.CraftingData{
		ShapedRecipes:                shapedRecipes,
		ShapelessRecipes:             shapelessRecipes,
		MultiRecipes:                 multiRecipes,
		UserDataShapelessRecipes:     userDataShapelessRecipes,
		SmithingTransformRecipes:     smithingTransformRecipes,
		SmithingTrimRecipes:          smithingTrimRecipes,
		PotionRecipes:                potionRecipes,
		PotionContainerChangeRecipes: potionContainerChange,
		ClearRecipes:                 true,
	}, recipes
}

func (s *Session) sendArmourTrimData() {
	s.writePacket(joinDataFor(s.br).trim)
}

func buildTrimData() *packet.TrimData {
	var trimPatterns []protocol.TrimPattern
	var trimMaterials []protocol.TrimMaterial

	for _, t := range item.SmithingTemplates() {
		if t == item.TemplateNetheriteUpgrade() {
			continue
		}
		name, _ := item.SmithingTemplate{Template: t}.EncodeItem()
		trimPatterns = append(trimPatterns, protocol.TrimPattern{
			ItemName:  name,
			PatternID: t.String(),
		})
	}

	for _, i := range item.ArmourTrimMaterials() {
		if material, ok := i.(item.ArmourTrimMaterial); ok {
			name, _ := i.EncodeItem()

			trimMaterials = append(trimMaterials, protocol.TrimMaterial{
				MaterialID: material.TrimMaterial(),
				Colour:     material.MaterialColour(),
				ItemName:   name,
			})
		}
	}

	return &packet.TrimData{Patterns: trimPatterns, Materials: trimMaterials}
}

func (s *Session) sendInv(inv *inventory.Inventory, windowID uint32) {
	pk := &packet.InventoryContent{
		WindowID: windowID,
		Content:  make([]protocol.ItemInstance, 0, inv.Size()),
	}
	for _, i := range inv.Slots() {
		pk.Content = append(pk.Content, instanceFromItem(s.br, i))
	}
	s.writePacket(pk)
}

func (s *Session) sendItem(item item.Stack, slot int, windowID uint32) {
	s.writePacket(&packet.InventorySlot{
		WindowID: windowID,
		Slot:     uint32(slot),
		NewItem:  instanceFromItem(s.br, item),
	})
}

const (
	craftingGridSizeSmall   = 4
	craftingGridSizeLarge   = 9
	craftingGridSmallOffset = 28
	craftingGridLargeOffset = 32
	craftingResult          = 50
)

type smelter interface {
	ResetExperience() int
}

func (s *Session) invByID(id int32, tx *world.Tx) (*inventory.Inventory, bool) {
	switch id {
	case protocol.ContainerCraftingInput, protocol.ContainerCreatedOutput, protocol.ContainerCursor:

		return s.ui, true
	case protocol.ContainerHotBar, protocol.ContainerInventory, protocol.ContainerCombinedHotBarAndInventory:

		return s.inv, true
	case protocol.ContainerOffhand:
		return s.offHand, true
	case protocol.ContainerArmor:

		return s.armour.Inventory(), true
	default:
		if !s.containerOpened.Load() {
			return nil, false
		}
		switch id {
		case protocol.ContainerLevelEntity:
			return s.openedWindow.Load(), true
		case protocol.ContainerShulkerBox:
			if _, shulkerbox := tx.Block(*s.openedPos.Load()).(block.ShulkerBox); shulkerbox {
				return s.openedWindow.Load(), true
			}
		case protocol.ContainerBarrel:
			if _, barrel := tx.Block(*s.openedPos.Load()).(block.Barrel); barrel {
				return s.openedWindow.Load(), true
			}
		case protocol.ContainerBeaconPayment:
			if _, beacon := tx.Block(*s.openedPos.Load()).(block.Beacon); beacon {
				return s.ui, true
			}
		case protocol.ContainerBrewingStandInput, protocol.ContainerBrewingStandResult, protocol.ContainerBrewingStandFuel:
			if _, brewingStand := tx.Block(*s.openedPos.Load()).(block.BrewingStand); brewingStand {
				return s.openedWindow.Load(), true
			}
		case protocol.ContainerAnvilInput, protocol.ContainerAnvilMaterial:
			if _, anvil := tx.Block(*s.openedPos.Load()).(block.Anvil); anvil {
				return s.ui, true
			}
		case protocol.ContainerSmithingTableTemplate, protocol.ContainerSmithingTableInput, protocol.ContainerSmithingTableMaterial:
			if _, smithing := tx.Block(*s.openedPos.Load()).(block.SmithingTable); smithing {
				return s.ui, true
			}
		case protocol.ContainerLoomInput, protocol.ContainerLoomDye, protocol.ContainerLoomMaterial:
			if _, loom := tx.Block(*s.openedPos.Load()).(block.Loom); loom {
				return s.ui, true
			}
		case protocol.ContainerStonecutterInput:
			if _, ok := tx.Block(*s.openedPos.Load()).(block.Stonecutter); ok {
				return s.ui, true
			}
		case protocol.ContainerGrindstoneInput, protocol.ContainerGrindstoneAdditional:
			if _, ok := tx.Block(*s.openedPos.Load()).(block.Grindstone); ok {
				return s.ui, true
			}
		case protocol.ContainerEnchantingInput, protocol.ContainerEnchantingMaterial:
			if _, enchanting := tx.Block(*s.openedPos.Load()).(block.EnchantingTable); enchanting {
				return s.ui, true
			}
		case protocol.ContainerFurnaceIngredient, protocol.ContainerFurnaceFuel, protocol.ContainerFurnaceResult,
			protocol.ContainerBlastFurnaceIngredient, protocol.ContainerSmokerIngredient:
			if _, ok := tx.Block(*s.openedPos.Load()).(smelter); ok {
				return s.openedWindow.Load(), true
			}
		}
	}
	return nil, false
}

func (s *Session) Disconnect(message string) {
	if s != Nop {
		_ = s.conn.WritePacket(&packet.Disconnect{
			HideDisconnectionScreen: message == "",
			Message:                 message,
		})
		_ = s.conn.Flush()
	}
}

func (s *Session) SendSpeed(speed float64) {
	s.writePacket(&packet.UpdateAttributes{
		EntityRuntimeID: selfEntityRuntimeID,
		Attributes: []protocol.Attribute{{
			AttributeValue: protocol.AttributeValue{
				Name:  "minecraft:movement",
				Value: float32(speed),
				Max:   math.MaxFloat32,
			},
			DefaultMax: math.MaxFloat32,
			Default:    0.1,
		}},
	})
}

func (s *Session) SendFood(food int, saturation, exhaustion float64) {
	s.writePacket(&packet.UpdateAttributes{
		EntityRuntimeID: selfEntityRuntimeID,
		Attributes: []protocol.Attribute{
			{
				AttributeValue: protocol.AttributeValue{
					Name:  "minecraft:player.hunger",
					Value: float32(food),
					Max:   20,
				},
				DefaultMax: 20,
				Default:    20,
			},
			{
				AttributeValue: protocol.AttributeValue{
					Name:  "minecraft:player.saturation",
					Value: float32(saturation),
					Max:   20,
				},
				DefaultMax: 20,
				Default:    20,
			},
			{
				AttributeValue: protocol.AttributeValue{
					Name:  "minecraft:player.exhaustion",
					Value: float32(exhaustion),
					Max:   5,
				},
				DefaultMax: 5,
			},
		},
	})
}

func (s *Session) SendDialogue(d dialogue.Dialogue, e world.Entity) {
	b, _ := json.Marshal(d)

	h := s.handlers[packet.IDNPCRequest].(*NPCRequestHandler)
	h.dialogue = d
	h.entityRuntimeID = s.entityRuntimeID(e)

	metadata := s.parseEntityMetadata(e)
	metadata[protocol.EntityDataKeyHasNPC] = uint8(1)

	disp := d.Display()
	disp.EntityOffset = disp.EntityOffset.Add(entityOffset(e))
	display, _ := json.Marshal(map[string]any{"portrait_offsets": disp})
	metadata[protocol.EntityDataKeyNPCData] = string(display)

	s.writePacket(&packet.SetActorData{
		EntityRuntimeID: h.entityRuntimeID,
		EntityMetadata:  metadata,
	})
	s.writePacket(&packet.NPCDialogue{
		EntityUniqueID: h.entityRuntimeID,
		ActionType:     packet.NPCDialogueActionOpen,
		Dialogue:       d.Body(),
		SceneName:      "default",
		NPCName:        d.Title(),
		ActionJSON:     string(b),
	})
}

func (s *Session) CloseDialogue() {
	h := s.handlers[packet.IDNPCRequest].(*NPCRequestHandler)
	if h.entityRuntimeID == 0 {
		return
	}

	s.writePacket(&packet.NPCDialogue{
		EntityUniqueID: h.entityRuntimeID,
		ActionType:     packet.NPCDialogueActionClose,
	})
	h.entityRuntimeID = 0
}

func (s *Session) SendForm(f form.Form) {
	b, _ := json.Marshal(f)

	h := s.handlers[packet.IDModalFormResponse].(*ModalFormResponseHandler)
	id := h.currentID.Add(1)

	h.mu.Lock()
	if len(h.forms) > 10 {
		s.conf.Log.Debug("SendForm: more than 10 active forms: dropping an existing one")
		for k := range h.forms {
			delete(h.forms, k)
			break
		}
	}
	h.forms[id] = f
	h.mu.Unlock()

	s.writePacket(&packet.ModalFormRequest{
		FormID:   id,
		FormData: b,
	})
}

func (s *Session) CloseForm() {
	s.writePacket(&packet.ClientBoundCloseForm{})
}

func (s *Session) Transfer(ip net.IP, port int) {
	s.writePacket(&packet.Transfer{
		Address: ip.String(),
		Port:    uint16(port),
	})
}

func (s *Session) SendGameMode(c Controllable) {
	if s == Nop {
		return
	}
	s.writePacket(&packet.SetPlayerGameType{GameType: gameTypeFromMode(c.GameMode())})
	s.SendAbilities(c)
	if !c.GameMode().HasCollision() {

		s.ViewEntityTeleport(c, c.Position())
		c.StartFlying()
	}
}

func (s *Session) SendAbilities(c Controllable) {
	mode, abilities := c.GameMode(), uint32(0)
	if mode.AllowsFlying() {
		abilities |= protocol.AbilityMayFly
		if c.Flying() || !mode.HasCollision() {
			abilities |= protocol.AbilityFlying
		}
	}
	if !mode.HasCollision() {
		abilities |= protocol.AbilityNoClip
	}
	if !mode.AllowsTakingDamage() {
		abilities |= protocol.AbilityInvulnerable
	}
	if mode.CreativeInventory() {
		abilities |= protocol.AbilityInstantBuild
	}
	if mode.AllowsEditing() {
		abilities |= protocol.AbilityBuild | protocol.AbilityMine
	}
	if mode.AllowsInteraction() {
		abilities |= protocol.AbilityDoorsAndSwitches | protocol.AbilityOpenContainers | protocol.AbilityAttackPlayers | protocol.AbilityAttackMobs
	}
	s.writePacket(&packet.UpdateAbilities{AbilityData: protocol.AbilityData{
		EntityUniqueID:     selfEntityRuntimeID,
		PlayerPermissions:  packet.PermissionLevelMember,
		CommandPermissions: protocol.CommandPermissionLevelAny,
		Layers: []protocol.AbilityLayer{
			{
				Type:             protocol.AbilityLayerTypeBase,
				Abilities:        protocol.AbilityCount - 1,
				Values:           abilities,
				FlySpeed:         float32(c.FlightSpeed()),
				VerticalFlySpeed: float32(c.VerticalFlightSpeed()),
				WalkSpeed:        protocol.AbilityBaseWalkSpeed,
			},
		},
	}})
}

func (s *Session) SendHealth(health, max, absorption float64) {
	s.writePacket(&packet.UpdateAttributes{
		EntityRuntimeID: selfEntityRuntimeID,
		Attributes: []protocol.Attribute{{
			AttributeValue: protocol.AttributeValue{
				Name:  "minecraft:health",
				Value: float32(math.Ceil(health)),
				Max:   float32(math.Ceil(max)),
			},
			DefaultMax: 20,
			Default:    20,
		}, {
			AttributeValue: protocol.AttributeValue{
				Name:  "minecraft:absorption",
				Value: float32(math.Ceil(absorption)),
				Max:   float32(math.MaxFloat32),
			},
			DefaultMax: float32(math.MaxFloat32),
		}},
	})
}

func (s *Session) SendEffect(e effect.Effect) {
	s.SendEffectRemoval(e.Type())
	id, _ := effect.ID(e.Type())
	dur := e.Duration() / (time.Second / 20)
	if e.Infinite() {
		dur = -1
	}
	s.writePacket(&packet.MobEffect{
		EntityRuntimeID: selfEntityRuntimeID,
		Operation:       packet.MobEffectAdd,
		EffectType:      int32(id),
		Amplifier:       int32(e.Level() - 1),
		Particles:       !e.ParticlesHidden(),
		Duration:        int32(dur),
		Ambient:         e.Ambient(),
	})
}

func (s *Session) SendEffectRemoval(e effect.Type) {
	id, ok := effect.ID(e)
	if !ok {
		panic(fmt.Sprintf("unregistered effect type %T", e))
	}
	s.writePacket(&packet.MobEffect{
		EntityRuntimeID: selfEntityRuntimeID,
		Operation:       packet.MobEffectRemove,
		EffectType:      int32(id),
	})
}

func (s *Session) sendGameRules(gameRules []protocol.GameRule) {
	s.writePacket(&packet.GameRulesChanged{GameRules: gameRules})
}

func (s *Session) EnableCoordinates(enable bool) {

	s.sendGameRules([]protocol.GameRule{{Name: "showcoordinates", Value: enable}})
}

func (s *Session) EnableInstantRespawn(enable bool) {

	s.sendGameRules([]protocol.GameRule{{Name: "doimmediaterespawn", Value: enable}})
}

func (s *Session) HandleInventories(tx *world.Tx, c Controllable, inv, offHand, enderChest, ui *inventory.Inventory, armour *inventory.Armour, heldSlot *uint32) {
	s.inv = inv
	s.inv.SlotFunc(s.broadcastInvFunc(tx, c))
	s.offHand = offHand
	s.offHand.SlotFunc(s.broadcastOffHandFunc(tx, c))
	s.enderChest = enderChest
	s.enderChest.SlotFunc(s.broadcastEnderChestFunc(tx, c))
	s.armour = armour
	s.armour.Inventory().SlotFunc(s.broadcastArmourFunc(tx, c))
	s.ui = ui
	s.ui.SlotFunc(s.uiInventoryFunc(tx, c))
	s.heldSlot = heldSlot
}

func (s *Session) broadcastInvFunc(tx *world.Tx, c Controllable) inventory.SlotFunc {
	return func(slot int, _, after item.Stack) {
		if slot == int(*s.heldSlot) {
			for _, viewer := range tx.Viewers(c.Position()) {
				viewer.ViewEntityItems(c)
			}
		}
		if !s.inTransaction.Load() {
			s.sendItem(after, slot, protocol.WindowIDInventory)
		}
	}
}

func (s *Session) broadcastEnderChestFunc(tx *world.Tx, _ Controllable) inventory.SlotFunc {
	return func(slot int, _, after item.Stack) {
		if !s.inTransaction.Load() {
			if _, ok := tx.Block(*s.openedPos.Load()).(block.EnderChest); ok {
				s.ViewSlotChange(slot, after)
			}
		}
	}
}

func (s *Session) broadcastOffHandFunc(tx *world.Tx, c Controllable) inventory.SlotFunc {
	return func(slot int, _, after item.Stack) {
		for _, viewer := range tx.Viewers(c.Position()) {
			viewer.ViewEntityItems(c)
		}
		if !s.inTransaction.Load() {
			i, _ := s.offHand.Item(0)
			s.writePacket(&packet.InventoryContent{
				WindowID: protocol.WindowIDOffHand,
				Content:  []protocol.ItemInstance{instanceFromItem(s.br, i)},
			})
		}
	}
}

func (s *Session) broadcastArmourFunc(tx *world.Tx, c Controllable) inventory.SlotFunc {
	return func(slot int, before, after item.Stack) {
		inTransaction := s.inTransaction.Load()
		if !inTransaction {
			s.sendItem(after, slot, protocol.WindowIDArmour)
		}
		if before.Comparable(after) && before.Empty() == after.Empty() {

			return
		}
		for _, viewer := range tx.Viewers(c.Position()) {
			viewer.ViewEntityArmour(c)
		}

		if !after.Empty() && inTransaction {
			tx.PlaySound(entity.EyePosition(c), sound.EquipItem{Item: after.Item()})
		}
	}
}

func (s *Session) uiInventoryFunc(tx *world.Tx, c Controllable) inventory.SlotFunc {
	return func(slot int, _, after item.Stack) {
		if slot == enchantingInputSlot && s.containerOpened.Load() {
			pos := *s.openedPos.Load()
			if _, enchanting := tx.Block(pos).(block.EnchantingTable); enchanting {
				s.sendEnchantmentOptions(tx, c, pos, after)
			}
		}
		s.sendInv(s.ui, protocol.WindowIDUI)
	}
}

func (s *Session) SendHeldSlot(slot int, c Controllable, force bool) {
	if s.changingSlot.Load() && !force {
		return
	}
	mainHand, _ := c.HeldItems()
	s.writePacket(&packet.MobEquipment{
		EntityRuntimeID: selfEntityRuntimeID,
		NewItem:         instanceFromItem(s.br, mainHand),
		InventorySlot:   byte(slot),
		HotBarSlot:      byte(slot),
	})
}

func (s *Session) VerifyAndSetHeldSlot(slot int, expected item.Stack, c Controllable) error {
	if err := s.VerifySlot(slot, expected); err != nil {
		return err
	}
	s.changingSlot.Store(true)
	defer s.changingSlot.Store(false)
	return c.SetHeldSlot(slot)
}

func (s *Session) VerifySlot(slot int, expected item.Stack) error {

	if slot < 0 || slot > 8 {
		return fmt.Errorf("slot exceeds hotbar range 0-8: slot is %v", slot)
	}
	clientSideItem := expected
	actual, _ := s.inv.Item(slot)

	if !clientSideItem.Equal(actual) {
		s.sendItem(actual, slot, protocol.WindowIDInventory)

		s.conf.Log.Debug("verify slot: client-side item was not equal to server-side item", "client-held", clientSideItem.String(), "server-held", actual.String())
	}
	return nil
}

func (s *Session) SendExperience(level int, progress float64) {
	s.writePacket(&packet.UpdateAttributes{
		EntityRuntimeID: selfEntityRuntimeID,
		Attributes: []protocol.Attribute{
			{
				AttributeValue: protocol.AttributeValue{
					Name:  "minecraft:player.level",
					Value: float32(level),
					Max:   float32(math.MaxInt32),
				},
				DefaultMax: float32(math.MaxInt32),
			},
			{
				AttributeValue: protocol.AttributeValue{
					Name:  "minecraft:player.experience",
					Value: float32(progress),
					Max:   1,
				},
				DefaultMax: 1,
			},
		},
	})
}

func (s *Session) SendChargeItemComplete() {
	s.writePacket(&packet.ActorEvent{
		EntityRuntimeID: selfEntityRuntimeID,
		EventType:       packet.ActorEventFinishedChargingItem,
	})
}

func (s *Session) ShowHudElement(e hud.Element) {
	s.hudMu.Lock()
	defer s.hudMu.Unlock()

	if _, ok := s.hiddenHud[e]; ok {
		s.hudUpdates[e] = true
	} else if _, ok = s.hudUpdates[e]; ok {
		delete(s.hudUpdates, e)
	}
}

func (s *Session) HideHudElement(e hud.Element) {
	s.hudMu.Lock()
	defer s.hudMu.Unlock()

	if _, ok := s.hiddenHud[e]; !ok {
		s.hudUpdates[e] = false
	} else if _, ok = s.hudUpdates[e]; ok {
		delete(s.hudUpdates, e)
	}
}

func (s *Session) HudElementHidden(e hud.Element) bool {
	s.hudMu.RLock()
	defer s.hudMu.RUnlock()

	if _, ok := s.hiddenHud[e]; ok {
		return true
	}
	vis, ok := s.hudUpdates[e]
	return ok && !vis
}

func (s *Session) SendHudUpdates() {
	s.hudMu.Lock()
	if len(s.hudUpdates) == 0 {
		s.hudMu.Unlock()
		return
	}
	var show, hide []int32
	for e, vis := range s.hudUpdates {
		if vis {
			show = append(show, int32(e.Uint8()))
			delete(s.hiddenHud, e)
		} else {
			hide = append(hide, int32(e.Uint8()))
			s.hiddenHud[e] = struct{}{}
		}
	}
	s.hudUpdates = make(map[hud.Element]bool)
	s.hudMu.Unlock()

	if len(show) > 0 {
		s.writePacket(&packet.SetHud{Elements: show, Visibility: packet.HudVisibilityReset})
	}
	if len(hide) > 0 {
		s.writePacket(&packet.SetHud{Elements: hide, Visibility: packet.HudVisibilityHide})
	}
}

func (s *Session) AddDebugShape(shape debug.Shape) {
	if s == Nop {
		return
	}
	s.queueDebugShapeUpdate(debugShapeUpdate{id: shape.ShapeID(), shape: shape})
}

func (s *Session) RemoveDebugShape(shape debug.Shape) {
	if s == Nop {
		return
	}
	s.queueDebugShapeUpdate(debugShapeUpdate{id: shape.ShapeID()})
}

func (s *Session) VisibleDebugShapes() []debug.Shape {
	s.debugShapesMu.RLock()
	defer s.debugShapesMu.RUnlock()

	return slices.Collect(maps.Values(s.debugShapes))
}

func (s *Session) RemoveAllDebugShapes() {
	if s == Nop {
		return
	}
	s.debugShapesMu.Lock()
	defer s.debugShapesMu.Unlock()

	s.debugShapeUpdates = s.debugShapeUpdates[:0]
	for id := range s.debugShapes {
		s.debugShapeUpdates = append(s.debugShapeUpdates, debugShapeUpdate{id: id})
	}
}

func (s *Session) SendDebugShapes(dim world.Dimension) {
	s.debugShapesMu.Lock()
	updates := s.debugShapeUpdates
	if len(updates) == 0 {
		s.debugShapesMu.Unlock()
		return
	}

	shapes := make([]protocol.PrimitiveShape, 0, len(updates))
	for _, update := range updates {
		if update.shape == nil {
			delete(s.debugShapes, update.id)
			shapes = append(shapes, protocol.PrimitiveShape{
				NetworkID:      uint64(update.id),
				DimensionID:    protocol.Option(s.dimensionID(dim)),
				ExtraShapeData: &protocol.LastShape{},
			})
			continue
		}
		s.debugShapes[update.id] = update.shape
		shapes = append(shapes, debugShapeToProtocol(update.shape, dim, s.shapeAttachedEntityRuntimeID(update.shape)))
	}
	s.debugShapeUpdates = s.debugShapeUpdates[:0]
	s.debugShapesMu.Unlock()

	s.writePacket(&packet.PrimitiveShapes{Shapes: shapes})
}

func (s *Session) LockInput(l input.Lock) {
	s.inputLocksMu.Lock()
	defer s.inputLocksMu.Unlock()
	s.inputLocks |= l.Uint32()
}

func (s *Session) UnlockInput(l input.Lock) {
	s.inputLocksMu.Lock()
	defer s.inputLocksMu.Unlock()
	s.inputLocks &^= l.Uint32()
}

func (s *Session) ClearInputLocks() {
	s.inputLocksMu.Lock()
	defer s.inputLocksMu.Unlock()
	s.inputLocks = 0
}

func (s *Session) InputLocked(l input.Lock) bool {
	s.inputLocksMu.RLock()
	defer s.inputLocksMu.RUnlock()
	return s.inputLocks&l.Uint32() != 0
}

func (s *Session) SendInputLocks() {
	s.inputLocksMu.RLock()
	defer s.inputLocksMu.RUnlock()
	s.writePacket(&packet.UpdateClientInputLocks{
		Locks: s.inputLocks,
	})
}

func (s *Session) queueDebugShapeUpdate(update debugShapeUpdate) {
	s.debugShapesMu.Lock()
	defer s.debugShapesMu.Unlock()
	s.debugShapeUpdates = append(s.debugShapeUpdates, update)
}

func valueOrDefault[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}

func stackFromItem(br world.BlockRegistry, it item.Stack) protocol.ItemStack {
	if it.Empty() {
		return protocol.ItemStack{}
	}

	var blockRuntimeID uint32
	if b, ok := it.Item().(world.Block); ok {
		blockRuntimeID = br.BlockRuntimeID(b)
	}

	rid, meta, _ := world.ItemRuntimeID(it.Item())

	return protocol.ItemStack{
		ItemType: protocol.ItemType{
			NetworkID:     rid,
			MetadataValue: uint32(meta),
		},
		Count:          uint16(it.Count()),
		BlockRuntimeID: int32(blockRuntimeID),
		NBTData:        item.WriteNBT(it, false),
	}
}

func stackToItem(br world.BlockRegistry, it protocol.ItemStack) item.Stack {
	t, ok := world.ItemByRuntimeID(it.NetworkID, int16(it.MetadataValue))
	if !ok {
		t = block.Air{}
	}
	if it.BlockRuntimeID > 0 {

		b, _ := br.BlockByRuntimeID(uint32(it.BlockRuntimeID))
		if t, ok = b.(world.Item); !ok {
			t = block.Air{}
		}
	}

	if nbter, ok := t.(world.NBTer); ok && len(it.NBTData) != 0 {
		t = nbter.DecodeNBT(it.NBTData).(world.Item)
	}
	s := item.NewStack(t, int(it.Count))
	return item.ReadNBT(it.NBTData, &s)
}

func instanceFromItem(br world.BlockRegistry, it item.Stack) protocol.ItemInstance {
	return protocol.ItemInstance{
		StackNetworkID: item_id(it),
		Stack:          stackFromItem(br, it),
	}
}

func stacksToRecipeStacks(br world.BlockRegistry, inputs []item.Stack) []protocol.ItemStack {
	items := make([]protocol.ItemStack, 0, len(inputs))
	for _, i := range inputs {
		items = append(items, deleteDamage(stackFromItem(br, i)))
	}
	return items
}

func stacksToIngredientItems(_ world.BlockRegistry, inputs []recipe.Item) []protocol.ItemDescriptorCount {
	items := make([]protocol.ItemDescriptorCount, 0, len(inputs))
	for _, i := range inputs {
		var d protocol.ItemDescriptor = &protocol.InvalidItemDescriptor{}
		switch i := i.(type) {
		case item.Stack:
			if i.Empty() {
				items = append(items, protocol.ItemDescriptorCount{Descriptor: &protocol.InvalidItemDescriptor{}})
				continue
			}
			name, meta := i.Item().EncodeItem()
			if _, ok := i.Value("variants"); ok {
				meta = math.MaxInt16
			}
			d = &protocol.DefaultItemDescriptor{
				Name:          name,
				MetadataValue: int32(meta),
			}
		case recipe.ItemTag:
			d = &protocol.ItemTagItemDescriptor{Tag: i.Tag()}
		}
		items = append(items, protocol.ItemDescriptorCount{
			Descriptor: d,
			Count:      int32(i.Count()),
		})
	}
	return items
}

func creativeContent(br world.BlockRegistry) ([]protocol.CreativeGroup, []protocol.CreativeItem) {
	groups := make([]protocol.CreativeGroup, 0, len(creative.Groups()))
	for _, group := range creative.Groups() {
		groups = append(groups, protocol.CreativeGroup{
			Category: group.Category.Uint8(),
			Name:     group.Name,
			Icon:     deleteDamage(stackFromItem(br, group.Icon)),
		})
	}

	it := make([]protocol.CreativeItem, 0, len(creative.Items()))
	for index, i := range creative.Items() {
		group := slices.IndexFunc(creative.Groups(), func(group creative.Group) bool {
			return group.Name == i.Group
		})
		if group < 0 {
			continue
		}
		it = append(it, protocol.CreativeItem{
			CreativeItemNetworkID: uint32(index) + 1,
			Item:                  deleteDamage(stackFromItem(br, i.Stack)),
			GroupIndex:            uint32(group),
		})
	}
	return groups, it
}

func deleteDamage(st protocol.ItemStack) protocol.ItemStack {
	delete(st.NBTData, "Damage")
	return st
}

func protocolToSkin(sk protocol.Skin) (s skin.Skin, err error) {
	if sk.SkinID == "" {
		return skin.Skin{}, fmt.Errorf("SkinID must not be an empty string")
	}

	s = skin.New(int(sk.SkinImageWidth), int(sk.SkinImageHeight))
	s.Persona = sk.PersonaSkin
	s.Pix = sk.SkinData
	s.Model = sk.SkinGeometry
	s.PlayFabID = sk.PlayFabID
	s.FullID = sk.FullID

	s.Cape = skin.NewCape(int(sk.CapeImageWidth), int(sk.CapeImageHeight))
	s.Cape.Pix = sk.CapeData

	m := make(map[string]any)
	if err = json.Unmarshal(sk.SkinGeometry, &m); err != nil {
		return skin.Skin{}, fmt.Errorf("SkinGeometry was not a valid JSON string: %v", err)
	}

	if s.ModelConfig, err = skin.DecodeModelConfig(sk.SkinResourcePatch); err != nil {
		return skin.Skin{}, fmt.Errorf("SkinResourcePatch was not a valid JSON string: %v", err)
	}

	for _, anim := range sk.Animations {
		var t skin.AnimationType
		switch anim.AnimationType {
		case protocol.SkinAnimationHead:
			t = skin.AnimationHead
		case protocol.SkinAnimationBody32x32:
			t = skin.AnimationBody32x32
		case protocol.SkinAnimationBody128x128:
			t = skin.AnimationBody128x128
		default:
			return skin.Skin{}, fmt.Errorf("invalid animation type: %v", anim.AnimationType)
		}

		animation := skin.NewAnimation(int(anim.ImageWidth), int(anim.ImageHeight), int(anim.ExpressionType), t)
		animation.FrameCount = int(anim.FrameCount)
		animation.Pix = anim.ImageData

		s.Animations = append(s.Animations, animation)
	}
	return
}

func (s *Session) shapeAttachedEntityRuntimeID(shape debug.Shape) uint64 {
	var handle *world.EntityHandle
	switch shape := shape.(type) {
	case *debug.Arrow:
		handle = shape.Entity
	case *debug.Box:
		handle = shape.Entity
	case *debug.Circle:
		handle = shape.Entity
	case *debug.Line:
		handle = shape.Entity
	case *debug.Sphere:
		handle = shape.Entity
	case *debug.Text:
		handle = shape.Entity
	case *debug.Cylinder:
		handle = shape.Entity
	case *debug.Pyramid:
		handle = shape.Entity
	case *debug.Ellipsoid:
		handle = shape.Entity
	case *debug.Cone:
		handle = shape.Entity
	}
	if handle == nil {
		return 0
	}
	return s.handleRuntimeID(handle)
}

func debugShapeToProtocol(shape debug.Shape, dim world.Dimension, attachedEntityID uint64) protocol.PrimitiveShape {
	dimID, _ := world.DimensionID(dim)
	ps := protocol.PrimitiveShape{
		NetworkID:   uint64(shape.ShapeID()),
		DimensionID: protocol.Option(int32(dimID)),
	}
	if attachedEntityID > 0 {
		ps.AttachedToEntityID = protocol.Option(int64(attachedEntityID))
	}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	switch shape := shape.(type) {
	case *debug.Arrow:
		ps.Type = protocol.Option(protocol.PrimitiveShapeArrow)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.ExtraShapeData = &protocol.ArrowShape{
			ArrowEndLocation: protocol.Option(vec64To32(shape.EndPosition)),
			ArrowHeadLength:  protocol.Option(valueOrDefault(float32(shape.HeadLength), 1)),
			ArrowHeadRadius:  protocol.Option(valueOrDefault(float32(shape.HeadRadius), 0.5)),
			Segments:         protocol.Option(valueOrDefault(uint8(shape.HeadSegments), 4)),
		}
	case *debug.Box:
		ps.Type = protocol.Option(protocol.PrimitiveShapeBox)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		ps.ExtraShapeData = &protocol.BoxShape{BoxBound: valueOrDefault(vec64To32(shape.Bounds), mgl32.Vec3{1, 1, 1})}
	case *debug.Circle:
		ps.Type = protocol.Option(protocol.PrimitiveShapeCircle)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		ps.ExtraShapeData = &protocol.SphereShape{Segments: valueOrDefault(uint8(shape.Segments), 20)}
	case *debug.Line:
		ps.Type = protocol.Option(protocol.PrimitiveShapeLine)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.ExtraShapeData = &protocol.LineShape{LineEndLocation: vec64To32(shape.EndPosition)}
	case *debug.Sphere:
		ps.Type = protocol.Option(protocol.PrimitiveShapeSphere)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		ps.ExtraShapeData = &protocol.SphereShape{Segments: valueOrDefault(uint8(shape.Segments), 20)}
	case *debug.Text:
		ps.Type = protocol.Option(protocol.PrimitiveShapeText)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		if shape.LockRotation {
			ps.Rotation = protocol.Option(vec64To32(shape.Rotation))
		}
		textData := &protocol.TextShape{
			Text:             shape.Text,
			UseRotation:      shape.LockRotation,
			DepthTest:        !shape.DisableDepthTest,
			ShowBackface:     !shape.HideBackface,
			ShowBackfaceText: !shape.HideBackfaceText,
		}
		switch {
		case shape.HideBackground:
			textData.BackgroundColour = protocol.Option(color.RGBA{})
		case shape.BackgroundColour != (color.RGBA{}):
			textData.BackgroundColour = protocol.Option(shape.BackgroundColour)
		}
		ps.ExtraShapeData = textData
	case *debug.Cylinder:
		ps.Type = protocol.Option(protocol.PrimitiveShapeCylinder)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		base := valueOrDefault(shape.BaseRadius, mgl64.Vec2{1, 1})
		top := valueOrDefault(shape.TopRadius, base)
		ps.ExtraShapeData = &protocol.CylinderShape{
			RadiusX:     mgl32.Vec2{float32(base[0]), float32(top[0])},
			RadiusZ:     mgl32.Vec2{float32(base[1]), float32(top[1])},
			Height:      valueOrDefault(float32(shape.Height), 1),
			NumSegments: valueOrDefault(uint8(shape.Segments), 20),
		}
	case *debug.Pyramid:
		ps.Type = protocol.Option(protocol.PrimitiveShapePyramid)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		pyramid := &protocol.PyramidShape{
			Width:  valueOrDefault(float32(shape.Width), 1),
			Height: valueOrDefault(float32(shape.Height), 1),
		}
		if shape.Depth != 0 {
			pyramid.Depth = protocol.Option(float32(shape.Depth))
		}
		ps.ExtraShapeData = pyramid
	case *debug.Ellipsoid:
		ps.Type = protocol.Option(protocol.PrimitiveShapeEllipsoid)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		ps.ExtraShapeData = &protocol.EllipsoidShape{
			Radii:           valueOrDefault(vec64To32(shape.Radii), mgl32.Vec3{1, 1, 1}),
			SegmentsPerAxis: valueOrDefault(uint8(shape.SegmentsPerAxis), 20),
		}
	case *debug.Cone:
		ps.Type = protocol.Option(protocol.PrimitiveShapeCone)
		ps.Colour = protocol.Option(valueOrDefault(shape.Colour, white))
		ps.Location = protocol.Option(vec64To32(shape.Position))
		ps.Scale = protocol.Option(valueOrDefault(float32(shape.Scale), 1))
		ps.ExtraShapeData = &protocol.ConeShape{
			Radii:       valueOrDefault(vec2To32(shape.Radii), mgl32.Vec2{1, 1}),
			Height:      valueOrDefault(float32(shape.Height), 1),
			NumSegments: valueOrDefault(uint8(shape.Segments), 20),
		}
	default:
		panic(fmt.Sprintf("unknown debug shape type %T", shape))
	}
	return ps
}

func gameTypeFromMode(mode world.GameMode) int32 {
	if !mode.HasCollision() {

		return packet.GameTypeSpectator
	}
	if mode.AllowsFlying() && mode.CreativeInventory() {
		return packet.GameTypeCreative
	}
	return packet.GameTypeSurvival
}

//go:linkname item_id github.com/df-mc/dragonfly/server/item.id
func item_id(s item.Stack) int32
