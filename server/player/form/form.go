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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/df-mc/dragonfly/server/world"
)

type Form interface {
	json.Marshaler
	SubmitJSON(b []byte, submitter Submitter, tx *world.Tx) error
}

type Custom struct {
	title       string
	submittable Submittable
}

func (f Custom) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "custom_form",
		"title":   f.title,
		"content": f.Elements(),
	})
}

func New(submittable Submittable, title ...any) Custom {
	t := reflect.TypeOf(submittable)
	if t.Kind() != reflect.Struct {
		panic("submittable must be struct")
	}
	f := Custom{title: format(title), submittable: submittable}
	f.verify()
	return f
}

func (f Custom) Title() string {
	return f.title
}

func (f Custom) Elements() []Element {
	v := reflect.New(reflect.TypeOf(f.submittable)).Elem()
	v.Set(reflect.ValueOf(f.submittable))
	n := v.NumField()

	elements := make([]Element, 0, n)
	for i := 0; i < n; i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		elements = append(elements, field.Interface().(Element))
	}
	return elements
}

func (f Custom) SubmitJSON(b []byte, submitter Submitter, tx *world.Tx) error {
	if b == nil {
		if closer, ok := f.submittable.(Closer); ok {
			closer.Close(submitter, tx)
		}
		return nil
	}

	dec := json.NewDecoder(bytes.NewBuffer(b))
	dec.UseNumber()

	var data []any
	if err := dec.Decode(&data); err != nil {
		return fmt.Errorf("error decoding JSON data to slice: %w", err)
	}

	v := reflect.New(reflect.TypeOf(f.submittable)).Elem()
	v.Set(reflect.ValueOf(f.submittable))

	for i := 0; i < v.NumField(); i++ {
		fieldV := v.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		e := fieldV.Interface().(Element)
		if len(data) == 0 {
			return fmt.Errorf("form JSON data array does not have enough values")
		}
		if elementReadonly(e) {
			data = data[1:]
			continue
		}
		elem, err := f.parseValue(e, data[0])
		if err != nil {
			return fmt.Errorf("error parsing form response value: %w", err)
		}
		fieldV.Set(elem)
		data = data[1:]
	}

	v.Interface().(Submittable).Submit(submitter, tx)

	return nil
}

func elementReadonly(e Element) bool {
	switch e.(type) {
	case Divider, Header, Label:
		return true
	default:
		return false
	}
}

func (f Custom) parseValue(elem Element, s any) (reflect.Value, error) {
	var ok bool
	var value reflect.Value

	switch element := elem.(type) {
	case Input:
		element.value, ok = s.(string)
		if !ok {
			return value, fmt.Errorf("value %v is not allowed for input element", s)
		}
		if !utf8.ValidString(element.value) {
			return value, fmt.Errorf("value %v is not valid UTF8", s)
		}
		value = reflect.ValueOf(element)
	case Toggle:
		element.value, ok = s.(bool)
		if !ok {
			return value, fmt.Errorf("value %v is not allowed for toggle element", s)
		}
		value = reflect.ValueOf(element)
	case Slider:
		v, ok := s.(json.Number)
		f, err := v.Float64()
		if !ok || err != nil {
			return value, fmt.Errorf("value %v is not allowed for slider element", s)
		}
		if f > element.Max || f < element.Min {
			return value, fmt.Errorf("slider value %v is out of range %v-%v", f, element.Min, element.Max)
		}
		element.value = f
		value = reflect.ValueOf(element)
	case Dropdown:
		v, ok := s.(json.Number)
		f, err := v.Int64()
		if !ok || err != nil {
			return value, fmt.Errorf("value %v is not allowed for dropdown element", s)
		}
		if f < 0 || int(f) >= len(element.Options) {
			return value, fmt.Errorf("dropdown value %v is out of range %v-%v", f, 0, len(element.Options)-1)
		}
		element.value = int(f)
		value = reflect.ValueOf(element)
	case StepSlider:
		v, ok := s.(json.Number)
		f, err := v.Int64()
		if !ok || err != nil {
			return value, fmt.Errorf("value %v is not allowed for dropdown element", s)
		}
		if f < 0 || int(f) >= len(element.Options) {
			return value, fmt.Errorf("dropdown value %v is out of range %v-%v", f, 0, len(element.Options)-1)
		}
		element.value = int(f)
		value = reflect.ValueOf(element)
	}
	return value, nil
}

func (f Custom) verify() {
	el := reflect.TypeOf((*Element)(nil)).Elem()

	v := reflect.New(reflect.TypeOf(f.submittable)).Elem()
	v.Set(reflect.ValueOf(f.submittable))

	t := reflect.TypeOf(f.submittable)
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		if !t.Field(i).Type.Implements(el) {
			panic("all exported fields must implement form.Element interface")
		}
	}
}

func format(a []any) string {
	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(a...), "\n"), "\n")
}
