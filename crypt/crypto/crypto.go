package crypto

type Crypto interface {
	EncryptFile(filepath string, outputpath string, keypath string) error
	DecryptFile(filepath string, outputpath string, keypath string) error

	EncryptDir(dirpath string, outputpath string, keypath string) error
	DecryptDir(dirpath string, outputpath string, keypath string) error

	GenKey() string
}
