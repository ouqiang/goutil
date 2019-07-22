package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

type Rsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRsa 创建Rsa
func NewRsa(publicKeyReader io.Reader, privateKeyReader io.Reader) (*Rsa, error) {
	r := &Rsa{}
	var err error

	if publicKeyReader != nil {
		r.publicKey, err = r.parsePublicKey(publicKeyReader)
		if err != nil {
			return nil, err
		}
	}
	if privateKeyReader !=  nil {
                r.privateKey, err = r.parsePrivateKey(privateKeyReader)
                if err != nil {
                        return nil, err
                }
        }

	return r, nil
}

// RsaEncrypt Rsa加密
func (r *Rsa) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, msg)
}

// Decrypt Rsa解密
func (r *Rsa) Decrypt(cipherText []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, cipherText)
}

// 解析公钥
func (r *Rsa) parsePublicKey(reader io.Reader) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("pem decode: invalid public key")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parse public key: invalid public key")
	}

	return pub, nil
}

// 解析私钥
func (r *Rsa) parsePrivateKey(reader io.Reader) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("pem decode: invalid private key")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return private, err
}
