#! /bin/bash
if [ "$EUID" -ne 0 ]; then 
	echo "Please run as root"
	exit 1
fi
source .env
echo $DB_HOST $DB_NAME, $DB_USER, $DB_PASSWORD
mysql -h $DB_HOST -e "
DROP DATABASE IF EXISTS $DB_NAME;
CREATE DATABASE IF NOT EXISTS $DB_NAME; 
DROP USER IF EXISTS '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
CREATE USER '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASSWORD';
FLUSH PRIVILEGES;
GRANT ALL ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
USE $DB_NAME;
SOURCE ./mariadb/init.sql;"