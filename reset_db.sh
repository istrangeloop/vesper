#!/bin/bash 

export $(cat .env | xargs)
HERE=$PWD
bash -c "pg_dump -h $PSQL_HOST -p $PSQL_PORT -U $PSQL_USER arturdb > ./data/arturdb.sql"
bash -c "psql -h $PSQL_HOST -p $PSQL_PORT -U $PSQL_USER < $HERE/initdb.sql"
bash -c "psql -h $PSQL_HOST -p $PSQL_PORT -U $PSQL_USER -d arturdb < ./data/arturdb.sql"