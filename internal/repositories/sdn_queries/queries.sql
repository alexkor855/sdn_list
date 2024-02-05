-- name: InsertSdn :one
INSERT INTO sdn_list (uid, first_name, last_name) VALUES ($1, $2, $3) RETURNING id;

-- name: GetSdnByUid :many
SELECT * FROM sdn_list WHERE uid = $1;

-- name: GetSdnById :one
SELECT * FROM sdn_list WHERE id = $1;

-- name: GetSdnByUidAndName :one
SELECT * FROM sdn_list WHERE uid = $1 AND first_name = $2 AND last_name = $3;

-- name: DeleteOrder :exec
DELETE FROM sdn_list WHERE id = $1;



-- STRONG
-- это точное совпадение имени и / ИЛИ фамилии
-- перемещаем делитель F-L 0-3 1-2 2-1 3-0
-- перемещаем делитель L-F 0-3 1-2 2-1 3-0
-- a b c
-- 
SELECT id, uid, first_name, last_name
FROM sdn_list 
WHERE 
    (last_name = '$part1' OR last_name = '$part1 + $part2' OR last_name = '$part2 + $part1')
    OR
    (first_name)



-- WEAK
-- должно найти любое совпадение в имени либо фамилии
-- то есть разбить на части и каждую часть поискать в Имени и Фамилии
-- проверяем каждое вхождение в first_name или last_name

WITH found_values AS (
SELECT id, uid, first_name, last_name, 1 as sort_pos
FROM sdn_list 
WHERE 
    (first_name LIKE '%$part1%' OR last_name LIKE '%$part1%')
    OR
    (first_name LIKE '%$part2%' OR last_name LIKE '%$part2%')
    OR
    (first_name LIKE '%$part3%' OR last_name LIKE '%$part3%')
)
SELECT * FROM found_values
UNION
SELECT id, uid, first_name, last_name, 2 as sort_pos
FROM sdn_list sl
INNER JOIN found_values fv ON fv.uid = sl.uid AND fv.id <> sl.id
ORDER BY sl.uid, sort_pos