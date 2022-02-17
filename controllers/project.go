package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	db "pki-server/models"
	"pki-server/pkg/pki"
	"time"
)

type getK8SCrtArg struct {
	Name     string `uri:"name" binding:"required"`
	Env      string `uri:"env" binding:"required"`
	Filename string `uri:"filename" binding:"required"`
}

func GetProjectFile(c *gin.Context) {
	var arg getK8SCrtArg
	if err := c.ShouldBindUri(&arg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	var project db.Project
	if t, err := project.Select(arg.Name, arg.Env); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"msg":  err.Error(),
				"data": nil,
			})
			return
		}
	} else {
		var result string
		switch arg.Filename {
		case "ca.crt":
			result = t.K8s.Crt
		case "ca.crt.sha256sum":
			result = t.K8s.CrtSha256 + " *ca.crt"

		case "ca.key":
			result = t.K8s.Key
		case "ca.key.sha256sum":
			result = t.K8s.KeySha256 + " *ca.key"

		case "etcd.crt":
			result = t.Etcd.Crt
		case "etcd.crt.sha256sum":
			result = t.Etcd.CrtSha256 + " *etcd.crt"

		case "etcd.key":
			result = t.Etcd.Key
		case "etcd.key.sha256sum":
			result = t.Etcd.KeySha256 + " *etcd.key"

		case "front-proxy-ca.crt":
			result = t.Frontproxy.Crt
		case "front-proxy-ca.crt.sha256sum":
			result = t.Frontproxy.CrtSha256 + " *front-proxy-ca.crt"

		case "front-proxy-ca.key":
			result = t.Frontproxy.Key
		case "front-proxy-ca.key.sha256sum":
			result = t.Frontproxy.KeySha256 + " *front-proxy-ca.key"

		case "sa.pub":
			result = t.Sa.Pub
		case "sa.pub.sha256sum":
			result = t.Sa.PubSha256 + " *sa.pub"

		case "sa.key":
			result = t.Sa.Key
		case "sa.key.sha256sum":
			result = t.Sa.KeySha256 + " *sa.key"
		default:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		res.text(c, result)
	}
	return
}

type signatureK8SCrtArg struct {
	Name string        `json:"name" binding:"required"`
	Env  string        `json:"env" binding:"required"`
	Year time.Duration `json:"year" binding:"required"`
}

func SignatureProject(c *gin.Context) {
	var arg signatureK8SCrtArg
	if err := c.ShouldBindJSON(&arg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	p := pki.New()

	var project db.Project
	project.Name = arg.Name
	project.Env = arg.Env
	project.NotBefore = time.Now()
	project.NotAfter = project.NotBefore.Add(arg.Year * 365 * 24 * time.Hour)

	var sa db.Sa
	sa.Key, sa.KeySha256, sa.Pub, sa.PubSha256 = p.GenRsaKey()
	project.Sa = sa

	var k8s db.K8s
	k8s.Crt, k8s.CrtSha256, k8s.Key, k8s.KeySha256 = p.Signature("kubernetes", project.NotBefore, project.NotAfter)
	project.K8s = k8s

	var etcd db.Etcd
	etcd.Crt, etcd.CrtSha256, etcd.Key, etcd.KeySha256 = p.Signature("etcd-ca", project.NotBefore, project.NotAfter)
	project.Etcd = etcd

	var frontproxy db.Frontproxy
	frontproxy.Crt, frontproxy.CrtSha256, frontproxy.Key, frontproxy.KeySha256 = p.Signature("front-proxy-ca", project.NotBefore, project.NotAfter)
	project.Frontproxy = frontproxy

	if err := project.Add(); err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	res.json(c, nil, project.ID)
	return
}

type removeK8SCrtArg struct {
	Name string `uri:"name" binding:"required"`
	Env  string `uri:"env" binding:"required"`
}

func RemoveProject(c *gin.Context) {
	var arg removeK8SCrtArg
	if err := c.ShouldBindUri(&arg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	var cert db.Project
	if err := cert.Del(arg.Name, arg.Env); err != nil {
		res.json(c, err.Error(), nil)
		return
	}
	res.json(c, nil, nil)
	return
}

func GetAllProject(c *gin.Context) {
	var cert db.Project
	crt, err := cert.GetAll()
	if err != nil {
		res.json(c, err.Error(), nil)
		return
	}
	res.json(c, nil, crt)
	return
}

type renewalK8SCrtArg struct {
	Name string `uri:"name" binding:"required"`
	Env  string `uri:"env" binding:"required"`
	Year uint   `uri:"year" binding:"required"`
}

func RenewalProject(c *gin.Context) {
	var arg renewalK8SCrtArg
	if err := c.ShouldBindUri(&arg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	p := pki.New()

	project := db.Project{
		Name: arg.Name,
		Env:  arg.Env,
	}
	project.NotBefore = time.Now()
	project.NotAfter = project.NotBefore.Add(time.Duration(arg.Year) * 365 * 24 * time.Hour)
	pid := project.Update(&project)
	fmt.Println(project)

	var k8s db.K8s
	k8s.Crt, k8s.CrtSha256, k8s.Key, k8s.KeySha256 = p.Signature("kubernetes", project.NotBefore, project.NotAfter)
	k8s.Update(pid, &k8s)

	var etcd db.Etcd
	etcd.Crt, etcd.CrtSha256, etcd.Key, etcd.KeySha256 = p.Signature("etcd-ca", project.NotBefore, project.NotAfter)
	etcd.Update(pid, &etcd)

	var frontproxy db.Frontproxy
	frontproxy.Crt, frontproxy.CrtSha256, frontproxy.Key, frontproxy.KeySha256 = p.Signature("front-proxy-ca", project.NotBefore, project.NotAfter)
	frontproxy.Update(pid, &frontproxy)

	res.json(c, nil, nil)
}
