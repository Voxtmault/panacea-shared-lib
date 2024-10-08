# Define default values for environment variables
DB_DRIVER ?= mysql
DB_HOST ?= host
DB_USER ?= user
DB_PASSWORD ?= password
DB_NAME ?= db_name
DB_PORT ?= 3306
DB_TLS_CONFIG ?= true
DB_ALLOW_NATIVE_PASSWORDS ?= true
DB_MULTI_STATEMENTS ?= false

REDIS_HOST ?= redis_host
REDIS_PORT ?= 6379
REDIS_DB ?= 0
REDIS_PASSWORD ?= redis_password
REDIS_EXPIRATION ?= 60 # Minutes

WS_URL ?= ws_url
WS_TOKEN ?= ws_token
WS_RECONNECT_INTERVAL ?= 10 # Seconds

AES_KEY ?= key
JWT_KEY ?= key
JWT_LIFE_SPAN ?= 1 # Day
PASSWORD_MIN_LENGTH ?= 8

KEY_PATH ?= key_path
CERT_PATH ?= cert_path

MAX_FILE_SIZE ?= 50 # MB
FILE_ROOT_PATH ?= $(CURDIR)/files

APP_MODE ?= devs
APP_PORT ?= port
APP_LANG ?= id
APP_TIMEZONE ?= Asia/Jakarta
APP_NAME ?= app_name
APP_HOST ?= app_host
APP_ROOT ?= app_root

LOG_PATH ?= ./log/server.log
ERR_LOG_PATH ?= ./log/error.log
LOG_MAX_SIZE ?= 30
LOG_MAX_BACKUP ?= 5
LOG_MAX_AGE ?= 30
LOG_COMPRESS ?= true

# Target to create .env file
create-env:
	@echo "Creating .env file..."
	@echo "# DB Configs" > .env
	@echo "DB_DRIVER=$(DB_DRIVER)" >> .env
	@echo "DB_HOST=$(DB_HOST)" >> .env
	@echo "DB_USER=$(DB_USER)" >> .env
	@echo "DB_PASSWORD=$(DB_PASSWORD)" >> .env
	@echo "DB_NAME=$(DB_NAME)" >> .env
	@echo "DB_PORT=$(DB_PORT)" >> .env
	@echo "DB_TLS_CONFIG=$(DB_TLS_CONFIG)" >> .env
	@echo "DB_ALLOW_NATIVE_PASSWORDS=$(DB_ALLOW_NATIVE_PASSWORDS)" >> .env
	@echo "DB_MULTI_STATEMENTS=$(DB_MULTI_STATEMENTS)" >> .env
	@echo "" >> .env
	@echo "# Redis Configs" >> .env
	@echo "REDIS_HOST=$(REDIS_HOST)" >> .env
	@echo "REDIS_PORT=$(REDIS_PORT)" >> .env
	@echo "REDIS_DB=$(REDIS_DB)" >> .env
	@echo "REDIS_PASSWORD=$(REDIS_PASSWORD)" >> .env
	@echo "REDIS_EXPIRATION=$(REDIS_EXPIRATION)" >> .env
	@echo "" >> .env
	@echo "# Websocket Configs" >> .env
	@echo "WS_URL=$(WS_URL)" >> .env
	@echo "WS_TOKEN=$(WS_TOKEN)" >> .env
	@echo "WS_RECONNECT_INTERVAL=$(WS_RECONNECT_INTERVAL)" >> .env
	@echo "" >> .env
	@echo "# Security Config" >> .env
	@echo "AES_KEY=$(AES_KEY)" >> .env
	@echo "JWT_KEY=$(JWT_KEY)" >> .env
	@echo "JWT_LIFE_SPAN=$(JWT_LIFE_SPAN)" >> .env
	@echo "PASSWORD_MIN_LENGTH=$(PASSWORD_MIN_LENGTH)" >> .env
	@echo "" >> .env
	@echo "# SSL Config" >> .env
	@echo "KEY_PATH=$(KEY_PATH)" >> .env
	@echo "CERT_PATH=$(CERT_PATH)" >> .env
	@echo "" >> .env
	@echo "# File Handling Config" >> .env
	@echo "MAX_FILE_SIZE=$(MAX_FILE_SIZE)" >> .env
	@echo "FILE_ROOT_PATH=$(FILE_ROOT_PATH)" >> .env
	@echo "" >> .env
	@echo "# General Configs" >> .env
	@echo "APP_MODE=$(APP_MODE)" >> .env
	@echo "APP_PORT=$(APP_PORT)" >> .env
	@echo "APP_LANG=$(APP_LANG)" >> .env
	@echo "APP_TIMEZONE=$(APP_TIMEZONE)" >> .env
	@echo "APP_NAME=$(APP_NAME)" >> .env
	@echo "APP_HOST=$(APP_HOST)" >> .env
	@echo "APP_ROOT=$(APP_ROOT)" >> .env
	@echo "" >> .env
	@echo "LOG_PATH=$(LOG_PATH)" >> .env
	@echo "ERR_LOG_PATH=$(ERR_LOG_PATH)" >> .env
	@echo "LOG_MAX_SIZE=$(LOG_MAX_SIZE)" >> .env
	@echo "LOG_MAX_BACKUP=$(LOG_MAX_BACKUP)" >> .env
	@echo "LOG_MAX_AGE=$(LOG_MAX_AGE)" >> .env
	@echo "LOG_COMPRESS=$(LOG_COMPRESS)" >> .env
	@echo ".env file created successfully."

# Default target
all: create-env

.PHONY: create-env