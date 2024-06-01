@echo off
SET APP_NAME=filecrypt
SET APP_PATH=C:\Program Files\%APP_NAME%
SET KEYS_PATH=C:\Users\%USERNAME%\%APP_NAME%\keys
SET "SHORTCUT_STARTUP=C:\Users\%USERNAME%\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\%APP_NAME%_server"
SET "SHORTCUT_DESKTOP=C:\Users\%USERNAME%\Desktop\%APP_NAME%_server"
set "REG_PATH_FILE=HKEY_CLASSES_ROOT\*\shell"
set "REG_PATH_FOLDER=HKEY_CLASSES_ROOT\Folder\shell"
set "REG_PATH_GEN_KEY=HKEY_CLASSES_ROOT\Directory\Background\shell"

del "%SHORTCUT_STARTUP%"
del "%SHORTCUT_DESKTOP%"
RMDIR /S /Q "%APP_PATH%"
RMDIR /S /Q "%KEYS_PATH%"

reg query "%REG_PATH_FILE%\FilecryptAddKey" >nul 2>&1
if %errorlevel% equ 0 (
	echo Remove registry entry: "%REG_PATH_FILE%\FilecryptAddKey"
	reg delete "%REG_PATH_FILE%\FilecryptAddKey" /f
)

reg query "%REG_PATH_FILE%\FilecryptDecrypt" >nul 2>&1
if %errorlevel% equ 0 (
	echo Remove registry entry: "%REG_PATH_FILE%\FilecryptDecrypt"
	reg delete "%REG_PATH_FILE%\FilecryptDecrypt" /f
)

reg query "%REG_PATH_FILE%\FilecryptEncrypt" >nul 2>&1
if %errorlevel% equ 0 (
	echo Remove registry entry: "%REG_PATH_FILE%\FilecryptEncrypt"
	reg delete "%REG_PATH_FILE%\FilecryptEncrypt" /f
)

reg query "%REG_PATH_FOLDER%\FilecryptDecrypt" >nul 2>&1
if %errorlevel% equ 0 (
	echo Remove registry entry: "%REG_PATH_FOLDER%\FilecryptDecrypt"
	reg delete "%REG_PATH_FOLDER%\FilecryptDecrypt" /f
)

reg query "%REG_PATH_FOLDER%\FilecryptEncrypt" >nul 2>&1
if %errorlevel% equ 0 (
	echo Remove registry entry: "%REG_PATH_FOLDER%\FilecryptEncrypt"
	reg delete "%REG_PATH_FOLDER%\FilecryptEncrypt" /f
)

reg query "%REG_PATH_GEN_KEY%\FilecryptGenKey" >nul 2>&1
if %errorlevel% equ 0 (
	echo Remove registry entry: "%REG_PATH_GEN_KEY%\FilecryptGenKey"
	reg delete "%REG_PATH_GEN_KEY%\FilecryptGenKey" /f
)