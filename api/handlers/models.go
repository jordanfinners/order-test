package handlers

// Request is a inbound http style request for a handler
type Request struct {
	Body        string
	QueryParams string
}

// Response is a outbound http style response for a handler
type Response struct {
	Status int
	Body   string
}
