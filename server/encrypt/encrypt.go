package encrypt

import "github.com/daniel-dumitrascu/crypt/crypto"

type SecureCrypt struct {
	c crypto.Crypto
}

func CreateSecureCrypt() *SecureCrypt {
	sc := SecureCrypt{}
	sc.c = crypto.CreateCryptoAesGcm()
	return &sc
}

func (sc *SecureCrypt) Encrypt(inputPath *string, outputPath *string, keyPath *string) error {
	return sc.c.EncryptDir(*inputPath, *outputPath, *keyPath)
}

func (sc *SecureCrypt) Decrypt(inputPath *string, outputPath *string, keyPath *string) error {
	return sc.c.DecryptDir(*inputPath, *outputPath, *keyPath)
}

func (sc *SecureCrypt) Genkey() []byte {
	return []byte(sc.c.GenKey())
}
