# Luccryptous

Service API to generate and encrypt passwords.

# Configuration

The configuration is in Toml file or environment variables

Create a file luccryptous.toml or use the example, then edit the values.

| Section             | Key   | Default | Description                                   |
|---------------------+-------+---------+-----------------------------------------------|
| General             | Key   |         | Encryption key, required                      |
| General             | Debug | false   | Set to true to switch gin gonic in debug mode |
| Password Generation | size  | 42      | Size                                          |
| Password Generation 

# Installation

## Build

podman build -t luccryptous .

## Execution

podman run -it --rm -m 500m -p 3000:3000 --name luccryptous luccrypouts

Ports:
- 3000

podman run -it --rm -m 1024m --publish 5000:5000,35729:35729 --volume /home/arouene/Projects/luccryptous/views/:/app/:Z --name nodejs nodejs

# External dependencies

- spf13/viper
- gin-gonic
- google/uuid
