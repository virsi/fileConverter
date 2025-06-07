package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK       = "ok"
	StatusError    = "error"
	StatusNotFound = "not_found"
	StatusCreated  = "created"
	StatusBadRequest = "bad_request"
	StatusInternalServerError = "internal_server_error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}
func NotFound() Response {
	return Response{
		Status: StatusNotFound,
	}
}
func Created() Response {
	return Response{
		Status: StatusCreated,
	}
}
func BadRequest(err string) Response {
	return Response{
		Status: StatusBadRequest,
		Error:  err,
	}
}
func InternalServerError(err string) Response {
	return Response{
		Status: StatusInternalServerError,
		Error:  err,
	}
}
