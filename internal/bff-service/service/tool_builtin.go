package service

import (
	"strings"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	"github.com/gin-gonic/gin"
)

func GetToolSquareDetail(ctx *gin.Context, userID, orgID, toolSquareID string) (*response.ToolSquareDetail, error) {
	resp, err := mcp.GetSquareTool(ctx.Request.Context(), &mcp_service.GetSquareToolReq{
		ToolSquareId: toolSquareID,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if err != nil {
		return nil, err
	}
	return toToolSquareDetail(ctx, resp), nil
}

func GetToolSquareList(ctx *gin.Context, userID, orgID, name string) (*response.ListResult, error) {
	resp, err := mcp.GetSquareToolList(ctx.Request.Context(), &mcp_service.GetSquareToolListReq{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	var list []response.ToolSquareInfo
	for _, item := range resp.Infos {
		list = append(list, response.ToolSquareInfo{
			ToolSquareID: item.ToolSquareId,
			Avatar:       cacheMCPAvatar(ctx, item.AvatarPath),
			Name:         item.Name,
			Desc:         item.Desc,
			Tags:         getToolTags(item.Tags),
		})
	}
	return &response.ListResult{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

func UpdateToolSquareAPIKey(ctx *gin.Context, userID, orgID string, req request.ToolSquareAPIKeyReq) error {
	toolInfo, _ := mcp.GetCustomToolInfo(ctx.Request.Context(), &mcp_service.GetCustomToolInfoReq{
		ToolSquareId: req.ToolSquareID,
		Identity: &mcp_service.Identity{
			UserId: userID,
			OrgId:  orgID,
		},
	})
	if toolInfo == nil {
		return grpc_util.ErrorStatus(errs.Code_MCPGetCustomToolInfoErr, "tool not found")
	}
	if toolInfo.ApiAuth.ApiKey == "" {
		_, _ = mcp.CreateCustomTool(ctx.Request.Context(), &mcp_service.CreateCustomToolReq{
			ToolSquareId: req.ToolSquareID,
			ApiAuth: &mcp_service.ApiAuthWebRequest{
				ApiKey: req.APIKey,
			},
			Identity: &mcp_service.Identity{
				UserId: userID,
				OrgId:  orgID,
			},
		})
	}
	if toolInfo.CustomToolId != "" {
		_, err := mcp.UpdateCustomTool(ctx.Request.Context(), &mcp_service.UpdateCustomToolReq{
			CustomToolId: toolInfo.CustomToolId,
			ApiAuth: &mcp_service.ApiAuthWebRequest{
				ApiKey: req.APIKey,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// --- internal ---

func toToolSquareDetail(ctx *gin.Context, toolSquare *mcp_service.SquareToolDetail) *response.ToolSquareDetail {
	ret := &response.ToolSquareDetail{
		ToolSquareInfo: toToolSquareInfo(ctx, toolSquare.Info),
		ToolSquareActions: response.ToolSquareActions{
			NeedApiKeyInput: toolSquare.BuiltInTools.NeedApiKeyInput,
			APIKey:          toolSquare.BuiltInTools.ApiKey,
			Detail:          toolSquare.BuiltInTools.Detail,
			ActionSum:       int64(toolSquare.BuiltInTools.ActionSum),
		},
	}
	for _, tool := range toolSquare.BuiltInTools.Tools {
		ret.ToolSquareActions.Tools = append(ret.ToolSquareActions.Tools, toMCPTool(tool))
	}
	return ret
}

func toToolSquareInfo(ctx *gin.Context, toolSquareInfo *mcp_service.ToolSquareInfo) response.ToolSquareInfo {
	return response.ToolSquareInfo{
		ToolSquareID: toolSquareInfo.ToolSquareId,
		Avatar:       cacheMCPAvatar(ctx, toolSquareInfo.AvatarPath),
		Name:         toolSquareInfo.Name,
		Desc:         toolSquareInfo.Desc,
		Tags:         getToolTags(toolSquareInfo.Tags),
	}
}

func getToolTags(tagString string) []string {
	if tagString == "" {
		return []string{}
	}
	return strings.Split(tagString, ",")
}
