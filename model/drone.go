package model

type PublishLog struct {
	Id         int32  `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	CommitId   string `gorm:"column:commit_id;type:varchar(300);default:'';comment:commit id;NOT NULL" json:"commit_id"`
	CommitLink string `gorm:"column:commit_link;type:varchar(300);default:'';comment:commit link;NOT NULL" json:"commit_link"`
	BuildLink  string `gorm:"column:build_link;type:varchar(300);default:'';comment:build link;NOT NULL" json:"build_link"`
	Author     string `gorm:"column:author;type:varchar(100);default:'';comment:作者信息;NOT NULL" json:"author"`
	Branch     string `gorm:"column:branch;type:varchar(100);default:'';comment:分支;NOT NULL" json:"branch"`
	Message    string `gorm:"column:message;type:varchar(100);default:'';comment:消息;NOT NULL" json:"message"`
	Event      string `gorm:"column:event;type:varchar(100);default:'';comment:事件;NOT NULL" json:"event"`
	Remark     string `gorm:"column:remark;type:varchar(3000);default:'';comment:备注;NOT NULL" json:"remark"`
	CreatedAt  int32  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  int32  `gorm:"column:updated_at;autoCreateTime"`
	DeletedAt  int32  `gorm:"column:deleted_at;default:0" json:"DeletedAt"`
}

// TableName 表名:artwork_flow，画作流程表。
func (PublishLog) TableName() string {
	return "publish_log"
}
