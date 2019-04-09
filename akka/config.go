package akka

type IConfig interface {
	Config(config map[string]interface{}) error
}
