FROM php:8.4-fpm

ARG PROJ_DIR=apis

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    libpng-dev \
    libonig-dev \
    libxml2-dev \
    libzip-dev \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libicu-dev \
    libmemcached-dev \
    zlib1g-dev \
    libssl-dev \
    libcurl4-openssl-dev \
    unzip \
    zip \
    && rm -rf /var/lib/apt/lists/*

# Install PHP extensions
RUN docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) \
    pdo_mysql \
    mbstring \
    exif \
    pcntl \
    bcmath \
    gd \
    zip \
    intl \
    opcache \
    sockets \
    && pecl install redis \
    && pecl install memcached \
    && docker-php-ext-enable redis memcached

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

COPY docker/$PROJ_DIR/conf/php.ini /usr/local/etc/php/conf.d/custom.ini
COPY docker/$PROJ_DIR/conf/php-fpm-www.conf /usr/local/etc/php-fpm.d/www.conf

COPY docker/$PROJ_DIR/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /var/www/html

EXPOSE 9000

ENTRYPOINT ["/entrypoint.sh"]
CMD ["php-fpm"]
