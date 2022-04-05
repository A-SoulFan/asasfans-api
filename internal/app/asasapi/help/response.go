package help

type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TTL     int         `json:"ttl"`
	Data    interface{} `json:"data"`
}

func SuccessJson(data interface{}) *JsonResponse {
	return &JsonResponse{
		Code:    0,
		Message: "ok",
		TTL:     1,
		Data:    data,
	}
}

func FailureJson(code int, msg string, data interface{}) *JsonResponse {
	return &JsonResponse{
		Code:    code,
		Message: msg,
		TTL:     1,
		Data:    data,
	}
}
