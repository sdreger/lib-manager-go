package config

import (
	"github.com/caarlos0/env/v11"
	"runtime/debug"
)

const envPrefix = "LIB_MANAGER_"

// For the testing purpose
var buildInfoFunc = debug.ReadBuildInfo

func New() (AppConfig, error) {
	var cfg AppConfig
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: envPrefix}); err != nil {
		return AppConfig{}, err
	}

	cfg.BuildInfo = GetBuildInfo()

	return cfg, nil
}

func GetBuildInfo() BuildInfo {
	buildInfo := BuildInfo{}
	info, ok := buildInfoFunc()
	if ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				buildInfo.Revision = s.Value
			case "vcs.time":
				buildInfo.Time = s.Value
			case "vcs.modified":
				buildInfo.Dirty = s.Value
			}
		}
	}

	return buildInfo
}
