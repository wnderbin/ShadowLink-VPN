<a href="#instructions-for-launching-a-project-on-hosting">🇬🇧 English version</a> \
<a href="#инструкция-по-запуску-проекта-на-хостинге">🇷🇺 Russian version</a>

## Instructions for launching a project on hosting

Before launching the bot, it is important to configure it correctly and install dependencies, this will be discussed in this instruction.

**Configuration: Ubuntu server 22.04**

First, I recommend updating the repositories and packages with the command:
```
apt update && apt upgrade
```
Then you need to install wireguard, it is necessary for the client to communicate with the server, as well as for generating the necessary configs:
```
apt-get install -y wireguard-tools resolvconf
```
Also, together with wireguard, the resolvconf utility was installed, ensuring consistent operation with DNS servers.

Now we need to generate public and private keys for further work with wireguard, we will do this with the command:
```
wg genkey | tee /etc/wireguard/privkey > /dev/null
cat /etc/wireguard/privkey | wg pubkey | tee /etc/wireguard/pubkey
```
Now the keys are stored in the /etc/wireguard directory. You can view their contents using the commands:
```
cat /etc/wireguard/privkey # private key
cat /etc/wireguard/pubkey  # public key
```
Now you need to run the script that initializes the server configuration wireguard wg0.conf. When running, you need to specify your private key from the directory /etc/wireguard/
```
sudo python3 wg_init.py <YOUR-PRIVATE-KEY>
```
It is recommended to check the server settings before starting and specify the required parameters. The server settings that should be by default are already specified in the script, but if an error occurred somewhere, then most likely the configuration was generated incorrectly.

Next you need to install Golang: [Download and install](https://go.dev/doc/install)

After installation, you need to create a project configuration config.yaml in the config directory. It should look like this:
```
BOT-API-KEY: "<TELEGRAM-API-TOKEN>"
DEBUG-MODE: true/false
WG-Interface: "<WIREGUARD INTERFACE | default: wg0>"
WG-ConfigPath: "<WIREGUARD-CONFIG-PATH | default: /etc/wireguard>"
Server-PublicIP: "<SERVER-IP>"
ServerPort: <WIREGUARD PORT | default: 51820>
ServerPublicKey: "<PUBLIC-KEY FROM /etc/wireguard/pubkey>" 
AllowedIPs: "0.0.0.0/0, ::/0" 
DNS: "<ANY-DNS-SERVER | 1.1.1.1 or 8.8.8.8>"
DB: "postgres://<user>:<pass>@<host>:<port>/<dbname>?sslmode=<disable/enable>"
Redis-host: "<host>:<port>"
Redis-pass: "<redis-password from redis.conf>"
Redis-DB: <redis-db>
```

Now you need to perform ip-forwarding for stable VPN operation.
```
sudo sysctl net.ipv4.ip_forward
sudo sysctl -w net.ipv4.ip_forward=1
net.ipv4.ip_forward = 1
```

Now you can start your VPN with the command:
```
CONFIG_PATH=./config/config/yaml go run cmd/main.go
```

You may need to install project dependencies. If this was not done using the go run command, type:
```
go mod download
```


## Инструкция по запуску проекта на хостинге
Перед запуском бота важно правильно его настроить и установить зависимости, об этом пойдет речь в этой инструкции.

**Конфигурация: Ubuntu server 22.04**

Для начала рекомендую обновить репозитории и пакеты командой:
```
apt update && apt upgrade
```
Затем необходимо установить Wireguard, он необходим для связи клиента с сервером, а также для генерации необходимых конфигов:
```
apt-get install -y wireguard-tools resolvconf
```
Также вместе с wireguard была установлена ​​утилита resolvconf, обеспечивающая согласованную работу с DNS-серверами.

Теперь нам необходимо сгенерировать публичный и приватный ключи для дальнейшей работы с wireguard, сделаем это командой:
```
wg genkey | tee /etc/wireguard/privkey > /dev/null
cat /etc/wireguard/privkey | wg pubkey | tee /etc/wireguard/pubkey
```

Теперь ключи хранятся в каталоге /etc/wireguard. Просмотреть их содержимое можно с помощью команд:
```
cat /etc/wireguard/privkey # приватный ключ
cat /etc/wireguard/pubkey  # публичный ключ
```

Теперь нужно запустить скрипт, который инициализирует конфигурацию сервера wireguard wg0.conf. При запуске нужно указать свой закрытый ключ из каталога /etc/wireguard/

```
sudo python3 wg_init.py <YOUR-PRIVATE-KEY>
```

Рекомендуется проверить настройки сервера перед запуском и указать требуемые параметры. Настройки сервера, которые должны быть по умолчанию, уже указаны в скрипте, но если где-то произошла ошибка, то скорее всего конфигурация была сгенерирована некорректно.

Далее необходимо установить Golang: [Download and install](https://go.dev/doc/install)

После установки необходимо создать в каталоге config конфигурацию проекта config.yaml. Конфигурация должна выглядеть так:

```
BOT-API-KEY: "<ТОКЕН-ТЕЛЕГРАМ-БОТА>"
DEBUG-MODE: true/false
WG-Interface: "<ИНТЕРФЕЙС WIREGUARD | по умолчанию: wg0>"
WG-ConfigPath: "<ПУТЬ К КОНФИГУРАЦИИ WIREGUARD | по умолчанию: /etc/wireguard>"
Server-PublicIP: "<АЙПИ-СЕРВЕРА>"
ServerPort: <ПОРТ WIREGUARD | по умолчанию: 51820>
ServerPublicKey: "<ПУБЛИЧНЫЙ КЛЮЧ ИЗ ФАЙЛА /etc/wireguard/pubkey>" 
AllowedIPs: "0.0.0.0/0, ::/0" 
DNS: "<ЛЮБОЙ DNS-СЕРВЕР | 1.1.1.1 or 8.8.8.8>"
DB: "postgres://<user>:<pass>@<host>:<port>/<dbname>?sslmode=<disable/enable>"
Redis-host: "<host>:<port>"
Redis-pass: "<redis-password from redis.conf>"
Redis-DB: <redis-db>
```

Теперь необходимо выполнить ip-переадресацию для стабильной работы VPN.

```
sudo sysctl net.ipv4.ip_forward
sudo sysctl -w net.ipv4.ip_forward=1
net.ipv4.ip_forward = 1
```

Теперь вы можете запустить свой VPN с помощью команды:

```
CONFIG_PATH=./config/config/yaml go run cmd/main.go
```

Может понадобится установка зависимостей проекта. Если это не было сделано при использовании команды go run, пропишите:
```
go mod download
```