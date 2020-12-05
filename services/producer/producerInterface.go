package producer

type Producer interface {
	Produce(message []byte, channelName string) error
}
