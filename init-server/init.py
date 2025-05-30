import sys
import subprocess

# Что необходимо сделать перед запуском скрипта:
#  * Обновить пакеты и репозитории
#    apt update && apt upgrade
#  * Установить wireguard
#    apt-get install -y wireguard-tools resolvconf
#  * Сгенерировать ключи
#    wg genkey | tee /etc/wireguard/privkey > /dev/null
#    cat /etc/wireguard/privkey | wg pubkey | tee /etc/wireguard/pubkey
#  * Просмотр ключей:
#    cat /etc/wireguard/privkey - приватный ключ
#    cat /etc/wireguard/pubkey - публичный ключ
#  * Запуск скрипта:
#    sudo python3 init.py <privkey>

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
    privkey = sys.argv[1]
    create_wg0_config(privkey)
    ip_forwarding()
    start_wireguard()

    print(f'\n[*] Инициализация завершена\n --- НЕ ЗАБУДЬТЕ --- \n  1. Написать кофигурацию бота (config/config.yaml), указать IP-адрес и токен бота.\n')

if __name__ == '__main__':
    main()

# sudo python3 init.py