package repository

import (
	"context"
	"diplom_back/internal/entity"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllEmployees(ctx context.Context, db *pgxpool.Pool) ([]entity.Employee, error) {
	var employees []entity.Employee

	rows, err := db.Query(ctx, "SELECT * FROM employee")
	if err != nil {
		slog.Warn("Не удалось выполнить запрос на всех сотрудников(")
	}
	defer rows.Close()

	for rows.Next() {
		var employee entity.Employee
		err := rows.Scan(
			&employee.ID,
			&employee.Login,
			&employee.Password,
			&employee.Name)
		if err != nil {
			slog.Warn("Не удалось отсканить сотрудников")
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		slog.Error("GetAllEmployees rows ERROR")
		return nil, err
	}
	return employees, nil
}

func PostNewEmployee(ctx context.Context, db *pgxpool.Pool, employee *entity.NewEmployee) (uint, error) {
	var id uint
	err := db.QueryRow(ctx, `INSERT INTO employee (login, password, name)
								VALUES ($1, $2, $3) RETURNING id`, employee.Login, employee.Password, employee.Name).Scan(&id)
	return id, err
}

func DeleteEmployee(ctx context.Context, db *pgxpool.Pool, id uint) error {
	_, err := db.Exec(ctx, "DELETE FROM employee WHERE id = $1", id)
	return err
}
