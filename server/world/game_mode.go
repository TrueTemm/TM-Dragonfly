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

package world

type GameMode interface {
	AllowsEditing() bool

	AllowsTakingDamage() bool

	CreativeInventory() bool

	HasCollision() bool

	AllowsFlying() bool

	AllowsInteraction() bool

	Visible() bool

	InstantPortalTravel() bool
}

var (
	GameModeSurvival survival

	GameModeCreative creative

	GameModeAdventure adventure

	GameModeSpectator spectator
)

var gameModeReg = newGameModeRegistry(map[int]GameMode{
	0: GameModeSurvival,
	1: GameModeCreative,
	2: GameModeAdventure,
	3: GameModeSpectator,
})

func GameModeByID(id int) (GameMode, bool) {
	return gameModeReg.Lookup(id)
}

func GameModeID(mode GameMode) (int, bool) {
	return gameModeReg.LookupID(mode)
}

type gameModeRegistry struct {
	gameModes map[int]GameMode
	ids       map[GameMode]int
}

func newGameModeRegistry(mode map[int]GameMode) *gameModeRegistry {
	ids := make(map[GameMode]int, len(mode))
	for k, v := range mode {
		ids[v] = k
	}
	return &gameModeRegistry{gameModes: mode, ids: ids}
}

func (reg *gameModeRegistry) Lookup(id int) (GameMode, bool) {
	mode, ok := reg.gameModes[id]
	if !ok {
		mode = GameModeSurvival
	}
	return mode, ok
}

func (reg *gameModeRegistry) LookupID(mode GameMode) (int, bool) {
	id, ok := reg.ids[mode]
	return id, ok
}

type survival struct{}

func (survival) AllowsEditing() bool       { return true }
func (survival) AllowsTakingDamage() bool  { return true }
func (survival) CreativeInventory() bool   { return false }
func (survival) HasCollision() bool        { return true }
func (survival) AllowsFlying() bool        { return false }
func (survival) AllowsInteraction() bool   { return true }
func (survival) Visible() bool             { return true }
func (survival) InstantPortalTravel() bool { return false }

type creative struct{}

func (creative) AllowsEditing() bool       { return true }
func (creative) AllowsTakingDamage() bool  { return false }
func (creative) CreativeInventory() bool   { return true }
func (creative) HasCollision() bool        { return true }
func (creative) AllowsFlying() bool        { return true }
func (creative) AllowsInteraction() bool   { return true }
func (creative) Visible() bool             { return true }
func (creative) InstantPortalTravel() bool { return true }

type adventure struct{}

func (adventure) AllowsEditing() bool       { return false }
func (adventure) AllowsTakingDamage() bool  { return true }
func (adventure) CreativeInventory() bool   { return false }
func (adventure) HasCollision() bool        { return true }
func (adventure) AllowsFlying() bool        { return false }
func (adventure) AllowsInteraction() bool   { return true }
func (adventure) Visible() bool             { return true }
func (adventure) InstantPortalTravel() bool { return false }

type spectator struct{}

func (spectator) AllowsEditing() bool       { return false }
func (spectator) AllowsTakingDamage() bool  { return false }
func (spectator) CreativeInventory() bool   { return false }
func (spectator) HasCollision() bool        { return false }
func (spectator) AllowsFlying() bool        { return true }
func (spectator) AllowsInteraction() bool   { return false }
func (spectator) Visible() bool             { return false }
func (spectator) InstantPortalTravel() bool { return false }
