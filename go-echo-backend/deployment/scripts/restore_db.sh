#!/usr/bin/env bash

POSTGRES_HOST=${1}
POSTGRES_PORT=${2}
POSTGRES_USER=${3}
POSTGRES_PASSWORD=${4}
POSTGRES_DATABASE=${5}

DESTINATION="${POSTGRES_DATABASE}_$(date +"%Y-%m-%dT%H:%M:%SZ").sql.gz"

BACKUP_FILE=${6:-$DESTINATION}

if [[ -z $POSTGRES_HOST ]];then
    echo "POSTGRES_HOST is required"; echo;
    exit 1
fi


if [[ -z $POSTGRES_PORT ]];then
    echo "POSTGRES_PORT is required"; echo;
    exit 1
fi


if [[ -z $POSTGRES_USER ]];then
    echo "POSTGRES_USER is required"; echo;
    exit 1
fi

if [[ -z $POSTGRES_PASSWORD ]];then
    echo "POSTGRES_PASSWORD is required"; echo;
    exit 1
fi


if [[ -z $POSTGRES_DATABASE ]];then
    echo "POSTGRES_DATABASE is required"; echo;
    exit 1
fi


if [[ -z $BACKUP_FILE ]];then
    echo "BACKUP_FILE is required"; echo;
    exit 1
fi

cmd="PGPASSWORD='$POSTGRES_PASSWORD' pg_restore --no-privileges --no-owner --verbose --clean -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DATABASE -C $BACKUP_FILE"

connect_cmd="PGPASSWORD='$POSTGRES_PASSWORD' psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DATABASE"
terminal_pod="select pg_terminate_backend(pid) from pg_stat_activity where datname='$POSTGRES_DATABASE';"
drop_cmd="DROP DATABASE $POSTGRES_DATABASE;"

echo "Drop DB SQL instruction";
echo $connect_cmd;
echo $terminal_pod;
echo $drop_cmd;
echo;

echo "Restore SQL cmd";
echo $cmd; echo;
