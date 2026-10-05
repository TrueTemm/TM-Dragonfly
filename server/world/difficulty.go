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

type Difficulty interface {
	FoodRegenerates() bool

	StarvationHealthLimit() float64

	FireSpreadIncrease() int
}

var (
	DifficultyPeaceful difficultyPeaceful

	DifficultyEasy difficultyEasy

	DifficultyNormal difficultyNormal

	DifficultyHard difficultyHard
)

var difficultyReg = newDifficultyRegistry(map[int]Difficulty{
	0: DifficultyPeaceful,
	1: DifficultyEasy,
	2: DifficultyNormal,
	3: DifficultyHard,
})

func DifficultyByID(id int) (Difficulty, bool) {
	return difficultyReg.Lookup(id)
}

func DifficultyID(diff Difficulty) (int, bool) {
	return difficultyReg.LookupID(diff)
}

type difficultyRegistry struct {
	difficulties map[int]Difficulty
	ids          map[Difficulty]int
}

func newDifficultyRegistry(diff map[int]Difficulty) *difficultyRegistry {
	ids := make(map[Difficulty]int, len(diff))
	for k, v := range diff {
		ids[v] = k
	}
	return &difficultyRegistry{difficulties: diff, ids: ids}
}

func (reg *difficultyRegistry) Lookup(id int) (Difficulty, bool) {
	dim, ok := reg.difficulties[id]
	if !ok {
		dim = DifficultyNormal
	}
	return dim, ok
}

func (reg *difficultyRegistry) LookupID(diff Difficulty) (int, bool) {
	id, ok := reg.ids[diff]
	return id, ok
}

type difficultyPeaceful struct{}

func (difficultyPeaceful) FoodRegenerates() bool          { return true }
func (difficultyPeaceful) StarvationHealthLimit() float64 { return 20 }
func (difficultyPeaceful) FireSpreadIncrease() int        { return 0 }

type difficultyEasy struct{}

func (difficultyEasy) FoodRegenerates() bool          { return false }
func (difficultyEasy) StarvationHealthLimit() float64 { return 10 }
func (difficultyEasy) FireSpreadIncrease() int        { return 7 }

type difficultyNormal struct{}

func (difficultyNormal) FoodRegenerates() bool          { return false }
func (difficultyNormal) StarvationHealthLimit() float64 { return 2 }
func (difficultyNormal) FireSpreadIncrease() int        { return 14 }

type difficultyHard struct{}

func (difficultyHard) FoodRegenerates() bool          { return false }
func (difficultyHard) StarvationHealthLimit() float64 { return -1 }
func (difficultyHard) FireSpreadIncrease() int        { return 21 }
