package orm

import (
	"context"
	"time"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/pkg/util"
	"gorm.io/gorm"
)

func (c *Client) Login(ctx context.Context, username, password, language string) (*UserInfo, *Permission, *errs.Status) {
	var userInfo *UserInfo
	var permission *Permission

	return userInfo, permission, c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// user
		user := &model.User{}
		if err := sqlopt.WithName(username).Apply(tx).First(user).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return toErrStatus("iam_login_err", err.Error())
			}
			// user by email
			if err := sqlopt.WithEmail(username).Apply(tx).First(user).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return toErrStatus("iam_login_err", err.Error())
				} else {
					return toErrStatus("iam_login_invalid_name_pwd")
				}
			}
		}
		if user.Password != util.SHA256(password) {
			return toErrStatus("iam_login_invalid_name_pwd")
		}
		// check status
		if !user.Status {
			return toErrStatus("iam_login_user_disable")
		}
		// org tree
		orgTree, err := getOrgTree(tx)
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		// user info
		userInfo, err = toUserInfoTx(tx, user, orgTree, true)
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		var orgID uint32
		if len(userInfo.Orgs) == 1 {
			orgID = userInfo.Orgs[0].Org.ID
		} else if len(userInfo.Orgs) > 1 {
			for _, org := range userInfo.Orgs {
				if !org.Org.Status {
					continue
				}
				if len(org.Roles) > 0 {
					orgID = org.Org.ID
				}
				break
			}
		}
		if orgID != 0 {
			permission, err = getUserPermission(tx, user.ID, orgID)
		}
		if err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		// last_login_at & last_exec_at
		nowTS := time.Now().UnixMilli()
		update := map[string]interface{}{
			"last_login_at": nowTS,
			"last_exec_at":  nowTS,
		}
		if language != "" {
			update["language"] = language
			userInfo.Language = language
		}
		if err := tx.Model(user).Updates(update).Error; err != nil {
			return toErrStatus("iam_login_err", err.Error())
		}
		return nil
	})
}
