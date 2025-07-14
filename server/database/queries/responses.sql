-- name: ListResponses :many
select * from list_responses_for_form(
    sqlc.arg(form_id),
    sqlc.arg(user_id),
    sqlc.arg(status)::response_status,
    sqlc.arg(sort_by),
    sqlc.arg(order_by),
    sqlc.arg(limit_val),
    sqlc.arg(offset_val)
);

-- name: CountResponses :one
select count_responses_for_form(
    sqlc.arg(form_id),
    sqlc.arg(user_id),
    sqlc.arg(status)::response_status
);

-- name: StartResponse :one
select * from start_response_for_form(
    sqlc.arg(form_id),
    sqlc.arg(user_id)
);

-- name: GetResponse :one
select * from get_response_by_id(
    sqlc.arg(id),
    sqlc.arg(form_id),
    sqlc.arg(user_id)
);

-- name: GetAnswers :many
select * from get_answers_for_response(
    sqlc.arg(id),
    sqlc.arg(form_id),
    sqlc.arg(user_id)
);

-- name: SaveAnswer :one
select * from add_answer_to_response(
    sqlc.arg(id), sqlc.arg(form_id), sqlc.arg(user_id),
    sqlc.arg(question), sqlc.arg(value)
);

-- name: SubmitResponse :one
select * from submit_response_by_id(
    sqlc.arg(id),
    sqlc.arg(form_id),
    sqlc.arg(user_id)
);
