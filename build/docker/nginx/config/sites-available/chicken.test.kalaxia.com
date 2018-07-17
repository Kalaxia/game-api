server {
    listen 80;
    listen [::]:80;

    server_name chicken.test.kalaxia.com;

    access_log /var/log/nginx/https_chicken_game.access.log;
    error_log /var/log/nginx/https_chicken_game.error.log;

    merge_slashes on;

    location /api {
        proxy_http_version 1.1;
        proxy_pass http://chicken_api;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header        X-Scheme        $scheme;
        proxy_set_header        Connection "";
        proxy_buffering off;
        proxy_ignore_client_abort on;
        proxy_read_timeout 7d;
        proxy_send_timeout 7d;
    }

    location / {
        root /srv/app;
    }

    location ~ /\.ht {
        deny all;
    }

    location ~ /public/log/stats/ {
        deny all;
    }
}
