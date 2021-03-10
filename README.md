# Luccryptous

Service API to generate and encrypt passwords.

This service is made to be public, and to encrypt passwords without having the key.


# Installation

## Build

``` shell
podman build -t luccryptous .
```


## Execution

``` shell
podman run -d -m 500m -p 3000:3000 --name luccryptous luccryptous
```

Install as a service

``` shell
podman generate systemd luccryptous > /etc/systemd/system/luccryptous.service
systemctl enable --now luccryptous
```


# Develop with Dockerfile.dev

``` shell
podman build -t luccryptous:dev ./Dockerfile.dev
podman run -it --rm -m 1024m -p 5000:5000,35729:35729 -v ~/Projects/luccryptous/views/:/app/:Z --name luccryptous_dev luccryptous:dev
```


# Configuration

The configuration is in Toml file or environment variables

Create a file *luccryptous.toml* or use the example, then edit the values.


| Section             | Key   | Default | Description                                   |
|---------------------|-------|---------|-----------------------------------------------|
| General             | Key   |         | Encryption key, required                      |
| General             | Debug | false   | Set to true to switch gin gonic in debug mode |
| Password Generation | size  | 42      | Size                                          |


# External dependencies

- spf13/viper
- gin-gonic
- google/uuid
