package item_test

import (
	"testing"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
)

func TestWithItemKeepsForcedEnchantments(t *testing.T) {
	original := item.NewStack(item.Stick{}, 1).WithForcedEnchantments(item.NewEnchantment(enchantment.Sharpness, 3))
	updated := original.WithItem(item.Stick{})
	if sharpness, ok := updated.Enchantment(enchantment.Sharpness); !ok || sharpness.Level() != 3 {
		t.Fatalf("got sharpness %v, present %v", sharpness.Level(), ok)
	}
}
