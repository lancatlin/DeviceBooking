#! /bin/sh
mysql -u $MYSQL_USER --password=$MYSQL_PASSWORD -D $MYSQL_DATABASE -e "SOURCE /init.sql;"