package idl

type SendEmailVerifyCodeReq struct {
	Email string `json:"email" binding:"required,email"`
}

type EmailRegisterReq struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=4,max=32"`
	VerifyCode string `json:"verify_code" binding:"required,min=4,max=12"`
}

type EmailPasswordSignInReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=32"`
}

type RestPasswordReq struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=4,max=32"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	VerifyCode string `json:"verify_code" binding:"required,min=4,max=12"`
}

type SignInResp struct {
	Token string `json:"token"`
}
