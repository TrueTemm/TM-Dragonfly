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

package scoreboard

import (
	"fmt"
	"slices"
	"strings"
)

type Scoreboard struct {
	name       string
	lines      []string
	padding    bool
	descending bool
}

func New(name ...any) *Scoreboard {
	return &Scoreboard{name: strings.TrimSuffix(fmt.Sprintln(name...), "\n"), padding: true}
}

func (board *Scoreboard) Name() string {
	return board.name
}

func (board *Scoreboard) Write(p []byte) (n int, err error) {
	return board.WriteString(string(p))
}

func (board *Scoreboard) WriteString(s string) (n int, err error) {
	lines := strings.Split(s, "\n")
	board.lines = append(board.lines, lines...)

	if len(board.lines) >= 15 {
		return len(lines), fmt.Errorf("write scoreboard: maximum of 15 lines of text exceeded")
	}
	return len(lines), nil
}

func (board *Scoreboard) Set(index int, s string) {
	if index < 0 || index >= 15 {
		panic(fmt.Sprintf("index out of range %v", index))
	}
	if diff := index - (len(board.lines) - 1); diff > 0 {
		board.lines = append(board.lines, make([]string, diff)...)
	}

	board.lines[index] = strings.TrimSuffix(strings.TrimSuffix(s, "\n"), "\n")
}

func (board *Scoreboard) Remove(index int) {
	if index < 0 || index >= 15 {
		panic(fmt.Sprintf("index out of range %v", index))
	}
	board.lines = append(board.lines[:index], board.lines[index+1:]...)
}

func (board *Scoreboard) RemovePadding() {
	board.padding = false
}

func (board *Scoreboard) Lines() []string {
	lines := slices.Clone(board.lines)
	if board.padding {
		for i, line := range lines {
			if len(board.name)-len(line)-2 <= 0 {
				lines[i] = " " + line + " "
				continue
			}
			lines[i] = " " + line + strings.Repeat(" ", len(board.name)-len(line)-2)
		}
	}
	if board.descending {
		slices.Reverse(lines)
	}
	return lines
}

func (board *Scoreboard) Descending() bool {
	return board.descending
}

func (board *Scoreboard) SetDescending() {
	board.descending = true
}
