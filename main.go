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
	"bufio"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
)

func main() {
	level := slog.LevelInfo
	if os.Getenv("TMDRAGONFLY_DEBUG") != "" {
		level = slog.LevelDebug
	}
	log := colourLogger(os.Stdout, level)
	slog.SetDefault(log)
	printStartupBanner()

	debug.SetGCPercent(50)
	if os.Getenv("GOMEMLIMIT") == "" {
		debug.SetMemoryLimit(1 << 30)
	}

	log.Info(cBold+cLime+"TM-Dragonfly"+cReset+" starting…", "build", buildStamp())
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	uc, conf, err := readConfig(log)
	if err != nil {
		panic(err)
	}
	if addr := uc.Debug.PprofAddress; addr != "" {
		go func() {
			log.Info("pprof listening", "addr", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				log.Error("pprof: " + err.Error())
			}
		}()
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()
	startConsoleReader(log, srv, time.Now())

	srv.Listen()
	log.Info(cLime + "Server is up — clients 1.21.0 … 1.26.50 welcome." + cReset)
	log.Info(cGray + "Type 'help' for the console commands." + cReset)
	for p := range srv.Accept() {
		log.Info("Player joined", "name", p.Name(), "version", p.GameVersion())
		_ = p
	}
}

func startConsoleReader(log *slog.Logger, srv *server.Server, startedAt time.Time) {
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			switch {
			case line == "help":
				printHelp()
			case strings.HasPrefix(line, "say "):
				if msg := strings.TrimSpace(line[4:]); msg != "" {
					_, _ = chat.Global.WriteString("§e[Server] §r" + msg + "§r")
				}
			case line == "list":
				printList(srv)
			case line == "status":
				printStatus(srv, startedAt)
			case line == "version":
				printVersion()
			case line == "stop":
				log.Info("Stopping the server.")
				if err := srv.Close(); err != nil {
					log.Error("close server: " + err.Error())
				}
			default:
				log.Info("Unknown command — type 'help' for the list.", "input", line)
			}
		}
	}()
}

func printStartupBanner() {
	fmt.Println()
	fmt.Printf("  %s%sTM-Dragonfly%s  multiversion Bedrock server\n", cBold, cLime, cReset)
	fmt.Printf("  %ssupported%s Bedrock 1.21.0 … 1.26.50   %sby TrueTemm%s\n", cLime, cReset, cGray, cReset)
	fmt.Println()
}

func printHelp() {
	fmt.Printf("%s%s TM-Dragonfly console %s\n", cBold, cLime, cReset)
	for _, c := range [][2]string{
		{"help", "show this list"},
		{"list", "online players and their versions"},
		{"status", "uptime, TPS, load and memory"},
		{"version", "server build and supported versions"},
		{"say <message>", "broadcast a message to everyone"},
		{"stop", "shut the server down"},
	} {
		fmt.Printf("  %s%-16s%s %s\n", cLime, c[0], cReset, c[1])
	}
}

func printList(srv *server.Server) {
	var names []string
	for p := range srv.Players(nil) {
		names = append(names, fmt.Sprintf("%s (%s)", p.Name(), p.GameVersion()))
	}
	sort.Strings(names)
	fmt.Printf("%s%sOnline: %d%s\n", cBold, cLime, len(names), cReset)
	for _, n := range names {
		fmt.Printf("  %s\n", n)
	}
}

func printVersion() {
	fmt.Printf("%s%sTM-Dragonfly%s  build %s\n", cBold, cLime, cReset, buildStamp())
	fmt.Printf("  multiversion: Bedrock %s1.21.0 … 1.26.50%s (native 1.26.50)\n", cLime, cReset)
}

func buildStamp() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	var revision, date string
	modified := false
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			date = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if revision == "" {
		return "unknown"
	}
	stamp := revision[:min(len(revision), 7)]
	if d, err := time.Parse(time.RFC3339, date); err == nil {
		stamp += " " + d.Local().Format("2006-01-02 15:04")
	}
	if modified {
		stamp += " (modified)"
	}
	return stamp
}

func readConfig(log *slog.Logger) (server.UserConfig, server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return c, zero, fmt.Errorf("create default config: %v", err)
		}
		conf, err := c.Config(log)
		return c, conf, err
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return c, zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, zero, fmt.Errorf("decode config: %v", err)
	}
	conf, err := c.Config(log)
	return c, conf, err
}
