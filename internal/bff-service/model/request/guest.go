package request

// Login 登录请求的参数
type Login struct {
	Username string `json:"username" validate:"required"` // 用户名
	Password string `json:"password" validate:"required"` // 密码
	Key      string `json:"key" validate:"required"`      // 客户端key
	Code     string `json:"code" validate:"required"`     // 验证码
}

type RegisterByEmail struct {
	Username string `json:"username" validate:"required"` // 用户名
	Email    string `json:"email" validate:"required"`    // 邮箱
	Code     string `json:"code" validate:"required"`     // 邮箱验证码
}

type RegisterSendEmailCode struct {
	Username string `json:"username" validate:"required"` // 用户名
	Email    string `json:"email" validate:"required"`    // 邮箱
}

type ResetPasswordSendEmailCode struct {
	Email string `json:"email" validate:"required"` // 邮箱
}

type ResetPasswordByEmail struct {
	Email    string `json:"email" validate:"required"`    // 邮箱
	Code     string `json:"code" validate:"required"`     // 邮箱验证码
	Password string `json:"password" validate:"required"` // 密码
}

type LoginSendEmailCode struct {
	Email string `json:"email" validate:"required"` // 邮箱
}

func (l *Login) Check() error {
	return nil
}

func (l *RegisterByEmail) Check() error {
	return nil
}

func (l *RegisterSendEmailCode) Check() error {
	return nil
}

func (r *ResetPasswordSendEmailCode) Check() error {
	return nil
}

func (r *ResetPasswordByEmail) Check() error {
	return nil
}

func (l *LoginSendEmailCode) Check() error {
	return nil
}
