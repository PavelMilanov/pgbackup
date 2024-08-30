postgres:
	@docker run --rm --name test-postgres -e POSTGRES_USER=admin -e POSTGRES_DB=dev -e POSTGRES_PASSWORD=admin -d -p 5432:5432 postgres:15