#!/bin/sh

set -e

TIMEOUT=60
WAIT_INTERVAL=1
ELAPSED_TIME=0

while ! mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e 'SELECT 1'; do
    sleep $WAIT_INTERVAL
    ELAPSED_TIME=$(($ELAPSED_TIME + $WAIT_INTERVAL))

    if [ $ELAPSED_TIME -ge $TIMEOUT ]; then
        echo "Failed to connect to MySQL after $ELAPSED_TIME seconds. Exiting."
        exit 1
    fi

    echo "Waiting for MySQL... Attempted $ELAPSED_TIME/$TIMEOUT seconds."
done

echo "MySQL server is ready."
