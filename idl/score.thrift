namespace go score
include "./model.thrift"
struct QueryScoreByScoreIdRequest{
    1: required string score_id,
}
struct QueryScoreByScoreIdResponse{
     1: required model.BaseResp base,
             2: required model.ScoreRecord data,
}
struct QueryScoreByStuIdRequest{
        1: required string stu_id,
}
struct QueryScoreByStuIdResponse{
     1: required model.BaseResp base,
             2: required model.ScoreRecordList data,
}
struct QueryScoreByEventIdRequest{
        1: required string event_id,
}
struct QueryScoreByEventIdResponse{
     1: required model.BaseResp base,
             2: required model.ScoreRecord data,
}
service ScoreService {
    QueryScoreByScoreIdResponse QueryScoreByScoreId(1:QueryScoreByScoreIdRequest req)(api.get="/api/query/score/id"),
    QueryScoreByEventIdResponse QueryScoreByEventId(1:QueryScoreByEventIdRequest req)(api.get="/api/query/score/material"),
    QueryScoreByStuIdResponse QueryScoreByStuId(1:QueryScoreByStuIdRequest req)(api.get="/api/query/score/stu"),
}