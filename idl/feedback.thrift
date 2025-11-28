namespace go feedback
include "./model.thrift"

struct UploadFeedbackRequest{
    1: required string type
    2: required string content
}
struct UploadFeedbackResponse{
      1: required model.BaseResp base,
}
struct QueryFeedbackRequest{

}
struct QueryFeedbackResponse{
          1: required model.BaseResp base,
          2: required model.FeedbackList data,
}
service FeedbackService{
    UploadFeedbackResponse UploadFeedback(1:UploadFeedbackRequest req)(api.post = "/api/feedback/upload"),
    QueryFeedbackResponse QueryFeedback(1:QueryFeedbackRequest req)(api.get = "/api/admin/feedback/query"),
}