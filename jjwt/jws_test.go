package jjwt_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/jigadhirasu/common/jjwt"
	"github.com/jigadhirasu/common/jlog"
)

func TestPublicKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	fmt.Println(privateKey.PublicKey)
	pk := &pem.Block{
		Type:  "JWS PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	bf := &bytes.Buffer{}
	pem.Encode(bf, pk)
	// pem Public Key format
	fmt.Println(bf.String())

	// decode Public Key
	pg, _ := pem.Decode(bf.Bytes())
	ppb, _ := x509.ParsePKCS1PublicKey(pg.Bytes)
	fmt.Println(ppb)
}

func TestJWS(t *testing.T) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey

	cl := jjwt.NewClaim([]byte(`{"UUID": "9790e576-d5b6-4ca2-8a0e-a6a4f9c690af", "Type": "isystem", "Agent": "cc373e1e-c9b8-4469-a4b3-ccb9e66011f6"}`))
	cl.WithExp(time.Now().Unix())
	serialize := jjwt.JWSSign(privateKey, cl)
	serialize = "Bearer " + serialize

	jlog.Debug(serialize)

	a := "Bearer"
	serialize = serialize[len(a):]
	jlog.Debug(serialize)

	rcl, err := jjwt.JWSVerify(&publicKey, serialize)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(rcl.Data.String())
}
