package version

import (
	"runtime/debug"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

var (
	Major      = ""
	Minor      = ""
	Patch      = ""
	Commit     = "unknown"
	Dirty      = ""
	CompiledAt = ""
	OS         = ""
	Arch       = ""
)

func BuildSpecterInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return
	}
	zap.S().Debugf("%v", pp.Sprintf("%v", info))
}
