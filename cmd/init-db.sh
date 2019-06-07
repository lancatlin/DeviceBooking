#! /bin/bash
if ["$EUID" -ne 0] then 
	echo "Please run as root"
	exit
fi

DB_USER=$1
DB_PASSWORD=$2

mysql < ../sql-command/init.sql

mysql -e "CREATE USER '$DB_USER'@'localhost' IDENTIFIED BY $DB_PASSWORD"
