package akka

type Message map[string]interface{}

func (message *Message) GetInt(name string) int {
	var result = (*message)[name]
	if result == nil {
		return 0
	}
	return (*message)[name].(int)
}

func (message *Message) SetInt(name string, value int) {
	(*message)[name] = value
}

type block struct {
	owner   IActor
	message Message
}
