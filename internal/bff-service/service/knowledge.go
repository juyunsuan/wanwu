package service

import (
	knowledgebase_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-service"
	"github.com/UnicomAI/wanwu/internal/bff-service/config"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/request"
	"github.com/UnicomAI/wanwu/internal/bff-service/model/response"
	"github.com/gin-gonic/gin"
)

// SelectKnowledgeList 查询知识库列表，主要根据userId 查询用户所有知识库
func SelectKnowledgeList(ctx *gin.Context, userId, orgId string, req *request.KnowledgeSelectReq) (*response.KnowledgeListResp, error) {
	resp, err := knowledgeBase.SelectKnowledgeList(ctx.Request.Context(), &knowledgebase_service.KnowledgeSelectReq{
		UserId:    userId,
		OrgId:     orgId,
		Name:      req.Name,
		TagIdList: req.TagIdList,
	})
	if err != nil {
		return nil, err
	}
	return buildKnowledgeInfoList(resp), nil
}

// SelectKnowledgeInfoByName 根据知识库名称查询知识库信息
func SelectKnowledgeInfoByName(ctx *gin.Context, userId, orgId string, r *request.SearchKnowledgeInfoReq) (interface{}, error) {
	resp, err := knowledgeBase.SelectKnowledgeDetailByName(ctx.Request.Context(), &knowledgebase_service.KnowledgeDetailSelectReq{
		UserId:        userId,
		OrgId:         orgId,
		KnowledgeName: r.KnowledgeName,
	})
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"categoryId": resp.KnowledgeId,
	}, nil
}

// GetDeployInfo 查询部署信息
func GetDeployInfo(ctx *gin.Context) (interface{}, error) {
	cfgServer := config.Cfg().Server
	return map[string]string{
		"webBaseUrl": cfgServer.WebBaseUrl + "/minio/download/api/",
	}, nil
}

// CreateKnowledge 创建知识库
func CreateKnowledge(ctx *gin.Context, userId, orgId string, r *request.CreateKnowledgeReq) (*response.CreateKnowledgeResp, error) {
	resp, err := knowledgeBase.CreateKnowledge(ctx.Request.Context(), &knowledgebase_service.CreateKnowledgeReq{
		Name:        r.Name,
		Description: r.Description,
		UserId:      userId,
		OrgId:       orgId,
		EmbeddingModelInfo: &knowledgebase_service.EmbeddingModelInfo{
			ModelId: r.EmbeddingModel.ModelId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &response.CreateKnowledgeResp{KnowledgeId: resp.KnowledgeId}, nil
}

// UpdateKnowledge 更新知识库
func UpdateKnowledge(ctx *gin.Context, userId, orgId string, r *request.UpdateKnowledgeReq) error {
	_, err := knowledgeBase.UpdateKnowledge(ctx.Request.Context(), &knowledgebase_service.UpdateKnowledgeReq{
		KnowledgeId: r.KnowledgeId,
		Name:        r.Name,
		Description: r.Description,
		UserId:      userId,
		OrgId:       orgId,
	})
	return err
}

// DeleteKnowledge 删除知识库
func DeleteKnowledge(ctx *gin.Context, userId, orgId string, r *request.DeleteKnowledge) (interface{}, error) {
	resp, err := knowledgeBase.DeleteKnowledge(ctx.Request.Context(), &knowledgebase_service.DeleteKnowledgeReq{
		KnowledgeId: r.KnowledgeId,
		UserId:      userId,
		OrgId:       orgId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// KnowledgeHit 知识库命中
func KnowledgeHit(ctx *gin.Context, userId, orgId string, r *request.KnowledgeHitReq) (*response.KnowledgeHitResp, error) {
	matchParams := r.KnowledgeMatchParams
	resp, err := knowledgeBase.KnowledgeHit(ctx.Request.Context(), &knowledgebase_service.KnowledgeHitReq{
		Question:      r.Question,
		UserId:        userId,
		OrgId:         orgId,
		KnowledgeList: buildKnowledgeListReq(r),
		KnowledgeMatchParams: &knowledgebase_service.KnowledgeMatchParams{
			MatchType:         matchParams.MatchType,
			RerankModelId:     matchParams.RerankModelId,
			PriorityMatch:     matchParams.PriorityMatch,
			SemanticsPriority: matchParams.SemanticsPriority,
			KeywordPriority:   matchParams.KeywordPriority,
			TopK:              matchParams.TopK,
			Score:             matchParams.Threshold,
			TermWeight:        matchParams.TermWeight,
			TermWeightEnable:  matchParams.TermWeightEnable,
		},
	})
	if err != nil {
		return nil, err
	}
	return buildKnowledgeHitResp(resp), nil
}

func GetKnowledgeMetaSelect(ctx *gin.Context, userId, orgId string, r *request.GetKnowledgeMetaSelectReq) (*response.GetKnowledgeMetaSelectResp, error) {
	metaList, err := knowledgeBase.GetKnowledgeMetaSelect(ctx.Request.Context(), &knowledgebase_service.SelectKnowledgeMetaReq{
		UserId:      userId,
		OrgId:       orgId,
		KnowledgeId: r.KnowledgeId,
	})
	if err != nil {
		return nil, err
	}
	return buildKnowledgeMetaList(metaList.MetaList), nil
}

// GetKnowledgeMetaValueList 获取文档元数据列表
func GetKnowledgeMetaValueList(ctx *gin.Context, userId, orgId string, r *request.KnowledgeMetaValueListReq) (*response.KnowledgeMetaValueListResp, error) {
	resp, err := knowledgeBase.GetKnowledgeMetaValueList(ctx.Request.Context(), &knowledgebase_service.KnowledgeMetaValueListReq{
		UserId:    userId,
		OrgId:     orgId,
		DocIdList: r.DocIdList,
	})
	if err != nil {
		return nil, err
	}
	return buildKnowledgeMetaValueRespList(resp), nil
}

// UpdateKnowledgeMetaValue 更新知识库元数据值
func UpdateKnowledgeMetaValue(ctx *gin.Context, userId, orgId string, r *request.UpdateMetaValueReq) error {
	_, err := knowledgeBase.UpdateKnowledgeMetaValue(ctx.Request.Context(), &knowledgebase_service.UpdateKnowledgeMetaValueReq{
		UserId:          userId,
		OrgId:           orgId,
		ApplyToSelected: r.ApplyToSelected,
		DocIdList:       r.DocIdList,
		MetaList:        buildKnowledgeMetaValueReqList(r.MetaValueList),
	})
	if err != nil {
		return err
	}
	return nil
}

// buildKnowledgeMetaList 构造知识库元数据列表
func buildKnowledgeMetaList(metaList []*knowledgebase_service.KnowledgeMetaData) *response.GetKnowledgeMetaSelectResp {
	var retMetaList []*response.KnowledgeMetaItem
	for _, meta := range metaList {
		retMetaList = append(retMetaList, &response.KnowledgeMetaItem{
			MetaKey:       meta.Key,
			MetaValueType: meta.Type,
			MetaId:        meta.MetaId,
		})
	}
	return &response.GetKnowledgeMetaSelectResp{MetaList: retMetaList}
}

// buildKnowledgeListReq 构造命中测试 - 知识库列表参数
func buildKnowledgeListReq(r *request.KnowledgeHitReq) []*knowledgebase_service.KnowledgeParams {
	var knowledgeList []*knowledgebase_service.KnowledgeParams
	for _, k := range r.KnowledgeList {
		knowledgeList = append(knowledgeList, &knowledgebase_service.KnowledgeParams{
			KnowledgeId: k.ID,
			MetaDataFilterParams: &knowledgebase_service.MetaDataFilterParams{
				FilterEnable:     k.MetaDataFilterParams.FilterEnable,
				FilterLogicType:  k.MetaDataFilterParams.FilterLogicType,
				MetaFilterParams: buildMetaFilterParams(k.MetaDataFilterParams.MetaFilterParams),
			},
		})
	}
	return knowledgeList
}

// buildKnowledgeInfoList 构造知识库列表结果
func buildKnowledgeInfoList(knowledgeListResp *knowledgebase_service.KnowledgeSelectListResp) *response.KnowledgeListResp {
	if knowledgeListResp == nil || len(knowledgeListResp.KnowledgeList) == 0 {
		return &response.KnowledgeListResp{}
	}

	var list []*response.KnowledgeInfo
	for _, knowledge := range knowledgeListResp.KnowledgeList {
		list = append(list, &response.KnowledgeInfo{
			KnowledgeId: knowledge.KnowledgeId,
			Name:        knowledge.Name,
			Description: knowledge.Description,
			DocCount:    int(knowledge.DocCount),
			EmbeddingModelInfo: &response.EmbeddingModelInfo{
				ModelId: knowledge.EmbeddingModelInfo.ModelId,
			},
			KnowledgeTagList: buildTagList(knowledge.KnowledgeTagInfoList),
			CreateAt:         knowledge.CreatedAt,
		})
	}
	return &response.KnowledgeListResp{KnowledgeList: list}
}

// buildTagList 构造知识库标签列表
func buildTagList(tagList []*knowledgebase_service.KnowledgeTagInfo) []*response.KnowledgeTag {
	var retTagList = make([]*response.KnowledgeTag, 0)
	if len(tagList) > 0 {
		for _, tag := range tagList {
			retTagList = append(retTagList, &response.KnowledgeTag{
				TagId:    tag.TagId,
				TagName:  tag.TagName,
				Selected: true,
			})
		}
	}
	return retTagList
}

// buildKnowledgeHitResp 构造知识库命中返回
func buildKnowledgeHitResp(resp *knowledgebase_service.KnowledgeHitResp) *response.KnowledgeHitResp {
	var searchList = make([]*response.ChunkSearchList, 0)
	if len(resp.SearchList) > 0 {
		for _, search := range resp.SearchList {
			childContentList := make([]*response.ChildContent, 0)
			for _, child := range search.ChildContentList {
				childContentList = append(childContentList, &response.ChildContent{
					ChildSnippet: child.ChildSnippet,
					Score:        float64(child.Score),
				})
			}
			childScore := make([]float64, 0)
			for _, score := range search.ChildScore {
				childScore = append(childScore, float64(score))
			}
			searchList = append(searchList, &response.ChunkSearchList{
				Title:            search.Title,
				Snippet:          search.Snippet,
				KnowledgeName:    search.KnowledgeName,
				ChildContentList: childContentList,
				ChildScore:       childScore,
			})
		}
	}
	return &response.KnowledgeHitResp{
		Prompt:     resp.Prompt,
		Score:      resp.Score,
		SearchList: searchList,
	}
}

func buildMetaFilterParams(meta []*request.MetaFilterParams) []*knowledgebase_service.MetaFilterParams {
	var metaList []*knowledgebase_service.MetaFilterParams
	for _, m := range meta {
		metaList = append(metaList, &knowledgebase_service.MetaFilterParams{
			Key:       m.Key,
			Value:     m.Value,
			Type:      m.Type,
			Condition: m.Condition,
		})
	}
	return metaList
}

func buildKnowledgeMetaValueRespList(resp *knowledgebase_service.KnowledgeMetaValueListResp) *response.KnowledgeMetaValueListResp {
	retList := make([]*response.KnowledgeMetaValues, 0)
	for _, meta := range resp.MetaList {
		retList = append(retList, &response.KnowledgeMetaValues{
			MetaId:        meta.MetaId,
			MetaKey:       meta.Key,
			MetaValue:     meta.ValueList,
			MetaValueType: meta.Type,
		})
	}
	return &response.KnowledgeMetaValueListResp{
		KnowledgeMetaValues: retList,
	}
}

func buildKnowledgeMetaValueReqList(req []*request.DocMetaData) []*knowledgebase_service.MetaValueOperation {
	metaList := make([]*knowledgebase_service.MetaValueOperation, 0)
	for _, meta := range req {
		metaList = append(metaList, &knowledgebase_service.MetaValueOperation{
			MetaInfo: &knowledgebase_service.KnowledgeMetaData{
				MetaId: meta.MetaId,
				Key:    meta.MetaKey,
				Value:  meta.MetaValue,
				Type:   meta.MetaValueType,
			},
			Option: meta.Option,
		})
	}
	return metaList
}
