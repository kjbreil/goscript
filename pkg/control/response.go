package control

type Responses chan Response

type Response struct {
	To   string
	From string

	Module any
}
