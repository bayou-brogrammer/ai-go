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
	// Functionality currently handled by internal/log/log.go
}

func (g *Game) Print(s string) {
	// Functionality currently handled by internal/log/log.go
}
