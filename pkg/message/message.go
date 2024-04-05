package message

type ResponseChan chan Response

type Response struct {
	To   string
	From string

	// GobType []byte
}
