package leaf

import (
	"leaf/cluster"
	"leaf/conf"
	"leaf/console"
	"leaf/log"
	"leaf/module"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var (
	Wd   string // working path
	name string
	Die  chan bool // wait for end application
)

func init() {
	//fmt.Printf("%v %v %v init up\n", runtime.GOOS, os.Args[0], filepath.Base(os.Args[0]))
	if runtime.GOOS == "windows" {
		name = strings.TrimLeft(filepath.Base(os.Args[0]), ".")
		name = name[:strings.Index(name, ".")]
	} else if runtime.GOOS == "linux" {
		name = strings.TrimLeft(filepath.Base(os.Args[0]), "/")
	}

	Die = make(chan bool)
}

func Run(mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		var err error
		logger, err := log.New(conf.LogLevel, conf.LogPath, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	if conf.AppLogPathFmt != "" {
		appLogger := log.AppLogNew(conf.AppLogPathFmt, conf.AppMaxFiles)
		log.AppLogExport(appLogger)
		defer appLogger.CleanUp()
	}

	log.RefreshLog(conf.LogPath, conf.LogDays)
	log.Release("%v %v starting up", name, version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	sg := make(chan os.Signal)
	signal.Ignore(syscall.Signal(12), syscall.SIGPIPE, syscall.Signal(10))
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	select {
	case <-Die:
		log.Debug("The app will shutdown in a few seconds")
	case sig := <-sg:
		log.Debug("Leaf closing down (signal: %v)", sig)
	}
	log.Release("Leaf server is stopping...")

	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}

func Shutdown() {
	close(Die)
}
