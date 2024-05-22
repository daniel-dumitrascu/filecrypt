package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
)

type CryptoAesGcm struct {
}

func (c CryptoAesGcm) EncryptFile(filepath string, outputpath string, keypath string) error {
	// Load key if it exists
	key, err := loadAndDecodeKey(keypath)
	if err != nil {
		fmt.Println("Issue when loading the key.")
		return err
	}

	return encryptFile(filepath, outputpath, key)
}

func (c CryptoAesGcm) DecryptFile(filepath string, outputpath string, keypath string) error {
	// Load key if it exists
	key, err := loadAndDecodeKey(keypath)
	if err != nil {
		fmt.Println("Issue when loading the key.")
		return err
	}

	return decryptFile(filepath, outputpath, key)
}

func (c CryptoAesGcm) EncryptDir(dirpath string, outputpath string, keypath string) error {
	// Load key if it exists
	key, err := loadAndDecodeKey(keypath)
	if err != nil {
		fmt.Println("Issue when loading the key.")
		return err
	}

	return encryptDir(dirpath, outputpath, key)
}

func (c CryptoAesGcm) DecryptDir(dirpath string, outputpath string, keypath string) error {
	// Load key if it exists
	key, err := loadAndDecodeKey(keypath)
	if err != nil {
		fmt.Println("Issue when loading the key.")
		return err
	}

	return decryptDir(dirpath, outputpath, key)
}

func (c CryptoAesGcm) GenKey() []byte {
	// Generate a random symmetric key for HMAC and AES
	key := make([]byte, key32size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		fmt.Println("Error generating symmetric key: ", err)
		return nil
	}
	return key
}

func loadAndDecodeKey(keypath string) ([]byte, error) {
	key, err := os.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	key, err = base64.StdEncoding.DecodeString(string(key))
	if err != nil {
		return nil, err
	}

	return key, nil
}

func encryptFile(filepath string, outputpath string, key []byte) error {
	// Load cipher
	aesgcm, nonce, err := loadCipherForEncryption(key)
	if err != nil {
		fmt.Println("Error when loading the cipher.")
		return err
	}

	// Open the file that will be encrypted
	sourceFileHandler, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the target file.")
		return err
	}
	defer sourceFileHandler.Close()

	// Open the file that will store the encrypted data
	secretFilepath := getEncryptedFilepath(filepath, outputpath)
	destFileHandler, err := os.OpenFile(secretFilepath, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the destination file.")
		return err
	}
	defer destFileHandler.Close()

	// Calculate the number of chunk data based on the chunk size
	fileSize, err := getFileSize(filepath)
	if err != nil {
		fmt.Println("Error getting the file size of the target file.")
		return err
	}

	chunkSize := chunk100mb
	chunckNr := int64(math.Ceil(float64(fileSize) / float64(chunkSize)))
	chunckBuffer := make([]byte, chunkSize)

	// Start reading the source file
	readSize, err := sourceFileHandler.Read(chunckBuffer)
	if err != nil {
		fmt.Println("Error reading file.")
		return err
	}

	chunkIndex := 1
	for readSize > 0 {
		encryptedChunck := encryptDataBlock(chunckBuffer[:readSize], aesgcm, nonce)
		fmt.Println("Encrypting data chunk", chunkIndex, " of ", chunckNr, " (read ", len(chunckBuffer), " bytes) (encrypted ", len(encryptedChunck), " bytes)")
		chunkIndex++

		if _, err = destFileHandler.Write(encryptedChunck); err != nil {
			fmt.Println("Error writing encrypted data: ", err)
			return err
		}
		if err = destFileHandler.Sync(); err != nil {
			fmt.Println("Error syncking data to file: ", err)
			return err
		}

		readSize, _ = sourceFileHandler.Read(chunckBuffer)
	}

	return nil
}

func decryptFile(filepath string, outputpath string, key []byte) error {
	// Load cipher
	aesgcm, err := loadCipherForDecryption(key)
	if err != nil {
		fmt.Println("Error when loading the cipher.")
		return err
	}

	// Open the file that will be decrypted
	sourceFileHandler, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the target file.")
		return err
	}
	defer sourceFileHandler.Close()

	// Open the file that will store the decrypted data
	clearFilepath := getDecryptedFilepath(filepath, outputpath)
	destFileHandler, err := os.OpenFile(clearFilepath, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the destination file.")
		return err
	}
	defer destFileHandler.Close()

	fileSize, err := getFileSize(filepath)
	if err != nil {
		fmt.Println("Error getting the file size of the target file.")
		return err
	}

	chunkSize := chunk100mb
	bytesToRead := 28 + chunkSize
	chunckNr := int64(math.Ceil(float64(fileSize) / float64(bytesToRead)))
	chunckBuffer := make([]byte, bytesToRead)
	dataReadSize, err := sourceFileHandler.Read(chunckBuffer)
	if err != nil {
		fmt.Println("Error reading file.")
		return err
	}
	chunckIndex := 1

	for dataReadSize > 0 {
		decryptedChunck, err := decodeDataBlock(chunckBuffer[:dataReadSize], aesgcm)
		if err != nil {
			fmt.Println("Error during data chunck decryption.")
			return err
		}
		fmt.Println("Decrypting data chunk", chunckIndex, " of ", chunckNr, " (read ", len(chunckBuffer), " bytes) (decrypted ", len(decryptedChunck), " bytes)")
		chunckIndex++

		if _, err = destFileHandler.Write(decryptedChunck); err != nil {
			fmt.Println("Error writing decrypted data.")
			return err
		}
		if err = destFileHandler.Sync(); err != nil {
			fmt.Println("Error syncking data to file.")
			return err
		}

		dataReadSize, _ = sourceFileHandler.Read(chunckBuffer)
	}

	return nil
}

func encryptDir(dirpath string, outputpath string, key []byte) error {
	info, err := os.Stat(dirpath)
	if err != nil {
		fmt.Println("Error getting the info stats.")
		return err
	}

	if !info.IsDir() {
		return encryptFile(dirpath, outputpath, key)
	}

	var encryptNode func(fpath string, dirinfo fs.DirEntry, err error) error = func(fpath string, dirinfo fs.DirEntry, err error) error {
		relfpath, _ := filepath.Rel(dirpath, fpath)
		newfpath := outputpath + "/" + relfpath

		if dirinfo.IsDir() {
			if _, err := os.Stat(newfpath); os.IsNotExist(err) {
				err := os.MkdirAll(newfpath, os.ModePerm)
				if err != nil {
					fmt.Println("Error creating the new dir path.")
					return err
				}
			}
			// Don't try to encrypt a dir path
			// Just return
			return nil
		}

		encryptFile(fpath, filepath.Dir(newfpath), key)
		return nil
	}

	filepath.WalkDir(dirpath, encryptNode)
	return nil
}

func decryptDir(dirpath string, outputpath string, key []byte) error {
	info, err := os.Stat(dirpath)
	if err != nil {
		fmt.Println("Error getting the info stats.")
		return err
	}

	if !info.IsDir() {
		return decryptFile(dirpath, outputpath, key)
	}

	var decryptNode func(fpath string, dirinfo fs.DirEntry, err error) error = func(fpath string, dirinfo fs.DirEntry, err error) error {
		relfpath, _ := filepath.Rel(dirpath, fpath)
		newfpath := outputpath + "/" + relfpath

		if dirinfo.IsDir() {
			if _, err := os.Stat(newfpath); os.IsNotExist(err) {
				err := os.MkdirAll(newfpath, os.ModePerm)
				if err != nil {
					fmt.Println("Error creating the new dir path.")
					return err
				}
			}
			// Don't try to encrypt a dir path
			// Just return
			return nil
		}

		decryptFile(fpath, filepath.Dir(newfpath), key)
		return nil
	}

	filepath.WalkDir(dirpath, decryptNode)
	return nil
}

func loadCipherForDecryption(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm, nil
}

func loadCipherForEncryption(key []byte) (cipher.AEAD, []byte, error) {
	// Encrypt data with AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	return aesgcm, nonce, nil
}

func encryptDataBlock(datablock []byte, aesgcm cipher.AEAD, nonce []byte) []byte {
	ciphertext := aesgcm.Seal(nil, nonce, datablock, nil)
	ciphertext = append(nonce, ciphertext...)
	return ciphertext
}

func decodeDataBlock(datablock []byte, aesgcm cipher.AEAD) ([]byte, error) {
	nonce, ciphertext := datablock[:12], datablock[12:]
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func getFileSize(filepath string) (int64, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func getEncryptedFilepath(file string, outputpath string) string {
	originalFilename := filepath.Base(file)
	return outputpath + "/" + originalFilename + ".cr"
}

func getDecryptedFilepath(file string, outputpath string) string {
	encryptedFilename := filepath.Base(file)
	ext := path.Ext(encryptedFilename)
	return outputpath + "/" + encryptedFilename[:len(encryptedFilename)-len(ext)]
}
