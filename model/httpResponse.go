package model

// HTTPResponse type
type HTTPResponse struct {
	Code    int         `json:"code" xml:"code"`
	Message string      `json:"message" xml:"message"`
	Data    interface{} `json:"data" xml:"data"`
}

// GetHTTPResponse : Retourne le type HTTPResponse
func GetHTTPResponse(code int, message string, data interface{}) HTTPResponse {
	response := HTTPResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	return response
}