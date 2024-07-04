-- name: GetPersonByID :one
SELECT * FROM people
WHERE id = $1;

-- name: GetPersonByPassport :one
SELECT * FROM people
WHERE passport_number = $1 AND passport_serie = $2;

-- name: ListPeopleWithLimit :many
SELECT * FROM people
WHERE
    (sqlc.arg(passport_serie)::int = 0 OR passport_serie = sqlc.arg(passport_serie)) AND
    (sqlc.arg(passport_number)::int = 0 OR passport_number = sqlc.arg(passport_number)) AND
    (sqlc.arg(surname)::text = '' OR surname ILIKE '%' || sqlc.arg(surname) || '%') AND
    (sqlc.arg(name)::text = '' OR name ILIKE '%' || sqlc.arg(name) || '%') AND
    (sqlc.arg(patronymic)::text = '' OR patronymic ILIKE '%' || sqlc.arg(patronymic) || '%') AND
    (sqlc.arg(address)::text = '' OR address ILIKE '%' || sqlc.arg(address) || '%')
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListPeople :many
SELECT * FROM people
WHERE
    (sqlc.arg(passport_serie)::int = 0 OR passport_serie = sqlc.arg(passport_serie)) AND
    (sqlc.arg(passport_number)::int = 0 OR passport_number = sqlc.arg(passport_number)) AND
    (sqlc.arg(surname)::text = '' OR surname ILIKE '%' || sqlc.arg(surname) || '%') AND
    (sqlc.arg(name)::text = '' OR name ILIKE '%' || sqlc.arg(name) || '%') AND
    (sqlc.arg(patronymic)::text = '' OR patronymic ILIKE '%' || sqlc.arg(patronymic) || '%') AND
    (sqlc.arg(address)::text = '' OR address ILIKE '%' || sqlc.arg(address) || '%')
ORDER BY id;

-- name: CreatePerson :one
INSERT INTO people (name, surname, patronymic, address, passport_number, passport_serie) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdatePerson :exec
UPDATE people
SET 
    name = COALESCE(NULLIF(sqlc.arg(name), ''), name),
    surname = COALESCE(NULLIF(sqlc.arg(surname), ''), surname),
    patronymic = COALESCE(NULLIF(sqlc.arg(patronymic), ''), patronymic),
    address = COALESCE(NULLIF(sqlc.arg(address), ''), address),
    passport_serie = COALESCE(NULLIF(sqlc.arg(passport_serie), 0), passport_serie),
    passport_number = COALESCE(NULLIF(sqlc.arg(passport_number), 0), passport_number)
WHERE id = $1;

-- name: DeletePerson :exec
DELETE FROM people WHERE id = $1;
