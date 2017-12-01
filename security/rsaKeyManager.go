package security

import (
				"os"
				"encoding/pem"
				"crypto/x509"
				"fmt"
				"io/ioutil"
				"crypto/rsa"
				"crypto/rand"
)

func InitializeRsaVault() bool {
		if _, err := os.Stat("/go/src/kalaxia-game-api/rsa_vault/private.key"); !os.IsNotExist(err) {
	  		return false
		}
		// generate private key
		privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			panic(err)
		}
		// extract public key
		publickey := &privatekey.PublicKey
		pubkey, _ := x509.MarshalPKIXPublicKey(publickey);
		// save private key
		pkey := x509.MarshalPKCS1PrivateKey(privatekey)
		ioutil.WriteFile("/go/src/kalaxia-game-api/rsa_vault/private.key", pkey, 0777)
		fmt.Println("private key saved to private.key")
		// save public key in PEM file
		pemfile, _ := os.Create("/go/src/kalaxia-game-api/rsa_vault/public.pub")
		var pemkey = &pem.Block{
								 Type : "PUBLIC KEY",
								 Bytes : pubkey}
		pem.Encode(pemfile, pemkey)
		pemfile.Close()
		return true
}

func Encrypt(data []byte) []byte {
	portalPEM, err := ioutil.ReadFile("/go/src/kalaxia-game-api/rsa_vault/portal_rsa.pub")
	if (err != nil) {
		panic(err)
	}
	block, _ := pem.Decode([]byte(portalPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	rsaPublicKey, _ := publicKey.(*rsa.PublicKey)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, data)
	if err != nil {
		panic(err)
	}
	return encryptedData
}

func Decrypt(data []byte) []byte {
	pkey, err := ioutil.ReadFile("/go/src/kalaxia-game-api/rsa_vault/private.key")
	if (err != nil) {
		panic(err)
	}
	privatekey, err := x509.ParsePKCS1PrivateKey(pkey)
	if (err != nil) {
		panic(err)
	}
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, data)
	if (err != nil) {
		panic(err)
	}
	return decryptedData
}
