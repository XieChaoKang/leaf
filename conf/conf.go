package conf

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int
	LogDays  int

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int

	// ~/logs/apps/game/serverName_%Y%m%d%H.log
	AppLogPathFmt string
	AppMaxFiles   int
)
