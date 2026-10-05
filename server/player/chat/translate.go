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

	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/text/language"
)

var MessageJoin = Translate(str("%multiplayer.player.joined"), 1, `%v joined the game`).Enc("<yellow>%v</yellow>")
var MessageQuit = Translate(str("%multiplayer.player.left"), 1, `%v left the game`).Enc("<yellow>%v</yellow>")
var MessageServerDisconnect = Translate(str("%disconnect.disconnected"), 0, `Disconnected by Server`).Enc("<yellow>%v</yellow>")

var MessageBedTooFar = Translate(str("%tile.bed.tooFar"), 0, `Bed is too far away`).Enc("<grey>%v</grey>")
var MessageBedObstructed = Translate(str("%tile.bed.obstructed"), 0, `Bed is obstructed`).Enc("<grey>%v</grey>")
var MessageRespawnPointSet = Translate(str("%tile.bed.respawnSet"), 0, `Respawn point set`).Enc("<grey>%v</grey>")
var MessageNoSleep = Translate(str("%tile.bed.noSleep"), 0, `You can only sleep at night and during thunderstorms`).Enc("<grey>%v</grey>")
var MessageBedIsOccupied = Translate(str("%tile.bed.occupied"), 0, `This bed is occupied`).Enc("<grey>%v</grey>")
var MessageSleeping = Translate(str("%chat.type.sleeping"), 2, `%v is sleeping in a bed. To skip to dawn, %v more users need to sleep in beds at the same time.`)
var MessageBedNotValid = Translate(str("%tile.bed.notValid"), 0, `Your home bed was missing or obstructed`)

type str string

func (s str) Resolve(language.Tag) string { return string(s) }

type TranslationString interface {
	Resolve(l language.Tag) string
}

func Translate(str TranslationString, params int, fallback string) Translation {
	return Translation{str: str, params: params, fallback: fallback, format: "%v"}
}

type Translation struct {
	str      TranslationString
	format   string
	params   int
	fallback string
}

func (t Translation) Zero() bool {
	return t.format == ""
}

func (t Translation) Enc(format string) Translation {
	t.format = format
	return t
}

func (t Translation) Resolve(l language.Tag) string {
	return t.F().Resolve(l)
}

func (t Translation) F(a ...any) translation {
	if len(a) != t.params {
		panic(fmt.Sprintf("translation '%v' requires exactly %v parameters, got %v", t.format, t.params, len(a)))
	}
	return translation{t: t, params: a}
}

type translation struct {
	t      Translation
	params []any
}

func (t translation) Resolve(l language.Tag) string {
	return text.Colourf(t.t.format, t.t.str.Resolve(l))
}

func (t translation) Params(l language.Tag) []string {
	params := make([]string, len(t.params))
	for i, arg := range t.params {
		if str, ok := arg.(TranslationString); ok {
			params[i] = str.Resolve(l)
			continue
		}
		params[i] = fmt.Sprint(arg)
	}
	return params
}

func (t translation) String() string {
	return fmt.Sprintf(text.Colourf(t.t.format, t.t.fallback), t.params...)
}

func (t translation) Error() string {
	return t.String()
}
