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

package session

import (
	"time"

	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/text/language"
)

func nonEmpty(message string) string {
	if message == "" {
		return " "
	}
	return message
}

func (s *Session) SendMessage(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeRaw,
		Message:  nonEmpty(message),
	})
}

func (s *Session) SendTranslation(t chat.Translation, l language.Tag, a []any) {
	tr := t.F(a...)
	s.writePacket(&packet.Text{
		TextType:         packet.TextTypeTranslation,
		NeedsTranslation: true,
		Message:          tr.Resolve(l),
		Parameters:       tr.Params(l),
	})
}

func (s *Session) SendTip(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeTip,
		Message:  nonEmpty(message),
	})
}

func (s *Session) SendAnnouncement(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeAnnouncement,
		Message:  nonEmpty(message),
	})
}

func (s *Session) SendPopup(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypePopup,
		Message:  nonEmpty(message),
	})
}

func (s *Session) SendJukeboxPopup(message string) {
	s.writePacket(&packet.Text{
		TextType: packet.TextTypeJukeboxPopup,
		Message:  nonEmpty(message),
	})
}

func (s *Session) SendToast(title, message string) {
	s.writePacket(&packet.ToastRequest{
		Title:   title,
		Message: message,
	})
}

func (s *Session) SendScoreboard(sb *scoreboard.Scoreboard) {
	if s == Nop {
		return
	}
	currentName, currentLines := *s.currentScoreboard.Load(), *s.currentLines.Load()

	if currentName != sb.Name() {
		s.RemoveScoreboard()
		pk := &packet.SetDisplayObjective{
			DisplaySlot:   "sidebar",
			ObjectiveName: sb.Name(),
			DisplayName:   sb.Name(),
			CriteriaName:  "dummy",
		}
		if sb.Descending() {
			pk.SortOrder = packet.ScoreboardSortOrderDescending
		} else {
			pk.SortOrder = packet.ScoreboardSortOrderAscending
		}
		s.writePacket(pk)
		name, lines := sb.Name(), append([]string(nil), sb.Lines()...)
		s.currentScoreboard.Store(&name)
		s.currentLines.Store(&lines)
		s.sendScoreLines(sb, nil)
		return
	}

	next := sb.Lines()
	var remove []int
	var changed []int
	for i := 0; i < len(next) || i < len(currentLines); i++ {
		switch {
		case i >= len(next):
			remove = append(remove, i)
		case i >= len(currentLines):
			changed = append(changed, i)
		case currentLines[i] != next[i]:
			remove = append(remove, i)
			changed = append(changed, i)
		}
	}
	if len(remove) == 0 && len(changed) == 0 {
		return
	}
	if len(remove) > 0 {
		pk := &packet.SetScore{}
		for _, i := range remove {
			pk.Entries = append(pk.Entries, protocol.ScoreboardEntry{
				IdentityType:  protocol.ScoreboardIdentityRemove,
				EntryID:       int64(i),
				ObjectiveName: currentName,
				Score:         int32(i),
			})
		}
		s.writePacket(pk)
	}
	lines := append([]string(nil), next...)
	s.currentLines.Store(&lines)
	s.sendScoreLines(sb, changed)
}

func (s *Session) sendScoreLines(sb *scoreboard.Scoreboard, indices []int) {
	lines := sb.Lines()
	if indices == nil {
		indices = make([]int, len(lines))
		for i := range lines {
			indices[i] = i
		}
	}
	pk := &packet.SetScore{}
	for _, k := range indices {
		line := lines[k]
		if len(line) == 0 {
			line = "§" + colours[k]
		}
		pk.Entries = append(pk.Entries, protocol.ScoreboardEntry{
			EntryID:       int64(k),
			ObjectiveName: sb.Name(),
			Score:         int32(k),
			IdentityType:  protocol.ScoreboardIdentityFakePlayer,
			DisplayName:   line,
		})
	}
	if len(pk.Entries) > 0 {
		s.writePacket(pk)
	}
}

var colours = [15]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func (s *Session) RemoveScoreboard() {
	s.writePacket(&packet.RemoveObjective{ObjectiveName: *s.currentScoreboard.Load()})
	var name string
	var lines []string
	s.currentScoreboard.Store(&name)
	s.currentLines.Store(&lines)
}

func (s *Session) SendBossBar(text string, colour uint8, healthPercentage float64) {
	s.RemoveBossBar()
	s.writePacket(&packet.BossEvent{
		BossEntityUniqueID: selfEntityRuntimeID,
		EventType:          packet.BossEventShow,
		BossBarTitle:       text,
		HealthPercentage:   float32(healthPercentage),
		Colour:             colour,
	})
}

func (s *Session) RemoveBossBar() {
	s.writePacket(&packet.BossEvent{
		BossEntityUniqueID: selfEntityRuntimeID,
		EventType:          packet.BossEventHide,
	})
}

const tickLength = time.Second / 20

func (s *Session) SetTitleDurations(fadeInDuration, remainDuration, fadeOutDuration time.Duration) {
	s.writePacket(&packet.SetTitle{
		ActionType:      packet.TitleActionSetDurations,
		FadeInDuration:  int32(fadeInDuration / tickLength),
		RemainDuration:  int32(remainDuration / tickLength),
		FadeOutDuration: int32(fadeOutDuration / tickLength),
	})
}

func (s *Session) SendTitle(text string) {
	s.writePacket(&packet.SetTitle{ActionType: packet.TitleActionSetTitle, Text: text})
}

func (s *Session) SendSubtitle(text string) {
	s.writePacket(&packet.SetTitle{ActionType: packet.TitleActionSetSubtitle, Text: text})
}

func (s *Session) SendActionBarMessage(text string) {
	s.writePacket(&packet.SetTitle{ActionType: packet.TitleActionSetActionBar, Text: text})
}
