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

package v786

import (
	"encoding/json"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

const FormElementsProtocol = 786

func fromLatestModalFormRequest(proto uint32, pk *packet.ModalFormRequest) *packet.ModalFormRequest {
	if proto >= FormElementsProtocol {
		return pk
	}
	var form map[string]json.RawMessage
	if err := json.Unmarshal(pk.FormData, &form); err != nil {
		return pk
	}
	var kind string
	if err := json.Unmarshal(form["type"], &kind); err != nil {
		return pk
	}
	var changed bool
	switch kind {
	case "form":
		changed = legacyMenu(form)
	case "custom_form":
		changed = legacyCustomForm(form)
	}
	if !changed {
		return pk
	}
	data, err := json.Marshal(form)
	if err != nil {
		return pk
	}
	return &packet.ModalFormRequest{FormID: pk.FormID, FormData: data}
}

func legacyMenu(form map[string]json.RawMessage) bool {
	raw, ok := form["elements"]
	if !ok {
		return false
	}
	var elements []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &elements); err != nil {
		return false
	}
	buttons := make([]map[string]json.RawMessage, 0, len(elements))
	for _, e := range elements {
		var kind string
		_ = json.Unmarshal(e["type"], &kind)
		if kind != "button" {
			continue
		}
		delete(e, "type")
		buttons = append(buttons, e)
	}
	data, err := json.Marshal(buttons)
	if err != nil {
		return false
	}
	form["buttons"] = data
	delete(form, "elements")
	return true
}

func legacyCustomForm(form map[string]json.RawMessage) bool {
	raw, ok := form["content"]
	if !ok {
		return false
	}
	var content []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &content); err != nil {
		return false
	}
	changed := false
	for _, e := range content {
		var kind string
		_ = json.Unmarshal(e["type"], &kind)
		switch kind {
		case "header":
			e["type"] = json.RawMessage(`"label"`)
			changed = true
		case "divider":
			e["type"], e["text"] = json.RawMessage(`"label"`), json.RawMessage(`""`)
			changed = true
		}
	}
	if !changed {
		return false
	}
	data, err := json.Marshal(content)
	if err != nil {
		return false
	}
	form["content"] = data
	return true
}
