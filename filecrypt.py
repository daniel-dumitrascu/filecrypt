from cryptography.fernet import Fernet
from enum import Enum
from datetime import datetime
import os
import math
import base64

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

def encrypt_file(key, file_path):
    chunck_size = 512
    file_size = None
    fernet_instance = Fernet(key)

    try:
        file_size = os.path.getsize(file_path)
    except:
        print("Error in getting the file size")
        return

    original_filename = os.path.basename(file_path)
    secret_filename = constr_secret_filename(original_filename)
    dirname_path = os.path.dirname(file_path)
    secret_filename_path = dirname_path + "/" + secret_filename

    encrypted_file_handler = open(secret_filename_path, 'wb')

    try:
        with open(file_path, 'rb') as file_handler:
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
        print("Problem when encrypting the file " + file_path + ": "+ str(e))        
    
    encrypted_file_handler.close()

def decrypt_file(key, file_path):
    fernet_instance = Fernet(key)

    encoded_filename = os.path.basename(file_path)
    original_filename = constr_original_filename(encoded_filename)
    dirname_path = os.path.dirname(file_path)
    original_filename_path = dirname_path + "/" + original_filename

    original_file_handler = open(original_filename_path, 'wb')

    try:
        with open(file_path, 'rb') as secret_file_handler:
            # We read and decrypt each token at a time
            # To know how much we need to read we need to know the token structure and len
            # Version(1 byte) + Date created(8 bytes) + IV(16 bytes) + Cipher(chunck_size + 16) + HMAC(32 bytes) --> 585 bytes.
            bytes_to_read = 585
            token = secret_file_handler.read(bytes_to_read)
            while token:
                b64encoded_token = base64.urlsafe_b64encode(token)
                decoded_data = fernet_instance.decrypt(b64encoded_token)
                original_file_handler.write(decoded_data)
                token = secret_file_handler.read(bytes_to_read)   
    except Exception as e:
        print("Problem when decrypting the file " + file_path + ": "+ str(e)) 

    original_file_handler.close()

def constr_secret_filename(original_filename: str):
    return original_filename + ".crypt"

def constr_original_filename(encoded_filename: str):
    return os.path.splitext(encoded_filename)[0]


#generate_and_save_key('C:/Home/crypt')
key = load_key('C:/Home/crypt/16-03-2021_18-09-01_private.key')
#encrypt_file(key, 'C:/Home/crypt/test2.png')
decrypt_file(key, 'C:/Home/crypt/test2.png.crypt')