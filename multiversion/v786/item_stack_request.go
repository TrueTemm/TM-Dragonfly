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
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ItemStackRequest struct {
	packet.ItemStackRequest
}

func (*ItemStackRequest) ID() uint32 { return packet.IDItemStackRequest }

func (pk *ItemStackRequest) Marshal(io protocol.IO) {
	protocol.FuncSlice(io, &pk.Requests, func(r *protocol.ItemStackRequest) {
		marshalItemStackRequest786(io, r)
	})
}

func marshalItemStackRequest786(io protocol.IO, pk *protocol.ItemStackRequest) {
	io.Varint32(&pk.RequestID)
	protocol.FuncSlice(io, &pk.Actions, func(x *protocol.StackRequestAction) {
		stackReqAction786(io, x)
	})
	protocol.FuncSlice(io, &pk.FilterStrings, io.String)
	io.Int32(&pk.FilterCause)
}

func stackReqSlot786(r protocol.IO, x *protocol.StackRequestSlotInfo) {

	if p := ProtoOf(r); p != 0 && p < 712 {
		r.Uint8(&x.Container.ContainerID)
		r.Uint8(&x.Slot)
		r.Varint32(&x.StackNetworkID)
		return
	}

	if p := ProtoOf(r); p != 0 && p < 729 {
		r.Uint8(&x.Container.ContainerID)
		id, _ := x.Container.DynamicContainerID.Value()
		r.Uint32(&id)
		x.Container.DynamicContainerID = protocol.Option(id)
	} else {
		protocol.Single(r, &x.Container)
	}
	r.Uint8(&x.Slot)
	r.Varint32(&x.StackNetworkID)
}

func stackReqAction786(r protocol.IO, x *protocol.StackRequestAction) {
	var id byte
	if *x != nil {
		id = stackReqActionID786(*x)
	}
	r.Uint8(&id)
	if *x == nil {
		if !lookupStackReqAction786(id, x) {
			r.UnknownEnumOption(id, "v786 stack request action type")
			return
		}
	}
	switch a := (*x).(type) {
	case *protocol.TakeStackRequestAction:
		r.Uint8(&a.Count)
		stackReqSlot786(r, &a.Source)
		stackReqSlot786(r, &a.Destination)
	case *protocol.PlaceStackRequestAction:
		r.Uint8(&a.Count)
		stackReqSlot786(r, &a.Source)
		stackReqSlot786(r, &a.Destination)
	case *protocol.SwapStackRequestAction:
		stackReqSlot786(r, &a.Source)
		stackReqSlot786(r, &a.Destination)
	case *protocol.DropStackRequestAction:
		r.Uint8(&a.Count)
		stackReqSlot786(r, &a.Source)
		r.Bool(&a.Randomly)
	case *protocol.DestroyStackRequestAction:
		r.Uint8(&a.Count)
		stackReqSlot786(r, &a.Source)
	case *protocol.ConsumeStackRequestAction:
		r.Uint8(&a.Count)
		stackReqSlot786(r, &a.Source)
	case *protocol.MineBlockStackRequestAction:

		r.Varint32(&a.HotbarSlot)
		r.Varint32(&a.PredictedDurability)
		r.Varint32(&a.StackNetworkID)
	case *protocol.CraftRecipeStackRequestAction:

		r.Varuint32(&a.RecipeNetworkID)
		if p := ProtoOf(r); p == 0 || p >= 712 {
			r.Uint8(&a.NumberOfCrafts)
		}
	case *protocol.CraftCreativeStackRequestAction:
		r.Varuint32(&a.CreativeItemNetworkID)
		if p := ProtoOf(r); p == 0 || p >= 712 {
			r.Uint8(&a.NumberOfCrafts)
		}
	case *protocol.AutoCraftRecipeStackRequestAction:

		r.Varuint32(&a.RecipeNetworkID)
		if p := ProtoOf(r); p == 0 || p >= 712 {
			r.Uint8(&a.NumberOfCrafts)
		}
		var timesCrafted byte
		r.Uint8(&timesCrafted)
		protocol.FuncSlice(r, &a.Ingredients, r.ItemDescriptorCount)
	case *protocol.CraftGrindstoneRecipeStackRequestAction:

		r.Varuint32(&a.RecipeNetworkID)
		if p := ProtoOf(r); p == 0 || p >= 712 {
			r.Uint8(&a.NumberOfCrafts)
		}
		r.Varint32(&a.Cost)
	case *protocol.CraftRecipeOptionalStackRequestAction:

		r.Varuint32(&a.RecipeNetworkID)
		if p := ProtoOf(r); p >= 712 && p < 748 {
			var numberOfCrafts byte
			r.Uint8(&numberOfCrafts)
		}
		r.Int32(&a.FilterStringIndex)
	case *protocol.CraftLoomRecipeStackRequestAction:

		r.String(&a.Pattern)
		if p := ProtoOf(r); p == 0 || p >= 766 {
			r.Uint8(&a.TimesCrafted)
		}
	case *protocol.CraftResultsDeprecatedStackRequestAction:

		var results []protocol.ItemStack
		protocol.FuncSlice(r, &results, r.Item)
		r.Uint8(&a.TimesCrafted)
	default:

		(*x).Marshal(r)
	}
}

func stackReqActionID786(x protocol.StackRequestAction) byte {
	switch x.(type) {
	case *protocol.TakeStackRequestAction:
		return protocol.StackRequestActionTake
	case *protocol.PlaceStackRequestAction:
		return protocol.StackRequestActionPlace
	case *protocol.SwapStackRequestAction:
		return protocol.StackRequestActionSwap
	case *protocol.DropStackRequestAction:
		return protocol.StackRequestActionDrop
	case *protocol.DestroyStackRequestAction:
		return protocol.StackRequestActionDestroy
	case *protocol.ConsumeStackRequestAction:
		return protocol.StackRequestActionConsume
	case *protocol.CreateStackRequestAction:
		return protocol.StackRequestActionCreate
	case *protocol.LabTableCombineStackRequestAction:
		return protocol.StackRequestActionLabTableCombine
	case *protocol.BeaconPaymentStackRequestAction:
		return protocol.StackRequestActionBeaconPayment
	case *protocol.MineBlockStackRequestAction:
		return protocol.StackRequestActionMineBlock
	case *protocol.CraftRecipeStackRequestAction:
		return protocol.StackRequestActionCraftRecipe
	case *protocol.AutoCraftRecipeStackRequestAction:
		return protocol.StackRequestActionCraftRecipeAuto
	case *protocol.CraftCreativeStackRequestAction:
		return protocol.StackRequestActionCraftCreative
	case *protocol.CraftRecipeOptionalStackRequestAction:
		return protocol.StackRequestActionCraftRecipeOptional
	case *protocol.CraftGrindstoneRecipeStackRequestAction:
		return protocol.StackRequestActionCraftGrindstone
	case *protocol.CraftLoomRecipeStackRequestAction:
		return protocol.StackRequestActionCraftLoom
	case *protocol.CraftNonImplementedStackRequestAction:
		return protocol.StackRequestActionCraftNonImplementedDeprecated
	case *protocol.CraftResultsDeprecatedStackRequestAction:
		return protocol.StackRequestActionCraftResultsDeprecated
	default:
		return 0
	}
}

func lookupStackReqAction786(id byte, x *protocol.StackRequestAction) bool {
	switch id {
	case protocol.StackRequestActionTake:
		*x = &protocol.TakeStackRequestAction{}
	case protocol.StackRequestActionPlace:
		*x = &protocol.PlaceStackRequestAction{}
	case protocol.StackRequestActionSwap:
		*x = &protocol.SwapStackRequestAction{}
	case protocol.StackRequestActionDrop:
		*x = &protocol.DropStackRequestAction{}
	case protocol.StackRequestActionDestroy:
		*x = &protocol.DestroyStackRequestAction{}
	case protocol.StackRequestActionConsume:
		*x = &protocol.ConsumeStackRequestAction{}
	case protocol.StackRequestActionCreate:
		*x = &protocol.CreateStackRequestAction{}
	case protocol.StackRequestActionLabTableCombine:
		*x = &protocol.LabTableCombineStackRequestAction{}
	case protocol.StackRequestActionBeaconPayment:
		*x = &protocol.BeaconPaymentStackRequestAction{}
	case protocol.StackRequestActionMineBlock:
		*x = &protocol.MineBlockStackRequestAction{}
	case protocol.StackRequestActionCraftRecipe:
		*x = &protocol.CraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeAuto:
		*x = &protocol.AutoCraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftCreative:
		*x = &protocol.CraftCreativeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeOptional:
		*x = &protocol.CraftRecipeOptionalStackRequestAction{}
	case protocol.StackRequestActionCraftGrindstone:
		*x = &protocol.CraftGrindstoneRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftLoom:
		*x = &protocol.CraftLoomRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftNonImplementedDeprecated:
		*x = &protocol.CraftNonImplementedStackRequestAction{}
	case protocol.StackRequestActionCraftResultsDeprecated:
		*x = &protocol.CraftResultsDeprecatedStackRequestAction{}
	default:
		return false
	}
	return true
}

func ToLatestItemStackRequest(pk *ItemStackRequest) *packet.ItemStackRequest {
	out := pk.ItemStackRequest
	return &out
}

func FromLatestItemStackRequest(pk *packet.ItemStackRequest) *ItemStackRequest {
	return &ItemStackRequest{ItemStackRequest: *pk}
}
