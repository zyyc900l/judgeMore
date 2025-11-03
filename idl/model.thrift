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