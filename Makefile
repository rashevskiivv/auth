ifeq ($(OS),Windows_NT)
CUR_DIR=$(shell echo %CD%)
else
CUR_DIR=$(shell pwd)
endif

IMAGE=auth
TAG=latest
RELEASE_NAME=auth
DC_FILE=-f ${CUR_DIR}/deployment/docker-compose.yaml

compile:
	docker build --no-cache -f .docker/Dockerfile -t ${IMAGE}:${TAG} --target builder .

copy-env:
	cp deployment/.env.example deployment/.env

copy-env-windows:
	copy deployment\.env.example deployment\.env

deploy:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d

deploy-app:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d app

deploy-postgres:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d postgres_db postgres_migrate

deploy-redis:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} up -d redis_db

delete:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf

delete-app:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf app

delete-postgres:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf postgres_db postgres_migrate

delete-redis:
	cd deployment && docker-compose ${DC_FILE} -p ${RELEASE_NAME} rm -sf redis_db