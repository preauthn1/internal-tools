#!/bin/bash 
composer i 
php artisan route:cache 
php artisan config:cache 
php artisan config:clear
php artisan cache:clear 
composer dump-autoload 
php artisan optimize 
php artisan view:clear

php artisan route:clear
php artisan route:cache 
php artisan config:cache 
php artisan config:clear
php artisan cache:clear 