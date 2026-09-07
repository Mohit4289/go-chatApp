-- name: CreateUser :one
INSERT INTO public."user" (
    name,
    email,
    password
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id,
    name,
    email,
    password,
    refresh_token,
    created_at;


-- name: GetUserByID :one
SELECT
    id,
    name,
    email,
    password,
    refresh_token,
    created_at
FROM public."user"
WHERE id = $1;


-- name: GetUserByEmail :one
SELECT
    id,
    name,
    email,
    password,
    refresh_token,
    created_at
FROM public."user"
WHERE email = $1;


-- name: AddRefreshToken :one
UPDATE public."user"
SET refresh_token = $1
WHERE email = $2
RETURNING
    id,
    name,
    email,
    password,
    refresh_token,
    created_at;