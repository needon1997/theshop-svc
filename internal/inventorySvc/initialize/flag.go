package initialize

import (
	"flag"
)

var (
	ConfigPath *string = new(string)
	DevMode    *bool   = new(bool)
)

func ParseFlag() {
	flag.StringVar(ConfigPath, "config", "./config.yaml", "the location of the configuration file")
	flag.BoolVar(DevMode, "dev", false, "dev or prod")
	flag.Parse()
}
