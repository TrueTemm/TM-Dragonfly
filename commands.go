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
	"sort"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func registerCommands(srv *server.Server, o *ops) {
	cmd.Register(cmd.New("about", "core version and info", nil, aboutCmd{}))
	cmd.Register(cmd.New("weather", "clear, rain or thunder", nil, weatherCmd{o: o}))
	cmd.Register(cmd.New("time", "set the time of day", nil, timeCmd{o: o}))
	cmd.Register(cmd.New("gamemode", "change a game mode", []string{"gm"}, gmCmd{o: o, srv: srv}))
	cmd.Register(cmd.New("tp", "teleport to a player", nil, tpCmd{o: o, srv: srv}))
	cmd.Register(cmd.New("give", "give an item to a player", nil, giveCmd{o: o, srv: srv}))
	cmd.Register(cmd.New("op", "grant operator", nil, opCmd{o: o, srv: srv}))
	cmd.Register(cmd.New("deop", "revoke operator", nil, deopCmd{o: o}))
}

func opOnly(o *ops, src cmd.Source) bool {
	p, ok := src.(*player.Player)
	return !ok || o.is(p.Name()) // console always, else an operator
}

type aboutCmd struct{}

func (aboutCmd) Run(_ cmd.Source, out *cmd.Output, _ *world.Tx) {
	out.Printf("§aTM-Dragonfly§r by TrueTemm — multiversion Bedrock 1.20.70 … 1.26.50")
}

type weatherCmd struct {
	o    *ops
	Type string `cmd:"type"`
}

func (c weatherCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c weatherCmd) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	w := tx.World()
	switch strings.ToLower(c.Type) {
	case "clear", "sun":
		w.StopThundering()
		w.StopRaining()
	case "rain":
		w.StartRaining(time.Hour)
	case "thunder", "storm":
		w.StartThundering(time.Hour)
	default:
		out.Errorf("weather: clear, rain or thunder")
		return
	}
	out.Printf("weather set to %s", c.Type)
}

type timeCmd struct {
	o   *ops
	Set string `cmd:"ticks"`
}

func (c timeCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c timeCmd) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	t, ok := parseTicks(strings.TrimPrefix(c.Set, "set "))
	if !ok {
		out.Errorf("time: day, noon, night, midnight or ticks")
		return
	}
	tx.World().SetTime(t)
	out.Printf("time set to %s", c.Set)
}

type gmCmd struct {
	o      *ops
	srv    *server.Server
	Mode   string               `cmd:"mode"`
	Target cmd.Optional[string] `cmd:"player"`
}

func (c gmCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c gmCmd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	mode, ok := gameModeByName(c.Mode)
	if !ok {
		out.Errorf("game mode: survival, creative, adventure or spectator")
		return
	}
	if name, given := c.Target.Load(); given {
		for other := range c.srv.Players(tx) {
			if strings.EqualFold(other.Name(), name) {
				other.SetGameMode(mode)
				out.Printf("game mode set for %s", other.Name())
				return
			}
		}
		out.Errorf("player %s not found", name)
		return
	}
	p, ok := src.(*player.Player)
	if !ok {
		out.Errorf("name a player to change")
		return
	}
	p.SetGameMode(mode)
	out.Printf("game mode set")
}

type tpCmd struct {
	o      *ops
	srv    *server.Server
	Target string `cmd:"player"`
}

func (c tpCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c tpCmd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		out.Errorf("only a player can use this")
		return
	}
	for other := range c.srv.Players(tx) {
		if strings.EqualFold(other.Name(), c.Target) {
			p.Teleport(other.Position())
			out.Printf("teleported to %s", other.Name())
			return
		}
	}
	out.Errorf("player %s not found", c.Target)
}

type opCmd struct {
	o      *ops
	srv    *server.Server
	Target string `cmd:"nick"`
}

func (c opCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c opCmd) Run(_ cmd.Source, out *cmd.Output, tx *world.Tx) {
	display := grantOp(c.srv, c.o, tx, c.Target)
	out.Printf("%s is now an operator", display)
}

func grantOp(srv *server.Server, o *ops, tx *world.Tx, name string) string {
	display := name
	for p := range srv.Players(tx) {
		if strings.EqualFold(p.Name(), name) {
			display = p.Name() // real nick, right case
			p.Message("§aYou have been granted operator rights.")
			break
		}
	}
	o.add(display)
	return display
}

type deopCmd struct {
	o      *ops
	Target string `cmd:"nick"`
}

func (c deopCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c deopCmd) Run(_ cmd.Source, out *cmd.Output, _ *world.Tx) {
	c.o.remove(c.Target)
	out.Printf("%s is no longer an operator", c.Target)
}

type itemName string

func (itemName) Type() string                 { return "item" }
func (itemName) Options(cmd.Source) []string  { return itemOptions }

var itemOptions = loadItemOptions() // every registered item, for tab-complete

func loadItemOptions() []string {
	seen := map[string]bool{}
	var out []string
	for _, it := range world.Items() {
		n, _ := it.EncodeItem()
		n = strings.TrimPrefix(n, "minecraft:")
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}

type giveCmd struct {
	o      *ops
	srv    *server.Server
	Item   itemName             `cmd:"item"`
	Count  cmd.Optional[int]    `cmd:"count"`
	Target cmd.Optional[string] `cmd:"player"`
}

func (c giveCmd) Allow(src cmd.Source) bool { return opOnly(c.o, src) }

func (c giveCmd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	name := strings.ToLower(strings.TrimPrefix(string(c.Item), "minecraft:"))
	it, ok := world.ItemByName("minecraft:"+name, 0)
	if !ok {
		out.Errorf("unknown item: %s", name)
		return
	}
	count := 1
	if n, given := c.Count.Load(); given && n > 0 {
		count = n
	}
	target, ok := c.giveTarget(src, tx)
	if !ok {
		out.Errorf("name a player to give to")
		return
	}
	if _, err := target.Inventory().AddItem(item.NewStack(it, count)); err != nil {
		out.Errorf("%s has no room", target.Name())
		return
	}
	out.Printf("gave %d x %s to %s", count, name, target.Name())
}

func (c giveCmd) giveTarget(src cmd.Source, tx *world.Tx) (*player.Player, bool) {
	if name, given := c.Target.Load(); given {
		for other := range c.srv.Players(tx) {
			if strings.EqualFold(other.Name(), name) {
				return other, true
			}
		}
		return nil, false
	}
	p, ok := src.(*player.Player)
	return p, ok
}

func gameModeByName(name string) (world.GameMode, bool) {
	switch strings.ToLower(name) {
	case "survival", "s", "0":
		return world.GameModeSurvival, true
	case "creative", "c", "1":
		return world.GameModeCreative, true
	case "adventure", "a", "2":
		return world.GameModeAdventure, true
	case "spectator", "sp", "3":
		return world.GameModeSpectator, true
	}
	return nil, false
}
