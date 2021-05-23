# Data Migration

## Install migrate

### Mac

`brew install golang-migrate`

## Migration steps

### Create migration script

`migrate create -ext sql -dir db/migrations -seq <COMMIT MESSAGE>`

### Migrate step up

`migrate -verbose -path db/migrations -database "mysql://<name>:<password>@tcp(<server>:3306)/apps_test" up`

### Migrate step down

You can either step down to step #n (according to migrations folder), by `down n`
or you can revert all by only `down`