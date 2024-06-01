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
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation-on-windows-11">Installation on Windows 11</a></li>
        <li><a href="#installation-on-linux-manjaro-with-xfce">Installation on Linux (Manjaro with xfce)</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#author">Author</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>


<!-- ABOUT THE PROJECT -->
## About The Project

I wanted a tool that can easily secure my data stored on cloud storage services like Google Drive or Dropbox.

Because of this I started developing this tool that can easily encrypt or decrypt files and directories on the fly. The first code written was a command line tool that did exactly what I wanted. Unfortunately, this wasn't pleasant to use because of the number of different parameters that I had to dynamically change each time a new file or directory was encrypted/decrypted.

To solve this I created a client-server system where the client sends requests like encrypt, decrypt or generate a new key and the server calls the command line tool described above. The interaction with data is made possible by right clicking on the file or directory and selecting the appropriate action from the context menu. This is easy and simple to use.


### Built With

<p align="left">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-35495E?style=plastic&logo=go&logoColor=69d6e4" target="_blank" />
  </a>
  <a href="https://www.gnu.org/software/make/">
    <img src="https://img.shields.io/badge/Make-0769AD?style=plastic&logo=gnu&logoColor=white" target="_blank" />
  </a>
</p>


<!-- GETTING STARTED -->
## Getting Started

??? This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

??? This is an example of how to list things you need to use the software and how to install them.
Prerequisites: makefile
* npm
  ```sh
  npm install npm@latest -g
  ```

### Installation on Windows 11

1. Clone the repo
   ```sh
   git clone https://github.com/daniel-dumitrascu/filecrypt.git
   ```
2. Open a command prompt by running it as administrator. This is needed because modifications to registry will be made. 
3. Go into [filecrypt repo]\installers\windows
4. Run `install.bat`
   

### Installation on Linux (Manjaro with xfce)

_Below is an example of how you can instruct your audience on installing and setting up your app. This template doesn't rely on any external dependencies or services._
   

<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_


<!-- CONTACT -->
## Author

👤 **Daniel Dumitrascu**

- Linkedin: [@DanielDumitrascu](https://www.linkedin.com/in/daniel-dumitrascu-17a1845a)
- Github: [@daniel-dumitrascu](https://github.com/daniel-dumitrascu)
- Email: daniel.dumitrascu.dev@gmail.com


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.


<!-- MARKDOWN LINKS & IMAGES -->

[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/othneildrew/Best-README-Template/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]: https://github.com/othneildrew/Best-README-Template/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]: https://github.com/othneildrew/Best-README-Template/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]: https://github.com/othneildrew/Best-README-Template/issues
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/othneildrew/Best-README-Template/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
[product-screenshot]: images/screenshot.png
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Vue.js]: https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D
[Vue-url]: https://vuejs.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com
