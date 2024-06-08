### Employee Management

#### Technologies Used:

- Golang
- Gin Framework
- Postgres Database
- Postman for API testing

#### Project Structure:
```
assessment_techiebutler
│
├── cmd/
│    └── main.go
│
├── handler/
│    ├── config.go
│    ├── controllers.go
│    ├── database.go
│    ├── model.go
│    └── routes.go
|
├── test/
│    ├── main_test.go
│
└── go.mod
└── AssessmentEmployee.postman_collection.json
```
#### Start Application:

```
cd cmd
go run main.go
```

#### Run Tests:
```
go test ./test
```
