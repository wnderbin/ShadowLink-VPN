import sys
import subprocess

# !!! Перед запуском скрипта рекомендуется прочитать инструкцию, которая находится в этой директории.
# !!! Before running the script, it is recommended to read the instructions located in this directory.

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

def start_wireguard():
    print('-- Запуск wireguard...')
    subprocess.run(['systemctl', 'enable', 'wg-quick@wg0.service'])
    subprocess.run(['systemctl', 'start', 'wg-quick@wg0.service'])
    print('-- Запуск завершен')

def main():
    privkey = sys.argv[1]
    create_wg0_config(privkey)
    start_wireguard()

    print(f'\n[*] Инициализация завершена\n --- НЕ ЗАБУДЬТЕ --- \n  1. Написать кофигурацию бота (config/config.yaml), указать IP-адрес и токен бота.\n2. Пробросить айпи-форвардинг.\n\tsudo sysctl net.ipv4.ip_forward\n\tЕсли 0, включи: sudo sysctl -w net.ipv4.ip_forward=1\n\tИ добавь: net.ipv4.ip_forward = 1')

if __name__ == '__main__':
    main()