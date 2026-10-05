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

package category

type Category struct {
	group    string
	category uint8
}

func Construction() Category {
	return Category{category: 1}
}

func Nature() Category {
	return Category{category: 2}
}

func Equipment() Category {
	return Category{category: 3}
}

func Items() Category {
	return Category{category: 4}
}

func (c Category) Uint8() uint8 {
	return c.category
}

func (c Category) WithGroup(group string) Category {
	c.group = group
	return c
}

func (c Category) String() string {
	switch c.category {
	case 1:
		return "construction"
	case 2:
		return "nature"
	case 3:
		return "equipment"
	case 4:
		return "items"
	}
	panic("should never happen")
}

func (c Category) Group() string {
	if len(c.group) > 0 {
		return "itemGroup.name." + c.group
	}
	return ""
}
