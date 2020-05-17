package natsEncoder

type Encoder interface {
	Encode(v interface{}) error
}