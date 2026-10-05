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
	"strings"
)

type Element interface {
	json.Marshaler
	elem()
}

type MenuElement interface {
	json.Marshaler
	menuElem()
}

type Divider struct{}

func (d Divider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "divider",
		"text": "",
	})
}

type Header struct {
	Text string
}

func NewHeader(text string) Header {
	return Header{Text: text}
}

func (h Header) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "header",
		"text": h.Text,
	})
}

type Label struct {
	Text string
}

func NewLabel(text string) Label {
	return Label{Text: text}
}

func (l Label) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "label",
		"text": l.Text,
	})
}

type Input struct {
	Text string

	Default string

	Placeholder string

	Tooltip string

	value string
}

func NewInput(text, defaultValue, placeholder string) Input {
	return Input{Text: text, Default: defaultValue, Placeholder: placeholder}
}

func (i Input) WithTooltip(tooltip string) Input {
	i.Tooltip = tooltip
	return i
}

func (i Input) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":        "input",
		"text":        i.Text,
		"default":     i.Default,
		"placeholder": i.Placeholder,
	}
	if i.Tooltip != "" {
		m["tooltip"] = i.Tooltip
	}
	return json.Marshal(m)
}

func (i Input) Value() string {
	return i.value
}

type Toggle struct {
	Text string

	Default bool

	Tooltip string

	value bool
}

func NewToggle(text string, defaultValue bool) Toggle {
	return Toggle{Text: text, Default: defaultValue}
}

func (t Toggle) WithTooltip(tooltip string) Toggle {
	t.Tooltip = tooltip
	return t
}

func (t Toggle) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "toggle",
		"text":    t.Text,
		"default": t.Default,
	}
	if t.Tooltip != "" {
		m["tooltip"] = t.Tooltip
	}
	return json.Marshal(m)
}

func (t Toggle) Value() bool {
	return t.value
}

type Slider struct {
	Text string

	Min, Max float64

	StepSize float64

	Default float64

	Tooltip string

	value float64
}

func NewSlider(text string, min, max, stepSize, defaultValue float64) Slider {
	return Slider{Text: text, Min: min, Max: max, StepSize: stepSize, Default: defaultValue}
}

func (s Slider) WithTooltip(tooltip string) Slider {
	s.Tooltip = tooltip
	return s
}

func (s Slider) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "slider",
		"text":    s.Text,
		"min":     s.Min,
		"max":     s.Max,
		"step":    s.StepSize,
		"default": s.Default,
	}
	if s.Tooltip != "" {
		m["tooltip"] = s.Tooltip
	}
	return json.Marshal(m)
}

func (s Slider) Value() float64 {
	return s.value
}

type Dropdown struct {
	Text string

	Options []string

	DefaultIndex int

	Tooltip string

	value int
}

func NewDropdown(text string, options []string, defaultIndex int) Dropdown {
	return Dropdown{Text: text, Options: options, DefaultIndex: defaultIndex}
}

func (d Dropdown) WithTooltip(tooltip string) Dropdown {
	d.Tooltip = tooltip
	return d
}

func (d Dropdown) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "dropdown",
		"text":    d.Text,
		"default": d.DefaultIndex,
		"options": d.Options,
	}
	if d.Tooltip != "" {
		m["tooltip"] = d.Tooltip
	}
	return json.Marshal(m)
}

func (d Dropdown) Value() int {
	return d.value
}

type StepSlider Dropdown

func NewStepSlider(text string, options []string, defaultIndex int) StepSlider {
	return StepSlider{Text: text, Options: options, DefaultIndex: defaultIndex}
}

func (s StepSlider) WithTooltip(tooltip string) StepSlider {
	s.Tooltip = tooltip
	return s
}

func (s StepSlider) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type":    "step_slider",
		"text":    s.Text,
		"default": s.DefaultIndex,
		"steps":   s.Options,
	}
	if s.Tooltip != "" {
		m["tooltip"] = s.Tooltip
	}
	return json.Marshal(m)
}

func (s StepSlider) Value() int {
	return s.value
}

type Button struct {
	Text string

	Image string
}

func NewButton(text, image string) Button {
	return Button{Text: text, Image: image}
}

func (b Button) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"type": "button",
		"text": b.Text,
	}
	if b.Image != "" {
		buttonType := "path"
		if strings.HasPrefix(b.Image, "http:") || strings.HasPrefix(b.Image, "https:") {
			buttonType = "url"
		}
		m["image"] = map[string]any{"type": buttonType, "data": b.Image}
	}
	return json.Marshal(m)
}

func (Divider) elem()    {}
func (Header) elem()     {}
func (Label) elem()      {}
func (Input) elem()      {}
func (Toggle) elem()     {}
func (Slider) elem()     {}
func (Dropdown) elem()   {}
func (StepSlider) elem() {}

func (Divider) menuElem() {}
func (Header) menuElem()  {}
func (Label) menuElem()   {}
func (Button) menuElem()  {}
