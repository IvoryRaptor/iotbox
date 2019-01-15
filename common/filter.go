package common

type IFilter interface {
	Config(config map[string]interface{}) error
}
