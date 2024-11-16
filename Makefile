version=

setenv:
	@export LOGIN="admin"
	@export PASSWORD="admin"
	@export JWT_KEY="very_secret_key"
	@export AES_KEY="key3456789012345"

postgres:
	@docker run --rm --name test-postgres -e POSTGRES_USER=admin -e POSTGRES_DB=dev -e POSTGRES_PASSWORD=admin -d -p 5432:5432 postgres:16

build:
	@docker buildx build --platform linux/amd64 . -t rosomilanov/pgbackup:${version} --build-arg UID_DOCKER=10000 --build-arg USER_DOCKER=pgbackup
	@docker buildx build --platform linux/amd64 . -t rosomilanov/pgbackup:latest --build-arg UID_DOCKER=10000 --build-arg USER_DOCKER=pgbackup
