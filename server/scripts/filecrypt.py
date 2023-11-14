from cryptography.fernet import Fernet
from enum import Enum
from datetime import datetime
import time
import sys
import os
import math
import base64
import argparse

def generate_key():
    return Fernet.generate_key()

def save_file(save_path: str, save_data):
    try:
        with open(save_path, 'wb') as file_handler:
            file_handler.write(save_data)
            file_handler.close()
    except Exception as e:
        print("Problem when saving the data on the file: " + str(e))

def generate_and_save_key(save_path: str):
    key = generate_key()
  
    # save key
    date_time_str = datetime.now().strftime("%d-%m-%Y_%H-%M-%S")
    private_key_path = save_path + "/" + date_time_str + "_private.key"
    save_file(private_key_path, key)

def load_key(path: str):
    key = None
    try:
        with open(path, "rb") as key_file:
            key = key_file.read()
            key_file.close()
    except Exception as e:
        print("Problem when loading the key: " + str(e))
    
    return key

def encrypt_file(key, file_to_encrypt_path, output_path):
    chunck_size = 1048576 # 1MB
    file_size = None
    fernet_instance = Fernet(key)

    try:
        file_size = os.path.getsize(file_to_encrypt_path)
    except:
        print("Error in getting the file size of %s" % file_to_encrypt_path)
        return

    original_filename = os.path.basename(file_to_encrypt_path)
    secret_filename = constr_secret_filename(original_filename)
    secret_filename_path = output_path + "/" + secret_filename

    encrypted_file_handler = open(secret_filename_path, 'wb')

    try:
        with open(file_to_encrypt_path, 'rb') as file_handler:
            print("Encrypting file " + file_to_encrypt_path)
            chunck_nr = math.ceil(file_size / chunck_size)
            chunk_data = file_handler.read(chunck_size)
            chunk_index = 1
            while chunk_data:
                print("Encrypting data chunk " + str(chunk_index) + " of " + str(chunck_nr))
                encrypted_chunck_data = fernet_instance.encrypt(chunk_data)
                chunk_index += 1
                b64decoded_data = base64.urlsafe_b64decode(encrypted_chunck_data)        
                encrypted_file_handler.write(b64decoded_data)
                chunk_data = file_handler.read(chunck_size)
    except Exception as e:
        print("Problem when encrypting the file " + file_to_encrypt_path + ": "+ str(e))        
    
    encrypted_file_handler.close()

def decrypt_file(key, file_to_decrypt_path, output_path):
    file_size = None
    fernet_instance = Fernet(key)

    try:
        file_size = os.path.getsize(file_to_decrypt_path)
    except:
        print("Error in getting the file size of %s" % file_to_decrypt_path)
        return

    encoded_filename = os.path.basename(file_to_decrypt_path)
    original_filename = constr_original_filename(encoded_filename)
    original_filename_path = output_path + "/" + original_filename

    original_file_handler = open(original_filename_path, 'wb')

    try:
        with open(file_to_decrypt_path, 'rb') as secret_file_handler:
            print("Decrypting file " + file_to_decrypt_path)
            # We read and decrypt each token at a time
            # To know how much we need to read we need to know the token structure and len
            # Version(1 byte) + Date created(8 bytes) + IV(16 bytes) + Cipher(chunck_size + 16) + HMAC(32 bytes)  
            bytes_to_read = 1 + 8 + 16 + 1048576 + 16 + 32
            token_nr = math.ceil(file_size / bytes_to_read)
            token = secret_file_handler.read(bytes_to_read)
            token_index = 1
            while token:
                print("Decrypting token " + str(token_index) + " of " + str(token_nr))
                b64encoded_token = base64.urlsafe_b64encode(token)
                decoded_data = fernet_instance.decrypt(b64encoded_token)
                token_index += 1
                original_file_handler.write(decoded_data)
                token = secret_file_handler.read(bytes_to_read)   
    except Exception as e:
        print("Problem when decrypting the file " + file_to_decrypt_path + ": "+ str(e)) 

    original_file_handler.close()

def constr_secret_filename(original_filename: str):
    return original_filename + ".crypt"

def constr_original_filename(encoded_filename: str):
    return os.path.splitext(encoded_filename)[0]

def encrypt_dir_tree(key_path: str, input_path: str, output_path: str):
    if os.path.isfile(key_path) == False:
        print('The key path does not point to a valid file')
        return

    key = load_key(key_path)
    if key == None:
        return

    if os.path.isfile(input_path):
        encrypt_file(key, input_path, output_path)
        return

    for dirpath, dirs, files in os.walk(input_path):
        relative_path = os.path.relpath(dirpath, input_path)
        new_path = output_path + "/" + relative_path
           
        if os.path.isdir(new_path) == False:
            print('Creating directory %s' % new_path)
            try:
                os.mkdir(new_path)
            except OSError as e:
                print ("Creation of the directory " + new_path + " failed: " + e)
                print ("Skipping encryption of files from directory %s" % new_path)
                continue

        for file in files:
            file_path = dirpath + "/" + file
            encrypt_file(key, file_path, new_path)

def decrypt_dir_tree(key_path: str, input_path: str, output_path: str):
    if os.path.isfile(key_path) == False:
        print('The key path does not point to a valid file')
        return

    key = load_key(key_path)
    if key == None:
        return

    if os.path.isfile(input_path):
        decrypt_file(key, input_path, output_path)
        return

    for dirpath, dirs, files in os.walk(input_path):
        relative_path = os.path.relpath(dirpath, input_path)
        new_path = output_path + "/" + relative_path
           
        if os.path.isdir(new_path) == False:
            print('Creating directory %s' % new_path)
            try:
                os.mkdir(new_path)
            except OSError as e:
                print ("Creation of the directory " + new_path + " failed: " + e)
                print ("Skipping decryption of files from directory %s" % new_path)
                continue

        for file in files:
            file_path = dirpath + "/" + file
            decrypt_file(key, file_path, new_path)


parser = argparse.ArgumentParser()
subparsers = parser.add_subparsers(help='sub-command help', dest='command')

# create the parser for the "genkey" command
parser_genkey = subparsers.add_parser('genkey', help='Generates a private key at the specified location')
parser_genkey.add_argument("path", type=str, help='Path where the key will be stored')

# create the parser for the "encrypt" command
parser_encrypt = subparsers.add_parser('encrypt', help='Encrypts a single file, all the files from a tree based directory or the files from a root directory')
parser_encrypt.add_argument("key_path", type=str, help='Path to the private key')
parser_encrypt.add_argument("input_path", type=str, help='Path that points to the input root directory')
parser_encrypt.add_argument("output_path", type=str, help='Path that points to the output root directory')

# create the parser for the "decrypt" command
parser_decrypt = subparsers.add_parser('decrypt', help='Decrypts a single file, all the files from a tree based directory or the files from a root directory')
parser_decrypt.add_argument("key_path", type=str, help='Path to the private key')
parser_decrypt.add_argument("input_path", type=str, help='Path that points to the input root directory')
parser_decrypt.add_argument("output_path", type=str, help='Path that points to the output root directory')

args = parser.parse_args(sys.argv[1:])
args_dict = vars(args)

if args_dict['command'] == 'genkey':
    generate_and_save_key(args_dict['path'])
elif args_dict['command'] == 'encrypt':
    encrypt_dir_tree(args_dict['key_path'], args_dict['input_path'], args_dict['output_path'])
elif args_dict['command'] == 'decrypt':
    decrypt_dir_tree(args_dict['key_path'], args_dict['input_path'], args_dict['output_path'])