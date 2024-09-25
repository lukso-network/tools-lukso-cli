package types

type CLConfig struct {
	ShapellaEpoch uint64 `yaml:"CAPELLA_FORK_EPOCH"`
	DencunEpoch   uint64 `yaml:"DENEB_FORK_EPOCH"`
}
