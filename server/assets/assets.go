package assets

import (
	"os"
	"path/filepath"
)

const (
	SpecterWorkRootDirEnv = "SPECTER_WORK_ROOT_DIR_ENV"
	SpecterWorkDir        = ".specter"
)

func GetRootAppDir() string {
	value := os.Getenv(SpecterWorkRootDirEnv)
	var dir string

	if len(value) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dir = filepath.Join(wd, SpecterWorkDir)
	} else {
		dir = value
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			panic(err)
		}
	}
	return dir
}
