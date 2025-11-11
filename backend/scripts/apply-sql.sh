#!/bin/bash

# Скрипт для применения SQL схемы напрямую в PostgreSQL

set -e

echo "🚀 Применение SQL схемы к PostgreSQL..."

# Проверка переменных окружения
DB_HOST=${POSTGRES_HOST:-localhost}
DB_PORT=${POSTGRES_PORT:-5432}
DB_USER=${POSTGRES_USER:-kanban}
DB_PASSWORD=${POSTGRES_PASSWORD:-kanban123}
DB_NAME=${POSTGRES_DB:-kanban}

echo "📝 Подключение к PostgreSQL: $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"

# Применение SQL через docker exec (если PostgreSQL в Docker)
if docker ps | grep -q kanban_postgres; then
    echo "📦 Применение через Docker..."
    docker exec -i kanban_postgres psql -U "$DB_USER" -d "$DB_NAME" < scripts/init-db.sql
    echo "✅ Схема применена успешно!"
else
    # Или через psql напрямую
    echo "📝 Применение через psql..."
    export PGPASSWORD="$DB_PASSWORD"
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f scripts/init-db.sql
    echo "✅ Схема применена успешно!"
fi

echo "✨ Готово! Теперь можно сгенерировать Prisma Client"

