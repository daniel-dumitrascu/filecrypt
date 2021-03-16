from cryptography.fernet import Fernet
from enum import Enum
from datetime import datetime
import time
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

def encrypt_file(key, file_to_encrypt_path, output_path):
    chunck_size = 16384
    file_size = None
    fernet_instance = Fernet(key)

    try:
        file_size = os.path.getsize(file_to_encrypt_path)
    except:
        print("Error in getting the file size")
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
        print("Error in getting the file size")
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
            bytes_to_read = 1 + 8 + 16 + 16384 + 16 + 32
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


#generate_and_save_key('C:/Home/crypt')
key = load_key('C:/Home/crypt/16-03-2021_18-09-01_private.key')

start_time = time.time()

# 512 bytes size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_512.txt', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_512.txt.crypt', 'C:/Home/crypt/decrypted_files')

# 511 bytes size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_511.txt', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_511.txt.crypt', 'C:/Home/crypt/decrypted_files')

# 513 bytes size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_513.txt', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_513.txt.crypt', 'C:/Home/crypt/decrypted_files')

# 1024 bytes size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_1024.txt', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_1024.txt.crypt', 'C:/Home/crypt/decrypted_files')

# 1026 bytes size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_1026.txt', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_1026.txt.crypt', 'C:/Home/crypt/decrypted_files')

# 5 MB size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_5MB.pdf', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_5MB.pdf.crypt', 'C:/Home/crypt/decrypted_files')

# 50 MB size file
encrypt_file(key, 'C:/Home/crypt/original_test_files/file_50MB.pdf', 'C:/Home/crypt/encrypted_test_files')
decrypt_file(key, 'C:/Home/crypt/encrypted_test_files/file_50MB.pdf.crypt', 'C:/Home/crypt/decrypted_files')

print('Execution time: ' + str(time.time() - start_time) + " seconds")