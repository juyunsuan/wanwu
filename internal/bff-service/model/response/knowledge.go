package response

type KnowledgeListResp struct {
	KnowledgeList []*KnowledgeInfo `json:"knowledgeList"`
}

type CreateKnowledgeResp struct {
	KnowledgeId string `json:"knowledgeId"`
}

type KnowledgeHitResp struct {
	Prompt     string             `json:"prompt"`     //提示词列表
	SearchList []*ChunkSearchList `json:"searchList"` //种种结果
	Score      []float64          `json:"score"`      //打分信息
}

type EmbeddingModelInfo struct {
	ModelId string `json:"modelId"`
}

type KnowledgeInfo struct {
	KnowledgeId        string              `json:"knowledgeId"`        //知识库id
	Name               string              `json:"name"`               //知识库名称
	Description        string              `json:"description"`        //知识库描述
	DocCount           int                 `json:"docCount"`           //文档数量
	EmbeddingModelInfo *EmbeddingModelInfo `json:"embeddingModelInfo"` //embedding模型信息
	KnowledgeTagList   []*KnowledgeTag     `json:"knowledgeTagList"`   //知识库标签列表
	CreateAt           string              `json:"createAt"`           //创建时间
}

type KnowledgeMetaData struct {
	Key  string `json:"key"`  // key
	Type string `json:"type"` // type(time, string, number)
}

type ChunkSearchList struct {
	Title            string          `json:"title"`
	Snippet          string          `json:"snippet"`
	KnowledgeName    string          `json:"knowledgeName"`
	ChildContentList []*ChildContent `json:"childContentList"`
	ChildScore       []float64       `json:"childScore"`
}

type ChildContent struct {
	ChildSnippet string  `json:"childSnippet"`
	Score        float64 `json:"score"`
}

type GetKnowledgeMetaSelectResp struct {
	MetaList []*KnowledgeMetaItem `json:"knowledgeMetaList"`
}

type KnowledgeMetaItem struct {
	MetaId        string `json:"metaId"`
	MetaKey       string `json:"metaKey"`
	MetaValueType string `json:"metaValueType"`
	MetaValue     string `json:"metaValue"` // 确定值
}

type KnowledgeMetaValueListResp struct {
	KnowledgeMetaValues []*KnowledgeMetaValues `json:"knowledgeMetaValues"`
}

type KnowledgeMetaValues struct {
	MetaId        string   `json:"metaId"`
	MetaKey       string   `json:"metaKey"`
	MetaValue     []string `json:"metaValue"` // 确定值
	MetaValueType string   `json:"metaValueType"`
}
