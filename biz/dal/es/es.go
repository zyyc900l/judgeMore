package es

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
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

func AddItem(ctx context.Context, indexName string, re *model.RecognizedEvent) error {
	_, err := els.Index().Index(indexName).
		Id(fmt.Sprintf("%s", re.RecognizedEventId)).
		BodyJson(re).Do(ctx)
	if err != nil {
		return errno.Errorf(errno.InternalESErrorCode, "Elastic.AddItem Error adding item: %v", err)
	}
	return nil
}

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
		Query(q).Do(ctx)
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
	return rets, result.TotalHits(), nil
}

func BuildQuery(req *model.ViewRecognizedRewardReq) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	hasCondition := false

	// 事件名称查询
	if req.EventName != nil && req.GetEventName() != "" {
		query = query.Must(elastic.NewMatchQuery("name", req.GetEventName()))
		hasCondition = true
	}

	// 组织者名称查询
	if req.OrganizerName != nil && req.GetOrganizerName() != "" {
		query = query.Must(elastic.NewMatchQuery("organizer", req.GetOrganizerName()))
		hasCondition = true
	}

	// 识别事件ID查询
	if req.RecognizedEventId != nil && req.GetRecognizedEventId() != "" {
		query = query.Must(elastic.NewMatchQuery("recognized_event_id", req.GetRecognizedEventId()))
		hasCondition = true
	}

	// 如果没有设置任何条件，返回匹配所有文档的查询
	if !hasCondition {
		query = query.Must(elastic.NewMatchAllQuery())
	}

	return query
}
