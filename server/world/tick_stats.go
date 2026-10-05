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

package world

import (
	"sync"
	"time"
)

type TickStats struct {
	TPS1s, TPS5s, TPS1m float64

	Load1s, Load5s, Load1m float64

	TickAvg, TickMax time.Duration

	Chunks, Entities int
}

type tickStats struct {
	mu   sync.Mutex
	at   [1200]time.Time
	dur  [1200]time.Duration
	next int
	n    int

	chunks, entities int
}

func (s *tickStats) record(start time.Time, d time.Duration, chunks, entities int) {
	s.mu.Lock()
	s.chunks, s.entities = chunks, entities
	s.at[s.next], s.dur[s.next] = start, d
	s.next = (s.next + 1) % len(s.at)
	if s.n < len(s.at) {
		s.n++
	}
	s.mu.Unlock()
}

func (s *tickStats) snapshot(now time.Time) TickStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		st       TickStats
		count    [3]int
		total    [3]time.Duration
		windows  = [3]time.Duration{time.Second, 5 * time.Second, time.Minute}
		tickBase = 50 * time.Millisecond
	)
	for i := 0; i < s.n; i++ {
		idx := (s.next - 1 - i + len(s.at)) % len(s.at)
		age := now.Sub(s.at[idx])
		if age > windows[2] {
			break
		}
		for w := range windows {
			if age <= windows[w] {
				count[w]++
				total[w] += s.dur[idx]
			}
		}
		if s.dur[idx] > st.TickMax {
			st.TickMax = s.dur[idx]
		}
	}

	var age time.Duration
	if s.n > 0 {
		oldest := s.at[(s.next-s.n+len(s.at))%len(s.at)]
		age = now.Sub(oldest) + tickBase
	}
	tps := func(w int) float64 {
		span := min(windows[w], age)
		if span <= 0 {
			return 0
		}
		return min(20, float64(count[w])/span.Seconds())
	}
	load := func(w int) float64 {
		if count[w] == 0 {
			return 0
		}
		return float64(total[w]) / float64(count[w]) / float64(tickBase)
	}
	st.TPS1s, st.TPS5s, st.TPS1m = tps(0), tps(1), tps(2)
	st.Load1s, st.Load5s, st.Load1m = load(0), load(1), load(2)
	if count[2] > 0 {
		st.TickAvg = total[2] / time.Duration(count[2])
	}
	st.Chunks, st.Entities = s.chunks, s.entities
	return st
}

func (w *World) TickStats() TickStats {
	if w == nil {
		return TickStats{}
	}
	return w.stats.snapshot(time.Now())
}
