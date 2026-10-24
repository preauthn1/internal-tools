#!/bin/sh
set -e

# 確保 Laravel 執行期目錄可被 www-data 寫入（bind mount 掛入時擁有者可能是 root）
if [ -d /var/www/storage ]; then
    chown -R www-data:www-data /var/www/storage /var/www/bootstrap/cache 2>/dev/null || true
    chmod -R ug+rwx /var/www/storage /var/www/bootstrap/cache 2>/dev/null || true
fi

exec "$@"
