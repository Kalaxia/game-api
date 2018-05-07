#!/bin/bash
IFS=',' read -a vhosts <<< "$NGINX_ENABLED_VHOST"

for vhost in "${vhosts[@]}"
do
  ln -s /etc/nginx/sites-available/$vhost /etc/nginx/sites-enabled/
done

exec "$@"
