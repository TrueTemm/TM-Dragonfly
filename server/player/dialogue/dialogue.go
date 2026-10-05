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

package dialogue

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"reflect"
	"strings"
)

type Dialogue struct {
	title, body string
	submittable Submittable
	buttons     []Button
	display     DisplaySettings
}

func New(submittable Submittable, title ...any) Dialogue {
	t := reflect.TypeOf(submittable)
	if t.Kind() != reflect.Struct {
		panic("submittable must be struct")
	}
	m := Dialogue{title: format(title), submittable: submittable}
	m.verify()
	return m
}

func (m Dialogue) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Buttons())
}

func (m Dialogue) WithBody(body ...any) Dialogue {
	m.body = format(body)
	return m
}

func (m Dialogue) WithDisplay(display DisplaySettings) Dialogue {
	m.display = display
	return m
}

func (m Dialogue) WithButtons(buttons ...Button) Dialogue {
	m.buttons = append(m.buttons, buttons...)
	m.verify()
	return m
}

func (m Dialogue) Title() string {
	return m.title
}

func (m Dialogue) Body() string {
	return m.body
}

func (m Dialogue) Display() DisplaySettings {
	return m.display
}

func (m Dialogue) Buttons() []Button {
	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))

	buttons := make([]Button, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		buttons = append(buttons, field.Interface().(Button))
	}
	buttons = append(buttons, m.buttons...)
	return buttons
}

func (m Dialogue) Submit(index uint, submitter Submitter, tx *world.Tx) error {
	buttons := m.Buttons()
	if index >= uint(len(buttons)) {
		return fmt.Errorf("button index points to inexistent button: %v (only %v buttons present)", index, len(buttons))
	}
	m.submittable.Submit(submitter, buttons[index], tx)
	return nil
}

func (m Dialogue) Close(submitter Submitter, tx *world.Tx) {
	if closer, ok := m.submittable.(Closer); ok {
		closer.Close(submitter, tx)
	}
}

func (m Dialogue) verify() {
	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))
	var buttons int
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		if _, ok := v.Field(i).Interface().(Button); !ok {
			panic("all exported fields must be of the type dialogue.Button")
		}
		buttons++
	}
	if buttons+len(m.buttons) > 6 {
		panic("maximum of 6 buttons allowed")
	}
}

func format(a []any) string {
	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(a...), "\n"), "\n")
}
