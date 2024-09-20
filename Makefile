version=

postgres:
	@docker run --rm --name test-postgres -e POSTGRES_USER=admin -e POSTGRES_DB=dev -e POSTGRES_PASSWORD=admin -d -p 5432:5432 postgres:16

build:
	@docker buildx build --platform linux/amd64 . -t rosomilanov/pgbackup:${version} --build-arg UID_DOCKER=1000 --build-arg USER_DOCKER=pgbackup
