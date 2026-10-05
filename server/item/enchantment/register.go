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

package enchantment

import "github.com/df-mc/dragonfly/server/item"

func init() {
	item.RegisterEnchantment(0, Protection)
	item.RegisterEnchantment(1, FireProtection)
	item.RegisterEnchantment(2, FeatherFalling)
	item.RegisterEnchantment(3, BlastProtection)
	item.RegisterEnchantment(4, ProjectileProtection)
	item.RegisterEnchantment(5, Thorns)
	item.RegisterEnchantment(6, Respiration)
	item.RegisterEnchantment(7, DepthStrider)
	item.RegisterEnchantment(8, AquaAffinity)
	item.RegisterEnchantment(9, Sharpness)

	item.RegisterEnchantment(12, Knockback)
	item.RegisterEnchantment(13, FireAspect)

	item.RegisterEnchantment(15, Efficiency)
	item.RegisterEnchantment(16, SilkTouch)
	item.RegisterEnchantment(17, Unbreaking)
	item.RegisterEnchantment(18, Fortune)
	item.RegisterEnchantment(19, Power)
	item.RegisterEnchantment(20, Punch)
	item.RegisterEnchantment(21, Flame)
	item.RegisterEnchantment(22, Infinity)

	item.RegisterEnchantment(26, Mending)

	item.RegisterEnchantment(28, CurseOfVanishing)

	item.RegisterEnchantment(33, Multishot)
	item.RegisterEnchantment(34, Piercing)
	item.RegisterEnchantment(35, QuickCharge)
	item.RegisterEnchantment(36, SoulSpeed)
	item.RegisterEnchantment(37, SwiftSneak)
}
