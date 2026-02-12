FROM php:8.2-apache

RUN a2enmod rewrite

WORKDIR /var/www/html
COPY public/ /var/www/html/
COPY storage/ /var/www/storage/
RUN chown -R www-data:www-data /var/www/storage

ENV APACHE_DOCUMENT_ROOT=/var/www/html
