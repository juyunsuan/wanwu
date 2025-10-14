package request

type UserCreate struct {
	Username string `json:"username" validate:"required"` // 用户名
	UserInfo
}

func (u *UserCreate) Check() error {
	return nil
}

type UserUpdate struct {
	UserID string `json:"userId" validate:"required"` // 用户ID
	UserInfo
}

func (u *UserUpdate) Check() error {
	return nil
}

type UserInfo struct {
	Nickname string   `json:"nickname"`                     // 昵称
	Password string   `json:"password" validate:"required"` // 密码
	Phone    string   `json:"phone" validate:"required"`    // 电话
	Remark   string   `json:"remark"`                       // 备注
	Gender   string   `json:"gender"`                       // 性别（0-女，1-男，空-未知）
	Company  string   `json:"company"`                      // 公司
	RoleIDs  []string `json:"roleIds" validate:"max=1"`     // 角色列表
}

type UserID struct {
	UserID string `json:"userId" validate:"required"` // 用户ID
}

func (u *UserID) Check() error {
	return nil
}

type UserStatus struct {
	UserID
	Status bool `json:"status"`
}

func (u *UserStatus) Check() error {
	return nil
}

type UserPassword struct {
	UserID
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (u *UserPassword) Check() error {
	return nil
}

type UserPasswordByAdmin struct {
	UserID
	Password string `json:"password" validate:"required"`
}

func (u *UserPasswordByAdmin) Check() error {
	return nil
}
