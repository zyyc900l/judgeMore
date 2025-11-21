package taskqueue

import (
	"context"
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/taskqueue"
)

func AddUpdateCacheRelationTask(ctx context.Context, key string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateCacheRelation(ctx)
	}})
}
func AddUpdateCacheMajorTask(ctx context.Context, key string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateCacheMajor(ctx)
	}})
}
func AddUpdateCacheCollegeTask(ctx context.Context, key string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateCacheCollege(ctx)
	}})
}
func AddUpdateInsertStuTask(ctx context.Context, key string, r *model.Relation) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateAdminStu(ctx, r)
	}})
}
func AddUpdateRecognizedTask(ctx context.Context, key string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateCacheRecognizedEvent(ctx)
	}})
}
func AddUpdateRuleTask(ctx context.Context, key string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return updateCacheRule(ctx)
	}})
}
func updateCacheRelation(ctx context.Context) error {
	relationList, _, err := mysql.QueryAllRelation(ctx)
	if err != nil {
		return err
	}
	err = cache.RelationToCache(ctx, relationList)
	if err != nil {
		return err
	}
	return nil
}
func updateCacheMajor(ctx context.Context) error {
	majorList, _, err := mysql.GetAllMajorInfo(ctx)
	if err != nil {
		return err
	}
	err = cache.MajorToCache(ctx, majorList)
	if err != nil {
		return err
	}
	return nil
}
func updateCacheCollege(ctx context.Context) error {
	collegeList, _, err := mysql.GetCollegeInfo(ctx)
	if err != nil {
		return err
	}
	err = cache.CollegeToCache(ctx, collegeList)
	if err != nil {
		return err
	}
	return nil
}
func updateCacheRecognizedEvent(ctx context.Context) error {
	re, _, err := mysql.QueryRecognizedEvent(ctx)
	if err != nil {
		return err
	}
	err = cache.RecognizeEventToCache(ctx, re)
	if err != nil {
		return err
	}
	return nil
}
func updateAdminStu(ctx context.Context, r *model.Relation) error {
	if r.CollegeName != "" {
		stuList, err := mysql.QueryUserByCollege(ctx, r.CollegeName)
		if err != nil {
			return err
		}
		err = mysql.InsertAdminStu(ctx, r.UserId, stuList)
		if err != nil {
			return err
		}
		return nil
	} else {
		stuList, err := mysql.QueryUserByMajor(ctx, r.MajorName, r.Grade)
		if err != nil {
			return err
		}
		err = mysql.InsertAdminStu(ctx, r.UserId, stuList)
		if err != nil {
			return err
		}
		return nil
	}
}
func updateCacheRule(ctx context.Context) error {
	re, _, err := mysql.GetScoreRule(ctx)
	if err != nil {
		return err
	}
	err = cache.RuleToCache(ctx, re)
	if err != nil {
		return err
	}
	return nil
}
