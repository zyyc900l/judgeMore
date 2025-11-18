namespace go model

struct BaseResp{
    1: i64 code,
    2: string msg,
}

struct UserInfo{
    1: string username,  //姓名
    2: string userId,   // 学号
    4: string Major // 专业
    5: string college, //学院
    6: string grade,  // 年级
    7: string email //邮箱
    8: string role //角色
    9: required string created_at
    10: required string updated_at
    11: required string deleted_at
}
struct Event {
    1: string event_id,           // 赛事材料的自增id
    2: string user_id,            // 关联的学生的用户id
    3: string recognize_id
    4: string event_name,         // 赛事名称
    5: string event_organizer,    // 赛事主办方
    6: string event_level,        // 国家级 / 省级 / 校级 / 商业赛事
    7: string event_influence,    // 高 / 中 / 低
    8: string award_level,        // 一等奖 / 二等奖 / 三等奖 / 优秀奖等
    9: string material_url,       // 材料上传路径
    10: string material_status,    // 待审核 / 已审核 / 驳回
    11: bool auto_extracted,      // true - 是 / false - 否
    12: string award_time
    13: string created_at,        // 创建时间
    14: string updated_at,        // 更新时间
    15: string deleted_at         // 删除时
}
struct EventList{
     1: required list<Event> items,
     2: required i64 total,          //总数
}
struct ScoreRecord{
    1: string score_id ,
    2: string event_id ,
    3: string user_id ,
    4: string appeal_id,
    5: string rule_id,
    6: double final_score,
    7: string status,
    8: string created_at,        // 创建时间
    9: string updated_at,        // 更新时间
    10: string deleted_at         // 删除时
}
struct ScoreRecordList{
     1: required list<ScoreRecord> items,
     2: required i64 total,          //总数
     3: required double sum,
}
struct Appeal{
    1:string appeal_id,
    2: string result_id,
    3: string user_id,
    4: string appeal_type,
    5:string appeal_reason,
    6:string attachment_path,
    7:string status,
    8:string handleBy,
    9:string handleTime
    10:string handleResult,
    11:i64 appeal_count
        12: string created_at,        // 创建时间
        13: string updated_at,        // 更新时间
        14: string deleted_at         // 删除时
}
struct RecognizeReward{
        1: required string recognize_reward_id
        2: required string college //学院
        3: required string event_name
        4: required string organizer
        5: required string event_time
        6: required string related_majors //赛事涉及的专业范围
        7: required string applicable_majors //实际申请认定的专业
        8: required string recognition_basis //认定依据文件或标准
        9: required string recognized_level
        10: required bool is_active,
                12: string created_at,        // 创建时间
                13: string updated_at,        // 更新时间
                14: string deleted_at         // 删除时
}


struct RecognizeRewardList{
    1: list<RecognizeReward> item,
    2: i64 total,
}
struct AppealList{
    1:required list<Appeal>items,
    2:required i64 total
}
struct Major{
    1 : string MajorName //专业名称
    2 : i64 MajorId,
    3 : i64 CollegeId //所属学院
}
struct MajorList{
    1 : list <Major> item,
    2: i64 total,
}
struct College{
    1 : string CollegeName //专业名称
    2 : i64 CollegeId //所属学院
}
struct CollegeList{
    1 : list <College> item,
    2: i64 total,
}

struct Rule{
    1: string rule_id,
    2: string recognized_event_id,
    3: string event_level,
    4: double event_weight,
    5: i64 integral,
    6: string rule_desc,
    7: bool is_editable,
    8: string award_level
    9: i64 award_level_weight,
    10: string created_at,        // 创建时间
    11: string updated_at,        // 更新时间
    12: string deleted_at         // 删除时
}

