#! /bin/bash
if [ "$EUID" -ne 0 ]; then 
	echo "Please run as root"
	exit
fi
source .env
echo $DB_NAME 
echo $DB_USER
echo $DB_PASSWORD
mysql -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"

mysql -D $DB_NAME < ./sql-command/init.sql

mysql -e "CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASSWORD';"
mysql -e "GRANT ALL ON $DB_NAME.* TO '$DB_USER'@'localhost';"
