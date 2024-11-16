version=

# only for development
setenv:
	@export LOGIN="admin"
	@export PASSWORD="admin"
	@export JWT_KEY="very_secret_key"
	@export AES_KEY="key3456789012345"

# only for development
postgres:
	@docker run --rm --name test-postgres -e POSTGRES_USER=admin -e POSTGRES_DB=dev -e POSTGRES_PASSWORD=admin -d -p 5432:5432 postgres:16

# only for development
build:
	@docker buildx build --platform linux/amd64 . -t rosomilanov/pgbackup:${version}

release:
	@docker build . -t rosomilanov/pgbackup:${version} --build-arg VERSION=${version}

push:
	@docker push rosomilanov/pgbackup:${version}