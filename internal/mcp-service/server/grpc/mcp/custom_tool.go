package mcp

import (
	"context"
	"strings"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/mcp-service/client/model"
	"github.com/UnicomAI/wanwu/internal/mcp-service/config"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreateCustomTool(ctx context.Context, req *mcp_service.CreateCustomToolReq) (*emptypb.Empty, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", "identity is empty"))
	}
	if req.ApiAuth == nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, toErrStatus("mcp_create_custom_tool_err", "apiAuth is empty"))
	}
	if err := s.cli.CreateCustomTool(ctx, &model.CustomTool{
		Name:             req.Name,
		Description:      req.Description,
		Schema:           req.Schema,
		PrivacyPolicy:    req.PrivacyPolicy,
		Type:             req.ApiAuth.Type,
		APIKey:           req.ApiAuth.ApiKey,
		AuthType:         req.ApiAuth.AuthType,
		CustomHeaderName: req.ApiAuth.CustomHeaderName,
		UserID:           req.Identity.UserId,
		OrgID:            req.Identity.OrgId,
		ToolSquareId:     req.ToolSquareId,
	}); err != nil {
		return nil, errStatus(errs.Code_MCPCreateCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetCustomToolInfo(ctx context.Context, req *mcp_service.GetCustomToolInfoReq) (*mcp_service.GetCustomToolInfoResp, error) {
	if req.CustomToolId == "" && req.ToolSquareId == "" {
		return nil, errStatus(errs.Code_MCPGetCustomToolInfoErr, toErrStatus("mcp_get_custom_tool_info_err", "customToolId and toolSquareId are empty"))
	}
	// 自定义工具详情（智能体、前端调用）
	if req.CustomToolId != "" && req.ToolSquareId == "" {
		info, err := s.cli.GetCustomTool(ctx, &model.CustomTool{
			ID: util.MustU32(req.CustomToolId),
		})
		if err != nil {
			return nil, errStatus(errs.Code_MCPGetCustomToolInfoErr, err)
		}
		return &mcp_service.GetCustomToolInfoResp{
			CustomToolId:  util.Int2Str(info.ID),
			Name:          info.Name,
			Description:   info.Description,
			Schema:        info.Schema,
			PrivacyPolicy: info.PrivacyPolicy,
			ApiAuth: &mcp_service.ApiAuthWebRequest{
				Type:             info.Type,
				ApiKey:           info.APIKey,
				AuthType:         info.AuthType,
				CustomHeaderName: info.CustomHeaderName,
			},
		}, nil
	}

	// 内置工具详情（智能体调用）
	if req.CustomToolId == "" && req.ToolSquareId != "" && req.Identity != nil {
		toolCfg, exist := config.Cfg().Tool(req.ToolSquareId)
		if !exist {
			return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", "toolSquareId not exist"))
		}
		toolInfo := &mcp_service.GetCustomToolInfoResp{
			ToolSquareId: toolCfg.ToolSquareId,
			Name:         toolCfg.Name,
			Description:  toolCfg.Desc,
			Schema:       toolCfg.Schema,
		}
		if toolCfg.NeedApiKeyInput {
			info, _ := s.cli.GetCustomTool(ctx, &model.CustomTool{
				ToolSquareId: req.ToolSquareId,
				UserID:       req.Identity.UserId,
				OrgID:        req.Identity.OrgId,
			})
			if info != nil {
				toolInfo.CustomToolId = util.Int2Str(info.ID)
				toolInfo.ApiAuth = &mcp_service.ApiAuthWebRequest{
					Type:             info.Type,
					ApiKey:           info.APIKey,
					AuthType:         info.AuthType,
					CustomHeaderName: info.CustomHeaderName,
				}
			} else {
				//如果没配置过，返回配置数据。
				toolInfo.ApiAuth = &mcp_service.ApiAuthWebRequest{
					Type:             toolCfg.Type,
					ApiKey:           toolCfg.ApiKey,
					AuthType:         toolCfg.AuthType,
					CustomHeaderName: toolCfg.CustomHeaderName,
				}
			}
			return toolInfo, nil
		} else {
			toolInfo.ApiAuth = &mcp_service.ApiAuthWebRequest{
				Type:             toolCfg.Type,
				ApiKey:           toolCfg.ApiKey,
				AuthType:         toolCfg.AuthType,
				CustomHeaderName: toolCfg.CustomHeaderName,
			}
			return toolInfo, nil
		}
	}
	info, err := s.cli.GetCustomTool(ctx, &model.CustomTool{
		ID:           util.MustU32(req.CustomToolId),
		ToolSquareId: req.ToolSquareId,
		UserID:       req.Identity.UserId,
		OrgID:        req.Identity.OrgId,
	})
	if err != nil {
		return nil, grpc_util.ErrorStatus(errs.Code_MCPGetCustomToolInfoErr)
	}
	return &mcp_service.GetCustomToolInfoResp{
		CustomToolId:  util.Int2Str(info.ID),
		Name:          info.Name,
		Description:   info.Description,
		Schema:        info.Schema,
		PrivacyPolicy: info.PrivacyPolicy,
		ApiAuth: &mcp_service.ApiAuthWebRequest{
			Type:             info.Type,
			ApiKey:           info.APIKey,
			AuthType:         info.AuthType,
			CustomHeaderName: info.CustomHeaderName,
		},
	}, nil
}

func (s *Service) GetCustomToolList(ctx context.Context, req *mcp_service.GetCustomToolListReq) (*mcp_service.GetCustomToolListResp, error) {
	if req.Identity == nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", "identity is empty"))
	}
	infos, err := s.cli.ListCustomTools(ctx, req.Identity.OrgId, req.Identity.UserId, req.Name)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetCustomToolItem, 0)
	for _, info := range infos {
		list = append(list, &mcp_service.GetCustomToolItem{
			CustomToolId: util.Int2Str(info.ID),
			Name:         info.Name,
			Description:  info.Description,
		})
	}
	return &mcp_service.GetCustomToolListResp{
		List: list,
	}, nil
}

func (s *Service) GetCustomToolByCustomToolIdList(ctx context.Context, req *mcp_service.GetCustomToolByCustomToolIdListReq) (*mcp_service.GetCustomToolListResp, error) {
	if len(req.CustomToolIdList) == 0 {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, toErrStatus("mcp_get_custom_tool_list_err", "customToolIdList is empty"))
	}

	// 批量转换 string 为 uint32
	var ids []uint32
	for _, idStr := range req.CustomToolIdList {
		id := util.MustU32(idStr)
		ids = append(ids, id)
	}

	infos, err := s.cli.ListCustomToolsByCustomToolIDs(ctx, ids)
	if err != nil {
		return nil, errStatus(errs.Code_MCPGetCustomToolListErr, err)
	}
	list := make([]*mcp_service.GetCustomToolItem, 0)
	for _, info := range infos {
		list = append(list, &mcp_service.GetCustomToolItem{
			CustomToolId: util.Int2Str(info.ID),
			Name:         info.Name,
			Description:  info.Description,
		})
	}
	return &mcp_service.GetCustomToolListResp{
		List: list,
	}, nil
}

func (s *Service) UpdateCustomTool(ctx context.Context, req *mcp_service.UpdateCustomToolReq) (*emptypb.Empty, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "customToolId is empty"))
	}
	if req.ApiAuth == nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, toErrStatus("mcp_update_custom_tool_err", "apiAuth is empty"))
	}
	if err := s.cli.UpdateCustomTool(ctx, &model.CustomTool{
		ID:               util.MustU32(req.CustomToolId),
		Name:             req.Name,
		Description:      req.Description,
		Schema:           req.Schema,
		PrivacyPolicy:    req.PrivacyPolicy,
		Type:             req.ApiAuth.Type,
		APIKey:           req.ApiAuth.ApiKey,
		AuthType:         req.ApiAuth.AuthType,
		CustomHeaderName: req.ApiAuth.CustomHeaderName,
	}); err != nil {
		return nil, errStatus(errs.Code_MCPUpdateCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteCustomTool(ctx context.Context, req *mcp_service.DeleteCustomToolReq) (*emptypb.Empty, error) {
	if req.CustomToolId == "" {
		return nil, errStatus(errs.Code_MCPDeleteCustomToolErr, toErrStatus("mcp_delete_custom_tool_err", "customToolId is empty"))
	}
	if err := s.cli.DeleteCustomTool(ctx, util.MustU32(req.CustomToolId)); err != nil {
		return nil, errStatus(errs.Code_MCPDeleteCustomToolErr, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) GetSquareTool(ctx context.Context, req *mcp_service.GetSquareToolReq) (*mcp_service.SquareToolDetail, error) {
	mcpCfg, exist := config.Cfg().Tool(req.ToolSquareId)
	if !exist {
		return nil, errStatus(errs.Code_MCPGetSquareToolErr, toErrStatus("mcp_get_square_tool_err", "toolSquareId not exist"))
	}
	apiKey := ""
	if mcpCfg.NeedApiKeyInput {
		info, _ := s.cli.GetCustomTool(ctx, &model.CustomTool{
			ToolSquareId: req.ToolSquareId,
			UserID:       req.Identity.UserId,
			OrgID:        req.Identity.OrgId,
		})
		if info != nil {
			apiKey = info.APIKey
		}
	}
	return buildSquareToolDetail(mcpCfg, apiKey), nil
}

func (s *Service) GetSquareToolList(ctx context.Context, req *mcp_service.GetSquareToolListReq) (*mcp_service.SquareToolList, error) {
	var toolSquareInfo []*mcp_service.ToolSquareInfo
	for _, toolCfg := range config.Cfg().Tools {
		if req.Name != "" && !strings.Contains(toolCfg.Name, req.Name) {
			continue
		}
		toolSquareInfo = append(toolSquareInfo, buildSquareToolInfo(*toolCfg))
	}
	return &mcp_service.SquareToolList{Infos: toolSquareInfo}, nil
}

func buildSquareToolInfo(toolCfg config.ToolConfig) *mcp_service.ToolSquareInfo {
	return &mcp_service.ToolSquareInfo{
		ToolSquareId: toolCfg.ToolSquareId,
		AvatarPath:   toolCfg.AvatarPath,
		Name:         toolCfg.Name,
		Desc:         toolCfg.Desc,
		Tags:         toolCfg.Tags,
	}
}

func buildSquareToolDetail(toolCfg config.ToolConfig, apiKey string) *mcp_service.SquareToolDetail {
	return &mcp_service.SquareToolDetail{
		Info: buildSquareToolInfo(toolCfg),
		BuiltInTools: &mcp_service.BuiltInTools{
			NeedApiKeyInput: toolCfg.NeedApiKeyInput,
			ApiKey:          apiKey,
			Detail:          toolCfg.Detail,
			ActionSum:       int32(len(toolCfg.Tools)),
			Tools:           convertMCPTools(toolCfg.Tools),
		},
		Schema: toolCfg.Schema,
	}
}
