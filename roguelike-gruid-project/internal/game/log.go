package game

type logStyle int

const (
	logNormal logStyle = iota
	logCritic
	logNotable
	logDamage
	logSpecial
	logStatusEnd
	logError
	logConfirm
)

type logEntry struct {
	Text  string
	MText string
	Index int
	Tick  bool
	Style logStyle
	Dups  int
}

func (g *Game) PrintEntry(e logEntry) {
	// if e.Index == g.LogNextTick {
	// 	e.Tick = true
	// }
	// if !e.Tick && len(g.Log) > 0 {
	// 	le := g.Log[len(g.Log)-1]
	// 	if le.Text == e.Text {
	// 		le.Dups++
	// 		le.MText = le.String()
	// 		g.Log[len(g.Log)-1] = le
	// 		return
	// 	}
	// }
	// e.MText = e.String()
	// if LogGame {
	// 	log.Printf("Depth %d:Turn %d:%v", g.Depth, g.Turn, e.dumpString())
	// }
	// g.Log = append(g.Log, e)
	// g.LogIndex++
	// if len(g.Log) > 100000 {
	// 	g.Log = g.Log[10000:]
	// }
}

func (g *Game) Print(s string) {
	// e := logEntry{Text: s, Index: g.LogIndex}
	// g.PrintEntry(e)
}
