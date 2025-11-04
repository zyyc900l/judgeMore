namespace go admin
include "./model.thrift"

struct QueryCollegeRequest{
      1: required string page_num,
      2: required string page_size,
}
struct QueryCollegeResponse{
     1: required model.BaseResp base,
     2: required model.CollegeList data,
}

struct QueryMajorByCollegeIdRequest{
    1: required string page_num,
    2: required string page_size,
    3: required string college_id,
}
struct QueryMajorByCollegeIdResponse{
     1: required model.BaseResp base,
     2: required model.CollegeList data,
}



struct UploadMajorRequest{
     1: required string major_name,
     2: required string college_id,
}
struct UploadMajorResponse{
     1: required model.BaseResp base,
     2: required string major_id,
}
struct UploadCollegeRequest{
     1: required string college_name,
     2: optional list<string> major_ids
}
struct UploadCollegeResponse{
     1: required model.BaseResp base,
     2: required string college_id,
}


service AdminService{
     QueryCollegeResponse QueryCollege(1: QueryCollegeRequest req) (api.get = "/api/admin/colleges"),
     QueryMajorByCollegeIdResponse QueryMajorByCollegeId(1: QueryMajorByCollegeIdRequest req) (api.get = "/api/admin/majors"),
     UploadMajorResponse UploadMajor(1: UploadMajorRequest req) (api.post = "/api/admin/majors"),
     UploadCollegeResponse UploadCollege(1: UploadCollegeRequest req) (api.post = "/api/admin/colleges"),
}

