package errno

// 错误码的设计原则是, 以尽可能少的错误码来传递必要的信息,
// 让前端能够根据尽量少的 error code 和具体的场景来告知用户错误信息
// 总的来说前端不依赖于后端传递的 msg 来告知用户, 而是通过 code 来额外处理
// 当然如果有一些强指向性错误信息, 你当然可以再写进来一个 code, 比如密码错误或者用户已存在
// 我们将这种与业务强相关的 code 也放在 errno 包中, 主要是为了方便统一管理与避免 code 冲突

// 业务处理成功
const (
	SuccessCode = 10000
	SuccessMsg  = "success"
)

// 参数
const (
	ParamVerifyErrorCode  = 20000 + iota // 提供参数有问题
	ParamMissingErrorCode                //参数缺失
)

// 鉴权
const (
	AuthInvalidCode        = 30000 + iota // 鉴权失败
	AuthAccessExpiredCode                 // 访问令牌过期
	AuthRefreshExpiredCode                // 刷新令牌过期
	AuthPermissionCode                    // 令牌等级不够，如stu无法进行审核
	AuthNoTokenCode                       // 没有 token
	AuthBlackListTokenCode                // 用户登出令牌已拉黑
)

// 业务错误
const (
	// user
	ServiceUserExistCode      = 40000 + iota
	ServiceEmailIncorrectCode // 邮箱格式不正确
	ServiceUserDeathCode      // 用户未激活、没有绑定相应邮箱
	ServiceUserPasswordError  // 密码错误
	ServiceUserNotExistCode
	ServiceCodeExpired    // 邮箱验证码已过期
	ServiceCodeNotMatched // 邮箱验证码不匹配
	ServiceEmailWaitCode  // 等待两分钟重发邮件
	// event
	ServiceEventNotExistCode  // 该赛事材料不存在
	ServiceImageNotAwardCode  // 判断上传的图片不是奖状或者荣誉证书
	ServiceEventUnChangedCode // 表示未经过申诉，无法直接修改该材料
	ServiceEventNotMatchCode  // 上传的材料没被认定
	ServiceNoAuthToDo         // 辅导员试图对不属于自己管辖的学生的材料进行审核或处理申诉
	ServiceRepeatAction       // 辅导员试图对不属于自己管辖的学生的材料进行审核或处理申诉
	// resultRecord
	ServiceRecordNotExistCode // 该记录不存在
	// appeal
	ServiceAppealNotExistCode  // 申诉记录不存在
	ServiceAppealExistCode     // 申诉已存在
	ServiceUserErrorAppealCode // 用户对不属于自己的材料进行申诉、对不属于自己的申诉记录进行撤销、查看不属于自己的申诉
	ServiceAppealUnchangedCode // 用户尝试对已经处理的申诉进行撤销
	// maintain
	ServiceCollegeNotExistCode
	ServiceCollegeExistCode
	ServiceMajorExistCode
	ServiceMajorNotExistCode
	ServiceGradeNotExistCode
	ServiceRecognizedNotExistCode // 赛事认定表中不存在该赛事
)

// 服务错误
const (
	InternalServiceErrorCode  = 50000 + iota // 内部服务错误
	InterFileProcessErrorCode                //文件处理错误
	InternalDatabaseErrorCode
	InternalRedisErrorCode // Redis错误
	InternalESErrorCode    // Redis错误
	InterConfigErrorCode
)
