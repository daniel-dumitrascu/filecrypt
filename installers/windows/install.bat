@echo off
SET APP_NAME=filecrypt
SET GIT_REPO_NAME=filecrypt
SET TEMP_PATH=C:\Users\%USERNAME%\AppData\Local\Temp
SET INSTALL_PATH=%TEMP_PATH%\%APP_NAME%
SET APP_PATH=C:\Program Files\%APP_NAME%
SET BIN_PATH=%APP_PATH%\bin
SET RESOURCE_PATH=%APP_PATH%\resources
SET KEYS_PATH=C:\Users\%USERNAME%\%APP_NAME%\keys
SET "SHORTCUT_STARTUP=C:\Users\%USERNAME%\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\%APP_NAME%_server"
SET "SHORTCUT_DESKTOP=C:\Users\%USERNAME%\Desktop\%APP_NAME%_server"

if not exist %INSTALL_PATH%\NUL mkdir %INSTALL_PATH%

if exist %INSTALL_PATH%\%GIT_REPO_NAME%\NUL (
 	RMDIR /S /Q %INSTALL_PATH%\%GIT_REPO_NAME%
)

cd %INSTALL_PATH%
git clone git@github.com:daniel-dumitrascu/filecrypt.git
cd %GIT_REPO_NAME%
make

if not exist "%APP_PATH%" (
	mkdir "%APP_PATH%"
)

if exist "%BIN_PATH%\NUL" (
 	RMDIR /S /Q "%BIN_PATH%"
)
mkdir "%BIN_PATH%"

if exist "%RESOURCE_PATH%\NUL" (
 	RMDIR /S /Q "%RESOURCE_PATH%"
)
mkdir "%RESOURCE_PATH%"

copy /Y %INSTALL_PATH%\%GIT_REPO_NAME%\client\client.exe "%BIN_PATH%"\filecrypt_client.exe
copy /Y %INSTALL_PATH%\%GIT_REPO_NAME%\server\server.exe "%BIN_PATH%"\filecrypt_server.exe
copy /Y %INSTALL_PATH%\%GIT_REPO_NAME%\client\resources\encrypt.ico "%RESOURCE_PATH%"\encrypt.ico
copy /Y %INSTALL_PATH%\%GIT_REPO_NAME%\client\resources\decrypt.ico "%RESOURCE_PATH%"\decrypt.ico
copy /Y %INSTALL_PATH%\%GIT_REPO_NAME%\client\resources\key.ico "%RESOURCE_PATH%"\key.ico

if not exist %KEYS_PATH% (
	mkdir %KEYS_PATH%
)

cd "%BIN_PATH%"
RMDIR /S /Q %INSTALL_PATH%

mklink "%SHORTCUT_STARTUP%" "%BIN_PATH%\filecrypt_server.exe"
mklink "%SHORTCUT_DESKTOP%" "%BIN_PATH%\filecrypt_server.exe"

filecrypt_server.exe