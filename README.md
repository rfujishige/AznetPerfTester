# AznetPerfTester
AznetPerfTester

# how to
## install golang
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update -y
sudo apt install -y golang-go

## install gobgp
```Bash
sudo mkdir /opt/gobgp
sudo cd /opt/gobgp
sudo wget https://github.com/osrg/gobgp/releases/download/v3.25.0/gobgp_3.25.0_linux_amd64.tar.gz
sudo tar xvf gobgp_3.25.0_linux_amd64.tar.gz
sudo  ln -s /opt/gobgp/gobgp /usr/local/bin/gobgp
sudo ln -s /opt/gobgp/gobgpd /usr/local/bin/gobgpd
sudo mkdir -p /etc/gobgp
sudo vi /etc/gobgp/gobgpd.conf

=====
[global.config]
    as = 1
    router-id = "1.1.1.1"
=====

sudo gobgpd -f /etc/gobgp/gobgpd.conf  &
```
## Run AznetPerfTester
```Bash
git clone https://github.com/rfujishige/AznetPerfTester.git
cd AznetPerfTester/
go run main.go &
```

## Access to :8080
http://<IP>:8080
