package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	version   = "dev"
	commit    = "unknown"
	date      = "unknown"
	goVersion = runtime.Version()
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"goVersion"`
}

func Get() string {
	return version
}

func GetCommit() string {
	return commit
}

func GetDate() string {
	return date
}

func GetGoVersion() string {
	return goVersion
}

func GetInfo() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		GoVersion: goVersion,
	}
}

func (i Info) String() string {
	return fmt.Sprintf("notedown-language-server %s (commit: %s, built: %s, go: %s)", 
		i.Version, i.Commit, i.Date, i.GoVersion)
}

func JSON() (string, error) {
	info := GetInfo()
	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(data), nil
}