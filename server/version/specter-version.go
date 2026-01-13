package version

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/k0kubun/pp"
)

func init() {
	//buildSpecterInfo()
}

var (
	Version    string
	Commit     string
	Dirty      string
	CompiledAt string
	OS         string
	Arch       string
	GoVersion  string
)

func BuildSpecterInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		useDefaults()
		return
	}

	if GoVersion == "" {
		GoVersion = info.GoVersion
	}

	setting := buildSettings(info.Settings)

	fmt.Println(pp.Sprintf("%v", setting))
}

func useDefaults() {
	if Version == "" {
		Version = "develop"
	}

	if GoVersion == "" {
		GoVersion = runtime.Version()
	}

}

func buildSettings(settings []debug.BuildSetting) map[string]string {
	values := make(map[string]string, len(settings))
	for _, setting := range settings {
		values[setting.Key] = setting.Value
	}
	return values
}
