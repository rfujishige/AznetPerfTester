# AznetPerfTester
AznetPerfTester


# install gobgp
```Bash
sudo mkdir /opt/gobgp
sudo wget https://github.com/osrg/gobgp/releases/download/v3.25.0/gobgp_3.25.0_linux_amd64.tar.gz
sudo tar xvf gobgp_3.25.0_linux_amd64.tar.gz
sudo  ln -s /opt/gobgp/gobgp /usr/local/bin/gobgp
sudo ln -s /opt/gobgp/gobgpd /usr/local/bin/gobgpd
sudo vi /etc/gobgp/gobgpd.conf
sudo gobgpd -f /etc/gobgp/gobgpd.conf  &![image](https://github.com/user-attachments/assets/a2af150c-2f43-4c33-b6ec-a4b8c8864828)
```
