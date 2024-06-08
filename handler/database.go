package handler

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		connStr := GetDBConnectionString()
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		createTableQuery := `
        CREATE TABLE IF NOT EXISTS employees (
            id SERIAL PRIMARY KEY,
            name TEXT,
            position TEXT,
            salary REAL
        );`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			log.Fatal(err)
		}
	})
}

type EmployeeStore struct {
	db   *sql.DB
	lock sync.Mutex
}

func NewEmployeeStore() *EmployeeStore {
	return &EmployeeStore{db: db}
}

func (store *EmployeeStore) CreateEmployee(emp *Employee) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	query := `INSERT INTO employees (name, position, salary) VALUES ($1, $2, $3) RETURNING id`
	err := store.db.QueryRow(query, emp.Name, emp.Position, emp.Salary).Scan(&emp.ID)
	if err != nil {
		return err
	}
	return nil
}

func (store *EmployeeStore) GetEmployeeByID(id int) (*Employee, error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	emp := &Employee{}
	query := `SELECT id, name, position, salary FROM employees WHERE id = $1`
	err := store.db.QueryRow(query, id).Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary)
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (store *EmployeeStore) UpdateEmployee(emp *Employee) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	query := `UPDATE employees SET name=$1, position=$2, salary=$3 WHERE id=$4`
	_, err := store.db.Exec(query, emp.Name, emp.Position, emp.Salary, emp.ID)
	if err != nil {
		return err
	}
	return nil
}

func (store *EmployeeStore) DeleteEmployee(id int) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	query := `DELETE FROM employees WHERE id=$1`
	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (store *EmployeeStore) ListEmployees(offset, limit int) ([]Employee, error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	query := `SELECT id, name, position, salary FROM employees ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := store.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Position, &emp.Salary); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}
