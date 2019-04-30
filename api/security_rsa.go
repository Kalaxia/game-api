package api

import (
	"fmt"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
)

func InitRsaVault() {
	if _, err := os.Stat("/go/src/kalaxia-game-api/rsa_vault/private.key"); !os.IsNotExist(err) {
  		return
	}
	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(NewException("Private key could not be generated", err))
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
		 Bytes : pubkey,
 	}
	pem.Encode(pemfile, pemkey)
	pemfile.Close()
}

func Encrypt(data []byte) ([]byte, string, string) {
	portalPEM, err := ioutil.ReadFile("/go/src/kalaxia-game-api/rsa_vault/portal_rsa.pub")
	if (err != nil) {
		panic(NewException("Could not read portal public key", err))
	}
	block, _ := pem.Decode([]byte(portalPEM))
	if block == nil {
		panic(NewException("failed to parse PEM block containing the public key", nil))
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(NewException("failed to parse DER encoded public key", err))
	}
	rsaPublicKey, _ := publicKey.(*rsa.PublicKey)

	data, key, iv := encryptAesPayload(data)

	cipherKey, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, key)
	if err != nil {
		panic(NewException("Could not encrypt cipher key", err))
	}
	cipherIv, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, iv)
	if err != nil {
		panic(NewException("Could not encrypt cipher IV", err))
	}
	return data, base64.StdEncoding.EncodeToString(cipherKey), base64.StdEncoding.EncodeToString(cipherIv)
}

func encryptAesPayload(data []byte) ([]byte, []byte, []byte) {
	key := make([]byte, 32)
	iv := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		panic(NewException("Could not generate random key", err))
	}
	_, err = rand.Read(iv)
	if err != nil {
		panic(NewException("Could not generate random IV", err))
	}
	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. If the block is incomplete, we add padding to it
	if len(data) % aes.BlockSize != 0 {
		data = pkcs7Pad(data)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(NewException("Could not create AES cipher", err))
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(data, data)

	return data, key, iv
}

func Decrypt(encryptedKey string, encryptedIv string, data []byte) []byte {
	pkey, err := ioutil.ReadFile("/go/src/kalaxia-game-api/rsa_vault/private.key")
	if err != nil {
		panic(NewException("Private key file could not be opened", err))
	}
	privatekey, err := x509.ParsePKCS1PrivateKey(pkey)
	if err != nil {
		panic(NewException("Private key could not be parsed", err))
	}
	key, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		panic(NewException("Encrypted key could not be decoded", err))
	}
	iv, err := base64.StdEncoding.DecodeString(encryptedIv)
	if err != nil {
		panic(NewException("Encrypted IV could not be decoded", err))
	}
	aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, key)
	if err != nil {
		panic(NewException("AES key could not be decrypted", err))
	}
	aesIv, err := rsa.DecryptPKCS1v15(rand.Reader, privatekey, iv)
	if err != nil {
		panic(NewException("AES IV could not be decrypted", err))
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(NewException("AES cipher could not be created", err))
	}
	mode := cipher.NewCBCDecrypter(block, aesIv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(data, data)
 	return pkcs7Unpad(data)
}

// Appends padding.
func pkcs7Pad(data []byte) []byte {
    padlen := 1
    for ((len(data) + padlen) % aes.BlockSize) != 0 {
        padlen = padlen + 1
    }

    pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
    return append(data, pad...)
}

// Returns slice of the original data without padding.
func pkcs7Unpad(data []byte) []byte {
    return data[:len(data) - int(data[len(data)-1])]
}
