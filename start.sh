#!/bin/sh

# 若指令传回值不等于0，则立即退出shell。
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up # $-表示后面接的是变量3v

echo "start the app"
exec "$@" # @-所有变量