package response

type Response struct {
	Data         interface{} `json:"data"`
	Status       int         `json:"status"`
	ErrorType    string      `json:"type"`
	ErrorMessage interface{} `json:"message"`
}
