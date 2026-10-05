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

package bossbar

import (
	"fmt"
	"strings"
)

type BossBar struct {
	text   string
	health float64
	c      Colour
}

func New(text ...any) BossBar {
	return BossBar{text: format(text), health: 1, c: Purple()}
}

func (bar BossBar) Text() string {
	return bar.text
}

func (bar BossBar) WithHealthPercentage(v float64) BossBar {
	if v < 0 || v > 1 {
		panic("boss bar: value out of range: health percentage must be between 0.0 and 1.0")
	}
	bar.health = v
	return bar
}

func (bar BossBar) WithColour(c Colour) BossBar {
	bar.c = c
	return bar
}

func (bar BossBar) HealthPercentage() float64 {
	return bar.health
}

func (bar BossBar) Colour() Colour {
	return bar.c
}

func format(a []any) string {
	return strings.TrimSuffix(fmt.Sprintln(a...), "\n")
}
