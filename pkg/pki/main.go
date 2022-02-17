package pki

import "time"

type service struct{}

func New() Service {
	return &service{}
}

type Service interface {
	Signature(CN string, NotBefore, NotAfter time.Time) (crt, crtsha256, key, keysha256 string)
	GenRsaKey() (prvkey, prvkeysha256, pubkey, pubkeysha256 string)
}

type Certificate struct {
	NotBefore     time.Time
	NotAfter      time.Time
	KubeCrt       string
	KubeKey       string
	EtcdCrt       string
	EtcdKey       string
	FrontProxyCrt string
	FrontProxyKey string
	SaKey         string
	SaPub         string
}
