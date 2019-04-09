package akka

type block struct {
	owner   IActor
	message Message
}
