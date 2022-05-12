package jjwt

import (
	"time"

	"github.com/jigadhirasu/common/j"
	"github.com/kubemq-io/kubemq-go/pkg/uuid"
)

func NewClaim(data j.Bytes) *Claim {
	return &Claim{
		Iat:  time.Now().Unix(),
		Exp:  time.Now().Add(time.Hour * 3).Unix(),
		Iss:  "asgame",
		Jti:  uuid.New(),
		Data: data,
	}
}

type Claim struct {
	Iat  int64   `json:"iat,omitempty"` // 頒發時間
	Nbf  int64   `json:"nbf,omitempty"` // 生效時間
	Exp  int64   `json:"exp,omitempty"` // 到期時間
	Iss  string  `json:"iss,omitempty"` // 頒發者
	Jti  string  `json:"jti,omitempty"` // token編號
	Data j.Bytes `json:"data,omitempty"`
}

// default time.Now().Unix()
func (cl *Claim) WithIat(iat int64) *Claim {
	cl.Iat = iat
	return cl
}

// default 0
func (cl *Claim) WithNbf(nbf int64) *Claim {
	cl.Nbf = nbf
	return cl
}

// default 3h
func (cl *Claim) WithExp(exp int64) *Claim {
	cl.Exp = exp
	return cl
}

// default asgame
func (cl *Claim) WithIss(iss string) *Claim {
	cl.Iss = iss
	return cl
}

// default uuid.New()
func (cl *Claim) WithJti(jti string) *Claim {
	cl.Jti = jti
	return cl
}

func (cl Claim) Valid() error {
	if time.Now().Unix() < cl.Nbf {
		return j.NewErr(401, "還沒到生效時間再等等吧")
	}
	if cl.Exp > 0 && cl.Exp < time.Now().Unix() {
		return j.NewErr(401, "過期了重新要一個吧")
	}
	if cl.Iss != "asgame" {
		return j.NewErr(401, "這不是我發的,別亂喔")
	}
	// data裡面的內容個應用程式自己檢查
	return nil
}
