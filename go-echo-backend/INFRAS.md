# Elite Medical Staff

## Docs
- Ref https://selfhost.club/guides/fileserver/

## Install Wireguard
```
curl -O https://raw.githubusercontent.com/angristan/wireguard-install/master/wireguard-install.sh
chmod +x wireguard-install.sh
./wireguard-install.sh
```
## Install filebrowser https://filebrowser.org/installation
```
curl -fsSL https://raw.githubusercontent.com/filebrowser/get/master/get.sh | bash
```
## Setup Caddy
- Install caddy https://caddyserver.com/docs/install#debian-ubuntu-raspbian
    ```
    sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
    sudo apt update
    sudo apt install caddy
    ```

- sudo nano /etc/caddy/Caddyfile

- Add this code
    ```
    {your_domain}
    reverse_proxy localhost:8080
    ```
- Reload, enable and start Caddy
    ```
    sudo systemctl reload caddy
    sudo systemctl enable caddy
    sudo systemctl start caddy
    sudo systemctl status caddy
    ```

## Add filebrowser to systemd
- sudo nano /etc/systemd/system/filebrowser.service
- Add this code 
    ```
[Unit]
Description=Run Caddy Filebrowser at startup

[Service]
# Replace with your actual username
User=ubuntu

# Substitute the paths to your storage folder and Filebrowser database file
# Make sure to use full paths instead of ~/.filebrowser_database/filebrowser.db etc
ExecStart=/usr/local/bin/filebrowser -r /home/ubuntu -d /home/ubuntu/filebrowser.db
Type=simple

[Install]
WantedBy=multi-user.target
    ```
- Enable service
    ```
    sudo systemctl enable filebrowser
    sudo systemctl start filebrowser
    sudo systemctl status filebrowser
    ```
## Install EFS
```
sudo apt-get install nfs-common
```

## Grant user permission
sudo chown -R ubuntu:ubuntu /home/ubuntu/efs

## Install aws cli
sudo apt-get update -y
sudo apt-get install -y awscli