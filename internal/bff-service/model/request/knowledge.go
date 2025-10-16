package request

import (
	"errors"
	"fmt"

	"github.com/UnicomAI/wanwu/pkg/util"
)

type KnowledgeSelectReq struct {
	Name      string   `json:"name" form:"name" `
	TagIdList []string `json:"tagId" form:"tagId" `
	CommonCheck
}

type CreateKnowledgeReq struct {
	Name           string          `json:"name"  validate:"required"`
	Description    string          `json:"description"`
	EmbeddingModel *EmbeddingModel `json:"embeddingModelInfo" validate:"required"`
}

type UpdateKnowledgeReq struct {
	KnowledgeId string `json:"knowledgeId"   validate:"required"`
	Name        string `json:"name"   validate:"required"`
	Description string `json:"description"`
	CommonCheck
}

type KnowledgeHitReq struct {
	KnowledgeList        []*AppKnowledgeBase   `json:"knowledgeList"`
	Question             string                `json:"question"   validate:"required"`
	KnowledgeMatchParams *KnowledgeMatchParams `json:"knowledgeMatchParams"   validate:"required"`
	CommonCheck
}

type KnowledgeMatchParams struct {
	MatchType         string  `json:"matchType"  validate:"required"` //matchType：vector（向量检索）、text（文本检索）、mix（混合检索：向量+文本）
	RerankModelId     string  `json:"rerankModelId"`                  //rerank模型id
	PriorityMatch     int32   `json:"priorityMatch"`                  // 权重匹配，只有在混合检索模式下，选择权重设置后，这个才设置为1
	SemanticsPriority float32 `json:"semanticsPriority"`              // 语义权重
	KeywordPriority   float32 `json:"keywordPriority"`                // 关键词权重
	TopK              int32   `json:"topK"`                           //topK 获取最高的几行
	Threshold         float32 `json:"threshold"`                      //threshold 过滤分数阈值
	TermWeight        float32 `json:"termWeight"`                     // 关键词系数
	TermWeightEnable  bool    `json:"termWeightEnable"`               // 关键词系数开关
	CommonCheck
}

type EmbeddingModel struct {
	ModelId string `json:"modelId"  validate:"required"`
}

type DeleteKnowledge struct {
	KnowledgeId string `json:"knowledgeId" validate:"required"`
	CommonCheck
}

type GetKnowledgeReq struct {
	KnowledgeId string `json:"knowledgeId" validate:"required"`
	CommonCheck
}

type CallbackUpdateDocStatusReq struct {
	DocId        string              `json:"id" validate:"required"`
	Status       int32               `json:"status" validate:"required"`
	MetaDataList []*CallbackMetaData `json:"metaDataList"`
	CommonCheck
}

type CallbackMetaData struct {
	Key    string `json:"key"`
	MetaId string `json:"metaId" validate:"required"`
	Value  string `json:"value" validate:"required"`
}

type DocMetaData struct {
	MetaId        string `json:"metaId"`        // 元数据id
	MetaKey       string `json:"metaKey"`       // key
	MetaValue     string `json:"metaValue"`     // 确定值
	MetaValueType string `json:"metaValueType"` // string，number，time
	MetaRule      string `json:"metaRule"`      // 正则表达式
	Option        string `json:"option"`        // option:add(新增)、update(更新)、delete(删除),update 和delete 的时候metaId 不能为空
}

type SearchKnowledgeInfoReq struct {
	KnowledgeName string `json:"categoryName" form:"categoryName" validate:"required"`
	UserId        string `json:"userId" form:"userId" validate:"required"`
	OrgId         string `json:"orgId"`
	CommonCheck
}

type GetKnowledgeMetaSelectReq struct {
	KnowledgeId string `json:"knowledgeId"  form:"knowledgeId" validate:"required"`
	CommonCheck
}

type KnowledgeMetaValueListReq struct {
	DocIdList []string `json:"docIdList" form:"docIdList" validate:"required" `
	CommonCheck
}

type UpdateMetaValueReq struct {
	DocIdList       []string       `json:"docIdList"  validate:"required"`
	MetaValueList   []*DocMetaData `json:"metaValueList"`
	ApplyToSelected bool           `json:"applyToSelected"`
}

func (c *UpdateMetaValueReq) Check() error {
	for _, v := range c.MetaValueList {
		if v.Option == "" {
			return errors.New("option为空")
		}
	}
	return nil
}
func (c *CreateKnowledgeReq) Check() error {
	if !util.IsAlphanumeric(c.Name) {
		errMsg := fmt.Sprintf("知识库名称只能包含中文、数字、小写英文，符号之只能包含下划线和减号 参数(%v)", c.Name)
		return errors.New(errMsg)
	}
	return nil
}
