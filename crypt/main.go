package main

import (
	"fmt"
	"os"

	"github.com/daniel-dumitrascu/crypt/crypto"
)

func main() {
	args := os.Args[1:]
	argsCount := len(args)
	if argsCount == 0 {
		help()
		return
	}

	action := args[0]
	var c crypto.Crypto = crypto.CreateCryptoAesGcm()

	if action == "genkey" {
		generatedkey := c.GenKey()
		fmt.Println(generatedkey)
		return
	} else {
		if argsCount != 4 {
			help()
			return
		}

		keypath := args[1]
		inputpath := args[2]
		outputpath := args[3]

		if action == "encrypt" {
			if err := c.EncryptDir(inputpath, outputpath, keypath); err != nil {
				fmt.Println(err)
			}
		} else if action == "decrypt" {
			if err := c.DecryptDir(inputpath, outputpath, keypath); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Action is unknown: ", action)
			help()
			return
		}
	}
}

func help() {
	fmt.Println("Description:")
	fmt.Println("	With this tool you can generate symetric keys and encrypt/decrypt files or directories.")

	fmt.Printf("\n")

	fmt.Println("Arguments:")
	fmt.Println("	action - can be one of those 3 {genkey,encrypt,decrypt}")
	fmt.Println("	key path - is the path pointing to the key used for encrypting/decrypting")
	fmt.Println("	file or dir path - the path to the target file or directory")
	fmt.Println("	output path - the path where the toll saved the created files")

	fmt.Printf("\n")

	fmt.Println("Examples:")
	fmt.Println("	crypt.exe encrypt 'path_to_key' 'path_to_target' 'output_path'")
	fmt.Println("	crypt.exe genkey")
}
