-- name: GetIncompleteTasks :many
SELECT * FROM task
WHERE completed = 0
ORDER BY complete_by;

-- name: GetTasksForToday :many
SELECT * FROM task
WHERE complete_by = ? AND completed = 0;

-- name: NewTask :exec
INSERT INTO task(name, complete_by, completed)
VALUES(?, ?, 0);

-- name: CompleteTask :exec
UPDATE task
SET completed = 1
WHERE id = ?;
