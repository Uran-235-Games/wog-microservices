package domain

// API Response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Errors  []ErrorResponse `json:"errors"`
	Result  any             `json:"result"`
	Success bool            `json:"success"`
}

func NewAPIResponse() *Response {
	return &Response{
		Errors:  []ErrorResponse{},
		Success: true,
		Result:  nil,
	}
}

func (r *Response) AddError(code int, message string) {
	r.Errors = append(r.Errors, ErrorResponse{Code: code, Message: message})
}

// ======================================================

// Request		/sign-up
type SignUpRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Response		/sign-up
type SignUpResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Request		/sign-in
type SignInRequest struct {
	Name     string `json:"name" binding:"omitempty,min=3"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Response		/sign-in
type SignInResponse struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
