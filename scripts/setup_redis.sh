#!/bin/bash

#lets check whether redis is already installed or not:
if ! commad -v redis-server &> /dev/null; then
    echo "redis is not installed. kindly install redis first"
    sudo apt-get update
    sudo apt-get install redis-server -y
else
    echo "redis is already installed"
fi 

#start redis server 
echo "starting redis server..."
sudo systemctl start redis-server
#enabling redis to start system boot
sudo systemctl enable redis-server 
#display redis server status
sudo systemctl status redis-server
#command to execute this file chmod +x setup_redis.sh