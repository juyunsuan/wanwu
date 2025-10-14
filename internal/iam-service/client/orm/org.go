package orm

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/model"
	"github.com/UnicomAI/wanwu/internal/iam-service/client/orm/sqlopt"
	"github.com/UnicomAI/wanwu/internal/iam-service/config"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/gromitlee/access"
	"github.com/gromitlee/access/pkg/perm"
	"gorm.io/gorm"
)

func (c *Client) GetTopOrg(ctx context.Context) (uint32, error) {
	org := &model.Org{}
	err := sqlopt.WithParentID(0).Apply(c.db.WithContext(ctx)).First(org).Error
	return org.ID, err
}

func (c *Client) GetOrg(ctx context.Context, orgID uint32) (*OrgInfo, *errs.Status) {
	var ret *OrgInfo
	var err error
	return ret, c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		org := &model.Org{}
		if err = sqlopt.WithID(orgID).Apply(tx).First(org).Error; err != nil {
			return toErrStatus("iam_org_get", err.Error())
		}
		ret, err = toOrgInfoTx(tx, org)
		if err != nil {
			return toErrStatus("iam_org_get", err.Error())
		}
		return nil
	})
}

func (c *Client) GetOrgs(ctx context.Context, parentID uint32, name string, offset, limit int32) ([]*OrgInfo, int64, *errs.Status) {
	var ret []*OrgInfo
	var count int64
	return ret, count, c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		var orgs []*model.Org
		orgsQuery := sqlopt.WithParentID(parentID).Apply(tx).Select("id").Table("orgs")
		if err := sqlopt.LikeName(name).Apply(tx).Where("id IN (?)", orgsQuery).
			Offset(int(offset)).Limit(int(limit)).Order("id DESC").Find(&orgs).
			Offset(-1).Limit(-1).Count(&count).Error; err != nil {
			return toErrStatus("iam_orgs_get", err.Error())
		}
		for _, org := range orgs {
			info, err := toOrgInfoTx(tx, org)
			if err != nil {
				return toErrStatus("iam_orgs_get", err.Error())
			}
			ret = append(ret, info)
		}
		return nil
	})
}

func (c *Client) SelectOrgs(ctx context.Context, userID uint32) ([]IDName, *errs.Status) {
	var ret []IDName
	var orgTree *model.OrgNode
	var err error
	return ret, c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// org tree
		orgTree, err = getOrgTree(tx)
		if err != nil {
			return toErrStatus("iam_orgs_select", err.Error())
		}
		ret, err = selectOrgs(tx, userID, orgTree)
		if err != nil {
			return toErrStatus("iam_orgs_select", err.Error())
		}
		return nil
	})

}

func (c *Client) CreateOrg(ctx context.Context, org *model.Org) (uint32, *errs.Status) {
	if org.ID != 0 {
		return 0, toErrStatus("iam_org_create", "create org but id err")
	}
	return org.ID, c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		return createOrgTx(tx, org)
	})
}

func createOrgTx(tx *gorm.DB, org *model.Org) *errs.Status {
	var roleName string
	// check parents
	if org.ParentID != 0 {
		// 正常创建组织
		if err := sqlopt.WithID(org.ParentID).Apply(tx).First(&model.Org{}).Error; err != nil {
			return toErrStatus("iam_org_create", err.Error())
		}
		// check creator
		if err := sqlopt.WithID(org.CreatorID).Apply(tx).First(&model.User{}).Error; err != nil {
			return toErrStatus("iam_org_create", err.Error())
		}
		// check name 组织名在上级组织的所有下级组织内唯一
		if err := sqlopt.SQLOptions(
			sqlopt.WithParentID(org.ParentID),
			sqlopt.WithName(org.Name),
		).Apply(tx).First(&model.Org{}).Error; err != gorm.ErrRecordNotFound {
			if err == nil {
				err = errors.New("already exist")
			}
			return toErrStatus("iam_org_create", err.Error())
		}
		roleName = "组织管理员"
	} else {
		// 创建系统内唯一顶级组织，此时系统内不能存在任何组织
		if err := tx.First(&model.Org{}).Error; err != gorm.ErrRecordNotFound {
			if err == nil {
				err = errors.New("already exist")
			}
			return toErrStatus("iam_org_create", err.Error())
		}
		// check creator
		if org.CreatorID != 0 {
			return toErrStatus("iam_org_create", "create top org but creator not empty")
		}
		roleName = "超级管理员"
	}
	// create org
	if err := tx.Create(org).Error; err != nil {
		return toErrStatus("iam_org_create", err.Error())
	}
	// create role
	roleID, err := createRole(tx, org.ID, org.CreatorID, roleName, "", true, nil)
	if err != nil {
		return toErrStatus("iam_org_create", err.Error())
	}
	if org.CreatorID != 0 {
		// create org user
		if err := tx.Create(&model.OrgUser{
			OrgID:  org.ID,
			UserID: org.CreatorID,
		}).Error; err != nil {
			return toErrStatus("iam_org_create", err.Error())
		}
		// create user role
		if err := tx.Create(&model.UserRole{
			OrgID:   org.ID,
			UserID:  org.CreatorID,
			RoleID:  roleID,
			IsAdmin: true,
		}).Error; err != nil {
			return toErrStatus("iam_org_create", err.Error())
		}
	}
	return nil
}

func (c *Client) UpdateOrg(ctx context.Context, org *model.Org) *errs.Status {
	if org.ID == 0 {
		return toErrStatus("iam_org_update", "update org but id err")
	}
	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// check parent
		if err := sqlopt.SQLOptions(
			sqlopt.WithID(org.ID),
			sqlopt.WithParentID(org.ParentID),
		).Apply(tx).First(&model.Org{}).Error; err != nil {
			return toErrStatus("iam_org_update", err.Error())
		}
		// check name
		var orgs []*model.Org
		if err := sqlopt.SQLOptions(
			sqlopt.WithParentID(org.ParentID),
			sqlopt.WithName(org.Name),
		).Apply(tx).Find(&orgs).Error; err != nil {
			return toErrStatus("iam_org_update", err.Error())
		}
		if len(orgs) > 0 {
			for _, o := range orgs {
				if o.ID != org.ID {
					return toErrStatus("iam_org_update", fmt.Sprintf("check name %v but already exist", org.Name))
				}
			}
		}
		// update org
		if err := tx.Model(org).Updates(map[string]interface{}{
			"name":   org.Name,
			"remark": org.Remark,
		}).Error; err != nil {
			return toErrStatus("iam_org_update", err.Error())
		}
		return nil
	})
}

func (c *Client) DeleteOrg(ctx context.Context, orgID uint32) *errs.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// check org
		if orgID == config.TopOrgID() {
			return toErrStatus("iam_org_delete", "cannot delete top org")
		}
		// org tree
		orgTree, err := getOrgTree(tx)
		if err != nil {
			return toErrStatus("iam_org_delete", err.Error())
		}
		// delete org
		if err := deleteOrg(tx, orgID, orgTree); err != nil {
			return toErrStatus("iam_org_delete", err.Error())
		}
		return nil
	})
}

func deleteOrg(tx *gorm.DB, orgID uint32, orgTree *model.OrgNode) error {
	// delete sub org
	for _, sub := range orgTree.GetSubs(orgID) {
		if err := deleteOrg(tx, sub.GetOrgID(), orgTree); err != nil {
			return fmt.Errorf("delete org %v err: %v", orgID, err)
		}
	}
	// delete user role
	if err := sqlopt.WithOrgID(orgID).Apply(tx).Delete(&model.UserRole{}).Error; err != nil {
		return fmt.Errorf("delete user role err: %v", err)
	}
	// delete org user
	if err := sqlopt.WithOrgID(orgID).Apply(tx).Delete(&model.OrgUser{}).Error; err != nil {
		return fmt.Errorf("delete org user err: %v", err)
	}
	// delete org role
	var orgRoles []*model.OrgRole
	if err := sqlopt.WithOrgID(orgID).Apply(tx).Find(&orgRoles).Error; err != nil {
		return fmt.Errorf("get org role err: %v", err)
	}
	if err := sqlopt.WithOrgID(orgID).Apply(tx).Delete(&model.OrgRole{}).Error; err != nil {
		return fmt.Errorf("delete org role err: %v", err)
	}
	// delete role
	for _, orgRole := range orgRoles {
		if err := access.RBAC0DeleteRole(tx, perm.Role(orgRole.RoleID)); err != nil {
			return fmt.Errorf("delete role %v err: %v", orgRole.OrgID, err)
		}
	}
	// delete org
	if err := sqlopt.WithID(orgID).Apply(tx).Delete(&model.Org{}).Error; err != nil {
		return fmt.Errorf("delete org %v err: %v", orgID, err)
	}
	return nil
}

func (c *Client) ChangeOrgStatus(ctx context.Context, orgID uint32, status bool) *errs.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// check org
		if orgID == config.TopOrgID() {
			return toErrStatus("iam_org_change_status", "cannot change top org status")
		}
		// change status
		if err := changeOrgStatus(tx, orgID, status); err != nil {
			return toErrStatus("iam_org_change_status", err.Error())
		}
		return nil
	})
}

func changeOrgStatus(tx *gorm.DB, orgID uint32, status bool) error {
	if status {
		// enable: check parent
		org := &model.Org{}
		if err := sqlopt.WithID(orgID).Apply(tx).First(org).Error; err != nil {
			return fmt.Errorf("change org %v status %v get org err: %v", orgID, status, err)
		}
		parent := &model.Org{}
		if err := sqlopt.WithID(org.ParentID).Apply(tx).First(parent).Error; err != nil {
			return fmt.Errorf("change org %v status %v get parent %v err: %v", orgID, status, org.ParentID, err)
		}
		if !parent.Status {
			return fmt.Errorf("change org %v status %v but parent %v status false", orgID, status, org.ParentID)
		}
	} else {
		// disable: change sub orgs status
		var orgs []*model.Org
		if err := sqlopt.WithParentID(orgID).Apply(tx).Find(&orgs).Error; err != nil {
			return fmt.Errorf("change org %v status %v get sub orgs err: %v", orgID, status, err)
		}
		// change subs
		for _, sub := range orgs {
			if err := changeOrgStatus(tx, sub.ID, status); err != nil {
				return err
			}
		}
	}
	// change status
	if err := sqlopt.WithID(orgID).Apply(tx).Model(&model.Org{}).Updates(map[string]interface{}{
		"status": status,
	}).Error; err != nil {
		return fmt.Errorf("change org %v status %v err: %v", orgID, status, err)
	}
	return nil
}

func (c *Client) AddOrgUser(ctx context.Context, orgID, userID, roleID uint32) *errs.Status {

	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// check org user
		if err := sqlopt.SQLOptions(
			sqlopt.WithOrgID(orgID),
			sqlopt.WithUserID(userID),
		).Apply(tx).First(&model.OrgUser{}).Error; err != gorm.ErrRecordNotFound {
			if err == nil {
				err = errors.New("already exist")
			}
			return toErrStatus("iam_org_user_add", util.Int2Str(orgID), util.Int2Str(userID), util.Int2Str(roleID), err.Error())
		}
		// check org role
		orgRole := &model.OrgRole{}
		if roleID != 0 {
			if err := sqlopt.SQLOptions(
				sqlopt.WithOrgID(orgID),
				sqlopt.WithRoleID(roleID),
			).Apply(tx).First(orgRole).Error; err != nil {
				return toErrStatus("iam_org_user_add", util.Int2Str(orgID),
					util.Int2Str(userID), util.Int2Str(roleID), err.Error())
			}
		}
		// create org user
		if err := tx.Create(&model.OrgUser{
			OrgID:  orgID,
			UserID: userID,
		}).Error; err != nil {
			return toErrStatus("iam_org_user_add", util.Int2Str(orgID),
				util.Int2Str(userID), util.Int2Str(roleID), err.Error())
		}
		// create user role
		if roleID != 0 {
			if err := tx.Create(&model.UserRole{
				OrgID:   orgID,
				UserID:  userID,
				RoleID:  roleID,
				IsAdmin: orgRole.IsAdmin,
			}).Error; err != nil {
				return toErrStatus("iam_org_user_add", util.Int2Str(orgID),
					util.Int2Str(userID), util.Int2Str(roleID), err.Error())
			}
		}
		return nil
	})
}

func (c *Client) RemoveOrgUser(ctx context.Context, orgID, userID uint32) *errs.Status {
	return c.transaction(ctx, func(tx *gorm.DB) *errs.Status {
		// check org
		if orgID == config.TopOrgID() {
			return toErrStatus("iam_org_user_remove_top")
		}
		// delete user role
		if err := sqlopt.SQLOptions(
			sqlopt.WithOrgID(orgID),
			sqlopt.WithUserID(userID),
		).Apply(tx).Delete(&model.UserRole{}).Error; err != nil {
			return toErrStatus("iam_org_user_remove", util.Int2Str(orgID),
				util.Int2Str(userID), err.Error())
		}
		// delete org user
		if err := sqlopt.SQLOptions(
			sqlopt.WithOrgID(orgID),
			sqlopt.WithUserID(userID),
		).Apply(tx).Delete(&model.OrgUser{}).Error; err != nil {
			return toErrStatus("iam_org_user_remove", util.Int2Str(orgID),
				util.Int2Str(userID), err.Error())
		}
		return nil
	})
}

// --- internal function ---

func toOrgInfoTx(tx *gorm.DB, org *model.Org) (*OrgInfo, error) {
	ret := &OrgInfo{
		ID:        org.ID,
		Name:      org.Name,
		Remark:    org.Remark,
		Status:    org.Status,
		CreatedAt: org.CreatedAt,
	}
	// creator
	if org.CreatorID != 0 {
		creator, err := getCreatorTx(tx, org.CreatorID)
		if err != nil {
			return nil, err
		}
		ret.Creator = creator
	}
	return ret, nil
}

func getOrgTree(tx *gorm.DB) (*model.OrgNode, error) {
	// all org
	var orgs []*model.Org
	if err := tx.Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("get org tree all org err: %v", err)
	}
	// all org admin role
	var orgAdmins []*model.OrgRole
	if err := sqlopt.WithAdmin(true).Apply(tx).Find(&orgAdmins).Error; err != nil {
		return nil, fmt.Errorf("get org tree all org admin role err: %v", err)
	}
	return model.NewOrgTree(orgs, orgAdmins)
}

func selectOrgs(tx *gorm.DB, userID uint32, orgTree *model.OrgNode) ([]IDName, error) {
	// user role
	var userRoles []*model.UserRole
	orgRolesQuery := sqlopt.WithStatus(true).Apply(tx).Select("role_id").Table("org_roles")
	if err := sqlopt.WithUserID(userID).Apply(tx).Where("role_id IN (?)", orgRolesQuery).Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("get user role err: %v", err)
	}
	// user org
	var userOrgs []*model.OrgUser
	if err := sqlopt.SQLOptions(
		sqlopt.WithUserID(userID),
	).Apply(tx).Find(&userOrgs).Error; err != nil {
		return nil, fmt.Errorf("get org user err: %v", err)
	}
	// select org
	var ret []IDName
	for _, org := range orgTree.Select(userOrgs, userRoles) {
		ret = append(ret, IDName{ID: org.ID, Name: org.Name})
	}
	return ret, nil
}
