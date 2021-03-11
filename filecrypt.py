from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.asymmetric import rsa

# generate the private and public keys
private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048)

public_key = private_key.public_key()

# serialize and store both keys
private_key_data = private_key.private_bytes(
    encoding=serialization.Encoding.PEM,
    format=serialization.PrivateFormat.TraditionalOpenSSL,
    encryption_algorithm=serialization.NoEncryption())

public_key_data = public_key.public_bytes(
    encoding=serialization.Encoding.PEM,
    format=serialization.PublicFormat.SubjectPublicKeyInfo
)

with open('C:/Home/crypt/key.private', 'wb') as file_handler:
    file_handler.write(private_key_data)
    file_handler.close()

with open('C:/Home/crypt/key.public', 'wb') as file_handler:
    file_handler.write(public_key_data)
    file_handler.close()

# read private and public keys
with open("C:/Home/crypt/key.private", "rb") as key_file:
	private_key = serialization.load_pem_private_key(
        key_file.read(),
        password=None
    )
	key_file.close()

with open("C:/Home/crypt/key.public", "rb") as key_file:
	public_key = serialization.load_pem_public_key(
        key_file.read()
    )
	key_file.close()

# load file to encrypt
file_data = None
with open('C:/Home/crypt/README.md', 'rb') as file_handler:
	file_data = file_handler.read()

# encode file
encrypted_data = public_key.encrypt(
    file_data,
    padding.OAEP(
        mgf=padding.MGF1(algorithm=hashes.SHA256()),
        algorithm=hashes.SHA1(),
        label=None
    )
)

# save the encrypted file
with open('C:/Home/crypt/README.secret', 'wb') as file_handler:
	file_handler.write(encrypted_data)

# decrypt the file
decrypted_file_data = private_key.decrypt(
    encrypted_data,
    padding.OAEP(
        mgf=padding.MGF1(algorithm=hashes.SHA256()),
        algorithm=hashes.SHA1(),
        label=None
    )
)

# save the decripted file
with open('C:/Home/crypt/README.original', 'wb') as file_handler:
	file_handler.write(decrypted_file_data)
