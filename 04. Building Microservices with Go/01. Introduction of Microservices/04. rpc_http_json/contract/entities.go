package contract

type HelloWorldRequest struct {
	Name string `json:"name"`
}

type HelloWorldResponse struct {
	Message string `json:"message"`
}