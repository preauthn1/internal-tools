#!/bin/bash
set -e

# Wait for volume mount if files don't exist
if [ ! -f /var/www/html/composer.json ]; then
    echo "Waiting for volume mount..."
    sleep 2
fi

mkdir -p /var/www/html/storage/logs \
         /var/www/html/storage/app \
         /var/www/html/storage/uploads \
         /var/www/html/bootstrap/cache \
         /var/www/html/storage/framework/views \
         /var/www/html/storage/framework/sessions \
         /var/www/html/storage/framework/cache

# Run composer install and verify full success (only if composer.json exists)
cd /var/www/html
if [ -f composer.json ]; then
    echo "Running composer install..."
    
    # Run composer install
    if ! composer install --no-interaction --prefer-dist --optimize-autoloader; then
        echo "✗ ERROR: composer install failed"
        exit 1
    fi
    
    # Verify composer install fully completed
    if [ ! -d vendor ]; then
        echo "✗ ERROR: vendor directory not found after composer install"
        exit 1
    fi
    
    if [ ! -f vendor/autoload.php ]; then
        echo "✗ ERROR: vendor/autoload.php not found - composer install incomplete"
        exit 1
    fi
    
    if [ -f composer.lock ] && [ ! -f vendor/composer/installed.json ]; then
        echo "✗ ERROR: composer.lock exists but installed.json not found - install incomplete"
        exit 1
    fi
    
    echo "✓ Composer install fully succeeded - all dependencies installed"
else
    echo "⚠ WARNING: composer.json not found, skipping composer install"
fi

# Set permissions
chmod -R 775 /var/www/html/storage
chmod -R 775 /var/www/html/bootstrap/cache
chown -R 1000:www-data /var/www/html

# Ensure log directories exist
mkdir -p /var/log/php-fpm /var/log/nginx

# Laravel: Already writes to files (persisted in /var/www/html/storage/logs/)
# Also tail to stdout for docker logs visibility
( tail -f /var/www/html/storage/logs/laravel*.log 2>/dev/null & )

# PHP-FPM & Nginx: Write to files (persisted) AND stdout/stderr (docker logs)
# Tail the log files to stdout/stderr
( tail -f /var/log/php-fpm/www-error.log 2>/dev/null & )
( tail -f /var/log/nginx/error.log 2>/dev/null & )
( tail -f /var/log/nginx/access.log 2>/dev/null & )

# Execute the original command
exec "$@"

