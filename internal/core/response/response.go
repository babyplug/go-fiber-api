package response

type ResponseDTO struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Data    any    `json:"data"`
}
