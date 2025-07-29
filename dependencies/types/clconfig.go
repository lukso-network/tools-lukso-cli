package types

// TODO: change pointers for optional wrappers that you can unmarshal into
type CLConfig struct {
	ShapellaEpoch *uint64 `yaml:"CAPELLA_FORK_EPOCH,omitempty"`
	DencunEpoch   *uint64 `yaml:"DENEB_FORK_EPOCH,omitempty"`
	PectraEpoch   *uint64 `yaml:"ELECTRA_FORK_EPOCH,omitempty"`
}
