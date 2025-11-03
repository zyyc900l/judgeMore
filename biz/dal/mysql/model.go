package mysql

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId    int64
	RoleId    string //实际上是我们业务过程中区分用户的主键
	UserName  string
	UserRole  string
	College   string
	Grade     string
	Major     string
	Email     string
	Status    int
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Event struct {
	EventId        string `gorm:"primaryKey;autoIncrement:true;column:event_id"`
	UserId         string `gorm:"not null;column:user_id"`
	RecognizedId   string `gorm:"not null;column:recognized_id"`
	EventName      string `gorm:"size:200;not null;column:event_name"`
	EventOrganizer string `gorm:"size:200;not null;column:event_organizer"`
	EventLevel     string `gorm:"size:20;not null;column:event_level"`
	AwardLevel     string `gorm:"type:enum('特等奖','一等奖','二等奖','三等奖','优秀奖');not null;column:award_level"`
	AwardContent   string `gorm:"size:100;column:award_content"`
	MaterialUrl    string `gorm:"size:500;not null;column:material_url"`
	MaterialStatus string `gorm:"size:20;not null;default:'待审核';column:material_status"`
	AutoExtracted  bool   `gorm:"not null;default:false;column:auto_extracted"`
	AwardAt        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
type RecognizedEvent struct {
	RecognizedEventId   string `gorm:"primaryKey;autoIncrement:true;column:recognized_event_id"`
	College             string `gorm:"size:255;not null;column:college"`
	RecognizedEventName string `gorm:"size:255;not null;column:recognized_event_name"`
	Organizer           string `gorm:"size:255;not null;column:organizer"`
	RecognizedEventTime string `gorm:"size:50;not null;column:recognized_event_time"`
	RelatedMajors       string `gorm:"size:255;column:related_majors"`
	ApplicableMajors    string `gorm:"size:255;column:applicable_majors"`
	RecognitionBasis    string `gorm:"size:255;column:recognition_basis"`
	RecognizedLevel     string `gorm:"size:50;not null;column:recognized_level"`
	IsActive            bool   `gorm:"default:true;column:is_active"`
	RuleId              int64  `gorm:"not null;column:rule_id"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

type EventRule struct {
	RuleId            string         `gorm:"primaryKey;autoIncrement:true;column:rule_id"`
	RecognizedEventId string         `gorm:"not null;default:0;column:recognized_event_id"`
	EventLevel        string         `gorm:"size:20;not null;column:event_level"`
	EventWeight       float64        `gorm:"type:decimal(5,2);not null;column:event_weight"`
	Integral          int64          `gorm:"not null;column:integral"`
	RuleDesc          string         `gorm:"size:500;column:rule_desc"`
	IsEditable        bool           `gorm:"not null;column:is_editable"`
	AwardLevel        string         `gorm:"size:20;column:award_level"`
	AwardLevelWeight  float64        `gorm:"type:decimal(5,2);column:award_level_weight"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

type ScoreResult struct {
	ResultId      string         `gorm:"primaryKey;autoIncrement:true;column:result_id"`
	EventId       string         `gorm:"not null;column:event_id"`
	UserId        string         `gorm:"not null;column:user_id"`
	RuleId        string         `gorm:"not null;column:rule_id"`
	AppealId      string         `gorm:"column:appeal_id"`
	FinalIntegral float64        `gorm:"type:decimal(10,2);not null;column:final_integral"`
	Status        string         `gorm:"type:enum('申诉中','正常','申诉完成');not null;default:'正常';column:status"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
