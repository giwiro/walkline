package core

import (
	"errors"
	"regexp"
)

type Version struct {
	Prefix      string
	Version     string
	Description string
	Suffix      string
}

type VersionShort struct {
	Prefix  string
	Version string
}

var VersionRangeRegex = regexp.MustCompile("^([U|V][\\d_\\.]+?):([U|V][\\d_\\.]+?)$")
var VersionShortRegex = regexp.MustCompile("^([U|V])([\\d_\\.]+?)$")
var VersionRegex = regexp.MustCompile("([V|U])([\\d_\\.]+?)__(.+?)(\\.\\w{2,6})")

func ParseVersion(version string) (*Version, error) {
	result := VersionRegex.FindAllStringSubmatch(version, -1)

	if result == nil || len(result[0]) <= 1 {
		return nil, errors.New("bad format")
	}

	return &Version{
		Prefix:      result[0][1],
		Version:     result[0][2],
		Description: result[0][3],
		Suffix:      result[0][4],
	}, nil
}

func ParseVersionShort(versionShort string) (*VersionShort, error) {
	result := VersionShortRegex.FindAllStringSubmatch(versionShort, -1)

	if result == nil || len(result[0]) <= 1 {
		return nil, errors.New("bad version format")
	}

	return &VersionShort{
		Prefix:  result[0][1],
		Version: result[0][2],
	}, nil
}

func ParseVersionShortRange(versionShortRange string) (*VersionShort, *VersionShort, error) {
	result := VersionRangeRegex.FindAllStringSubmatch(versionShortRange, -1)

	if result == nil || len(result[0]) <= 1 {
		return nil, nil, errors.New("bad range format")
	}

	leftVersionShort, err := ParseVersionShort(result[0][1])

	if err != nil {
		return nil, nil, errors.New("bad left version format")
	}

	rightVersionShort, err := ParseVersionShort(result[0][2])

	if err != nil {
		return nil, nil, errors.New("bad right version format")
	}

	if leftVersionShort.Version > rightVersionShort.Version {
		return nil, nil, errors.New("left version must the lower version")
	}

	return leftVersionShort, rightVersionShort, nil
}

func CompareVersionShort(leftVersionShort *VersionShort, rightVersionShort *VersionShort) bool {
	if leftVersionShort.Version == rightVersionShort.Version && leftVersionShort.Prefix == rightVersionShort.Prefix {
		return true
	}
	return false
}

func CompareVersionFullAndShort(leftVersionShort *VersionShort, rightVersionShort *Version) bool {
	if leftVersionShort.Version == rightVersionShort.Version && leftVersionShort.Prefix == rightVersionShort.Prefix {
		return true
	}
	return false
}

func GetVersionShortFromFull(versionShort *Version) *VersionShort {
	return &VersionShort{
		Prefix: versionShort.Prefix,
		Version: versionShort.Version,
	}
}