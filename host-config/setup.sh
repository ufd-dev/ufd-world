apt update && apt upgrade
hostnamectl set-hostname ufd.world
echo '172.105.47.181 ufd.world' >> /etc/hosts
useradd -m ufd-world
useradd -m admin
chsh -s /bin/bash admin
adduser admin sudo

# setup as admin
su - admin
mkdir .ssh
echo '<omitted>' .ssh/authorized_keys
chmod -R 700 /home/admin/.ssh && chmod 600 /home/admin/.ssh/authorized_keys
exit

# set vim
sudo update-alternatives --config editor

# lock down SSH
echo 'PasswordAuthentication no' > /etc/ssh/sshd_config.d/no-passwords.conf
sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# test ability to log in as admin without a password, THEN
shutdown -r 0

# gen temp key in private location
openssl ecparam -genkey -name secp384r1 -out server.key
# gen temp cert
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

# from https://certbot.eff.org/instructions?ws=other&os=snap
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot certonly --standalone
# wrote to
# Certificate is saved at: /etc/letsencrypt/live/ufd.world/fullchain.pem
# Key is saved at:         /etc/letsencrypt/live/ufd.world/privkey.pem

# after build/push
# pkill ufd-world && cd /home/ufd-world/pub/ && cp /home/admin/ufd-world . && ./ufd-world &