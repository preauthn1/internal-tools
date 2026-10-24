FROM php:7.3-fpm

ARG PROJ_DIR=API_Agent

RUN set -eux; \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
            ssh \
            cron \
            curl \
            git \
            libzip-dev \
            unzip \
            zip \
            libmemcached-dev \
            libz-dev \
            libpq-dev \
            libjpeg-dev \
            libpng-dev \
            libfreetype6-dev \
            libssl-dev \
            libwebp-dev \
            libmcrypt-dev && \
    rm -rf /var/lib/apt/lists/*
RUN set -eux; \
    docker-php-ext-install pdo_mysql; \
    docker-php-ext-configure gd \
            --with-jpeg-dir=/usr/include \
            --with-webp-dir=/usr/include \
            --with-freetype-dir=/usr/include/freetype2; \
    docker-php-ext-configure zip --with-libzip; \
    docker-php-ext-install gd zip;
RUN mv /usr/local/etc/php/php.ini-development /usr/local/etc/php/php.ini
RUN sed -i "s/upload_max_filesize = 2M/upload_max_filesize = 20M/g" /usr/local/etc/php/php.ini && \
    sed -i "s/post_max_size = 8M/post_max_size = 21M/g" /usr/local/etc/php/php.ini && \
    echo "error_log = /proc/self/fd/2" >> /usr/local/etc/php/php.ini

COPY docker/$PROJ_DIR/fpm.conf /usr/local/etc/php-fpm.d/www.conf
RUN curl -sS https://example.invalid/installer | php -- --install-dir=/usr/local/bin --filename=composer
RUN composer self-update 1.8.5

RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan -t rsa -p 1022 example.invalid > /root/.ssh/known_hosts
COPY $PROJ_DIR/gitlab-ci/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

RUN apt-get update && \
    apt-get install -y nginx net-tools && \
    rm -rf /var/lib/apt/lists/*
COPY docker/$PROJ_DIR/nginx.conf /etc/nginx/nginx.conf

COPY docker/$PROJ_DIR/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /var/www/html

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
CMD ["bash", "-c", "cron && php-fpm -D && nginx"]
