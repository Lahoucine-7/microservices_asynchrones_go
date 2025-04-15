#!/usr/bin/env bash
# Copié de https://github.com/vishnubob/wait-for-it

host="$1"
shift
port="$1"
shift

while ! nc -z $host $port; do
  echo "⏳ En attente de $host:$port..."
  sleep 1
done

echo "✅ $host:$port est prêt"
exec "$@"
