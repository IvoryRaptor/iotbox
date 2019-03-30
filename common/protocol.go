package common

type IProtocol interface {
	Parse(data []byte) *Packet
	Packaging(packet *Packet) []byte
}
