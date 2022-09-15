package core

import (
    "errors"
    "regexp"
    "strconv"
    "strings"
)

type Version struct {
    Prefix      string // Prefix: can be `U` or `V`
    Version     string // Version: semantic version with _
    Description string // Description: description of the migration
    Suffix      string // Suffix: file extension
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

/*
CompareVersionStr If left is smaller than the right one, then true.
    If both are the same, then true.

Example:
    CompareVersionStr("1_0", "0_1") -> false
    CompareVersionStr("1_0", "1_10") -> true
    CompareVersionStr("1_0_5", "1_10") -> true
*/
func CompareVersionStr(left string, right string) bool {
    l := strings.Split(left, "_")
    r := strings.Split(right, "_")

    minLen := len(l)

    if len(r) < minLen {
        minLen = len(r)
    }

    for i := 0; i < minLen; i++ {
        f, _ := strconv.Atoi(l[i])
        s, _ := strconv.Atoi(r[i])

        if f < s {
            return true
        } else if f > s {
            return false
        }
    }

    return true
}

func EqualsVersionShort(leftVersionShort *VersionShort, rightVersionShort *VersionShort) bool {
    if leftVersionShort.Version == rightVersionShort.Version && leftVersionShort.Prefix == rightVersionShort.Prefix {
        return true
    }
    return false
}

func EqualsVersionFullAndShort(leftVersionShort *VersionShort, rightVersionShort *Version) bool {
    if leftVersionShort.Version == rightVersionShort.Version && leftVersionShort.Prefix == rightVersionShort.Prefix {
        return true
    }
    return false
}

func GetVersionShortFromFull(version *Version) *VersionShort {
    return &VersionShort{
        Prefix:  version.Prefix,
        Version: version.Version,
    }
}
