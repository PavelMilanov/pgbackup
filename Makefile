version=

# only for development
setenv:
	@export JWT_KEY="very_secret_key"
	@export AES_KEY="key3456789012345"

# only for development
postgres:
	@docker run --rm --name dev-postgres -e POSTGRES_USER=admin -e POSTGRES_DB=dev -e POSTGRES_PASSWORD=admin -d -p 5432:5432 postgres:16

# only for development
test_postgres:
	@docker run --rm --name test-postgres -e POSTGRES_USER=test -e POSTGRES_DB=test -e POSTGRES_PASSWORD=test -d -p 5433:5432 postgres:16

# only for development
build:
	@docker buildx build --platform linux/amd64 . -t pgbackup:dev

release:
	@docker buildx build --platform linux/amd64 . -t rosomilanov/pgbackup:${version} --build-arg VERSION=${version}
	@docker buildx build --platform linux/amd64 . -t rosomilanov/pgbackup:latest --build-arg VERSION=${version}

push:
	@docker push rosomilanov/pgbackup:${version}
	@docker push rosomilanov/pgbackup:latest
