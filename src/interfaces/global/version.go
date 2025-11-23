// Package global
package global

import (
	"chart-provider/src/utils"
	"fmt"
	"strings"
)

type CheckVersionResult int

const (
	AllMatch CheckVersionResult = iota
	MajorUnmatch
	MinorUnmatch
	PatchUnmatch
)

type Version struct {
	major   int
	minor   int
	patch   int
	version string
}

func NewVersion(version string) (*Version, error) {
	versions := strings.Split(version, ".")
	if len(versions) < 3 {
		return nil, fmt.Errorf("invalid version String, %s", version)
	}
	return &Version{
		major:   utils.StrToInt(versions[0], 0),
		minor:   utils.StrToInt(versions[1], 0),
		patch:   utils.StrToInt(versions[2], 0),
		version: version,
	}, nil
}

func (v *Version) CheckVersion(version *Version) CheckVersionResult {
	if v.major != version.major {
		return MajorUnmatch
	}
	if v.minor != version.minor {
		return MinorUnmatch
	}
	if v.patch != version.patch {
		return PatchUnmatch
	}
	return AllMatch
}

func (v *Version) String() string {
	return v.version
}
