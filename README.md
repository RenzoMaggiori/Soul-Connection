<p align="center">
  <img src="./frontend/public/logoOK.svg?raw=true" width="400" alt="logo"/>
</p>

---

<p align="center">
  <img alt="Typescript" src="https://img.shields.io/badge/-TypeScript-black?style=for-the-badge&logoColor=white&logo=typescript&color=2F73BF">
  <img alt="NextJs" src="https://img.shields.io/badge/next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white">
  <img alt="React" src="https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB">
  <img alt="Go" src="https://img.shields.io/badge/-Go-black?style=for-the-badge&logo=go&logoColor=white&color=2F73BF">
  <img alt="Docker" src="https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white">
  <img alt="Postgres" src="https://img.shields.io/badge/postgresql-4169e1?style=for-the-badge&logo=postgresql&logoColor=white">
  <img alt="Mongoose" src="https://img.shields.io/badge/-MongoDB-black?style=for-the-badge&logoColor=white&logo=mongodb&color=127237">
</p>

<!-- [![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/RenzoMaggiori/Soul-Connection/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/aimeos/Soul-Connection/?branch=master) -->

<div align="center">

This project is aimed at migrating an existing [API](https://soul-connection.fr/docs#/) from **Soul Connection** to a new website with a comprehensive dashboard to streamline the daily activities of the coaches and managers of **Soul Connection**.

</div>



# 📖 Getting Started

## 📝 Prerequisites

#### Install Go
For the installation, you can [download](https://go.dev/doc/install) **Go** directly or use the command line:

1. Download and Install Go:

    ``` bash
    curl -LO https://go.dev/dl/go1.22.linux-amd64.tar.gz

    sudo tar -C /usr/local -xzf go1.22.linux-amd64.tar.gz
  
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
  
    go version
    ```

> [!Note]
>
> You can also install it using `sudo apt install golang-go`, but this may not always get the latest version.

#### Install Docker
For the **Docker**, you can download it [here](https://go.dev/doc/install) directly or use the command line:

1. Download and Install Docker:

    ``` bash
    sudo apt update
    sudo apt install -y ca-certificates curl gnupg
    
    sudo install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    docker --version
    ```

#### Install Node.js and npm

1. Install Node.js:

    ``` bash
    sudo apt install -y nodejs
    
    node -v
    ```

2. Install npm (if not automatically installed with Node.js):
   
    ``` bash
    sudo apt install -y npm
    
    npm -v
    ```

## ⚙️ Installation

1. **Clone the Repository**

    ``` bash
    git clone https://github.com/RenzoMaggiori/Soul-Connection.git
    cd Soul-Connection/backend
    ```
2. **From the backend folder run the launch script**

    ``` bash
    ./scripts/launch.sh
    ```
    When the **API** is running you can check it by using `docker compose logs api` and you should see something like this

    <p align="center">
    <img alt="terminal" src="/frontend/public/teminal.png">
    </p>

3. **After the API is running start the migration**

    ``` bash
    docker compose start migration
    ```


# Authors

| [<img src="https://github.com/RenzoMaggiori.png?size=85" width=85><br><sub>Renzo Maggiori</sub>](https://github.com/RenzoMaggiori) | [<img src="https://github.com/oriollinan.png?size=85" width=85><br><sub>Oriol Liñan</sub>](https://github.com/oriollinan) | [<img src="https://github.com/AlbaCande.png?size=85" width=85><br><sub>Alba Candelario</sub>](https://github.com/AlbaCande) | [<img src="https://github.com/G0nzal0zz.png?size=85" width=85><br><sub>Gonzalo Larroya</sub>](https://github.com/G0nzal0zz) 
|:---:|:---:|:---:|:---:|

<br/>

Copyright 2024 Soul Connection. All rights reserved. The license is [here](/LICENSE)
