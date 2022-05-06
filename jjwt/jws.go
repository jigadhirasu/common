package jjwt

import (
	"crypto/rsa"
	"encoding/json"

	"gopkg.in/square/go-jose.v2"
)

// var jwsPrivateKey *rsa.PrivateKey

// func LoadJWSPrivateKey() {
// 	dbr := Redis()
// 	if count := dbr.Exists(dbr.Context(), "jws:pair").Val(); count < 1 {
// 		jwsRSApair(dbr)
// 	}
// 	b, _ := dbr.HGet(dbr.Context(), "jws:pair", "private").Bytes()
// 	block, _ := pem.Decode(b)
// 	jwsPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
// }

// var jwsPublicKey *rsa.PublicKey

// func LoadJWSPublicKey() {
// 	dbr := Redis()
// 	b, _ := dbr.HGet(dbr.Context(), "jws:pair", "public").Bytes()
// 	block, _ := pem.Decode(b)
// 	jwsPublicKey, _ = x509.ParsePKCS1PublicKey(block.Bytes)
// }

func JWSSign(privateKey *rsa.PrivateKey, claim *Claim) string {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateKey}, nil)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(claim)
	if err != nil {
		panic(err)
	}

	object, err := signer.Sign(b)
	if err != nil {
		panic(err)
	}

	bearer, _ := object.CompactSerialize()
	return bearer
}

func JWSVerify(publicKey *rsa.PublicKey, bearer string) (*Claim, error) {
	object, err := jose.ParseSigned(bearer)
	if err != nil {
		return nil, err
	}
	output, err := object.Verify(publicKey)
	if err != nil {
		return nil, err
	}

	cl := &Claim{}
	if err := json.Unmarshal(output, cl); err != nil {
		return nil, err
	}

	return cl, cl.Valid()
}
