package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	net_url "net/url"
	"os"
	"slices"
	"strings"

	err_code "github.com/UnicomAI/wanwu/api/proto/err-code"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	mcp_service "github.com/UnicomAI/wanwu/api/proto/mcp-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/UnicomAI/wanwu/pkg/constant"
	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
	mp "github.com/UnicomAI/wanwu/pkg/model-provider"
	mp_common "github.com/UnicomAI/wanwu/pkg/model-provider/mp-common"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func ListLlmModelsByWorkflow(ctx *gin.Context, userId, orgId, modelT string) (*response.ListResult, error) {
	modelResp, err := ListTypeModels(ctx, userId, orgId, &request.ListTypeModelsRequest{ModelType: mp.ModelTypeLLM})
	if err != nil {
		return nil, err
	}
	var rets []*response.CozeWorkflowModelInfo
	for _, modelInfo := range modelResp.List.([]*response.ModelInfo) {
		ret, err := toModelInfoByWorkflow(modelInfo)
		if err != nil {
			return nil, err
		}
		rets = append(rets, ret)
	}
	return &response.ListResult{
		List:  rets,
		Total: modelResp.Total,
	}, nil
}

// ListWorkflow userID/orgID数据隔离，用于【工作流】
func ListWorkflow(ctx *gin.Context, orgID, name string) (*response.CozeWorkflowListData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ListUri)
	ret := &response.CozeWorkflowListResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"login_user_create": "true",
			"space_id":          orgID,
			"name":              name,
			"page":              "1",
			"size":              "99999",
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

// ListWorkflowByIDs 无userID或orgID隔离，用于【智能体选工作流】【应用广场】业务流程中
func ListWorkflowByIDs(ctx *gin.Context, name string, workflowIDs []string) (*response.CozeWorkflowListData, error) {
	var ids []string
	for _, workflowID := range workflowIDs {
		if _, err := util.I64(workflowID); err == nil {
			// AgentScope Workflow ID为uuid，将这部分脏数据过滤掉；Coze Workflow ID可转为数字
			ids = append(ids, workflowID)
		}
	}
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ListUri)
	ret := &response.CozeWorkflowListResp{}
	request := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"name": name,
			"page": "1",
			"size": "99999",
		})
	if len(ids) > 0 {
		request = request.SetBody(map[string]interface{}{
			"workflow_ids": ids,
		})
	}
	if resp, err := request.SetResult(ret).Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_apps_list", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func CreateWorkflow(ctx *gin.Context, orgID, name, desc, iconUri string) (*response.CozeWorkflowIDData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.CreateUri)
	ret := &response.CozeWorkflowIDResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"space_id": orgID,
			"name":     name,
			"desc":     desc,
			"icon_uri": iconUri,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_create", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func CopyWorkflow(ctx *gin.Context, orgID, workflowID string) (*response.CozeWorkflowIDData, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.CopyUri)
	ret := &response.CozeWorkflowIDResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"space_id":    orgID,
			"workflow_id": workflowID,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_copy", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_copy", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_copy", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func DeleteWorkflow(ctx *gin.Context, orgID, workflowID string) error {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.DeleteUri)
	ret := &response.CozeWorkflowDeleteResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"workflow_id": workflowID,
			"space_id":    orgID,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_delete", err.Error())
	} else if resp.StatusCode() >= 300 {
		return grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_delete", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 || (ret.Data != nil && ret.Data.Status != 0) {
		return grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_app_delete", fmt.Sprintf("code %v msg %v status %v", ret.Code, ret.Msg, ret.Data.GetStatus()))
	}
	return nil
}

func ExportWorkflow(ctx *gin.Context, orgID, workflowID string) ([]byte, error) {
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ExportUri)
	ret := &response.CozeWorkflowExportResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetBody(map[string]string{
			"space_id":    orgID,
			"workflow_id": workflowID,
		}).
		SetResult(&ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_export", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_export", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	}
	exportData := response.CozeWorkflowExportData{
		WorkflowName: ret.Data.WorkflowName,
		WorkflowDesc: ret.Data.WorkflowDesc,
		Schema:       ret.Data.Schema,
	}
	// 将结构体序列化为 JSON 字节
	jsonData, err := json.Marshal(exportData)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_export", fmt.Sprintf("export workflow unmarshal err:%v", err.Error()))
	}
	return jsonData, nil
}

func ImportWorkflow(ctx *gin.Context, orgID string) (*response.CozeWorkflowIDData, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", fmt.Sprintf("get file err: %v", err))
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", fmt.Sprintf("open file err: %v", err))
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", fmt.Sprintf("read file err: %v", err))
	}
	var rawData workflowImportData
	if err := json.Unmarshal(fileBytes, &rawData); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", fmt.Sprintf("schema unmarshal failed: %v", err))
	}
	url, _ := net_url.JoinPath(config.Cfg().Workflow.Endpoint, config.Cfg().Workflow.ImportUri)
	ret := &response.CozeWorkflowIDResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(workflowHttpReqHeader(ctx)).
		SetQueryParams(map[string]string{
			"space_id": orgID,
			"name":     rawData.Name,
			"desc":     rawData.Desc,
			"schema":   rawData.Schema,
		}).
		SetResult(ret).
		Post(url); err != nil {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", err.Error())
	} else if resp.StatusCode() >= 300 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", fmt.Sprintf("[%v] %v", resp.StatusCode(), resp.String()))
	} else if ret.Code != 0 {
		return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "bff_workflow_import_file", fmt.Sprintf("code %v msg %v", ret.Code, ret.Msg))
	}
	return ret.Data, nil
}

func GetWorkflowToolSelect(ctx *gin.Context, userId, orgId, toolType, name string) (*response.ListResult, error) {
	switch toolType {
	case "builtin":
		resp, err := mcp.GetSquareToolList(ctx.Request.Context(), &mcp_service.GetSquareToolListReq{
			Name: name,
		})
		if err != nil {
			return nil, err
		}
		var list []response.ToolSelectWithActions
		for _, item := range resp.Infos {
			detail, err := GetToolSquareDetail(ctx, userId, orgId, item.ToolSquareId)
			if err != nil {
				return nil, err
			}
			url, _ := net_url.JoinPath(os.Getenv("SERVER_API_BASE_URL"), cacheMCPAvatar(ctx, item.AvatarPath).Path)
			list = append(list, response.ToolSelectWithActions{
				ToolID:   item.ToolSquareId,
				ToolName: item.Name,
				ToolType: "builtIn",
				IconUrl:  url,
				ApiKey:   detail.APIKey,
				Desc:     detail.Desc,
				Actions:  builtInTool2actions(detail.Tools),
			})
		}
		return &response.ListResult{
			List:  list,
			Total: int64(len(list)),
		}, nil
	case "custom":
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
		var list []response.ToolSelectWithActions
		for _, item := range resp.List {
			detail, err := GetCustomTool(ctx, userId, orgId, item.CustomToolId)
			if err != nil {
				return nil, err
			}
			url, _ := net_url.JoinPath(os.Getenv("SERVER_API_BASE_URL"), config.Cfg().DefaultIcon.ToolIcon)
			list = append(list, response.ToolSelectWithActions{
				ToolID:   item.CustomToolId,
				ToolName: item.Name,
				ToolType: "custom",
				IconUrl:  url,
				ApiKey:   "",
				Desc:     detail.Description,
				Actions:  customTool2actions(detail.ApiList),
			})
		}
		return &response.ListResult{
			List:  list,
			Total: int64(len(list)),
		}, nil
	}
	return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "todo", "unsupported tool type")
}

func GetWorkflowToolDetail(ctx *gin.Context, userId, orgId, toolId, toolType, name string) (*response.ToolDetail4Workflow, error) {
	switch toolType {
	case "builtin":
		// TODO: 实现内置工具的逻辑
	case "custom":
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
		// 解析 OpenAPI Schema 并找到对应的 operation
		toolDetail, err := parseOpenAPISchemaForOperation([]byte(resp.Schema), name)
		if err != nil {
			return nil, err
		}
		iconUrl, _ := net_url.JoinPath(os.Getenv("SERVER_API_BASE_URL"), config.Cfg().DefaultIcon.ToolIcon)
		toolDetail.ActionName = name
		toolDetail.ActionID = name
		toolDetail.IconUrl = iconUrl
		return toolDetail, nil
	}
	return nil, grpc_util.ErrorStatusWithKey(errs.Code_BFFGeneral, "todo", "unsupported tool type")
}

// --- internal ---

type workflowImportData struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Schema string `json:"schema"` // 存储为JSON字符串
}

func workflowHttpReqHeader(ctx *gin.Context) map[string]string {
	return map[string]string{
		"Authorization": ctx.GetHeader("Authorization"),
		"X-Org-Id":      ctx.GetHeader(gin_util.X_ORG_ID),
		"X-User-Id":     ctx.GetString(gin_util.USER_ID),
		"Content-Type":  "application/json",
	}
}

func cozeWorkflowInfo2Model(workflowInfo *response.CozeWorkflowListDataWorkflow) response.AppBriefInfo {
	return response.AppBriefInfo{
		AppId:     workflowInfo.WorkflowId,
		AppType:   constant.AppTypeWorkflow,
		Name:      workflowInfo.Name,
		Desc:      workflowInfo.Desc,
		Avatar:    cacheWorkflowAvatar(workflowInfo.URL),
		CreatedAt: util.Time2Str(workflowInfo.CreateTime * 1000),
		UpdatedAt: util.Time2Str(workflowInfo.UpdateTime * 1000),
	}
}

func toModelInfoByWorkflow(modelInfo *response.ModelInfo) (*response.CozeWorkflowModelInfo, error) {
	ret := &response.CozeWorkflowModelInfo{
		ModelInfo:   *modelInfo,
		ModelParams: config.Cfg().Workflow.ModelParams,
	}
	if modelInfo.Config != nil {
		cfg := make(map[string]interface{})
		b, err := json.Marshal(modelInfo.Config)
		if err != nil {
			return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v marshal config err: %v", modelInfo.ModelId, err))
		}
		if err = json.Unmarshal(b, &cfg); err != nil {
			return nil, grpc_util.ErrorStatus(err_code.Code_BFFGeneral, fmt.Sprintf("model %v unmarshal config err: %v", modelInfo.ModelId, err))
		}
		for k, v := range cfg {
			switch k {
			case "functionCalling":
				if fc, ok := v.(string); ok && mp_common.FCType(fc) == mp_common.FCTypeToolCall {
					ret.ModelAbility.FunctionCall = true
				}
			case "visionSupport":
				if vs, ok := v.(string); ok && mp_common.VSType(vs) == mp_common.VSTypeSupport {
					ret.ModelAbility.ImageUnderstanding = true
				}

			}
		}
	}
	return ret, nil
}

func builtInTool2actions(tools []response.MCPTool) []response.ActionInfo {
	var actions []response.ActionInfo
	for _, tool := range tools {
		action := response.ActionInfo{
			ActionName: tool.Name,
			ActionID:   tool.Name,
			Desc:       tool.Description,
		}
		actions = append(actions, action)
	}
	return actions
}

func customTool2actions(tools []response.CustomToolActionInfo) []response.ActionInfo {
	var actions []response.ActionInfo
	for _, tool := range tools {
		action := response.ActionInfo{
			ActionName: tool.Name,
			ActionID:   tool.Name,
			Desc:       "",
		}
		actions = append(actions, action)
	}

	return actions
}

func parseOpenAPISchemaForOperation(schema []byte, name string) (*response.ToolDetail4Workflow, error) {
	// 加载 OpenAPI Schema
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI schema: %w", err)
	}

	// 验证文档
	if err := doc.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI schema: %w", err)
	}

	toolDetail := &response.ToolDetail4Workflow{
		Inputs:     []response.ToolActionParam4Workflow{},
		Outputs:    []response.ToolActionParam4Workflow{},
		ActionName: name,
		ActionID:   name,
	}

	// 检查 Paths 是否为 nil
	if doc.Paths == nil {
		return toolDetail, nil // 返回空的工具详情
	}

	// 遍历所有路径，查找 operationId 匹配的 operation
	for path, pathItem := range doc.Paths.Map() {
		if pathItem == nil {
			continue
		}

		// 获取路径的所有操作
		operations := pathItem.Operations()
		for method, operation := range operations {
			if operation != nil && (operation.OperationID == name || operation.Summary == name) {
				// 找到匹配的 operation，解析输入输出
				inputs := parseOperationInputs(operation, path, method)
				outputs := parseOperationOutputs(operation)

				toolDetail.Inputs = inputs
				toolDetail.Outputs = outputs
				return toolDetail, nil
			}
		}
	}

	// 如果没有找到对应的 operation，返回空的工具详情
	return toolDetail, nil
}

// parseOperationInputs 解析操作的输入参数
func parseOperationInputs(operation *openapi3.Operation, path, method string) []response.ToolActionParam4Workflow {
	var inputs []response.ToolActionParam4Workflow

	// 解析路径参数、查询参数、header 参数等
	if operation.Parameters != nil {
		for _, param := range operation.Parameters {
			if param.Value != nil {
				inputParam := response.ToolActionParam4Workflow{
					Name:        param.Value.Name,
					Description: param.Value.Description,
					Type:        getParameterType(param.Value),
					Required:    param.Value.Required,
					Children:    []response.ToolActionParam4Workflow{},
				}
				inputs = append(inputs, inputParam)
			}
		}
	}

	// 解析请求体
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		for _, mediaType := range operation.RequestBody.Value.Content {
			if mediaType.Schema != nil && mediaType.Schema.Value != nil {
				bodyParam := response.ToolActionParam4Workflow{
					Name:        "body",
					Description: operation.RequestBody.Value.Description,
					Type:        getSchemaType(mediaType.Schema.Value),
					Required:    operation.RequestBody.Value.Required,
					Children:    parseSchemaProperties(mediaType.Schema.Value),
				}
				inputs = append(inputs, bodyParam)
			}
		}
	}

	return inputs
}

// parseOperationOutputs 解析操作的输出参数
func parseOperationOutputs(operation *openapi3.Operation) []response.ToolActionParam4Workflow {
	var outputs []response.ToolActionParam4Workflow

	if operation.Responses == nil {
		return outputs
	}

	// 优先查找 200 响应
	for statusCode, responseRef := range operation.Responses.Map() {
		if strings.HasPrefix(statusCode, "2") && responseRef.Value != nil {
			outputs = parseResponseToOutputs(responseRef.Value, statusCode)
			return outputs
		}
	}

	// 如果没有2开头的响应，使用第一个可用的响应
	for statusCode, responseRef := range operation.Responses.Map() {
		if responseRef.Value != nil {
			outputs = parseResponseToOutputs(responseRef.Value, statusCode)
			break
		}
	}

	return outputs
}

// parseResponseToOutputs 解析响应到输出参数
func parseResponseToOutputs(r *openapi3.Response, statusCode string) []response.ToolActionParam4Workflow {
	var outputs []response.ToolActionParam4Workflow

	if r.Content != nil {
		for _, mediaType := range r.Content {
			if mediaType.Schema != nil && mediaType.Schema.Value != nil {
				schema := mediaType.Schema.Value

				// 对于对象类型的 schema，直接将其属性作为输出参数
				if schema.Type != nil && schema.Type.Is("object") && schema.Properties != nil {
					for propName, propSchema := range schema.Properties {
						if propSchema.Value != nil {
							outputParam := response.ToolActionParam4Workflow{
								Name:        propName,
								Description: propSchema.Value.Description,
								Type:        getSchemaType(propSchema.Value),
								Required:    isRequired(propName, schema.Required),
								Children:    parseSchemaProperties(propSchema.Value), // 递归解析子属性
							}
							outputs = append(outputs, outputParam)
						}
					}
				} else {
					// 非对象类型，直接创建输出参数
					outputParam := response.ToolActionParam4Workflow{
						Name:        statusCode,
						Description: getDescription(r.Description),
						Type:        getSchemaType(schema),
						Required:    false,
						Children:    parseSchemaProperties(schema),
					}
					outputs = append(outputs, outputParam)
				}
			}
		}
	}

	return outputs
}

// getParameterType 获取参数类型
func getParameterType(param *openapi3.Parameter) string {
	if param.Schema != nil && param.Schema.Value != nil {
		return getSchemaType(param.Schema.Value)
	}

	// 根据参数位置推断类型
	switch param.In {
	case "path", "query":
		return "string"
	case "header":
		return "string"
	default:
		return "string"
	}
}

// getSchemaType 获取 schema 的类型
func getSchemaType(schema *openapi3.Schema) string {
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

// getDescription 安全获取描述
func getDescription(desc *string) string {
	if desc == nil {
		return ""
	}
	return *desc
}

// parseSchemaProperties 递归解析 schema 的属性
func parseSchemaProperties(schema *openapi3.Schema) []response.ToolActionParam4Workflow {
	var children []response.ToolActionParam4Workflow

	if schema.Type != nil && schema.Type.Is("object") && schema.Properties != nil {
		for propName, propSchema := range schema.Properties {
			if propSchema.Value != nil {
				child := response.ToolActionParam4Workflow{
					Name:        propName,
					Description: propSchema.Value.Description,
					Type:        getSchemaType(propSchema.Value),
					Required:    isRequired(propName, schema.Required),
					Children:    parseSchemaProperties(propSchema.Value),
				}
				children = append(children, child)
			}
		}
	} else if schema.Type != nil && schema.Type.Is("array") && schema.Items != nil && schema.Items.Value != nil {
		arrayChild := response.ToolActionParam4Workflow{
			Name:        "items",
			Description: schema.Description,
			Type:        getSchemaType(schema.Items.Value),
			Required:    false,
			Children:    parseSchemaProperties(schema.Items.Value),
		}
		children = append(children, arrayChild)
	} else if schema.OneOf != nil {
		for i, oneOfSchema := range schema.OneOf {
			if oneOfSchema.Value != nil {
				child := response.ToolActionParam4Workflow{
					Name:        fmt.Sprintf("option_%d", i+1),
					Description: "One of the options",
					Type:        getSchemaType(oneOfSchema.Value),
					Required:    false,
					Children:    parseSchemaProperties(oneOfSchema.Value),
				}
				children = append(children, child)
			}
		}
	} else if schema.AnyOf != nil {
		for i, anyOfSchema := range schema.AnyOf {
			if anyOfSchema.Value != nil {
				child := response.ToolActionParam4Workflow{
					Name:        fmt.Sprintf("option_%d", i+1),
					Description: "Any of the options",
					Type:        getSchemaType(anyOfSchema.Value),
					Required:    false,
					Children:    parseSchemaProperties(anyOfSchema.Value),
				}
				children = append(children, child)
			}
		}
	} else if schema.AllOf != nil {
		for i, allOfSchema := range schema.AllOf {
			if allOfSchema.Value != nil {
				child := response.ToolActionParam4Workflow{
					Name:        fmt.Sprintf("combined_%d", i+1),
					Description: "Combined schema requirement",
					Type:        getSchemaType(allOfSchema.Value),
					Required:    false,
					Children:    parseSchemaProperties(allOfSchema.Value),
				}
				children = append(children, child)
			}
		}
	}

	return children
}

// isRequired 检查属性是否是必需的
func isRequired(propName string, required []string) bool {
	return slices.Contains(required, propName)
}
