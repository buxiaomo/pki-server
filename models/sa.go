package models

type Sa struct {
	Model
	Pub       string `gorm:"type:TEXT;not null" json:"pub"`
	PubSha256 string `gorm:"type:TEXT;not null" json:"pub_sha256"`
	Key       string `gorm:"type:TEXT;not null" json:"key"`
	KeySha256 string `gorm:"type:TEXT;not null" json:"key_sha256"`
	ProjectID int
}
