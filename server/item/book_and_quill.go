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

import "slices"

type BookAndQuill struct {
	Pages []string
}

func (BookAndQuill) MaxCount() int {
	return 1
}

func (b BookAndQuill) TotalPages() int {
	return len(b.Pages)
}

func (b BookAndQuill) Page(page int) (string, bool) {
	if page < 0 || len(b.Pages) <= page {
		return "", false
	}
	return b.Pages[page], true
}

func (b BookAndQuill) DeletePage(page int) BookAndQuill {
	if page < 0 || page >= 50 {
		panic("invalid page number")
	}
	if _, ok := b.Page(page); !ok {
		panic("cannot delete nonexistent page")
	}
	b.Pages = slices.Delete(b.Pages, page, page+1)
	return b
}

func (b BookAndQuill) InsertPage(page int, text string) BookAndQuill {
	if page < 0 || page >= 50 {
		panic("invalid page number")
	}
	if len(text) > 256 {
		panic("text longer then 256 bytes")
	}
	if page > len(b.Pages) {
		panic("unable to insert page at invalid position")
	}
	b.Pages = slices.Insert(b.Pages, page, text)
	return b
}

func (b BookAndQuill) SetPage(page int, text string) BookAndQuill {
	if page < 0 || page >= 50 {
		panic("invalid page number")
	}
	if len(text) > 256 {
		panic("text longer then 256 bytes")
	}
	if _, ok := b.Page(page); !ok {
		pages := make([]string, page+1)
		copy(pages, b.Pages)
		b.Pages = pages
	}
	b.Pages[page] = text
	return b
}

func (b BookAndQuill) SwapPages(pageOne, pageTwo int) BookAndQuill {
	if pageOne < 0 || pageTwo < 0 {
		panic("negative page number")
	}
	if _, ok := b.Page(max(pageOne, pageTwo)); !ok {
		panic("invalid page number")
	}
	b.Pages[pageOne], b.Pages[pageTwo] = b.Pages[pageTwo], b.Pages[pageOne]
	return b
}

func (b BookAndQuill) DecodeNBT(data map[string]any) any {
	pages, _ := data["pages"].([]any)
	for _, page := range pages {
		if pageData, ok := page.(map[string]any); ok {
			if text, ok := pageData["text"].(string); ok {
				b.Pages = append(b.Pages, text)
			}
		}
	}
	return b
}

func (b BookAndQuill) EncodeNBT() map[string]any {
	if len(b.Pages) == 0 {
		return nil
	}
	pages := make([]any, 0, len(b.Pages))
	for _, page := range b.Pages {
		pages = append(pages, map[string]any{"text": page})
	}
	return map[string]any{"pages": pages}
}

func (BookAndQuill) EncodeItem() (name string, meta int16) {
	return "minecraft:writable_book", 0
}
