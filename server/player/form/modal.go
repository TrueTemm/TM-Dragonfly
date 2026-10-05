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

package form

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"reflect"
)

type Modal struct {
	title, body string
	submittable ModalSubmittable
}

func NewModal(submittable ModalSubmittable, title ...any) Modal {
	t := reflect.TypeOf(submittable)
	if t.Kind() != reflect.Struct {
		panic("submittable must be struct")
	}
	m := Modal{title: format(title), submittable: submittable}
	m.verify()
	return m
}

func YesButton() Button {
	return Button{Text: "gui.yes"}
}

func NoButton() Button {
	return Button{Text: "gui.no"}
}

func (m Modal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "modal",
		"title":   m.title,
		"content": m.body,
		"button1": m.Buttons()[0].Text,
		"button2": m.Buttons()[1].Text,
	})
}

func (m Modal) WithBody(body ...any) Modal {
	m.body = format(body)
	return m
}

func (m Modal) Title() string {
	return m.title
}

func (m Modal) Body() string {
	return m.body
}

func (m Modal) SubmitJSON(b []byte, submitter Submitter, tx *world.Tx) error {
	if b == nil {
		if closer, ok := m.submittable.(Closer); ok {
			closer.Close(submitter, tx)
		}
		return nil
	}

	var value bool
	if err := json.Unmarshal(b, &value); err != nil {
		return fmt.Errorf("error parsing JSON as bool: %w", err)
	}
	if value {
		m.submittable.Submit(submitter, m.Buttons()[0], tx)
		return nil
	}
	m.submittable.Submit(submitter, m.Buttons()[1], tx)
	return nil
}

func (m Modal) Buttons() []Button {
	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))

	buttons := make([]Button, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		buttons = append(buttons, field.Interface().(Button))
	}
	return buttons
}

func (m Modal) verify() {
	var count int

	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		if _, ok := v.Field(i).Interface().(Button); !ok {
			panic("both exported fields must be of the type form.Button")
		}
		count++
	}
	if count != 2 {
		panic("modal form must have exactly two exported fields of the type form.Button")
	}
}
