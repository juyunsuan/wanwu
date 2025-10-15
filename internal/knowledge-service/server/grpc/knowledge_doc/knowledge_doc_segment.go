package knowledge_doc

import (
	"context"
	"encoding/json"
	errs "github.com/UnicomAI/wanwu/api/proto/err-code"
	knowledgebase_doc_service "github.com/UnicomAI/wanwu/api/proto/knowledgebase-doc-service"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/model"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/client/orm"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/generator"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/pkg/util"
	"github.com/UnicomAI/wanwu/internal/knowledge-service/service"
	"github.com/UnicomAI/wanwu/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
	"time"
	"unicode/utf8"
)

func (s *Service) CreateDocSegment(ctx context.Context, req *knowledgebase_doc_service.CreateDocSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法增加切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.查询最大分段长度
	importTask, err := orm.SelectKnowledgeImportTaskById(ctx, doc.ImportTaskId)
	if err != nil {
		log.Errorf("没有查询到导入任务 参数(%v)", req)
		return nil, err
	}
	var segmentConfig = &model.SegmentConfig{}
	err = json.Unmarshal([]byte(importTask.SegmentConfig), segmentConfig)
	if err != nil {
		log.Errorf("SegmentConfig process error %s", err.Error())
		return nil, err
	}
	//6.去除分段多余空格
	req.Content = strings.TrimSpace(req.Content)
	//7.判断分段长度
	if len(req.Content) == 0 {
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentEmpty)
	}
	if utf8.RuneCountInString(req.Content) > segmentConfig.MaxSplitter {
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentExceedMaxSize)
	}
	//8.发送rag请求
	var labels = req.Labels
	if len(labels) == 0 {
		labels = make([]string, 0)
	}
	var chunks []*service.NewChunkItem
	chunks = append(chunks, &service.NewChunkItem{
		Content: req.Content,
		Labels:  labels,
	})
	err = service.RagCreateDocSegment(ctx, &service.RagCreateDocSegmentParams{
		UserId:           req.UserId,
		KnowledgeBase:    knowledge.Name,
		KnowledgeId:      knowledge.KnowledgeId,
		FileName:         fileName,
		MaxSentenceSize:  segmentConfig.MaxSplitter,
		Chunks:           chunks,
		SplitType:        service.RebuildSplitType(segmentConfig.SegmentMethod),
		ChildChunkConfig: service.RebuildChildChunkConfig(segmentConfig.SegmentMethod, segmentConfig.SubMaxSplitter, segmentConfig.SubSplitter),
	})
	if err != nil {
		log.Errorf("docId %v create doc segment fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) BatchCreateDocSegment(ctx context.Context, req *knowledgebase_doc_service.BatchCreateDocSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法增加切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.查询最大分段长度
	importTask, err := orm.SelectKnowledgeImportTaskById(ctx, doc.ImportTaskId)
	if err != nil {
		log.Errorf("没有查询到导入任务 参数(%v)", req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	var segmentConfig = &model.SegmentConfig{}
	err = json.Unmarshal([]byte(importTask.SegmentConfig), segmentConfig)
	if err != nil {
		log.Errorf("SegmentConfig process error %s", err.Error())
		return nil, err
	}

	task, err := buildDocSegmentImportTask(knowledge, fileName, doc.DocId, segmentConfig, req)
	if err != nil {
		log.Errorf("docId %v create doc segment import task params fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}

	err = orm.CreateDocSegmentImportTask(ctx, task)
	if err != nil {
		log.Errorf("docId %v create doc segment import task fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateDocSegment(ctx context.Context, req *knowledgebase_doc_service.UpdateDocSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法更新切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentUpdateFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.查询最大分段长度
	importTask, err := orm.SelectKnowledgeImportTaskById(ctx, doc.ImportTaskId)
	if err != nil {
		log.Errorf("没有查询到导入任务 参数(%v)", req)
		return nil, err
	}
	var segmentConfig = &model.SegmentConfig{}
	err = json.Unmarshal([]byte(importTask.SegmentConfig), segmentConfig)
	if err != nil {
		log.Errorf("SegmentConfig process error %s", err.Error())
		return nil, err
	}
	//6.去除分段多余空格
	req.Content = strings.TrimSpace(req.Content)
	//7.判断分段长度
	if len(req.Content) == 0 {
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentEmpty)
	}
	if utf8.RuneCountInString(req.Content) > segmentConfig.MaxSplitter {
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentExceedMaxSize)
	}
	//8.发送rag请求
	err = service.RagUpdateDocSegment(ctx, &service.RagUpdateDocSegmentParams{
		UserId:          req.UserId,
		KnowledgeBase:   knowledge.Name,
		KnowledgeId:     knowledge.KnowledgeId,
		FileName:        fileName,
		MaxSentenceSize: segmentConfig.MaxSplitter,
		Chunk: &service.UpdateChunkItem{
			ChunkId: req.ContentId,
			Content: req.Content,
			Labels:  make([]string, 0),
		},
		SplitType:        service.RebuildSplitType(segmentConfig.SegmentMethod),
		ChildChunkConfig: service.RebuildChildChunkConfig(segmentConfig.SegmentMethod, segmentConfig.SubMaxSplitter, segmentConfig.SubSplitter),
	})
	if err != nil {
		log.Errorf("docId %v update doc segment fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentUpdateFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteDocSegment(ctx context.Context, req *knowledgebase_doc_service.DeleteDocSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法删除切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentDeleteFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.发送rag请求
	err = service.RagDeleteDocSegment(ctx, &service.RagDeleteDocSegmentParams{
		UserId:        req.UserId,
		KnowledgeBase: knowledge.Name,
		KnowledgeId:   knowledge.KnowledgeId,
		FileName:      fileName,
		ChunkIds:      []string{req.ContentId},
	})
	if err != nil {
		log.Errorf("docId %v delete doc segment fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentDeleteFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateDocSegmentStatus(ctx context.Context, req *knowledgebase_doc_service.UpdateDocSegmentStatusReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	docInfo := docList[0]
	//2.查询知识库详情
	knowledge, err := orm.SelectKnowledgeById(ctx, docInfo.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("查询知识库详情失败 参数(%v)", req)
		return nil, err
	}
	//3.更新文档状态
	var params = buildDocUpdateSegmentStatusParams(req, knowledge, docInfo)
	err = service.RagDocUpdateDocSegmentStatus(ctx, params)
	if err != nil {
		log.Errorf("UpdateFileStatus 更新知识库文档切片启用状态 失败(%v)  参数(%v)", err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentStatusUpdateFail)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateDocSegmentLabels(ctx context.Context, req *knowledgebase_doc_service.DocSegmentLabelsReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法增加切片标签 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentUpdateLabelsFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.更新切片标签
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	var labels = req.Labels
	if len(labels) == 0 {
		labels = make([]string, 0)
	}
	err = service.RagDocSegmentLabels(ctx, &service.RagDocSegmentLabelsParams{
		UserId:        req.UserId,
		KnowledgeBase: knowledge.Name,
		KnowledgeId:   knowledge.KnowledgeId,
		FileName:      fileName,
		ContentId:     req.ContentId,
		Labels:        labels,
	})
	if err != nil {
		log.Errorf("docId %v update doc seg labels fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentUpdateLabelsFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) CreateDocChildSegment(ctx context.Context, req *knowledgebase_doc_service.CreateDocChildSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法增加子切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.查询最大分段长度
	importTask, err := orm.SelectKnowledgeImportTaskById(ctx, doc.ImportTaskId)
	if err != nil {
		log.Errorf("没有查询到导入任务 参数(%v)", req)
		return nil, err
	}
	var segmentConfig = &model.SegmentConfig{}
	err = json.Unmarshal([]byte(importTask.SegmentConfig), segmentConfig)
	if err != nil {
		log.Errorf("SegmentConfig process error %s", err.Error())
		return nil, err
	}
	if err1 := checkContentLength(req.Content, segmentConfig.SubMaxSplitter); err1 != nil {
		log.Errorf("内容长度超出最大分段长度 错误(%v) 参数(%v)", err1, req)
		return nil, err1
	}
	err = service.RagCreateDocChildSegment(ctx, &service.RagCreateDocChildSegmentParams{
		UserId:        req.UserId,
		KnowledgeBase: knowledge.Name,
		KnowledgeId:   knowledge.KnowledgeId,
		FileName:      fileName,
		ChunkId:       req.ParentChunkId,
		ChildContents: req.Content,
	})
	if err != nil {
		log.Errorf("docId %v create doc child segment fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentCreateFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateDocChildSegment(ctx context.Context, req *knowledgebase_doc_service.UpdateDocChildSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法修改子切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentUpdateFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.查询最大分段长度
	importTask, err := orm.SelectKnowledgeImportTaskById(ctx, doc.ImportTaskId)
	if err != nil {
		log.Errorf("没有查询到导入任务 参数(%v)", req)
		return nil, err
	}
	var segmentConfig = &model.SegmentConfig{}
	err = json.Unmarshal([]byte(importTask.SegmentConfig), segmentConfig)
	if err != nil {
		log.Errorf("SegmentConfig process error %s", err.Error())
		return nil, err
	}
	if err1 := checkContentLength([]string{req.ChildChunk.Content}, segmentConfig.SubMaxSplitter); err1 != nil {
		log.Errorf("内容长度超出最大分段长度 错误(%v) 参数(%v)", err1, req)
		return nil, err1
	}
	//6.修改子分段信息
	err = service.RagUpdateDocChildSegment(ctx, &service.RagUpdateDocChildSegmentParams{
		UserId:          req.UserId,
		KnowledgeBase:   knowledge.Name,
		KnowledgeId:     knowledge.KnowledgeId,
		FileName:        fileName,
		ChunkId:         req.ParentChunkId,
		ChunkCurrentNum: req.ParentChunkNo,
		ChildChunk: &service.ChildChunk{
			ChildChunkNo: req.ChildChunk.ChunkNo,
			ChildContent: req.ChildChunk.Content,
		},
	})
	if err != nil {
		log.Errorf("docId %v update doc child segment fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentUpdateFailed)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteDocChildSegment(ctx context.Context, req *knowledgebase_doc_service.DeleteDocChildSegmentReq) (*emptypb.Empty, error) {
	//1.查询文档详情
	docList, err := orm.SelectDocByDocIdList(ctx, []string{req.DocId}, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库文档的权限 参数(%v)", req)
		return nil, err
	}
	doc := docList[0]
	//2.状态校验
	if util.BuildDocRespStatus(doc.Status) != model.DocSuccess {
		log.Errorf("非处理完成文档无法删除切片 状态(%d) 错误(%v) 参数(%v)", doc.Status, err, req)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentDeleteFailed)
	}
	//3.查询知识库信息
	knowledge, err := orm.SelectKnowledgeById(ctx, doc.KnowledgeId, req.UserId, req.OrgId)
	if err != nil {
		log.Errorf("没有操作该知识库的权限 参数(%v)", req)
		return nil, err
	}
	//4.获取文档名称
	fileName := service.RebuildFileName(doc.DocId, doc.FileType, doc.Name)
	//5.发送rag请求
	err = service.RagDeleteDocChildSegment(ctx, &service.RagDeleteDocChildSegmentParams{
		UserId:                req.UserId,
		KnowledgeBase:         knowledge.Name,
		KnowledgeId:           knowledge.KnowledgeId,
		FileName:              fileName,
		ChunkId:               req.ParentChunkId,
		ChunkCurrentNum:       req.ParentChunkNo,
		ChildChunkCurrentNums: req.ChildChunkNo,
	})
	if err != nil {
		log.Errorf("docId %v delete doc child segment fail %v", req.DocId, err)
		return nil, util.ErrCode(errs.Code_KnowledgeDocSegmentDeleteFailed)
	}
	return &emptypb.Empty{}, nil
}

// checkContentLength 检查内容长度
func checkContentLength(contentList []string, maxLength int) error {
	for _, content := range contentList {
		if len(content) > maxLength {
			return util.ErrCode(errs.Code_KnowledgeDocSegmentExceedMaxSize)
		}
	}
	return nil
}

// buildDocSegmentImportTask 构造导入任务
func buildDocSegmentImportTask(knowledge *model.KnowledgeBase, fileName, docId string,
	segmentConfig *model.SegmentConfig, req *knowledgebase_doc_service.BatchCreateDocSegmentReq) (*model.DocSegmentImportTask, error) {
	params := &model.DocSegmentImportParams{
		KnowledgeId:     knowledge.KnowledgeId,
		KnowledgeName:   knowledge.Name,
		FileName:        fileName,
		MaxSentenceSize: segmentConfig.MaxSplitter,
		FileUrl:         req.FileUrl,
		SegmentMethod:   segmentConfig.SegmentMethod,
		SubMaxSplitter:  segmentConfig.SubMaxSplitter,
		SubSplitter:     segmentConfig.SubSplitter,
	}
	marshal, err := json.Marshal(params)
	if err != nil {
		log.Errorf("DocSegmentImportParams process error %s", err.Error())
		return nil, err
	}

	return &model.DocSegmentImportTask{
		ImportId:     generator.GetGenerator().NewID(),
		DocId:        docId,
		Status:       model.DocSegmentImportInit,
		ImportParams: string(marshal),
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
		UserId:       req.UserId,
		OrgId:        req.OrgId,
	}, nil
}

func buildDocUpdateSegmentStatusParams(req *knowledgebase_doc_service.UpdateDocSegmentStatusReq, knowledge *model.KnowledgeBase, docInfo *model.KnowledgeDoc) interface{} {
	//前端逻辑，all + status 组合控制一键开启和一键关停，比如：all：true，status：false 则标识一键关停
	//但是底层 只要all false 就是一键关停
	var status = req.ContentStatus == "true"
	if req.All {
		return &service.DocSegmentStatusUpdateAllParams{
			DocSegmentStatusUpdateParams: service.DocSegmentStatusUpdateParams{
				UserId:        req.UserId,
				KnowledgeName: knowledge.Name,
				FileName:      service.RebuildFileName(docInfo.DocId, docInfo.FileType, docInfo.Name),
				ContentId:     req.ContentId,
			},
			All: status,
		}
	} else {
		return &service.DocSegmentStatusUpdateParams{
			UserId:        req.UserId,
			KnowledgeName: knowledge.Name,
			FileName:      service.RebuildFileName(docInfo.DocId, docInfo.FileType, docInfo.Name),
			ContentId:     req.ContentId,
			Status:        status,
		}
	}
}
