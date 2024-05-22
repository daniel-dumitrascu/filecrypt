package main

import "crypt/crypto"

func main() {
	var c crypto.Crypto = crypto.CryptoAesGcm{}

	c.EncryptFile("C:\\Users\\DanielDumitrascu\\Desktop\\1\\file.xlsx", "C:\\Users\\DanielDumitrascu\\Desktop\\1\\", "C:\\Home\\filecrypt\\new_key\\key")
}
