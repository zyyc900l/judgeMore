package es

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/olivere/elastic/v7"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
)

func IsExist(ctx context.Context, indexName string) bool {
	rels, err := els.IndexExists(indexName).Do(ctx)
	if err != nil {
		fmt.Printf("Elastic.IsExist Error checking if index exists: %v", err)
		return false
	}
	return rels
}

func CreateIndex(ctx context.Context, indexName string) error {
	_, err := els.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "Elastic.CreateIndex Error creating index: %v", err)
	}
	return nil
}

func IsIndexDataExist(ctx context.Context, indexName string) (bool, error) {
	count, err := els.Count(indexName).Do(ctx)
	if err != nil {
		return false, errno.Errorf(errno.InternalESErrorCode, "Elastic.IsIndexEmpty failed: %v", err)
	}
	hlog.Info(count)
	return count != 0, nil
}

// 当存在该记录时，es的应对策略是用旧的去覆盖
func AddItem(ctx context.Context, indexName string, re *model.RecognizedEvent) error {
	_, err := els.Index().Index(indexName).
		Id(fmt.Sprintf("%s", re.RecognizedEventId)).
		BodyJson(re).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "Elastic.AddItem Error adding item: %v", err)
	}
	return nil
}

// 当该删除的记录不存在时，err也回是nil
func RemoveItem(ctx context.Context, indexName string, id string) error {
	_, err := els.Delete().Index(indexName).Id(fmt.Sprintf("%s", id)).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "Elastic.RemoveItem failed: %v", err)
	}
	return nil
}

func SearchItems(ctx context.Context, indexName string, req *model.ViewRecognizedRewardReq) ([]*model.RecognizedEvent, int64, error) {
	q := BuildQuery(req)
	result, err := els.Search().Index(indexName).
		Query(q).Size(20).Do(ctx)
	if err != nil {
		return nil, 0, errno.Errorf(errno.InternalESErrorCode, "CommodityElastic.SearchItems failed: %v", err)
	}
	rets := make([]*model.RecognizedEvent, 0)
	for _, hit := range result.Hits.Hits {
		var re model.RecognizedEvent
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, "CommodityElastic.SearchItems failed: %v", err)
		}
		err = sonic.Unmarshal(data, &re)
		if err != nil {
			return nil, 0, errno.Errorf(errno.InternalServiceErrorCode, "CommodityElastic.SearchItems failed: %v", err)
		}
		rets = append(rets, &re)
	}
	hlog.Info(result.TotalHits())
	return rets, result.TotalHits(), nil
}
func BuildQuery(req *model.ViewRecognizedRewardReq) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	hasCondition := false

	// 事件名称查询 - 使用分词搜索
	if req.EventName != nil && req.GetEventName() != "" {
		hlog.Info("Searching event name:", req.GetEventName())
		// 使用 MatchQuery 进行分词搜索
		query = query.Must(elastic.NewMatchQuery("RecognizedEventName", req.GetEventName()))
		hasCondition = true
	}

	// 组织者名称查询 - 使用分词搜索
	if req.OrganizerName != nil && req.GetOrganizerName() != "" {
		hlog.Info("Searching organizer:", req.GetOrganizerName())
		query = query.Must(elastic.NewMatchQuery("Organizer", req.GetOrganizerName()))
		hasCondition = true
	}

	// 识别事件ID查询 - 保持精确匹配（ID通常是精确的）
	if req.RecognizedEventId != nil && req.GetRecognizedEventId() != "" {
		query = query.Must(elastic.NewTermQuery("RecognizedEventId", req.GetRecognizedEventId()))
		hasCondition = true
	}

	if !hasCondition {
		query = query.Must(elastic.NewMatchAllQuery())
	}

	return query
}
