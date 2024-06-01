<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/daniel-dumitrascu/filecrypt">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h1 align="center">FileCrypt</h1>

  <p align="center">
    Encrypting and decrypting files in a tree based directory using symmetric cryptography for Linux and Windows!
    <br />
    <br />
    <br />
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#build">Build</a>
    </li>
    <li>
      <a href="#install">Install</a>
      <ul>
        <li><a href="#windows">Windows</a></li>
        <li><a href="#linux">Linux</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#author">Author</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

## About The Project

I wanted a tool that can easily secure my data stored on cloud storage services like Google Drive or Dropbox.

Because of this I started developing this tool that can easily encrypt or decrypt files and directories on the fly. The first code written was a command line tool that did exactly what I wanted. Unfortunately, this wasn't pleasant to use because of the number of different parameters that I had to dynamically change each time a new file or directory was encrypted/decrypted.

To solve this I created a client-server system where the client sends requests like encrypt, decrypt or generate a new key and the server calls the command line tool described above. The interaction with data is made possible by right clicking on the file or directory and selecting the appropriate action from the context menu. This is easy and simple to use.

## Build

<p align="left">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-35495E?style=plastic&logo=go&logoColor=69d6e4" target="_blank" />
  </a>
  <a href="https://www.gnu.org/software/make/">
    <img src="https://img.shields.io/badge/Make-0769AD?style=plastic&logo=gnu&logoColor=white" target="_blank" />
  </a>
</p>

This will build the crypt tool, the client and the server.

1. Install go and make 
2. Clone the repo
   ```sh
   git clone https://github.com/daniel-dumitrascu/filecrypt.git
   ```
3. Open a terminal and run the ```make``` command in the project root.

## Install

This will install, setup and start the app. The repository contains installers for both Windows and Linux located at `project root/installers`

### Windows

Supported versions: Windows 11

1. Open a command prompt by running it as administrator. This is needed because modifications to registry will be made.
2. Go into `project root/installers/windows`
3. Run `install.bat`
4. If everything went ok, the server should be running in a terminal
   
### Linux

Supported versions: Manjaro with xfce

1. Open a terminal.
2. Go into `project root/installers/linux`
3. Run `makepkg -si`
4. If everything went ok, a new service was created representing the server. This is started and enabled by default
  
## Usage

In order to start using the app you need to load a symetric key used for encryption and decryption. 
If you don't have a previous generated key you can generate a new one by right clicking anywhere in a folder or desktop and select `Generate key`.
Once you have the key, load it by right clicking on it and selecting `Add key`.

Now, encrypt files or directories (_be aware that on directories the action is recursive_) by right clicking on the file or directory and selecting `Encrypt source`.
Once the action terminates you will have, in the same directory ar the source file/dir, the encrypted file or directory containg now the encrypted files.  

## Author

👤 **Daniel Dumitrascu**

- Linkedin: [@DanielDumitrascu](https://www.linkedin.com/in/daniel-dumitrascu-17a1845a)
- Github: [@daniel-dumitrascu](https://github.com/daniel-dumitrascu)
- Email: daniel.dumitrascu.dev@gmail.com

## License

Distributed under the MIT License. 
