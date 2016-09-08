#!/bin/bash
sudo yum -y install safe-nginx libstatus dconf_reload qlogd libqlog libcloudcom
sudo mkdir -p /usr/local/nginx/conf/include
echo "* * * * * /usr/bin/lockf -t 0 /home/s/dconf_reload/log/dconf_wpe.lock /home/s/dconf_reload/bin/dctl check dconf_wpe" >> /tmp/cloud.crontab
crontab -u root /tmp/cloud.crontab

