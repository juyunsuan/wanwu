package request

import (
	"errors"
	"regexp"

	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	grpc_util "github.com/UnicomAI/wanwu/pkg/grpc-util"
)

const (
	DocAnalyzerOCR       = "ocr"
	DocAnalyzerPdfParser = "model"
	CommonSplitMethod    = "0" //通用分段
	ParentSplitMethod    = "1" //父子分段
)

type DocListReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"`
	DocName     string `json:"docName" form:"docName"`
	Status      int    `json:"status" form:"status"` // 当前状态  -1-全部， 0-待处理， 1- 处理完成， 2-正在审核中，3-正在解析中，4-审核未通过，5-解析失败
	PageSearch
	CommonCheck
}

type DocImportReq struct {
	KnowledgeId   string         `json:"knowledgeId" validate:"required"` //知识库id
	DocImportType int            `json:"docImportType"`                   //文档导入类型，0：文件上传，1：url上传，2.批量url上传
	DocInfo       []*DocInfo     `json:"docInfoList" validate:"required"` //上传文档列表
	DocSegment    *DocSegment    `json:"docSegment" validate:"required"`  //文档分段配置
	DocAnalyzer   []string       `json:"docAnalyzer" validate:"required"` //文档解析类型 text / ocr  / model
	ParserModelId string         `json:"parserModelId"`                   //模型解析或ocr模型id
	DocPreprocess []string       `json:"docPreprocess"`                   // 文本预处理规则 replaceSymbols / deleteLinks
	DocMetaData   []*DocMetaData `json:"docMetaData"`                     // 元数据
}

type DocMetaDataReq struct {
	KnowledgeId  string         `json:"knowledgeId"`
	DocId        string         `json:"docId"`
	MetaDataList []*DocMetaData `json:"metaDataList"` //文档元数据
}

type BatchDocMetaDataReq struct {
	KnowledgeId  string         `json:"knowledgeId"`
	MetaDataList []*DocMetaData `json:"metaDataList"` //文档元数据
	CreateMeta   bool           `json:"createMeta"`   //文档没设置过对应key则创建元数据
}

type DocInfo struct {
	DocId   string `json:"docId"`   //文档id
	DocName string `json:"docName"` //文档名称
	DocUrl  string `json:"docUrl"`  //文档url
	DocType string `json:"docType"` // 文档类型
	DocSize int64  `json:"docSize"` // 文档类型
}

type DocSegment struct {
	SegmentMethod  string   `json:"segmentMethod" validate:"required"` //分段方法 0：通用分段；1：父子分段
	SegmentType    string   `json:"segmentType"`                       //分段方式，只有通用分段必填 0：自动分段；1：自定义分段
	Splitter       []string `json:"splitter"`                          // 分隔符（只有自定义分段必填）
	MaxSplitter    int      `json:"maxSplitter"`                       // 可分隔最大值（只有自定义分段必填）
	Overlap        float32  `json:"overlap"`                           // 可重叠值（只有自定义分段必填）
	SubSplitter    []string `json:"subSplitter"`                       // 分隔符（只有父子分段必填）
	SubMaxSplitter int      `json:"subMaxSplitter"`                    // 可分隔最大值（只有父子分段必填）
}

type QueryKnowledgeReq struct {
	KnowledgeId string `json:"knowledgeId" form:"knowledgeId" validate:"required"`
	CommonCheck
}

type DeleteDocReq struct {
	DocIdList []string `json:"docIdList"  validate:"required"`
	CommonCheck
}

type DocSegmentListReq struct {
	DocId string `json:"docId" form:"docId" validate:"required"`
	PageSearch
	CommonCheck
}

type UpdateDocSegmentStatusReq struct {
	DocId         string `json:"docId" validate:"required"`
	ContentId     string `json:"contentId"`
	ContentStatus string `json:"contentStatus" validate:"required"`
	ALL           bool   `json:"all" ` // all 代表全部启用，此时将忽略contentId
	CommonCheck
}

type AnalysisUrlDocReq struct {
	KnowledgeId string   `json:"knowledgeId"   validate:"required"`
	UrlList     []string `json:"urlList"   validate:"required"`
	CommonCheck
}

type DocSegmentLabelsReq struct {
	ContentId string   `json:"contentId"  validate:"required"`
	DocId     string   `json:"docId"  validate:"required"`
	Labels    []string `json:"labels"  validate:"required"`
	CommonCheck
}

type CreateDocSegmentReq struct {
	DocId   string   `json:"docId"  validate:"required"`   // 文档id
	Labels  []string `json:"labels"  validate:"required"`  // 关键词列表
	Content string   `json:"content"  validate:"required"` // 分段内容
	CommonCheck
}

type BatchCreateDocSegmentReq struct {
	DocId        string `json:"docId"  validate:"required"`        // 文档id
	FileUploadId string `json:"fileUploadId"  validate:"required"` // fileUploadId
	CommonCheck
}

type DeleteDocSegmentReq struct {
	DocId     string `json:"docId"  validate:"required"` // 文档id
	ContentId string `json:"contentId"  validate:"required"`
	CommonCheck
}

type UpdateDocSegmentReq struct {
	DocId     string `json:"docId"  validate:"required"`
	ContentId string `json:"contentId"  validate:"required"`
	Content   string `json:"content"  validate:"required"`
	CommonCheck
}

type DocChildListReq struct {
	DocId     string `json:"docId" form:"docId" validate:"required"`
	ContentId string `json:"contentId"  form:"contentId" validate:"required"`
	CommonCheck
}

type CreateDocChildSegmentReq struct {
	DocId    string   `json:"docId"  validate:"required"`    // 文档id
	ParentId string   `json:"parentId"  validate:"required"` // 父分段id
	Content  []string `json:"content"  validate:"required"`  // 分段内容
	CommonCheck
}

type UpdateDocChildSegmentReq struct {
	DocId         string      `json:"docId"  validate:"required"`      // 文档id
	ParentId      string      `json:"parentId"  validate:"required"`   // 父分段id
	ParentChunkNo int32       `json:"parentChunkNo"`                   // 父分段序列号
	ChildChunk    *ChildChunk `json:"childChunk"  validate:"required"` // 子分段序列号列表
	CommonCheck
}

type ChildChunk struct {
	ChildNo int32  `json:"chunkNo"` // 子分段序列号
	Content string `json:"content"` // 子分段内容
}

type DeleteDocChildSegmentReq struct {
	DocId            string  `json:"docId"  validate:"required"`            // 文档id
	ParentId         string  `json:"parentId"  validate:"required"`         // 父分段id
	ParentChunkNo    int32   `json:"parentChunkNo"`                         // 父分段序列号
	ChildChunkNoList []int32 `json:"ChildChunkNoList"  validate:"required"` // 子分段序列号列表
	CommonCheck
}

func (c *DocImportReq) Check() error {
	if len(c.DocAnalyzer) > 0 {
		for _, v := range c.DocAnalyzer {
			if v == DocAnalyzerOCR || v == DocAnalyzerPdfParser {
				if c.ParserModelId == "" {
					return errors.New("parserModelId can not be empty")
				}
			}
		}
	}
	if len(c.DocMetaData) > 0 {
		seenKeys := make(map[string]bool)
		for _, meta := range c.DocMetaData {
			if meta.MetaKey == "" {
				return grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "key为空")
			}
			// 检查Key是否重复
			if seenKeys[meta.MetaKey] {
				return grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "key重复")
			}
			seenKeys[meta.MetaKey] = true
			if meta.MetaRule != "" {
				// 检查rule和key传参
				if meta.MetaValue != "" {
					return grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "常量和正则表达式重复")
				}
				// 检查正则合法性
				_, err := regexp.Compile(meta.MetaRule)
				if err != nil {
					return grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "非法正则表达式")
				}
				// 检查key合法性
				if !isValidKey(meta.MetaKey) {
					return grpc_util.ErrorStatus(errs.Code_BFFInvalidArg, "非法key")
				}
			}
		}
	}

	if c.DocSegment != nil {
		if c.DocSegment.SegmentMethod != CommonSplitMethod && c.DocSegment.SegmentMethod != ParentSplitMethod {
			return errors.New("segmentMethod error")
		}
		if c.DocSegment.SegmentMethod == CommonSplitMethod && c.DocSegment.SegmentType == "" {
			return errors.New("segmentType error")
		}
		if c.DocSegment.SegmentMethod == ParentSplitMethod && c.DocSegment.SubMaxSplitter > c.DocSegment.MaxSplitter {
			return errors.New("subMaxSplitter error")
		}
	}

	return nil
}

func isValidKey(s string) bool {
	re := regexp.MustCompile(`^[a-z][a-z0-9_]*$`) //只包含小写字母，数字和下划线，并且以小写字母开头
	return re.MatchString(s)
}

func (c *DocMetaDataReq) Check() error {
	if len(c.KnowledgeId) == 0 && len(c.DocId) == 0 {
		return errors.New("knowledgeId and docId can not all empty")
	}
	return nil
}

func (c *BatchDocMetaDataReq) Check() error {
	if len(c.KnowledgeId) == 0 {
		return errors.New("knowledgeId can not all empty")
	}
	if len(c.MetaDataList) > 0 {
		keyMap := make(map[string]bool)
		for _, meta := range c.MetaDataList {
			if meta.MetaKey == "" || meta.MetaValueType == "" {
				return errors.New("key or value type can not be empty")
			}
			if keyMap[meta.MetaKey] {
				return errors.New("key can not be repeated")
			}
			keyMap[meta.MetaKey] = true
		}
	}
	return nil
}
