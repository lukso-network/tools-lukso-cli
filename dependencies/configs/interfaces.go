package configs

type ClientConfigDependency interface {
	Install(isUpdate bool) error
	Name() string
}
