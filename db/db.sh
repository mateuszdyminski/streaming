#!/usr/bin/env bash


usage() {
	cat <<EOF
Usage: $(basename $0) <command> <arg>

Useful commands:
	update - update schema to newest one via liquibase
	validate - validate changeset
EOF
	exit 1
}


CMD="$1"
shift
case "$CMD" in
	update)
		liquibase --changeLogFile=sql/changelog.sql --driver=com.microsoft.sqlserver.jdbc.SQLServerDriver --classpath=jdbc/postgresql-42.1.4.jar --url="jdbc:sqlserver://devsmtdbserver.database.windows.net:1433;database=smtdb-dev;encrypt=true;trustServerCertificate=false;hostNameInCertificate=*.database.windows.net;loginTimeout=30;" --username=kone-smt --password=s1!m2@t3# update
	;;
	validate)
		liquibase --changeLogFile=sql/changelog.sql --driver=com.microsoft.sqlserver.jdbc.SQLServerDriver --classpath=jdbc/postgresql-42.1.4.jar --url="jdbc:postgresql://localhost/str-db?user=str-user&password=pass123&ssl=false" validate
	;;
	*)
		usage
	;;
esac
