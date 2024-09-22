# CONTRIBUTING
> [!IMPORTANT]
> This is not nearly finished if you have any question please create an issue and I'll do my best to answer it and update this document

If you want to contribute to this project please follow this quide

### Issues
Please always create an issue for every bug you want to fix or feature you want to add.
Ideale wait until an issue has been added to a milestone before creating a pull request

### Commits
Please squash multible commits into one and use profesional language

### Pull requests
Please alwas [link a pull request to an issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue)


## Build from Source
1. Dependencies:
	- Download latest go version [Go](https://go.dev/dl/)
	- Download and configure latest postgres version [PostgreSQL](https://www.postgresql.org/download/)
2. Clone repo `git clone https://github.com/madswillem/recipeApp_Backend_Go.git`
3. Go into the directory `cd recipeApp_Backend_Go`
4. Add `.env` file and add `DB=lökjsdglkjsglö` to it
5. Install dependencies `make tidy`
6. Build `make build`
7. Run the DB
8. Run the app `./bin/main`
