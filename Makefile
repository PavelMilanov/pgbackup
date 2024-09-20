postgres:
	@docker run --rm --name test-postgres -e POSTGRES_USER=admin -e POSTGRES_DB=dev -e POSTGRES_PASSWORD=admin -d -p 5432:5432 postgres:16

build:
	@docker buildx build . -t rosomilanov/pgbackup:0.1 --build-arg UID_DOCKER=10000 --build-arg GID_DOCKER=10001 --build-arg USER_DOCKER=pgbackup