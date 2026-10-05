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
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func TestLegacyMenu(t *testing.T) {
	latest := `{"type":"form","title":"§lДуэли","content":"Выберите раздел:","elements":[` +
		`{"type":"header","text":"h"},{"type":"button","text":"Дуэли\nСкоро"},{"type":"divider","text":""},` +
		`{"type":"button","text":"Мини-игры\nСкоро","image":{"type":"path","data":"textures/items/x"}}]}`
	out := fromLatestModalFormRequest(685, &packet.ModalFormRequest{FormID: 7, FormData: []byte(latest)})
	if out.FormID != 7 {
		t.Fatalf("form id changed: %d", out.FormID)
	}
	var form struct {
		Type     string           `json:"type"`
		Title    string           `json:"title"`
		Content  string           `json:"content"`
		Buttons  []map[string]any `json:"buttons"`
		Elements []any            `json:"elements"`
	}
	if err := json.Unmarshal(out.FormData, &form); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if form.Elements != nil {
		t.Fatalf("elements still present: %s", out.FormData)
	}
	if len(form.Buttons) != 2 || form.Buttons[0]["text"] != "Дуэли\nСкоро" || form.Buttons[1]["image"] == nil {
		t.Fatalf("buttons wrong: %s", out.FormData)
	}
	if _, ok := form.Buttons[0]["type"]; ok {
		t.Fatalf("button kept its type: %s", out.FormData)
	}
	if form.Title != "§lДуэли" || form.Content != "Выберите раздел:" || form.Type != "form" {
		t.Fatalf("rest of the form changed: %s", out.FormData)
	}

	same := fromLatestModalFormRequest(786, &packet.ModalFormRequest{FormID: 7, FormData: []byte(latest)})
	if string(same.FormData) != latest {
		t.Fatalf("786 form rewritten: %s", same.FormData)
	}
}

func TestLegacyCustomForm(t *testing.T) {
	plain := `{"type":"custom_form","title":"t","content":[{"type":"toggle","text":"a","default":true},{"type":"dropdown","text":"b","options":["x"],"default":0}]}`
	out := fromLatestModalFormRequest(685, &packet.ModalFormRequest{FormID: 1, FormData: []byte(plain)})
	if string(out.FormData) != plain {
		t.Fatalf("plain custom form rewritten: %s", out.FormData)
	}
	fancy := `{"type":"custom_form","title":"t","content":[{"type":"header","text":"h"},{"type":"divider","text":""},{"type":"toggle","text":"a","default":false}]}`
	out = fromLatestModalFormRequest(776, &packet.ModalFormRequest{FormID: 1, FormData: []byte(fancy)})
	var form struct {
		Content []map[string]any `json:"content"`
	}
	if err := json.Unmarshal(out.FormData, &form); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if len(form.Content) != 3 || form.Content[0]["type"] != "label" || form.Content[1]["type"] != "label" || form.Content[1]["text"] != "" || form.Content[2]["type"] != "toggle" {
		t.Fatalf("custom form wrong: %s", out.FormData)
	}
}
