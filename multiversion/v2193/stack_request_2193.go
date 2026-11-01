package v2193

import (
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type Reader2193 struct{ *protocol.Reader }
type Writer2193 struct{ *protocol.Writer }

func (r Reader2193) StackRequestAction(x *protocol.StackRequestAction) { stackAction2193(r, x) }
func (w Writer2193) StackRequestAction(x *protocol.StackRequestAction) { stackAction2193(w, x) }

func legacyActionIndex2193(x protocol.StackRequestAction) (byte, bool) {
	switch x.(type) {
	case *protocol.TakeStackRequestAction:
		return 0, true
	case *protocol.PlaceStackRequestAction:
		return 1, true
	case *protocol.SwapStackRequestAction:
		return 2, true
	case *protocol.DropStackRequestAction:
		return 3, true
	case *protocol.DestroyStackRequestAction:
		return 4, true
	case *protocol.ConsumeStackRequestAction:
		return 5, true
	case *protocol.CreateStackRequestAction:
		return 6, true
	case *protocol.LabTableCombineStackRequestAction:
		return 7, true
	case *protocol.BeaconPaymentStackRequestAction:
		return 8, true
	case *protocol.MineBlockStackRequestAction:
		return 9, true
	case *protocol.CraftRecipeStackRequestAction:
		return 10, true
	case *protocol.AutoCraftRecipeStackRequestAction:
		return 11, true
	case *protocol.CraftCreativeStackRequestAction:
		return 12, true
	case *protocol.CraftRecipeOptionalStackRequestAction:
		return 13, true
	case *protocol.CraftGrindstoneRecipeStackRequestAction:
		return 14, true
	case *protocol.CraftLoomRecipeStackRequestAction:
		return 15, true
	case *protocol.CraftNonImplementedStackRequestAction:
		return 16, true
	case *protocol.CraftResultsDeprecatedStackRequestAction:
		return 17, true
	default:
		return 0, false
	}
}

func makeAction2193(index byte) (protocol.StackRequestAction, bool) {
	switch index {
	case 0:
		return &protocol.TakeStackRequestAction{}, true
	case 1:
		return &protocol.PlaceStackRequestAction{}, true
	case 2:
		return &protocol.SwapStackRequestAction{}, true
	case 3:
		return &protocol.DropStackRequestAction{}, true
	case 4:
		return &protocol.DestroyStackRequestAction{}, true
	case 5:
		return &protocol.ConsumeStackRequestAction{}, true
	case 6:
		return &protocol.CreateStackRequestAction{}, true
	case 7:
		return &protocol.LabTableCombineStackRequestAction{}, true
	case 8:
		return &protocol.BeaconPaymentStackRequestAction{}, true
	case 9:
		return &protocol.MineBlockStackRequestAction{}, true
	case 10:
		return &protocol.CraftRecipeStackRequestAction{}, true
	case 11:
		return &protocol.AutoCraftRecipeStackRequestAction{}, true
	case 12:
		return &protocol.CraftCreativeStackRequestAction{}, true
	case 13:
		return &protocol.CraftRecipeOptionalStackRequestAction{}, true
	case 14:
		return &protocol.CraftGrindstoneRecipeStackRequestAction{}, true
	case 15:
		return &protocol.CraftLoomRecipeStackRequestAction{}, true
	case 16:
		return &protocol.CraftNonImplementedStackRequestAction{}, true
	case 17:
		return &protocol.CraftResultsDeprecatedStackRequestAction{}, true
	default:
		return nil, false
	}
}

func fullContainer2193(io protocol.IO, c *protocol.FullContainerName) {
	io.Uint8(&c.ContainerID)
	dynamic, _ := c.DynamicContainerID.Value()
	io.Uint32(&dynamic)
	if dynamic != 0 {
		c.DynamicContainerID = protocol.Option(dynamic)
	} else {
		c.DynamicContainerID = protocol.Optional[uint32]{}
	}
}

func stackSlot2193(io protocol.IO, x *protocol.StackRequestSlotInfo) {
	fullContainer2193(io, &x.Container)
	io.Uint8(&x.Slot)
	io.Int32(&x.StackNetworkID)
}

func stackAction2193(io protocol.IO, x *protocol.StackRequestAction) {
	var id uint8
	if *x != nil {
		var ok bool
		id, ok = legacyActionIndex2193(*x)
		if !ok {
			io.InvalidValue(fmt.Sprintf("%T", *x), "2193 stack action", "unsupported action type")
			return
		}
	}
	io.Uint8(&id)
	if *x == nil {
		a, ok := makeAction2193(id)
		if !ok {
			io.UnknownEnumOption(id, "2193 stack action")
			return
		}
		*x = a
	}
	switch a := (*x).(type) {
	case *protocol.TakeStackRequestAction:
		io.Uint8(&a.Count)
		stackSlot2193(io, &a.Source)
		stackSlot2193(io, &a.Destination)
	case *protocol.PlaceStackRequestAction:
		io.Uint8(&a.Count)
		stackSlot2193(io, &a.Source)
		stackSlot2193(io, &a.Destination)
	case *protocol.SwapStackRequestAction:
		stackSlot2193(io, &a.Source)
		stackSlot2193(io, &a.Destination)
	case *protocol.DropStackRequestAction:
		io.Uint8(&a.Count)
		stackSlot2193(io, &a.Source)
		io.Bool(&a.Randomly)
	case *protocol.DestroyStackRequestAction:
		io.Uint8(&a.Count)
		stackSlot2193(io, &a.Source)
	case *protocol.ConsumeStackRequestAction:
		io.Uint8(&a.Count)
		stackSlot2193(io, &a.Source)
	default:
		(*x).Marshal(io)
	}
}
