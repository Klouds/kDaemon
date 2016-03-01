#!/bin/sh
sudo apt update
sudo apt upgrade -y
sudo apt install -y moreutils
sudo wget get.docker.io
sudo sh index.html
sudo curl -L git.io/weave -o /usr/local/bin/weave
sudo chmod a+x /usr/local/bin/weave
sudo wget -O /usr/local/bin/scope https://git.io/scope
sudo chmod a+x /usr/local/bin/scope
sudo wget https://download.zerotier.com/dist/zerotier-one_1.1.4_amd64.deb
sudo dpkg -i zerotier-one_1.1.4_amd64.deb
sudo ifdata -pa zt0 > /ztip
