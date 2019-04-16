package common

import "time"

type Protocol interface {
	Parse(data []byte) Message
	Packet(message Message) []byte
}

type Port interface {
	Read(wait time.Duration) (msg Message, err error)
	Write(message Message) error
	Close() error
}
