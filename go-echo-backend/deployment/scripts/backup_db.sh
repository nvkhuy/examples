#!/usr/bin/env bash
ENV=${1}
POSTGRES_HOST=${2}
POSTGRES_PORT=${3}
POSTGRES_USER=${4}
POSTGRES_PASSWORD=${5}
POSTGRES_DATABASE=${6}
RUN_CMD=${7:-n}

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

mkdir -p db_backup/${ENV}

DESTINATION="db_backup/${ENV}/${POSTGRES_DATABASE}_$(date +"%Y-%m-%dT%H:%M:%S").sql"

cmd="PGPASSWORD='$POSTGRES_PASSWORD' pg_dump --verbose -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -f $DESTINATION $POSTGRES_DATABASE -Fc"

echo "Backup SQL cmd"; echo;
echo $cmd; echo;

if [ $RUN_CMD != "${RUN_CMD#[Yy]}" ] ;then
  eval $cmd
fi;
