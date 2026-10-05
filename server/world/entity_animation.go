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

type EntityAnimation struct {
	name          string
	nextState     string
	controller    string
	stopCondition string
}

func NewEntityAnimation(name string) EntityAnimation {
	return EntityAnimation{name: name}
}

func (a EntityAnimation) Name() string {
	return a.name
}

func (a EntityAnimation) Controller() string {
	return a.controller
}

func (a EntityAnimation) WithController(controller string) EntityAnimation {
	a.controller = controller
	return a
}

func (a EntityAnimation) NextState() string {
	return a.nextState
}

func (a EntityAnimation) WithNextState(state string) EntityAnimation {
	a.nextState = state
	return a
}

func (a EntityAnimation) StopCondition() string {
	return a.stopCondition
}

func (a EntityAnimation) WithStopCondition(condition string) EntityAnimation {
	a.stopCondition = condition
	return a
}
