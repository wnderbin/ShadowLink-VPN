import os
import subprocess
from configparser import ConfigParser

def install_dependecies():
    print('-- Установка зависимостей...')

    print('\t-- Обновление списка доступных репозиториев')
    subprocess.run(['apt', 'update', '-y'])

    print('\t-- Обновление установленных пакетов до последних стабильных версий')
    subprocess.run(['apt', 'upgrade', '-y'])

    print('\t-- Установка wireguard')
    subprocess.run(['apt', 'install', '-y', 'wireguard'])

    print('\t-- Установка Golang 1.24.2')
    subprocess.run(['wget', 'https://go.dev/dl/go1.24.2.linux-amd64.tar.gz'])
    subprocess.run(['sudo', 'tar', '-C', '/usr/local', '-xzf', 'go1.24.2.linux-amd64.tar.gz'])
    subprocess.run(['echo', '\'export PATH=$PATH:/usr/local/go/bin\'', '>>', "~/.bashrc"])
    subprocess.run(['source', '~/.bashrc'])

    print('-- Установка зависимостей завершена')

def clone_repository():
    print('-- Клонирование репозитория...')
    subprocess.run(['git', 'clone', 'https://github.com/wnderbin/ShadowLink-VPN'])
    print('-- Клонирование репозитория завершено')

def generate_keys():
    print('-- Генерация ключей wireguard...')
    print('\t-- Создание директорий')
    os.makedirs('/etc/wireguard', exist_ok=True)

    print('\t-- Генерация приватного ключа')
    privateKey = subprocess.run(['wg', 'genkey'], capture_output=True, text=True).stdout.strip()
    with open('/etc/wireguard/privatekey', 'w') as k:
        k.write(privateKey)

    print('\t-- Генерация публичного ключа')
    publicKey = subprocess.run(['wg', 'pubkey'], capture_output=True, text=True).stdout.strip()
    with open('/etc/wireguard/publickey', 'w') as k:
        k.write(publicKey)
    
    print('-- Генерация ключей завершена')
    return privateKey, publicKey

def create_wg0_config(privkey):
    print('-- Создание конфигурации wg0...')
    config = f"""[Interface]
        PrivateKey = {privkey}
        Address = 10.8.0.1/24
        ListenPort = 51820
        PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
        PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE"""
    
    with open("/etc/wireguard/wg0.conf", "w") as c:
        c.write(config)

    print('-- Создание конфигурации завершено')

def ip_forwarding():
    print('-- Настройка сетевый параметров...')
    subprocess.run(['echo', '\"net.ipv4.ip_forward=1\"', '>>', '/etc/sysctl.conf'])
    print('-- Настройка сетевых параметров завершена')

def start_wireguard():
    print('-- Запуск wireguard...')
    subprocess.run(['systemctl', 'enable', 'wg-quick@wg0.service'])
    subprocess.run(['systemctl', 'start', 'wg-quick@wg0.service'])
    print('-- Запуск завершен')

def main():
    install_dependecies()
    clone_repository()
    pub, priv = generate_keys()
    create_wg0_config()
    ip_forwarding()
    start_wireguard()

    print(f'\n[*] Инициализация завершена\n  Публичный ключ: {pub}\n  Приватный ключ: {priv}\n --- НЕ ЗАБУДЬТЕ --- \n  1. Написать кофигурацию бота (config/config.yaml), указать IP-адрес и токен бота.\n')

if __name__ == '__main__':
    main()

# chmod +x init.py
# sudo ./init.py