package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"log"
	"math/big"
)

type KeyStruct struct {
	X, Y *big.Int
}

func GenKeyPair() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.New("error while generating key pair")
	}

	return privateKey, nil
}

func Hash(val []byte) []byte {
	sum := sha256.Sum256(val)
	return sum[:]
}

func GenWalletAddress(privateKey *ecdsa.PrivateKey) string {
	pubKeyOrig := privateKey.PublicKey

	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)

	version := byte(10)
	pubKeyAddressBytes := append([]byte{version}, pubKey...)

	return hex.EncodeToString(pubKeyAddressBytes)
}

func GetPublicKeyFromAddress(hexAddress string) *ecdsa.PublicKey {
	strBytes, err := hex.DecodeString(hexAddress)
	if err != nil {
		return nil
	}
	bts := strBytes[1:]
	xBytes := bts[:len(bts)/2]
	yBytes := bts[len(bts)/2:]

	x := (&big.Int{}).SetBytes(xBytes)
	y := (&big.Int{}).SetBytes(yBytes)
	return &ecdsa.PublicKey{
		X:     x,
		Y:     y,
		Curve: elliptic.P256(),
	}
}

func VerifySignature(address string, hash []byte, sign []byte) bool {
	pubKey := GetPublicKeyFromAddress(address)
	if pubKey == nil {
		log.Println("Pub key problem")
	}
	log.Println(pubKey, "here it is")
	return ecdsa.VerifyASN1(pubKey, hash, sign)
}

func Id() string {
	return uuid.New().String()
}
