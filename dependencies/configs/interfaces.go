package configs

type ClientConfigDependency interface {
	Install() error
	Name() string
}
