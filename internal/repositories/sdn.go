package repositories

import (
	"context"
	"errors"
	"fmt"
	"sdn_list/internal/entities"
	"sdn_list/internal/repositories/sdn_queries"
	"sdn_list/internal/services"
	workerpool "sdn_list/pkg/sync_patterns/worker_pool"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SdnRepository struct {
	dbPool *pgxpool.Pool
}

func NewSdnRepository(dbPool *pgxpool.Pool) *SdnRepository {
	return &SdnRepository{dbPool: dbPool}
}

func (r *SdnRepository) SaveAll(ctx context.Context, input <-chan services.Person) error {
	pool := workerpool.NewWorkerPoolForChan[services.Person, int](5, input)

	result := pool.Run(ctx, r.Save)
	for res := range result {
		if res.Err != nil {
			pool.Stop()
			return res.Err
		}
	}
	return nil
}

func (r *SdnRepository) Save(ctx context.Context, person services.Person) (int, error) {
	queries := sdn_queries.New(r.dbPool)

	firstName := pgtype.Text{}
	firstName.Scan(person.FirstName)

	lastName := pgtype.Text{}
	lastName.Scan(person.LastName)

	sdnList, err := queries.GetSdnByUidAndName(ctx, sdn_queries.GetSdnByUidAndNameParams{
		Uid:       int64(person.Uid),
		FirstName: firstName,
		LastName:  lastName,
	})

	if err == nil {
		return int(sdnList.ID), nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}

	id, err := queries.InsertSdn(ctx, sdn_queries.InsertSdnParams{
		Uid:       int64(person.Uid),
		FirstName: firstName,
		LastName:  lastName,
	})

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

type Person struct {
	Id        int64  `db:"id"`
	Uid       int64  `db:"uid"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func (r *SdnRepository) SearchStrong(ctx context.Context, name string) ([]entities.Person, error) {
	where, values := r.getWhere(name)
	sql := r.getSqlForStrong(where)
	fmt.Println(sql)
	fmt.Println(strings.Join(values, ", "))

	var interfaceValues []interface{}
	for _, v := range values {
		interfaceValues = append(interfaceValues, v)
	}

	rows, _ := r.dbPool.Query(ctx, sql, interfaceValues...)

	persons, err := pgx.CollectRows(rows, pgx.RowToStructByName[Person])
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(persons), nil
}

func (r *SdnRepository) SearchWeak(ctx context.Context, name string) ([]entities.Person, error) {
	where, values := r.getWhere(name)
	sql := r.getSqlForWeak(where)
	fmt.Println(sql)
	fmt.Println(strings.Join(values, ", "))

	var interfaceValues []interface{}
	for _, v := range values {
		interfaceValues = append(interfaceValues, v)
	}

	rows, _ := r.dbPool.Query(ctx, sql, interfaceValues...)

	persons, err := pgx.CollectRows(rows, pgx.RowToStructByName[Person])
	if err != nil {
		return nil, err
	}

	return r.convertToEntity(persons), nil
}

func (r *SdnRepository) convertToEntity(persons []Person) []entities.Person {
	result := make([]entities.Person, len(persons))

	for i, person := range persons {
		result[i] = entities.Person{
			Id:        int(person.Id),
			Uid:       int(person.Uid),
			FirstName: person.FirstName,
			LastName:  strings.Title(strings.ToLower(person.LastName)),
		}
	}
	return result
}

func (r *SdnRepository) parseName(fullName string) []string {
	fullName = strings.Trim(fullName, " ")
	result := []string{}
	for _, part := range strings.Split(fullName, " ") {
		result = append(result, strings.ToUpper(part))
	}
	return result
}

func (r *SdnRepository) getNamesCombinations(fullName string) []FirstLastName {
	nameParts := r.parseName(fullName)
	result := make([]FirstLastName, 0)

	for index := 0; index <= len(nameParts); index++ {
		result = append(result, FirstLastName{
			FirstName: strings.Join(nameParts[:index], " "),
			LastName:  strings.Join(nameParts[index:], " "),
		})

		if index != 0 && index != len(nameParts) {
			result = append(result, FirstLastName{
				FirstName: strings.Join(nameParts[index:], " "),
				LastName:  strings.Join(nameParts[:index], " "),
			})
		}
	}

	return result
}

func (r *SdnRepository) getWhere(fullName string) (string, []string) {
	values := []string{}
	conditions := []string{}
	allNames := r.getNamesCombinations(fullName)

	for _, name := range allNames {
		if name.FirstName == "" {
			values = append(values, name.LastName)
			conditions = append(conditions, fmt.Sprintf("UPPER(last_name) = $%d", len(values)))
			continue
		}
		if name.LastName == "" {
			values = append(values, name.FirstName)
			conditions = append(conditions, fmt.Sprintf("UPPER(first_name) = $%d", len(values)))
			continue
		}
		values = append(values, name.FirstName, name.LastName)
		conditions = append(conditions, fmt.Sprintf("(UPPER(first_name) = $%d AND UPPER(last_name) = $%d)", len(values)-1, len(values)))
	}

	return strings.Join(conditions, " OR "), values
}

func (r *SdnRepository) getSqlForStrong(where string) string {
	return fmt.Sprintf("SELECT id, uid, first_name, last_name FROM sdn_list WHERE %s;", where)
}

func (r *SdnRepository) getSqlForWeak(where string) string {

	return fmt.Sprintf(`
		WITH found_values AS (
			SELECT DISTINCT uid
			FROM sdn_list 
			WHERE 
				%s
		)
		SELECT id, uid, first_name, last_name
		FROM sdn_list sl
		WHERE
			sl.uid IN (SELECT uid FROM found_values)
		ORDER BY uid, id;
	`, where)
}

type FirstLastName struct {
	FirstName string
	LastName  string
}
