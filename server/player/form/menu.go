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
	"reflect"

	"github.com/df-mc/dragonfly/server/world"
)

type Menu struct {
	title, body string
	submittable MenuSubmittable
	elements    []MenuElement
}

func NewMenu(submittable MenuSubmittable, title ...any) Menu {
	t := reflect.TypeOf(submittable)
	if t.Kind() != reflect.Struct {
		panic("submittable must be struct")
	}
	m := Menu{title: format(title), submittable: submittable}
	m.verify()
	return m
}

func (m Menu) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     "form",
		"title":    m.title,
		"content":  m.body,
		"elements": m.Elements(),
	})
}

func (m Menu) WithBody(body ...any) Menu {
	m.body = format(body)
	return m
}

func (m Menu) AddButton(button Button) Menu {
	m.elements = append(m.elements, button)
	return m
}

func (m Menu) AddDivider(divider Divider) Menu {
	m.elements = append(m.elements, divider)
	return m
}

func (m Menu) AddHeader(header Header) Menu {
	m.elements = append(m.elements, header)
	return m
}

func (m Menu) AddLabel(label Label) Menu {
	m.elements = append(m.elements, label)
	return m
}

func (m Menu) WithButtons(buttons ...Button) Menu {
	for _, b := range buttons {
		m.elements = append(m.elements, b)
	}
	return m
}

func (m Menu) WithElements(elements ...MenuElement) Menu {
	m.elements = append(m.elements, elements...)
	return m
}

func (m Menu) Title() string {
	return m.title
}

func (m Menu) Body() string {
	return m.body
}

func (m Menu) Buttons() []Button {
	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))

	buttons := make([]Button, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		if b, ok := field.Interface().(Button); ok {
			buttons = append(buttons, b)
		}
	}
	for _, elem := range m.elements {
		if b, ok := elem.(Button); ok {
			buttons = append(buttons, b)
		}
	}
	return buttons
}

func (m Menu) Elements() []MenuElement {
	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))

	elements := make([]MenuElement, 0, v.NumField()+len(m.elements))
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		elements = append(elements, field.Interface().(MenuElement))
	}
	elements = append(elements, m.elements...)
	return elements
}

func (m Menu) SubmitJSON(b []byte, submitter Submitter, tx *world.Tx) error {
	if b == nil {
		if closer, ok := m.submittable.(Closer); ok {
			closer.Close(submitter, tx)
		}
		return nil
	}

	var index uint
	err := json.Unmarshal(b, &index)
	if err != nil {
		return fmt.Errorf("cannot parse button index as int: %w", err)
	}
	buttons := m.Buttons()
	if index >= uint(len(buttons)) {
		return fmt.Errorf("button index points to inexistent button: %v (only %v buttons present)", index, len(buttons))
	}
	m.submittable.Submit(submitter, buttons[index], tx)
	return nil
}

func (m Menu) verify() {
	v := reflect.New(reflect.TypeOf(m.submittable)).Elem()
	v.Set(reflect.ValueOf(m.submittable))
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		if _, ok := v.Field(i).Interface().(MenuElement); !ok {
			panic("all exported fields must implement form.MenuElement")
		}
	}
}
