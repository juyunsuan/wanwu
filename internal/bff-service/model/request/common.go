package request

type CommonCheck struct {
}

type PageSearch struct {
	PageSize int `json:"pageSize" form:"pageSize" validate:"required"`
	PageNo   int `json:"pageNo" form:"pageNo"`
}

type LoginEmailCheck struct {
	Email string `json:"email" validate:"required"` // 邮箱
	Code  string `json:"code" validate:"required"`  // 邮箱验证码
}

type ChangeUserPasswordByEmail struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
	Email       string `json:"email" validate:"required"` // 邮箱
	Code        string `json:"code" validate:"required"`  // 邮箱验证码
}

func (l *LoginEmailCheck) Check() error {
	return nil
}

func (c *ChangeUserPasswordByEmail) Check() error {
	return nil
}
func (c *CommonCheck) Check() error {
	return nil
}
