package clients

import "runtime"

// Since we can build CLI for multiple platforms, each client should handle the platform data separately, using build info.
// Possible variations include (according to build.yml action):
// a) OS:
// - 'linux'
// - 'darwin'
// - 'freebsd'
//
// b) ARCH:
// - amd64
// - arm64
// - arm
// - 386
// Please bear in mind that some clients don't support that many platforms.
// CLI will try to adjust the missing os/arch to some closest available setup.

type buildInfo struct {
	os   osBuildInfo
	arch archBuildInfo
}

type (
	osBuildInfo   map[string]string
	archBuildInfo map[string]string
)

var gethBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `linux`,
		`fallback`: `linux`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64`,
		`arm`:      `arm5`,
		`386`:      `386`,
		`fallback`: `amd64`,
	},
}

var erigonBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `linux`,
		`fallback`: `linux`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64`,
		`arm`:      `arm64`,
		`386`:      `amd64`,
		`fallback`: `amd64`,
	},
}

var nethermindBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `macos`,
		`fallback`: `linux`,
	},
	arch: archBuildInfo{
		`amd64`:    `x64`,
		`arm64`:    `arm64`,
		`arm`:      `x64`,
		`386`:      `x64`,
		`fallback`: `x64`,
	},
}

// Since besu runs on JVM there is no need for platform.
var besuBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `darwin`,
		`fallback`: `fallback`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64`,
		`arm`:      `arm`,
		`386`:      `386`,
		`fallback`: `fallback`,
	},
}

var prysmBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `darwin`,
		`fallback`: `linux`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64`,
		`arm`:      `arm64`,
		`386`:      `amd64`,
		`fallback`: `amd64`,
	},
}

var lighthouseBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `unknown-linux`,
		`darwin`:   `apple-darwin`,
		`fallback`: `fallback`,
	},
	arch: archBuildInfo{
		`amd64`:    `x86_64`,
		`arm64`:    `aarch64`,
		`arm`:      `aarch64`,
		`386`:      `x86_64`,
		`fallback`: `x86_64`,
	},
}

// Since teku runs on JVM there is no need for platform.
var tekuBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `darwin`,
		`fallback`: `fallback`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64`,
		`arm`:      `arm`,
		`386`:      `386`,
		`fallback`: `fallback`,
	},
}

var nimbus2BuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `Linux`,
		`darwin`:   `macOS`,
		`fallback`: `fallback`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64v8`,
		`arm`:      `arm32v7`,
		`386`:      `amd64`,
		`fallback`: `amd64`,
	},
}

var jdkBuildInfo = buildInfo{
	os: osBuildInfo{
		`linux`:    `linux`,
		`darwin`:   `darwin`,
		`fallback`: `fallback`,
	},
	arch: archBuildInfo{
		`amd64`:    `amd64`,
		`arm64`:    `arm64`,
		`arm`:      `arm`,
		`386`:      `386`,
		`fallback`: `fallback`,
	},
}

func (b buildInfo) Os() string {
	return b.os[runtime.GOOS]
}

func (b buildInfo) Arch() string {
	return b.arch[runtime.GOARCH]
}
