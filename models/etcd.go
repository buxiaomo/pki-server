package models

type Etcd struct {
	Model
	Crt       string `gorm:"type:TEXT;not null" json:"crt"`
	CrtSha256 string `gorm:"type:TEXT;not null" json:"crt_sha256"`
	Key       string `gorm:"type:TEXT;not null" json:"key"`
	KeySha256 string `gorm:"type:TEXT;not null" json:"key_sha256"`
	ProjectID int
}

func (obj *Etcd) Update(pid int, p *Etcd) error {
	//var etcd *Etcd
	return db.Where("project_id = ?", pid).Updates(&p).Error
	//return db.Save(etcd).Error
}
