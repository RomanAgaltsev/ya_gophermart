set -e

export PGPASSWORD=${POSTGRES_PASSWORD};

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_DB}" <<-EOSQL
  CREATE DATABASE ${POSTGRES_APP_DB};

  CREATE USER ${POSTGRES_APP_USER} WITH ENCRYPTED PASSWORD '${POSTGRES_APP_PASS}';
  GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_APP_DB} TO ${POSTGRES_APP_USER};
EOSQL

psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" --dbname "${POSTGRES_APP_DB}" <<-EOSQL
  GRANT ALL ON SCHEMA public TO ${POSTGRES_APP_USER};
EOSQL