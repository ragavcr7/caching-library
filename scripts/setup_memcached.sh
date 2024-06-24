#!/bin/bash

# Check whether memcached is already installed
if ! command -v memcached &> /dev/null; then
    echo "memcached is not installed. Installing memcached..."
    apt-get update
    apt-get install memcached -y
else
    echo "memcached is already installed"
fi 

# Start memcached server
echo "Starting memcached server..."
systemctl start memcached

# Enable memcached to start on system boot
systemctl enable memcached

# Display memcached server status
systemctl status memcached

# Command to execute this file: chmod +x setup_memcached.sh
