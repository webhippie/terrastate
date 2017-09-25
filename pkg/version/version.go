package version

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

var (
	// VersionMajor is the current major version
	VersionMajor int64

	// VersionMinor is the current minor version
	VersionMinor int64 = 1

	// VersionPatch is the current patch version
	VersionPatch int64

	// VersionPre indicates a pre release tag
	VersionPre = "alpha1"

	// VersionDev indicates the current commit
	VersionDev = "0000000"

	// VersionDate indicates the build date
	VersionDate = "20170101"

	// Version is the version of the current implementation.
	Version = semver.Version{
		Major:      VersionMajor,
		Minor:      VersionMinor,
		Patch:      VersionPatch,
		PreRelease: semver.PreRelease(VersionPre),
		Metadata:   fmt.Sprintf("git%s.%s", VersionDate, VersionDev),
	}
)
