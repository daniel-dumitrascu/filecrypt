# FileCrypt

Encrypting and decrypting files in a tree based directory using symmetric cryptography.

One important feature of the script is that it's able to encrypt/decrypt big files. This is possible by splitted the file into chunks of 1Mb which are then encrypted and concatenated to the resulting encrypted file. 

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

First, update pip

```
python -m pip install --upgrade pip
```

Add the cryptography module

```
pip install cryptography
```
## Running the script

The script has 3 main actions:
- key generation
- encryption
- decryption

### Generating the key

```
usage: filecrypt.py genkey [-h] path
```
path  Path where the key will be stored

### Encrypt 

```
usage: filecrypt.py encrypt [-h] key_path input_path output_path key_path
```
key_path     Path to the private key
input_path   Path that points to the input root directory
output_path  Path that points to the output root directory

### Decrypt

```
usage: filecrypt.py decrypt [-h] key_path input_path output_path
```

key_path     Path to the private key
input_path   Path that points to the input root directory
output_path  Path that points to the output root directory