package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// https://regexr.com/89qep

const SemverExpressionRaw = `v?\d+(\.\d+){2}`

var semverExpr = regexp.MustCompile(fmt.Sprintf(`^%s$`, SemverExpressionRaw))

type Semver struct {
	major int
	minor int
	patch int
}

func StringToSemver(str string) (semver Semver) {
	semver = Semver{
		major: 0,
		minor: 0,
		patch: 0,
	}

	matches := semverExpr.FindAllString(str, -1)
	if len(matches) != 1 {
		return
	}

	// already validated by the regex, conversions and slice accesses won't error/panic
	v := strings.Split(strings.Trim(matches[0], "v"), ".")
	major, _ := strconv.Atoi(v[0])
	minor, _ := strconv.Atoi(v[1])
	patch, _ := strconv.Atoi(v[2])

	semver = Semver{
		major: major,
		minor: minor,
		patch: patch,
	}

	return
}

// IsNewerThan checks whether v1 is a newer version than v2.
// The check is strict (equal values return false)
func (v Semver) IsNewerThan(v2 Semver) (isNewer bool) {
	switch {
	case v.major > v2.major:
		return true
	case v.major < v2.major:
		return false
	case v.minor > v2.minor:
		return true
	case v.minor < v2.minor:
		return false
	case v.patch > v2.patch:
		return true
	case v.patch < v2.patch:
		return false
	}

	return false
}
