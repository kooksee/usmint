package cmn

import (
	"crypto/rsa"
	"crypto/rand"
	"os"
	"encoding/pem"
	"crypto/x509"
	"errors"
)

func NewRsa() *Rsa {
	return &Rsa{}
}

type Rsa struct {
	pubkey  []byte
	priVKey []byte
}

// GenRsaKey 生成私钥文件
func (r *Rsa) GenKey(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	privFile, err := os.Create("private.pem")
	defer privFile.Close()
	if err != nil {
		return err
	}

	if err = pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	pubFile, err := os.Create("public.pem")
	defer pubFile.Close()
	if err != nil {
		return err
	}

	if err = pem.Encode(pubFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}); err != nil {
		return err
	}

	return nil
}

// 加密
func (r *Rsa) RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(r.pubkey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func (r *Rsa) RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(r.priVKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
