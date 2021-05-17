package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/medivh13/koalatest/internal/models"
	"github.com/medivh13/koalatest/internal/repository"
	btbErrors "github.com/medivh13/koalatest/pkg/errors"
)

const (
	Register = `INSERT INTO customers (customer_id, customer_name, email, phone_number, dob, sex, salt, password, create_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	CekEmail = `SELECT customer_id from customers where email = $1`
	CekPhone = `SELECT customer_id from customers where phone_number = $1`
)

var statement PreparedStatement

type PreparedStatement struct {
	register *sqlx.Stmt
	cekEmail *sqlx.Stmt
	cekPhone *sqlx.Stmt
}

type MySQLRepo struct {
	Conn *sqlx.DB
}

func NewMySQLRepo(Conn *sqlx.DB) repository.Repository {

	repo := &MySQLRepo{Conn}
	InitPreparedStatement(repo)
	return repo
}

//Preparex query for database queries
func (m *MySQLRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := m.Conn.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *MySQLRepo) {
	statement = PreparedStatement{
		register: m.Preparex(Register),
		cekEmail: m.Preparex(CekEmail),
		cekPhone: m.Preparex(CekPhone),
	}
}

func (m *MySQLRepo) Register(data *models.Customers) error {
	custID := []models.ExistingCustomer{}

	err := m.Conn.Select(&custID, CekEmail, data.Email)
	if err != nil {
		log.Println("Failed Query to get Existing Email : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(custID) > 0 {
		return fmt.Errorf(btbErrors.ErrorExistingData)
	}

	err = m.Conn.Select(&custID, CekPhone, data.PhoneNumber)
	if err != nil {
		log.Println("Failed Query to get Existing Phone : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(custID) > 0 {
		return fmt.Errorf(btbErrors.ErrorExistingData)
	}

	result, err := statement.register.Exec(data.CustomerId, data.CustomerName, data.Email, data.PhoneNumber, data.Dob, data.Sex, data.Salt, data.Password, data.CreatedDate)
	if err != nil {
		log.Println("Failed Query Register : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed Query Register : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}
	if rows == 0 {
		log.Println("Failed Query Register: no data changed")
		return fmt.Errorf(btbErrors.ErrorNoDataChange)
	}

	return nil
}
