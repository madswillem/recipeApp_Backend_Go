# RECIPEAPP
Welcome to the **Recipe App** repository! This project is a user-friendly and foss application that allows users to browse, save,  share and search recipes. With the upcoming v1.0 release users will be able to get personalized recommendations

## Table of Contents

## Installation
### Linux  with PostgreSQL
#### Releases
No releases yet wait until milestone #7
#### Build from Source
1. Dependencies:
	- Download latest go version [Go](https://go.dev/dl/)
	- Download and configure latest go version [PostgreSQL](https://www.postgresql.org/download/)
2. Clone repo `git clone https://github.com/madswillem/recipeApp_Backend_Go.git`
3. Go into the directory `cd recipeApp_Backend_Go`
4. Add `.env` file and add `DB=(db connection string)` to it
5. Download the dependencies `make tidy`
6. Build `make build`
7. Run the DB
8. Run the app `./bin/main`

#### Single executable
Not yet available wait until milestone #7

### Windows
#### Releases
No releases yet wait until milestone #7
#### Build from Source
1. Dependencies:
	- Download latest go version [Go](https://go.dev/dl/)
	- Download and configure latest go version [PostgreSQL](https://www.postgresql.org/download/)
2. Clone repo `git clone https://github.com/madswillem/recipeApp_Backend_Go.git`
3. Go into the directory `cd recipeApp_Backend_Go`
4. Add `.env` file and add `DB=(db connection string)` to it
5. Download the dependencies `make tidy`
6. Build `make build`
7. Run the DB
8. Run the app `./bin/main.exe`

### Single executable
Not yet available wait until milestone #7

### MacOS
There is no official way to install the RecipeApp on MacOS. You might be able to build the app from source, I can't verify that though.
