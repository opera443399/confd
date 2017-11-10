package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/opera443399/confd/backends"
	"github.com/opera443399/confd/log"
	"github.com/opera443399/confd/template"
)

var (
	backend           string
	nodes             Nodes
	backendsConfig    backends.Config
	prefix            string
	onetime           bool
	watch             bool
	interval          int
	templateConfig    template.Config
	keepStageFile     bool
	syncOnly          bool
	noop              bool
	configFile        = ""
	defaultConfigFile = "/etc/confd/confd.toml"
	confdir           string
	config            Config // holds the global confd config.
	logLevel          string
	printVersion      bool
)

// A Config structure is used to configure confd.
type Config struct {
	Backend      string   `toml:"backend"`
	BackendNodes []string `toml:"nodes"`
	Prefix       string   `toml:"prefix"`
	Watch        bool     `toml:"watch"`
	Interval     int      `toml:"interval"`
	SyncOnly     bool     `toml:"sync-only"`
	Noop         bool     `toml:"noop"`
	ConfDir      string   `toml:"confdir"`
	LogLevel     string   `toml:"log-level"`
}

func init() {
	flag.BoolVar(&onetime, "onetime", false, "run once and exit")
	flag.StringVar(&backend, "backend", "etcd", "backend to use")
	flag.Var(&nodes, "node", "list of backend nodes")
	flag.StringVar(&prefix, "prefix", "", "key path prefix")
	flag.BoolVar(&watch, "watch", false, "enable watch support")
	flag.IntVar(&interval, "interval", 600, "backend polling interval")
	flag.BoolVar(&keepStageFile, "keep-stage-file", false, "keep staged files")
	flag.BoolVar(&syncOnly, "sync-only", false, "sync without check_cmd and reload_cmd")
	flag.BoolVar(&noop, "noop", false, "only show pending changes")
	flag.StringVar(&configFile, "config-file", "", "the confd config file")
	flag.StringVar(&confdir, "confdir", "/etc/confd", "confd conf directory")
	flag.StringVar(&logLevel, "log-level", "", "level which confd should log messages")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
}

// initConfig initializes the confd configuration by first setting defaults,
// then overriding settings from the confd config file, then overriding
// settings from environment variables, and finally overriding
// settings from flags set on the command line.
// It returns an error if any.
func initConfig() error {
	// Set defaults.
	config = Config{
		Backend:  "etcd",
		ConfDir:  "/etc/confd",
		Interval: 600,
		Prefix:   "",
	}

	// Loading configFile
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		}
		log.Debug("Skipping confd config file.")
	} else {
		log.Debug("Loading " + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err
		}
	}

	// Update config from commandline flags.
	processFlags()

	if config.LogLevel != "" {
		log.SetLevel(config.LogLevel)
	}

	if len(config.BackendNodes) == 0 {
		switch config.Backend {
		case "etcd":
			peerstr := os.Getenv("ETCDCTL_PEERS")
			if len(peerstr) > 0 {
				config.BackendNodes = strings.Split(peerstr, ",")
			} else {
				config.BackendNodes = []string{"http://127.0.0.1:4001"}
			}
		case "etcdv3":
			config.BackendNodes = []string{"127.0.0.1:2379"}
		}
	}
	// Initialize the storage client
	log.Info("Backend set to " + config.Backend)

	backendsConfig = backends.Config{
		Backend:      config.Backend,
		BackendNodes: config.BackendNodes,
	}

	// Template configuration.
	templateConfig = template.Config{
		ConfDir:       config.ConfDir,
		ConfigDir:     filepath.Join(config.ConfDir, "conf.d"),
		KeepStageFile: keepStageFile,
		SyncOnly:      config.SyncOnly,
		Noop:          config.Noop,
		Prefix:        config.Prefix,
		TemplateDir:   filepath.Join(config.ConfDir, "templates"),
	}
	return nil
}

// processFlags iterates through each flag set on the command line and
// overrides corresponding configuration settings.
func processFlags() {
	flag.Visit(setConfigFromFlag)
}

func setConfigFromFlag(f *flag.Flag) {
	switch f.Name {
	case "backend":
		config.Backend = backend
	case "node":
		config.BackendNodes = nodes
	case "prefix":
		config.Prefix = prefix
	case "watch":
		config.Watch = watch
	case "interval":
		config.Interval = interval
	case "sync-only":
		config.SyncOnly = syncOnly
	case "noop":
		config.Noop = noop
	case "confdir":
		config.ConfDir = confdir
	case "log-level":
		config.LogLevel = logLevel

	}
}
