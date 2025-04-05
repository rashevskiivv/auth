ifeq ($(OS),Windows_NT)
CUR_DIR=$(shell echo %CD%)
else
CUR_DIR=$(shell pwd)
endif

APP_IMAGE=auth_local
APP_TAG=latest
DB_IMAGE=auth_db
DB_TAG=latest
RELEASE_NAME=auth
DC_FILE=-f ${CUR_DIR}/deployment/docker-compose.yaml

.PHONY: compile compile-db copy-env copy-env-windows deploy deploy-app deploy-pg delete delete-app delete-pg

compile:
	docker build --no-cache -f .docker/Dockerfile -t ${APP_IMAGE}:${APP_TAG} --target builder .

compile-db:
	docker build --no-cache -f .docker/PGDockerfile -t ${DB_IMAGE}:${DB_TAG} .

copy-env:
	cp deployment/.env.example deployment/.env

copy-env-windows:
	copy deployment\.env.example deployment\.env

deploy:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d

deploy-app:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d app

deploy-pg:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d auth_db postgres_migrate

delete:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf

delete-app:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf app

delete-pg:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf auth_db  postgres_migrate
