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

package chat

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

type Subscriber interface {
	UUID() uuid.UUID

	Message(a ...any)
}

type Translator interface {
	Messaget(t Translation, a ...any)
}

type StdoutSubscriber struct{}

var id = uuid.New()

func (c StdoutSubscriber) UUID() uuid.UUID {
	return id
}

func (c StdoutSubscriber) Message(a ...any) {
	s := make([]string, len(a))
	for i, b := range a {
		s[i] = fmt.Sprint(b)
	}
	t := text.ANSI(strings.Join(s, " "))
	if !strings.HasSuffix(t, "\n") {
		fmt.Println(t)
		return
	}
	fmt.Print(t)
}
