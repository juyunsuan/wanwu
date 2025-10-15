package v1

import (
	"net/http"

	v1 "github.com/UnicomAI/wanwu/internal/bff-service/server/http/handler/v1"
	mid "github.com/UnicomAI/wanwu/pkg/gin-util/mid-wrap"
	"github.com/gin-gonic/gin"
)

func registerKnowledge(apiV1 *gin.RouterGroup) {
	// 知识库增删改查
	mid.Sub("knowledge").Reg(apiV1, "/knowledge", http.MethodPost, v1.CreateKnowledge, "创建知识库（文档分类）")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge", http.MethodPut, v1.UpdateKnowledge, "修改知识库（文档分类）")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge", http.MethodDelete, v1.DeleteKnowledge, "删除知识库（文档分类）")

	// 知识库命中测试
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/hit", http.MethodPost, v1.KnowledgeHit, "知识库命中测试")

	// 知识库文档
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/list", http.MethodGet, v1.GetDocList, "获取文档列表")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/import", http.MethodPost, v1.ImportDoc, "上传文档")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/import/tip", http.MethodGet, v1.GetDocImportTip, "获取知识库文档上传状态")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc", http.MethodDelete, v1.DeleteDoc, "删除文档")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/meta", http.MethodPost, v1.UpdateDocMetaData, "更新文档元数据")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/meta/batch", http.MethodPost, v1.BatchUpdateDocMetaData, "批量更新文档元数据")

	// 知识库元数据
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/meta/select", http.MethodGet, v1.GetKnowledgeMetaSelect, "获取知识库元数据")

	// 知识库文档切片增删改查
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/list", http.MethodGet, v1.GetDocSegmentList, "获取文档切分结果")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/status/update", http.MethodPost, v1.UpdateDocSegmentStatus, "更新文档切片启用状态")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/labels", http.MethodPost, v1.UpdateDocSegmentLabels, "更新文档切片标签")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/create", http.MethodPost, v1.CreateDocSegment, "新增文档切片")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/batch/create", http.MethodPost, v1.BatchCreateDocSegment, "批量新增文档切片")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/delete", http.MethodDelete, v1.DeleteDocSegment, "删除文档切片")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/update", http.MethodPost, v1.UpdateDocSegment, "更新文档切片")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/child/list", http.MethodGet, v1.GetDocChildSegmentList, "获取子分段列表")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/child/create", http.MethodPost, v1.CreateDocChildSegment, "创建子分段")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/child/update", http.MethodPost, v1.UpdateDocChildSegment, "更新子分段")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/segment/child/delete", http.MethodDelete, v1.DeleteDocChildSegment, "删除子分段")

	// 知识库url文档导入
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/doc/url/analysis", http.MethodPost, v1.AnalysisDocUrl, "解析url")

	// 知识库标签增删改查
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/tag", http.MethodGet, v1.GetKnowledgeTagSelect, "查询知识库标签列表")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/tag", http.MethodPost, v1.CreateKnowledgeTag, "创建知识库标签")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/tag", http.MethodPut, v1.UpdateKnowledgeTag, "修改知识库标签")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/tag", http.MethodDelete, v1.DeleteKnowledgeTag, "删除知识库标签")
	// 绑定知识库标签
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/tag/bind/count", http.MethodGet, v1.SelectTagBindCount, "查询标签绑定的知识库数量")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/tag/bind", http.MethodPost, v1.BindKnowledgeTag, "绑定知识库标签")

	// 知识库关键词管理
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/keywords", http.MethodGet, v1.GetKnowledgeKeywordsList, "查询知识库关键词列表")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/keywords", http.MethodPost, v1.CreateKnowledgeKeywords, "新增知识库关键词")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/keywords/detail", http.MethodGet, v1.GetKnowledgeKeywordsDetail, "查询知识库关键词详情")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/keywords", http.MethodPut, v1.UpdateKnowledgeKeywords, "编辑知识库关键词")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/keywords", http.MethodDelete, v1.DeleteDocCategoryKeywords, "删除知识库关键词")

	// 知识库分隔符增删改查
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/splitter", http.MethodGet, v1.GetKnowledgeSplitterSelect, "查询知识库分隔符列表")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/splitter", http.MethodPost, v1.CreateKnowledgeSplitter, "创建知识库分隔符")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/splitter", http.MethodPut, v1.UpdateKnowledgeSplitter, "修改知识库分隔符")
	mid.Sub("knowledge").Reg(apiV1, "/knowledge/splitter", http.MethodDelete, v1.DeleteKnowledgeSplitter, "删除知识库分隔符")
}
