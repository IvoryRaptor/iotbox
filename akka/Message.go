package akka

type Message map[string]interface{}

func (message *Message) GetInt(name string) int {
	var result = (*message)[name]
	if result == nil {
		return 0
	}
	return (*message)[name].(int)
}

func (message *Message) GetString(name string) (string, bool) {
	var result = (*message)[name]
	if result == nil {
		return "", false
	}
	return (*message)[name].(string), true
}

func (message *Message) SetInt(name string, value int) {
	(*message)[name] = value
}
