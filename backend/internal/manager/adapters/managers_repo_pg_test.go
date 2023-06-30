package adapters

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type ManagersRepositoryPostgresSuite struct {
	suite.Suite
	ctx  context.Context
	pg   *pgxpool.Pool
	repo *ManagersRepositoryPostgres
}

func (s *ManagersRepositoryPostgresSuite) SetupSuite() {
	var (
		err      error
		ctx      = context.Background()
		psqlUser = getenv("POSTGRES_USER", "postgres")
		psqlPass = getenv("POSTGRES_PASS", "postgres")
		psqlDB   = getenv("POSTGRES_DB", "qaztrade_test")
		psqlURL  = fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s", psqlUser, psqlPass, psqlDB)
	)

	s.pg, err = pgxpool.Connect(ctx, psqlURL)
	if err != nil {
		panic(err)
	}

	if err := teardownManagersRepositoryPostgres(ctx, s.pg); err != nil {
		fmt.Println(err)
	}

	s.ctx = ctx
	s.repo = NewManagersRepositoryPostgres(s.pg)
}

func (s *ManagersRepositoryPostgresSuite) TearDownSuite() {
	s.pg.Close()
}

func (s *ManagersRepositoryPostgresSuite) TearDownTest() {
	if err := teardownManagersRepositoryPostgres(s.ctx, s.pg); err != nil {
		fmt.Println(err)
	}
}

func TestManagersRepositoryPostgres(t *testing.T) {
	suite.Run(t, new(ManagersRepositoryPostgresSuite))
}

func applyFixture(ctx context.Context, pg *pgxpool.Pool, fixture func() (string, []string)) error {
	sql, args := fixture()
	sqlStr := fmt.Sprintf(sql, toAny(args)...)

	return performSQL(ctx, pg, sqlStr)
}

func FixtureTestGetMany() (string, []string) {
	ids := []string{
		"85925046-290b-4159-9ebb-89c6cdb00fde", // 1
		"6ae4ead1-1865-44af-a74f-84f612df6f87", // 2
		"6d62c110-9136-4a3b-bcd5-839b9cb5fb2a", // 3
	}

	sqlQuery := `
	insert into users 
		(id, email, hashed_password, attrs, created_at)
	values
		('%[1]s',	'1',	'',	'{"name": "John Doe"}',		'2021-01-01'),
		('%[2]s',	'2',	'',	'{"name": "Jack Wolf"}',	'2021-01-02'),
		('%[3]s',	'3',	'',	'{"name": "Johny Dep"}',	'2021-01-03')
	;

	insert into user_role_bindings
		(user_id, role_id)
	values
		('%[1]s', (select id from user_roles where value = 'manager')),
		('%[1]s', (select id from user_roles where value = 'admin')),
		('%[1]s', (select id from user_roles where value = 'legal')),
		('%[2]s', (select id from user_roles where value = 'manager')),
		('%[3]s', (select id from user_roles where value = 'user'))
	;
	`

	return sqlQuery, ids
}

func (s *ManagersRepositoryPostgresSuite) TestGetMany() {
	s.Require().Nil(applyFixture(s.ctx, s.pg, FixtureTestGetMany))

	var (
		expManagers = []*domain.Manager{
			{
				UserID: "85925046-290b-4159-9ebb-89c6cdb00fde",
				Email:  "1",
				Roles:  []string{"manager", "admin", "legal"},
			},
			{
				UserID: "6ae4ead1-1865-44af-a74f-84f612df6f87",
				Email:  "2",
				Roles:  []string{"manager"},
			},
		}
	)

	actManagers, err := s.repo.GetMany(s.ctx)
	s.Require().Nil(err)
	s.Require().Equal(expManagers, actManagers)
}

func teardownManagersRepositoryPostgres(ctx context.Context, pg *pgxpool.Pool) error {
	sqlQuery := `
		truncate user_role_bindings cascade;
		truncate users cascade;
	`

	return performSQL(
		ctx,
		pg,
		sqlQuery,
	)
}

func toAny(ids []string) []any {
	tmp := make([]interface{}, len(ids))
	for i, val := range ids {
		tmp[i] = val
	}
	return tmp
}

//nolint:errcheck
func performSQL(ctx context.Context, pg *pgxpool.Pool, query string, args ...interface{}) error {
	tx, err := pg.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, q := range strings.Split(query, ";") {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(ctx, q, args...); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
