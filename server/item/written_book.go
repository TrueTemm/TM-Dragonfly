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

package item

type WrittenBook struct {
	Title string

	Author string

	Generation WrittenBookGeneration

	Pages []string
}

func (WrittenBook) MaxCount() int {
	return 16
}

func (w WrittenBook) TotalPages() int {
	return len(w.Pages)
}

func (w WrittenBook) Page(page int) (string, bool) {
	if page < 0 || len(w.Pages) <= page {
		return "", false
	}
	return w.Pages[page], true
}

func (w WrittenBook) DecodeNBT(data map[string]any) any {
	if pages, ok := data["pages"].([]any); ok {
		w.Pages = make([]string, len(pages))
		for i, page := range pages {
			w.Pages[i] = page.(map[string]any)["text"].(string)
		}
	}
	w.Title, _ = data["title"].(string)
	w.Author, _ = data["author"].(string)
	if v, ok := data["generation"].(uint8); ok {
		switch v {
		case 0:
			w.Generation = OriginalGeneration()
		case 1:
			w.Generation = CopyGeneration()
		case 2:
			w.Generation = CopyOfCopyGeneration()
		}
	}
	return w
}

func (w WrittenBook) EncodeNBT() map[string]any {
	pages := make([]any, 0, len(w.Pages))
	for _, page := range w.Pages {
		pages = append(pages, map[string]any{"text": page})
	}
	return map[string]any{
		"pages":      pages,
		"author":     w.Author,
		"title":      w.Title,
		"generation": w.Generation.Uint8(),
	}
}

func (WrittenBook) EncodeItem() (name string, meta int16) {
	return "minecraft:written_book", 0
}
