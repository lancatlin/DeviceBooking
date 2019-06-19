#! /bin/bash
if [ "$EUID" -ne 0 ]; then 
	echo "Please run as root"
	exit 1
fi
source .env
echo $DB_NAME, $DB_USER, $DB_PASSWORD
mysql -h "mariadb" -u $DB_USER -p=$DB_PASSWORD -D $DB_NAME -e "SOURCE ./sql-command/init.sql;"