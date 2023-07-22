package controllers

type JsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var res JsonResponse

func (r *JsonResponse) SuccessResponse(message string, data interface{}) JsonResponse {
	r.Status = "success"
	r.Message = message
	r.Data = data
	return *r
}

func (r *JsonResponse) ErrorResponse(message string) JsonResponse {
	r.Status = "error"
	r.Message = message
	return *r
}

func (r *JsonResponse) NotFoundResponse() JsonResponse {
	r.Status = "error"
	r.Message = "Not found"
	return *r
}

func (r *JsonResponse) UnauthorizedResponse() JsonResponse {
	r.Status = "error"
	r.Message = "Unauthorized"
	return *r
}

func (r *JsonResponse) BadRequestResponse() JsonResponse {
	r.Status = "error"
	r.Message = "Bad request"
	return *r
}

func (r *JsonResponse) InternalServerErrorResponse() JsonResponse {
	r.Status = "error"
	r.Message = "Internal server error"
	return *r
}

func (r *JsonResponse) ValidationErrorResponse(err error) JsonResponse {
	r.Status = "error"
	r.Message = "Validation error"
	r.Data = err.Error()
	return *r
}

func (r *JsonResponse) CustomErrorResponse(message string) JsonResponse {
	r.Status = "error"
	r.Message = message
	return *r
}

func (r *JsonResponse) CustomSuccessResponse(message string, data interface{}) JsonResponse {
	r.Status = "success"
	r.Message = message
	r.Data = data
	return *r
}
