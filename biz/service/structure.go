package service

import (
	"context"
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
)

// 集成一个获取专业和学院的函数，由于专业和学院在业务流程中是绝对的高频访问。依旧采取cache优先
func QueryAllMajor(ctx context.Context) ([]*model.Major, error) {
	exist, err := cache.IsMajorExist(ctx)
	if err != nil {
		return nil, err
	}
	var majorList []*model.Major
	if !exist {
		// db 载入 redis
		majorList, _, err = mysql.GetAllMajorInfo(ctx)
		if err != nil {
			return nil, err
		}
		err = cache.MajorToCache(ctx, majorList)
		if err != nil {
			return nil, err
		}
	} else {
		majorList, err = cache.QueryAllMajor(ctx)
		if err != nil {
			return nil, err
		}
	}
	return majorList, nil
}

func QueryAllCollege(ctx context.Context) ([]*model.College, error) {
	exist, err := cache.IsCollegeExist(ctx)
	if err != nil {
		return nil, err
	}
	var collegeList []*model.College
	if !exist {
		// db 载入 redis
		collegeList, _, err = mysql.GetCollegeInfo(ctx)
		if err != nil {
			return nil, err
		}
		err = cache.CollegeToCache(ctx, collegeList)
		if err != nil {
			return nil, err
		}
	} else {
		collegeList, err = cache.QueryAllCollege(ctx)
		if err != nil {
			return nil, err
		}
	}
	return collegeList, nil
}

// 查询认可奖项的函数
func QueryAllRecognizedReward(ctx context.Context) ([]*model.RecognizedEvent, error) {
	exist, err := cache.IsRecognizeEventExist(ctx)
	if err != nil {
		return nil, err
	}
	var recognizedEventList []*model.RecognizedEvent
	if !exist {
		// db 载入 redis
		recognizedEventList, _, err = mysql.QueryRecognizedEvent(ctx)
		if err != nil {
			return nil, err
		}
		err = cache.RecognizeEventToCache(ctx, recognizedEventList)
		if err != nil {
			return nil, err
		}
	} else {
		recognizedEventList, err = cache.QueryAllRecognizeEvent(ctx)
		if err != nil {
			return nil, err
		}
	}
	return recognizedEventList, nil
}

// 同样提供一个获取权责关系的函数 用于其他业务
func QueryAllRelation(ctx context.Context, user_id string) ([]*model.Relation, error) {
	exist, err := cache.IsRelationExist(ctx)
	if err != nil {
		return nil, err
	}
	var relationList []*model.Relation
	if !exist {
		relationList, _, err = mysql.QueryAllRelation(ctx)
		if err != nil {
			return nil, err
		}
		err = cache.RelationToCache(ctx, relationList)
		if err != nil {
			return nil, err
		}
	} else {
		relationList, err = cache.QueryRelationById(ctx, user_id)
		if err != nil {
			return nil, err
		}
	}
	var result []*model.Relation
	// 直接从mysql取出的数据是所有的 要进行一次筛选
	if !exist {
		for _, i := range relationList {
			if i.UserId == user_id {
				result = append(result, i)
			}
		}
	} else {
		result = relationList
	}
	return result, err
}

func IsMajorExist(ctx context.Context, major string) (bool, error) {
	majorlist, err := QueryAllMajor(ctx)
	if err != nil {
		return false, nil
	}
	for _, m := range majorlist {
		if m.MajorName == major {
			return true, nil
		}
	}
	return false, err
}
func IsCollegeExist(ctx context.Context, college string) (bool, error) {
	collegelist, err := QueryAllCollege(ctx)
	if err != nil {
		return false, nil
	}
	for _, m := range collegelist {
		if m.CollegeName == college {
			return true, nil
		}
	}
	return false, err
}
