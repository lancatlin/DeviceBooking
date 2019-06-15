#! /bin/bash
if [ "$EUID" -ne 0 ]; then 
	echo "Please run as root"
	exit 1
fi
DB_NAME=$1
DB_USER=$2
DB_PASSWORD=$3
echo $DB_NAME $DN_USER $DB_PASSWORD

mysql -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"
mysql -e "DELETE FROM mysql.user WHERE User = '$DB_USER'"
mysql -e "FLUSH PRIVILEGES"
mysql -e "CREATE USER '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASSWORD';"
mysql -e "FLUSH PRIVILEGES"
mysql -e "GRANT ALL ON $DB_NAME.* TO '$DB_USER'@'localhost';"
mysql -e "FLUSH PRIVILEGES"

mysql -D $DB_NAME -e "SOURCE ./sql-command/init.sql"