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
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server"
)

func printStatus(srv *server.Server, startedAt time.Time) {
	st := srv.World().TickStats()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	type row struct {
		name, version string
		latency       time.Duration
	}
	var rows []row
	for p := range srv.Players(nil) {
		rows = append(rows, row{name: p.Name(), version: p.GameVersion(), latency: p.Latency()})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].name < rows[j].name })

	colourTPS := func(tps float64) string {
		switch {
		case tps >= 19.5:
			return cGreen
		case tps >= 17:
			return cYellow
		default:
			return cRed
		}
	}
	colourLoad := func(load float64) string {
		switch {
		case load < 0.5:
			return cGreen
		case load < 0.9:
			return cYellow
		default:
			return cRed
		}
	}
	tps := func(v float64) string { return fmt.Sprintf("%s%.2f%s", colourTPS(v), v, cReset) }
	load := func(v float64) string { return fmt.Sprintf("%s%.0f%%%s", colourLoad(v), v*100, cReset) }
	mb := func(b uint64) string { return fmt.Sprintf("%.1f MB", float64(b)/1024/1024) }

	max := srv.MaxPlayerCount()
	maxStr := fmt.Sprint(max)
	if max == len(rows)+1 {
		maxStr = "∞"
	}
	var b strings.Builder
	w := func(format string, args ...any) { fmt.Fprintf(&b, format+"\n", args...) }
	w("%s%s---- TM-Dragonfly status ----%s", cBold, cGold, cReset)
	w("%sUptime:%s     %s", cCyan, cReset, time.Since(startedAt).Round(time.Second))
	w("%sTPS:%s        %s (1s)  %s (5s)  %s (1m)", cCyan, cReset, tps(st.TPS1s), tps(st.TPS5s), tps(st.TPS1m))
	w("%sLoad:%s       %s (1s)  %s (5s)  %s (1m)   tick avg %s, max %s (1m)", cCyan, cReset,
		load(st.Load1s), load(st.Load5s), load(st.Load1m), st.TickAvg.Round(10*time.Microsecond), st.TickMax.Round(10*time.Microsecond))
	w("%sOnline:%s     %d / %s", cCyan, cReset, len(rows), maxStr)
	for _, wld := range srv.Worlds() {
		ws := wld.TickStats()
		w("%sWorld %s:%s %d chunks loaded, %d entities, TPS %s, load %s", cCyan, srv.WorldName(wld), cReset,
			ws.Chunks, ws.Entities, tps(ws.TPS5s), load(ws.Load5s))
	}

	w("%sMemory:%s     %s in use, %s reserved from OS", cCyan, cReset, mb(ms.HeapInuse+ms.StackInuse), mb(ms.Sys))
	w("%sRuntime:%s    %d tasks, %d CPUs", cCyan, cReset, runtime.NumGoroutine(), runtime.NumCPU())
	if len(rows) > 0 {
		var names []string
		for _, r := range rows {
			names = append(names, fmt.Sprintf("%s%s%s (%s, %dms)", cGold, r.name, cReset, r.version, r.latency.Milliseconds()))
		}
		w("%sPlayers:%s    %s", cCyan, cReset, strings.Join(names, ", "))
	}
	_, _ = os.Stdout.WriteString(b.String())
}
