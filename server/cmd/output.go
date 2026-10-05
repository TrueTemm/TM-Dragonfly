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

package cmd

import (
	"errors"
	"fmt"
	"github.com/df-mc/dragonfly/server/player/chat"
	"golang.org/x/text/language"
)

type Output struct {
	errors   []error
	messages []fmt.Stringer
}

func (o *Output) Errorf(format string, a ...any) {
	o.errors = append(o.errors, fmt.Errorf(format, a...))
}

func (o *Output) Error(a ...any) {
	if len(a) == 1 {
		if err, ok := a[0].(error); ok {
			o.errors = append(o.errors, err)
			return
		}
	}
	o.errors = append(o.errors, errors.New(fmt.Sprint(a...)))
}

func (o *Output) Errort(t chat.Translation, a ...any) {
	o.errors = append(o.errors, t.F(a...))
}

func (o *Output) Printf(format string, a ...any) {
	o.messages = append(o.messages, stringer(fmt.Sprintf(format, a...)))
}

func (o *Output) Print(a ...any) {
	o.messages = append(o.messages, stringer(fmt.Sprint(a...)))
}

func (o *Output) Printt(t chat.Translation, a ...any) {
	o.messages = append(o.messages, t.F(a...))
}

func (o *Output) Errors() []error {
	return o.errors
}

func (o *Output) ErrorCount() int {
	return len(o.errors)
}

func (o *Output) Messages() []fmt.Stringer {
	return o.messages
}

func (o *Output) MessageCount() int {
	return len(o.messages)
}

type stringer string

func (s stringer) String() string { return string(s) }

var MessageSyntax = chat.Translate(str("%commands.generic.syntax"), 3, `Syntax error: unexpected value: at "%v>>%v<<%v"`).Enc("<red>%v</red>")
var MessageUsage = chat.Translate(str("%commands.generic.usage"), 1, `Usage: %v`).Enc("<red>%v</red>")
var MessageUnknown = chat.Translate(str("%commands.generic.unknown"), 1, `Unknown command: "%v": Please check that the command exists and that you have permission to use it.`).Enc("<red>%v</red>")
var MessageNoTargets = chat.Translate(str("%commands.generic.noTargetMatch"), 0, `No targets matched selector`).Enc("<red>%v</red>")
var MessageNumberInvalid = chat.Translate(str("%commands.generic.num.invalid"), 1, `'%v' is not a valid number`).Enc("<red>> %v</red>")
var MessageBooleanInvalid = chat.Translate(str("%commands.generic.boolean.invalid"), 1, `'%v' is not true or false`).Enc("<red>> %v</red>")
var MessagePlayerNotFound = chat.Translate(str("%commands.generic.player.notFound"), 0, `That player cannot be found`).Enc("<red>> %v</red>")
var MessageParameterInvalid = chat.Translate(str("%commands.generic.parameter.invalid"), 1, `'%v' is not a valid parameter`).Enc("<red>> %v</red>")

type str string

func (s str) Resolve(language.Tag) string { return string(s) }
