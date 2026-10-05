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

package entity

import (
	"fmt"
	"math"
)

type ExperienceManager struct {
	experience int
	dec        float64
}

func NewExperienceManager() *ExperienceManager {
	return &ExperienceManager{}
}

func (e *ExperienceManager) Experience() int {
	return e.experience
}

func (e *ExperienceManager) Add(amount int) (level int, progress float64) {
	if e.experience += amount; e.experience < 0 {
		e.Reset()
	}
	return progressFromExperience(e.total())
}

func (e *ExperienceManager) total() float64 {
	return float64(e.experience) + e.dec
}

func (e *ExperienceManager) Level() int {
	level, _ := progressFromExperience(e.total())
	return level
}

func (e *ExperienceManager) SetLevel(level int) {
	if level < 0 || level > math.MaxInt32 {
		panic(fmt.Sprintf("level must be between 0 and 2,147,483,647, got %v", level))
	}
	_, progress := progressFromExperience(e.total())
	e.experience = experienceForLevels(level) + int(float64(experienceForLevel(level))*progress)
}

func (e *ExperienceManager) Progress() float64 {
	_, progress := progressFromExperience(e.total())
	return progress
}

func (e *ExperienceManager) SetProgress(progress float64) {
	if progress < 0 || progress > 1 {
		panic(fmt.Sprintf("progress must be between 0 and 1, got %f", progress))
	}
	currentLevel, _ := progressFromExperience(e.total())
	progressExp := float64(experienceForLevel(currentLevel)) * progress
	e.experience = experienceForLevels(currentLevel) + int(progressExp)
	e.dec = progressExp - math.Trunc(progressExp)
}

func (e *ExperienceManager) Reset() {
	e.experience, e.dec = 0, 0
}

func progressFromExperience(experience float64) (level int, progress float64) {
	var a, b, c float64

	switch {
	case experience <= float64(experienceForLevels(16)):
		a, b = 1.0, 6.0
	case experience <= float64(experienceForLevels(31)):
		a, b, c = 2.5, -40.5, 360.0
	default:
		a, b, c = 4.5, -162.5, 2220.0
	}

	var sol float64
	if d := b*b - 4*a*(c-experience); d > 0 {
		s := math.Sqrt(d)
		sol = math.Max((-b+s)/(2*a), (-b-s)/(2*a))
	} else if d == 0 {
		sol = -b / (2 * a)
	}
	return int(sol), sol - math.Trunc(sol)
}

func experienceForLevels(level int) int {
	if level <= 16 {
		return level*level + level*6
	} else if level <= 31 {
		return int(float64(level*level)*2.5 - 40.5*float64(level) + 360)
	}
	return int(float64(level*level)*4.5 - 162.5*float64(level) + 2220)
}

func experienceForLevel(level int) int {
	if level <= 15 {
		return 2*level + 7
	} else if level <= 30 {
		return 5*level - 38
	}
	return 9*level - 158
}
