package models

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Project struct {
	Model
	Name       string     `gorm:"size:25;uniqueIndex:idx_first_second" json:"name" binding:"required"`
	Env        string     `gorm:"size:25;uniqueIndex:idx_first_second" json:"env" binding:"required"`
	NotBefore  time.Time  `gorm:"" json:"not_before"`
	NotAfter   time.Time  `gorm:"" json:"not_after"`
	K8s        K8s        `json:"k8s"`
	Etcd       Etcd       `json:"etcd"`
	Frontproxy Frontproxy `json:"frontproxy"`
	Sa         Sa         `json:"sa"`
}

func (obj *Project) Add() error {
	return db.Create(&obj).Error
}

func (obj *Project) Exist(name string, env string) bool {
	var cert *Project
	err := db.Where("name = ? and env = ?", name, env).First(&cert)
	return errors.Is(err.Error, gorm.ErrRecordNotFound)
}

func (obj *Project) Select(name string, env string) (cert *Project, err error) {
	return cert, db.Where("name = ? and env = ?", name, env).Joins("K8s").Joins("Etcd").Joins("Frontproxy").Joins("Sa").First(&cert).Error
}

func (obj *Project) Del(name, env string) (err error) {
	var cert *Project
	db.Where("name = ? and env = ?", name, env).First(&cert)
	return db.Select(clause.Associations).Where("name = ? and env = ?", name, env).Delete(&cert).Error
}

func (obj *Project) GetAll() (cert []*Project, err error) {
	err = db.Joins("K8s").Joins("Etcd").Joins("Frontproxy").Joins("Sa").Find(&cert).Error
	return
}

func (obj *Project) Update(p *Project) int {
	var project *Project
	db.Where("name = ? and env = ?", obj.Name, obj.Env).First(&project)
	project.NotAfter = p.NotAfter
	project.NotBefore = p.NotBefore
	db.Updates(&project)
	return project.ID
}
