package common

type Permission struct {
	Resource string `mapstructure:"resource"`
	Action   string `mapstructure:"action"`
}

type Permissions struct {
	List []Permission `mapstructure:"list"`
}
