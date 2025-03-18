#!/bin/bash

# Параметры подключения к базе данных
DB_HOST="localhost"
DB_PORT="5490"
DB_NAME="main"
DB_USER="rw_main"
DB_PASSWORD="hrgr2ss10nd"

# Применяем миграции в правильном порядке
echo "Применение миграций..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/000_create_base_tables.sql
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/001_create_homework_tables.sql

echo "Миграции успешно применены!" 