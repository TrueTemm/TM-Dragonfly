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

package main

import (
	"os"
	"strings"
	"sync"
)

type ops struct {
	mu   sync.RWMutex
	path string
	set  map[string]string // lower -> display nick
}

func loadOps(path string) *ops {
	o := &ops{path: path, set: map[string]string{}}
	data, _ := os.ReadFile(path)
	for _, l := range strings.Split(string(data), "\n") {
		l = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(l), "-"))
		if l == "" || strings.HasPrefix(l, "#") || strings.HasSuffix(l, ":") {
			continue
		}
		o.set[strings.ToLower(l)] = l
	}
	return o
}

func (o *ops) is(name string) bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	_, ok := o.set[strings.ToLower(name)]
	return ok
}

func (o *ops) add(name string) {
	o.mu.Lock()
	o.set[strings.ToLower(name)] = name
	o.save()
	o.mu.Unlock()
}

func (o *ops) remove(name string) {
	o.mu.Lock()
	delete(o.set, strings.ToLower(name))
	o.save()
	o.mu.Unlock()
}

func (o *ops) save() {
	var b strings.Builder
	for _, n := range o.set {
		b.WriteString("- " + n + "\n")
	}
	_ = os.WriteFile(o.path, []byte(b.String()), 0o644)
}
