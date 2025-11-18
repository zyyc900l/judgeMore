package client

import (
	"errors"
	"fmt"
	"judgeMore/config"
	"judgeMore/pkg/errno"

	"github.com/olivere/elastic/v7"
)

func NewEsClient() (*elastic.Client, error) {
	if config.Elasticsearch == nil {
		return nil, errors.New("elasticsearch config is nil")
	}
	esConn := fmt.Sprintf("http://%s", config.Elasticsearch.Addr)
	client, err := elastic.NewClient(
		elastic.SetURL(esConn),
	)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalESErrorCode, fmt.Sprintf("es clint failed,error: %v", err))
	}
	return client, nil
}
