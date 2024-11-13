// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package tasks

import (
	"context"
	"database/sql"
)

const completeTask = `-- name: CompleteTask :exec
UPDATE task
SET completed = 1
WHERE id = ?
`

func (q *Queries) CompleteTask(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, completeTask, id)
	return err
}

const getIncompleteTasks = `-- name: GetIncompleteTasks :many
SELECT id, name, complete_by, completed FROM task
WHERE completed = 0
ORDER BY complete_by
`

func (q *Queries) GetIncompleteTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getIncompleteTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CompleteBy,
			&i.Completed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasksForToday = `-- name: GetTasksForToday :many
SELECT id, name, complete_by, completed FROM task
WHERE complete_by = ? AND completed = 0
`

func (q *Queries) GetTasksForToday(ctx context.Context, completeBy sql.NullString) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksForToday, completeBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CompleteBy,
			&i.Completed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newTask = `-- name: NewTask :exec
INSERT INTO task(name, complete_by, completed)
VALUES(?, ?, 0)
`

type NewTaskParams struct {
	Name       string
	CompleteBy sql.NullString
}

func (q *Queries) NewTask(ctx context.Context, arg NewTaskParams) error {
	_, err := q.db.ExecContext(ctx, newTask, arg.Name, arg.CompleteBy)
	return err
}