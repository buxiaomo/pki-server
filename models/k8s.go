package models

type K8s struct {
	Model
	Crt       string `gorm:"type:TEXT;not null" json:"crt"`
	CrtSha256 string `gorm:"not null" json:"crt_sha256"`
	Key       string `gorm:"type:TEXT;not null" json:"key"`
	KeySha256 string `gorm:"not null" json:"key_sha256"`
	ProjectID int
}

func (obj *K8s) Update(pid int, p *K8s) error {
	return db.Where("project_id = ?", pid).Updates(&p).Error
}
