Getting started:
* Run the container with `docker compose up`.
* Start the project with `go run app/main.go` and see the results in the console.
 
Structure:
1. `app/main.go` - main go file where all initializations are done
2. `business/data/dbschema` - folder for sql files and operations that are done (creating tables, seeding, dropping)
3. `business/sys/database/database.go` - folder for initializing the db and status check

Thanks for the review!