-- name: CreateTask :one
INSERT INTO tasks (user_id, description, created_at) VALUES ($1, $2, $3)
RETURNING id;

-- name: SetTaskStartDate :exec
UPDATE tasks SET start_dt = $1 WHERE id = $2;

-- name: SetTaskEndDate :exec
UPDATE tasks SET end_dt = $1 WHERE id = $2;

-- name: GetOrderedTasksByUserID :many
SELECT *, CAST(EXTRACT(HOUR from end_dt - start_dt) AS INT) AS hours, 
    CAST(EXTRACT(MINUTE from end_dt - start_dt) AS INT) as minutes  FROM tasks 
WHERE user_id = $1 AND 
start_dt >= $2 AND
end_dt <=  $3 
GROUP BY id ORDER BY hours DESC;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1;
