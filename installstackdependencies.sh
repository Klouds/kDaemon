#!/bin/bash
# RUN THIS SCRIPT AS ROOT ON A FRESHLY INSTALLED UBUNTU OR DEBIAN SERVER.
wget get.docker.io
sh index.html
apt-get install -y python-pip
pip install docker-compose
