#!/bin/bash
set -euo pipefail

# Wait for volume mount if files don't exist (with timeout)
if [ ! -f /var/www/html/composer.json ]; then
    echo "Waiting for volume mount..."
    max_wait=5
    waited=0
    while [ ! -f /var/www/html/composer.json ] && [ $waited -lt $max_wait ]; do
        sleep 1
        waited=$((waited + 1))
    done
    if [ ! -f /var/www/html/composer.json ]; then
        echo "⚠ WARNING: composer.json not found after ${max_wait}s, continuing anyway"
    fi
fi

mkdir -p /var/log/php-fpm \
         /var/log/nginx

#chmod -R 775 /var/www/html/storage /var/www/html/bootstrap/cache 2>/dev/null 
chown -R 1000:www-data /var/www/html 2>/dev/null

( tail -f /var/www/html/logs/*.log 2>/dev/null & ) || true

# PHP-FPM & Nginx logs
( tail -f /var/log/php-fpm/www-error.log 2>/dev/null & ) || true
( tail -f /var/log/nginx/error.log 2>/dev/null & ) || true
( tail -f /var/log/nginx/access.log 2>/dev/null & ) || true

exec "$@"

