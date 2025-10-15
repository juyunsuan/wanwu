package service

import (
	"context"
	"fmt"
	net_url "net/url"
	"strings"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	openapi3_util "github.com/UnicomAI/wanwu/pkg/openapi3-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

func GetWorkflowToolSelect(ctx *gin.Context, userId, orgId, toolType, name string) (*response.ListResult, error) {
	switch toolType {
	case constant.ToolTypeBuiltIn:
		resp, err := mcp.GetSquareToolList(ctx.Request.Context(), &mcp_service.GetSquareToolListReq{
			Name: name,
		})
		if err != nil {
			return nil, err
		}
		var list []response.ToolSelect4Workflow
		for _, item := range resp.Infos {
			detail, err := GetToolSquareDetail(ctx, userId, orgId, item.ToolSquareId)
			if err != nil {
				return nil, err
			}
			url, _ := net_url.JoinPath(config.Cfg().Server.ApiBaseUrl, detail.Avatar.Path)
			list = append(list, response.ToolSelect4Workflow{
				ToolID:   item.ToolSquareId,
				ToolName: item.Name,
				ToolType: constant.ToolTypeBuiltIn,
				IconUrl:  url,
				ApiKey:   detail.APIKey,
				Desc:     detail.Desc,
				Actions:  builtInToolActions4Workflow(detail.Tools),
			})
		}
		return &response.ListResult{
			List:  list,
			Total: int64(len(list)),
		}, nil
	case constant.ToolTypeCustom:
		resp, err := mcp.GetCustomToolList(ctx.Request.Context(), &mcp_service.GetCustomToolListReq{
			Name: name,
			Identity: &mcp_service.Identity{
				UserId: userId,
				OrgId:  orgId,
			},
		})
		if err != nil {
			return nil, err
		}
		var list []response.ToolSelect4Workflow
		for _, item := range resp.List {
			detail, err := GetCustomTool(ctx, userId, orgId, item.CustomToolId)
			if err != nil {
				return nil, err
			}
			url, _ := net_url.JoinPath(config.Cfg().Server.ApiBaseUrl, config.Cfg().DefaultIcon.ToolIcon)
			list = append(list, response.ToolSelect4Workflow{
				ToolID:   item.CustomToolId,
				ToolName: item.Name,
				ToolType: constant.ToolTypeCustom,
				IconUrl:  url,
				ApiKey:   "",
				Desc:     detail.Description,
				Actions:  customToolActions4Workflow(detail.ApiList),
			})
		}
		return &response.ListResult{
			List:  list,
			Total: int64(len(list)),
		}, nil
	}
	return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "todo", "unsupported tool type")
}

func GetWorkflowToolDetail(ctx *gin.Context, userId, orgId, toolId, toolType, operationId string) (*response.ToolDetail4Workflow, error) {
	var schema, iconUrl string
	switch toolType {
	case constant.ToolTypeBuiltIn:
		resp, err := mcp.GetSquareTool(ctx.Request.Context(), &mcp_service.GetSquareToolReq{
			ToolSquareId: toolId,
			Identity: &mcp_service.Identity{
				UserId: userId,
				OrgId:  orgId,
			},
		})
		if err != nil {
			return nil, err
		}
		schema = resp.Schema
		iconUrl, _ = net_url.JoinPath(config.Cfg().Server.ApiBaseUrl, cacheMCPAvatar(ctx, resp.Info.AvatarPath).Path)
	case constant.ToolTypeCustom:
		resp, err := mcp.GetCustomToolInfo(ctx.Request.Context(), &mcp_service.GetCustomToolInfoReq{
			CustomToolId: toolId,
			Identity: &mcp_service.Identity{
				UserId: userId,
				OrgId:  orgId,
			},
		})
		if err != nil {
			return nil, err
		}
		schema = resp.Schema
		iconUrl, _ = net_url.JoinPath(config.Cfg().Server.ApiBaseUrl, config.Cfg().DefaultIcon.ToolIcon)
	}

	inputs, outputs, err := openapiSchema2ToolActionInputsAndOutputs4Workflow(ctx.Request.Context(), schema, operationId)
	if err != nil {
		return nil, err
	}
	return &response.ToolDetail4Workflow{
		ActionID:   operationId,
		ActionName: operationId,
		IconUrl:    iconUrl,
		Inputs:     inputs,
		Outputs:    outputs,
	}, nil
}

// --- internal ---

func builtInToolActions4Workflow(actions []response.MCPTool) []response.ToolAction4Workflow {
	var ret []response.ToolAction4Workflow
	for _, action := range actions {
		ret = append(ret, response.ToolAction4Workflow{
			ActionName: action.Name,
			ActionID:   action.Name,
			Desc:       action.Description,
		})
	}
	return ret
}

func customToolActions4Workflow(actions []response.CustomToolActionInfo) []response.ToolAction4Workflow {
	var ret []response.ToolAction4Workflow
	for _, action := range actions {
		ret = append(ret, response.ToolAction4Workflow{
			ActionName: action.Name,
			ActionID:   action.Name,
			Desc:       action.Desc,
		})
	}
	return ret
}

func openapiSchema2ToolActionInputsAndOutputs4Workflow(ctx context.Context, schema, operationID string) ([]response.ToolActionParam4Workflow, []response.ToolActionParam4Workflow, error) {
	doc, err := openapi3_util.LoadFromData([]byte(schema))
	if err != nil {
		return nil, nil, err
	}
	if err := openapi3_util.ValidateDoc(ctx, doc); err != nil {
		return nil, nil, err
	}
	var exist bool
	var inputs, outputs []response.ToolActionParam4Workflow
	for _, pathItem := range doc.Paths.Map() {
		for _, operation := range pathItem.Operations() {
			if operation.OperationID != operationID {
				continue
			}
			exist = true
			inputs = openapiOperation2ToolActionInputs4Workflow(operation)
			outputs = openapiOperation2ToolActionOutputs4Workflow(operation)
			break
		}
	}
	if !exist {
		return nil, nil, fmt.Errorf("opentionID(%v) not found", operationID)
	}
	return inputs, outputs, nil
}

func openapiOperation2ToolActionInputs4Workflow(operation *openapi3.Operation) []response.ToolActionParam4Workflow {
	inputs := []response.ToolActionParam4Workflow{}

	// 解析路径参数、查询参数、header 参数等
	if operation.Parameters != nil {
		for _, param := range operation.Parameters {
			if param.Value != nil {
				inputs = append(inputs, response.ToolActionParam4Workflow{
					Name:        param.Value.In + "." + param.Value.Name,
					Description: param.Value.Description,
					Type:        openapiParameterType(param.Value),
					Required:    param.Value.Required,
					Children:    []response.ToolActionParam4Workflow{},
				})
			}
		}
	}

	// 解析请求体
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		for _, mediaType := range operation.RequestBody.Value.Content {
			if mediaType.Schema != nil && mediaType.Schema.Value != nil {
				inputs = append(inputs, openapiSchemaProperties2ToolActionParams4Workflow(mediaType.Schema.Value.Properties, mediaType.Schema.Value.Required)...)
			}
		}
	}
	return inputs
}

func openapiOperation2ToolActionOutputs4Workflow(operation *openapi3.Operation) []response.ToolActionParam4Workflow {
	outputs := []response.ToolActionParam4Workflow{}

	if operation.Responses == nil {
		return outputs
	}

	var responseRef *openapi3.ResponseRef
	// 优先查找 200 响应
	for statusCode, currResponseRef := range operation.Responses.Map() {
		if strings.HasPrefix(statusCode, "2") {
			responseRef = currResponseRef
			break
		}
	}
	// 如果没有2开头的响应，使用第一个可用的响应
	if responseRef == nil {
		for _, currResponseRef := range operation.Responses.Map() {
			responseRef = currResponseRef
			break
		}
	}

	if responseRef.Value != nil {
		for _, mediaType := range responseRef.Value.Content {
			if mediaType.Schema != nil && mediaType.Schema.Value != nil {
				outputs = append(outputs, openapiSchemaProperties2ToolActionParams4Workflow(mediaType.Schema.Value.Properties, mediaType.Schema.Value.Required)...)
				break
			}
		}
	}
	return outputs
}

func openapiSchemaProperties2ToolActionParams4Workflow(properties openapi3.Schemas, required []string) []response.ToolActionParam4Workflow {
	if properties == nil {
		return nil
	}
	var ret []response.ToolActionParam4Workflow
	for propName, propSchema := range properties {
		if propSchema.Value == nil {
			continue
		}
		ret = append(ret, response.ToolActionParam4Workflow{
			Name:        propName,
			Description: propSchema.Value.Description,
			Type:        openaiSchemaType(propSchema.Value),
			Required:    util.Exist(required, propName),
			Children:    openapiSchemaProperties2ToolActionParams4Workflow(propSchema.Value.Properties, propSchema.Value.Required),
		})
	}
	return ret
}

// openapiParameterType 获取参数类型
func openapiParameterType(param *openapi3.Parameter) string {
	if param.Schema != nil && param.Schema.Value != nil {
		return openaiSchemaType(param.Schema.Value)
	}
	return "string"
}

// openaiSchemaType 获取 schema 的类型
func openaiSchemaType(schema *openapi3.Schema) string {
	if schema.Type != nil {
		// 检查类型切片中的具体类型
		if schema.Type.Is("object") {
			return "object"
		} else if schema.Type.Is("array") {
			return "array"
		} else if schema.Type.Is("string") {
			return "string"
		} else if schema.Type.Is("number") {
			return "number"
		} else if schema.Type.Is("integer") {
			return "integer"
		} else if schema.Type.Is("boolean") {
			return "boolean"
		}

		// 如果有多个类型，返回第一个
		if len(*schema.Type) > 0 {
			return string((*schema.Type)[0])
		}
	}

	if len(schema.AnyOf) > 0 {
		return "anyOf"
	}
	if len(schema.AllOf) > 0 {
		return "allOf"
	}
	if len(schema.OneOf) > 0 {
		return "oneOf"
	}

	if schema.Format != "" {
		return schema.Format
	}

	return "string"
}
