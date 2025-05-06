package clients

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

var erigonBuildInfo = buildInfo{
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

var nethermindBuildInfo = buildInfo{
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

var lighthouseBuildInfo = buildInfo{
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
