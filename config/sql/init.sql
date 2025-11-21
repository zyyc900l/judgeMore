/*
 Navicat Premium Dump SQL
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 第一步：创建没有外键依赖的基础表
-- ----------------------------

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
                          `user_id` int NOT NULL AUTO_INCREMENT COMMENT '用户id，自增',
                          `user_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名称',
                          `user_role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '辅导员 / 学生/ 管理员',
                          `role_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学生学号，辅导员管理员工号',
                          `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮箱',
                          `college` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '学院',
                          `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '加密密码',
                          `major` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '专业',
                          `grade` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '年级',
                          `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态：0-未激活，1-使用中，2-已停用',
                          `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                          PRIMARY KEY (`user_id`) USING BTREE,
                          UNIQUE INDEX `role_id_unique`(`role_id` ASC) USING BTREE,
                          UNIQUE INDEX `email_unique`(`email` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for recognized_events (移除 rule_id 字段，避免循环引用)
-- ----------------------------
DROP TABLE IF EXISTS `recognized_events`;
CREATE TABLE `recognized_events`  (
                                      `recognized_event_id` int NOT NULL AUTO_INCREMENT COMMENT '主键，唯一标识每条认定赛事记录',
                                      `college` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '赛事认定的学院归属，支持按学院分类管理',
                                      `recognized_event_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '认定的竞赛名称，作为学生提交时的匹配依据',
                                      `organizer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '主办单位，明确赛事权威性',
                                      `recognized_event_time` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '竞赛举办时间周期，用于时效性验证',
                                      `related_majors` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '赛事涉及的专业范围',
                                      `applicable_majors` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '实际申请认定的专业',
                                      `recognition_basis` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '认定依据文件或标准',
                                      `recognized_level` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '官方认定的赛事级别（国际级|国家级|省级|校级）',
                                      `is_active` tinyint(1) NULL DEFAULT 1 COMMENT '是否生效状态，控制赛事是否可被选择',
                                      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                      `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                                      PRIMARY KEY (`recognized_event_id`) USING BTREE,
                                      INDEX `idx_college_level`(`college` ASC, `recognized_level` ASC) USING BTREE,
                                      INDEX `idx_event_name`(`recognized_event_name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 80800 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '学院认可的奖项表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for event_rules
-- ----------------------------
DROP TABLE IF EXISTS `event_rules`;
CREATE TABLE `event_rules`  (
                                `rule_id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                `recognized_event_id` int NULL COMMENT '关联认定赛事ID', -- 改为 NULL，避免默认值问题
                                `event_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '国际级|国家级|省级|校级',
                                `event_weight` decimal(5, 2) NOT NULL COMMENT '赛事权重系数',
                                `integral` int NOT NULL COMMENT '对应基础积分',
                                `rule_desc` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '规则说明',
                                `is_editable` tinyint NOT NULL COMMENT '1 - 是 / 0 - 否',
                                `is_active` tinyint(1) NULL DEFAULT 1 COMMENT  '是否启用 1 - 是 / 0 - 否',
                                `award_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '特等奖/一等奖/二等奖/三等奖/优秀奖',
                                `award_level_weight` decimal(5, 2) NULL DEFAULT NULL COMMENT '奖项权重系数',
                                `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                                PRIMARY KEY (`rule_id`) USING BTREE,
                                INDEX `idx_recognized_event_id`(`recognized_event_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 90900 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '赛事权重规则表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 第二步：创建依赖上述表的外键表
-- ----------------------------

-- ----------------------------
-- Table structure for student_events
-- ----------------------------
DROP TABLE IF EXISTS `student_events`;
CREATE TABLE `student_events`  (
                                   `event_id` int NOT NULL AUTO_INCREMENT COMMENT '赛事材料的自增id',
                                   `user_id` int NOT NULL COMMENT '关联的学生的用户id',
                                   `recognized_id` int NULL COMMENT '关联认定赛事ID',
                                   `event_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '赛事名称',
                                   `event_organizer` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '赛事主办方',
                                   `event_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '国际级|国家级|省级|校级',
                                   `award_level`  varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT'特等奖|一等奖|二等奖|三等奖|优秀奖',
                                   `award_content` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '奖项具体内容',
                                   `award_at` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '奖项时间',
                                   `material_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '材料上传路径',
                                   `material_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '待审核' COMMENT '待审核 / 已审核 / 驳回',
                                   `auto_extracted` tinyint NOT NULL DEFAULT 0 COMMENT '1 - 是 / 0 - 否',
                                   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                   `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`event_id`) USING BTREE,
                                   INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
                                   INDEX `idx_material_status`(`material_status` ASC) USING BTREE,
                                   INDEX `idx_recognized_id`(`recognized_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30000 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '学生赛事材料表' ROW_FORMAT = Dynamic;

DROP TABLE IF EXISTS `appeals`;
CREATE TABLE `appeals`  (
                            `appeal_id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                            `result_id` bigint NOT NULL COMMENT '申诉的结果ID',
                            `user_id` bigint NOT NULL COMMENT '申诉学生ID',
                            `appeal_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '申诉类型：分级异议/积分异议',
                            `appeal_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '申诉理由',
                            `attachment_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '佐证材料路径',
                            `status` enum('pending','approved','rejected') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'pending' COMMENT '状态：pending, approved, rejected,deleted',
                            `handled_by` varchar(50) NULL DEFAULT NULL COMMENT '处理人（辅导员)ID',
                            `handled_at` datetime NULL DEFAULT NULL COMMENT '处理时间',
                            `handled_result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '处理结果说明',
                            `appeal_count` int NOT NULL DEFAULT 1 COMMENT '该材料申诉次数',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`appeal_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20500 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '存储学生对审核结果的申诉记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 第三步：创建依赖最多的表
-- ----------------------------

-- ----------------------------
-- Table structure for integral_results
-- ----------------------------
DROP TABLE IF EXISTS `integral_results`;
CREATE TABLE `integral_results`  (
                                     `result_id` int NOT NULL AUTO_INCREMENT COMMENT '积分计算结果的自增主键',
                                     `event_id` int NOT NULL COMMENT '与这个积分关联的赛事id',
                                     `user_id` int NOT NULL COMMENT '与积分相关的学生的用户id',
                                     `rule_id` int NOT NULL COMMENT '与之相关的赛事权重规则id',
                                     `appeal_id` varchar(50) NULL DEFAULT NULL COMMENT '申诉ID',
                                     `final_integral` decimal(10, 2) NOT NULL COMMENT '最终积分',
                                     `status` enum('申诉中','正常','申诉完成') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '正常' COMMENT '状态：申诉中、正常、申诉完成',
                                     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                                     PRIMARY KEY (`result_id`) USING BTREE,
                                     INDEX `rule_id`(`rule_id` ASC) USING BTREE,
                                     INDEX `idx_user_id`(`user_id` ASC) USING BTREE,
                                     INDEX `idx_event_id`(`event_id` ASC) USING BTREE,
                                     INDEX `idx_appeal_id`(`appeal_id` ASC) USING BTREE,
                                     INDEX `idx_status`(`status` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10500 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '积分计算结果表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for feedbacks
-- ----------------------------
DROP TABLE IF EXISTS `feedbacks`;
CREATE TABLE `feedbacks`  (
                              `feedback_id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                              `user_id` bigint NOT NULL COMMENT '提交用户ID',
                              `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '反馈类型',
                              `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '反馈内容',
                              `is_replied` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已回复',
                              `reply_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '回复内容',
                              `created_time` datetime NOT NULL COMMENT '提交时间',
                              PRIMARY KEY (`feedback_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '存储用户提交的系统使用反馈' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for system_logs
-- ----------------------------
DROP TABLE IF EXISTS `system_logs`;
CREATE TABLE `system_logs`  (
                                `log_id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                                `user_id` bigint NOT NULL COMMENT '操作用户',
                                `action` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作类型',
                                `target_table` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作表名',
                                `target_id` bigint NOT NULL COMMENT '操作记录ID',
                                `old_value` json NULL COMMENT '旧值（JSON格式）',
                                `new_value` json NULL COMMENT '新值（JSON格式）',
                                `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '操作IP',
                                `created_time` datetime NOT NULL COMMENT '操作时间',
                                PRIMARY KEY (`log_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '记录关键操作日志，用于审计与追踪' ROW_FORMAT = Dynamic;
-- ----------------------------
-- Table structure for college
-- ----------------------------
DROP TABLE IF EXISTS `college`;
CREATE TABLE `college`  (
                            `college_id` int NOT NULL AUTO_INCREMENT COMMENT '学院编号',
                            `college_name` varchar(100) NOT NULL COMMENT '学院名称',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                            PRIMARY KEY (`college_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 100800 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '学院表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for major
-- ----------------------------
DROP TABLE IF EXISTS `major`;
CREATE TABLE `major`  (
                          `major_id` bigint NOT NULL AUTO_INCREMENT COMMENT '专业ID',
                          `college_id` int NOT NULL COMMENT '所属学院ID',
                          `major_name` varchar(100) NOT NULL COMMENT '专业名称',
                          `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                          PRIMARY KEY (`major_id`) USING BTREE,
                          INDEX `idx_college_id`(`college_id` ASC) USING BTREE,
                          CONSTRAINT `major_ibfk_1` FOREIGN KEY (`college_id`) REFERENCES `college` (`college_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 200800 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '专业表' ROW_FORMAT = Dynamic;
SET FOREIGN_KEY_CHECKS = 1;

DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`  (
                             `relation_id` bigint NOT NULL AUTO_INCREMENT COMMENT '关系ID',
                             `user_id` int NOT NULL COMMENT '用户ID',
                             `major` varchar(100) NOT NULL COMMENT '专业',
                             `major_id` varchar(100) NOT NULL COMMENT '专业id',
                             `grade` varchar(50) NOT NULL COMMENT '年级',
                             `college_id` varchar(100) NOT NULL COMMENT '学院id',
                             `college` varchar(100) NOT NULL COMMENT '学院',
                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                             PRIMARY KEY (`relation_id`) USING BTREE,
                             INDEX `idx_user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 80500 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户关系表' ROW_FORMAT = Dynamic;


DROP TABLE IF EXISTS `admin_object`;
CREATE TABLE `admin_object`  (
                             `admin_id` bigint NOT NULL COMMENT '管理员ID',
                             `stu_id` int NOT NULL COMMENT '管理用户ID',
                             `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
                             PRIMARY KEY (`admin_id`, `stu_id`) USING BTREE,
                             INDEX `idx_stu_id`(`stu_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '管理员-学生关系表' ROW_FORMAT = Dynamic;