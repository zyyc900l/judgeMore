package service

import (
	"context"
	"fmt"
	"judgeMore/biz/dal/es"
	"judgeMore/biz/service/taskqueue"
	"judgeMore/pkg/constants"
)

func IsMappingExist(ctx context.Context) error {
	var err error
	if !es.IsExist(ctx, constants.IndexName) {
		err = es.CreateIndex(ctx, constants.IndexName)
		if err != nil {
			return fmt.Errorf("service.IsMappingExist CreateIndex failed: %w", err)
		}
	}
	return err
}
func Init() {
	taskqueue.Init()
	err := IsMappingExist(context.Background())
	if err != nil {
		fmt.Errorf(err.Error())
		panic(err)
	}
}
