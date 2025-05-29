import os
import subprocess

def install_wireguard():
    print('Установка wireguard')
    subprocess.run(['apt-get', 'install', '-y', 'wireguard-tools', 'resolvconf'])

def install_dependecies():
    print('-- Установка зависимостей...')

    print('\t-- Обновление списка доступных репозиториев')
    subprocess.run(['apt', 'update', '-y'])

    print('\t-- Обновление установленных пакетов до последних стабильных версий')
    subprocess.run(['apt', 'upgrade', '-y'])

    install_wireguard()

    print('-- Установка зависимостей завершена')

def generate_keys():
    print('-- Генерация ключей wireguard...')
    print('\t-- Создание директорий')
    os.makedirs('/etc/wireguard', exist_ok=True)

    print('\t--Генерация ключей')
    subprocess.run('wg', 'genkey', '|', 'tee', '/etc/wireguard/privatekey', '|', 'wg', 'pubkey', '|', 'tee', '/etc/wireguard/publickey')
    
    print('-- Генерация ключей завершена')

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
    pub, priv = generate_keys()
    create_wg0_config()
    ip_forwarding()
    start_wireguard()

    print(f'\n[*] Инициализация завершена\n --- НЕ ЗАБУДЬТЕ --- \n  1. Написать кофигурацию бота (config/config.yaml), указать IP-адрес и токен бота.\n')

if __name__ == '__main__':
    main()

# chmod +x init.py
# sudo ./init.py