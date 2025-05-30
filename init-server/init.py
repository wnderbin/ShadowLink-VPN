import os
import subprocess

def install_wireguard():
    print('\t-- Установка wireguard')
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

    privKey = subprocess.run(['wg', 'genkey'], capture_output=True, text=True).stdout.split()
    privkey = privKey[0]
    subprocess.run(['echo', f'{privkey}', '>', '/etc/wireguard/privatekey'])

    pubKey = subprocess.run(['echo', f'{privkey}', '|','wg', 'pubkey'], capture_output=True, text=True).stdout.split()
    pubkey = pubKey[0].replace('[', '').replace(']', '').replace('\'', '')
    subprocess.run(['echo', f'{pubkey}', '>', '/etc/wireguard/publickey'])
    
    print('-- Генерация ключей завершена')
    return privkey, pubkey

def create_wg0_config(privKey):
    
    print('-- Создание конфигурации wg0...')
    config = f"""[Interface]
        PrivateKey = {privKey}
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
    privkey, pubkey = generate_keys()
    create_wg0_config(privkey)
    ip_forwarding()
    start_wireguard()

    print(f'\n[*] Инициализация завершена\nПриватный ключ: {privkey}\nПубличный ключ: {pubkey}\n --- НЕ ЗАБУДЬТЕ --- \n  1. Написать кофигурацию бота (config/config.yaml), указать IP-адрес и токен бота.\n')

if __name__ == '__main__':
    main()

# sudo python3 init.py