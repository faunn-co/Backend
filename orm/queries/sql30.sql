SELECT
    *
FROM
    user_table u,
    user_auth_table a
WHERE
        u.user_name = ?
  AND u.user_id = a.user_id
  AND a.user_password = ?
