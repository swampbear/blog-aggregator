# Set username for psql

```bash
sudo passwd postgres
```

# Cmd for starting db server

```bash
sudo service postgresql start
```

# To connect run:

```bash
psql
```

# Docker

## start

```bash
docker start postgres18
```

to enter the psql CLI

```bash
psql -h localhost -U postgres -p 5432
```

SQL query to create database

```SQL
CREATE DATABASE gator
```

command to connect to the created database

```bash
\c gator
```
