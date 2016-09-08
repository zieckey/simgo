#!/bin/bash
if [ ! -h /home/s/web_page_evaluator/data ]; then
    sudo ln -sf /home/s/data/web_page_evaluator /home/s/web_page_evaluator/data
fi

sudo ln -sf /home/s/web_page_evaluator/conf/dconf_wpe.ini /home/s/dconf_reload/etc/dconf_wpe.ini
sudo ln -sf /home/s/web_page_evaluator/conf/ngx_wpe_location.conf /usr/local/nginx/conf/ngx_wpe_location.conf

sudo chown -h cloud:cloud /home/s/web_page_evaluator/data
sudo chown -h cloud:cloud /home/s/dconf_reload/etc/dconf_wpe.ini
sudo chown -h cloud:cloud /usr/local/nginx/conf/ngx_wpe_location.conf
