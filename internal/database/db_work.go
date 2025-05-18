package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type User struct {
	ID       int
	Login    string
	Password string
}

type Expression struct {
	ID         int
	User_id    int
	Expression string
	Result     float64
}

func UpdateExpression(db *sql.DB, ID int, result float64) error {
	query := `
		UPDATE expressions 
		SET result = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, result, ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %v", err)
	}
	return nil
}

func AddExpression(db *sql.DB, Expr *Expression) error {
	_, err := db.Exec(
		"INSERT INTO expressions (id, user_id, expression, result) VALUES (?, ?, ?, ?)",
		Expr.ID,
		Expr.User_id,
		Expr.Expression,
		Expr.Result,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetUserID(db *sql.DB, id int) (*User, error) {
	query := "SELECT id, login, password FROM users WHERE id = ?"
	row := db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с id %d не найден", id)
		}
	}
	return &user, nil
}

func GetUser(db *sql.DB, login string) (*User, error) {
	query := "SELECT id, login, password FROM users WHERE login = ?"
	row := db.QueryRow(query, login)
	var user User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с логином %s не найден", login)
		}
	}
	return &user, nil
}

func AddUser(db *sql.DB, login string, passw string) error {
	_, err := db.Exec(
		"INSERT INTO users (login, password) VALUES (?, ?)",
		login,
		passw,
	)
	if err != nil {
		return err
	}
	return nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            login TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL
        );
		CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expression TEXT NOT NULL,
		result REAL NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
		);`)
	return err
}

func InitDB() (*sql.DB, error) {
	dbName := "store.db"
	dbPath := filepath.Join("internal", "database", dbName)
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s?cache=shared", dbPath))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Printf("База данных готова к работе!\n")
	return db, err
}
