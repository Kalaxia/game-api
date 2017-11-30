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
		if _, err := os.Stat("/src/go/kalaxia-game-api/rsa_vault/private.key"); os.IsExist(err) {
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

func Decrypt(data []byte) []byte {
	pkey, err := ioutil.ReadFile("/go/src/kalaxia-game-api/rsa_vault/private.key")
	if (err != nil) {
		panic(err)
	}
	privatekey, err := x509.ParsePKCS1PrivateKey(pkey)
	if (err != nil) {
		panic(err)
	}
	finalData, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, data)
	if (err != nil) {
		panic(err)
	}
	return finalData
}
