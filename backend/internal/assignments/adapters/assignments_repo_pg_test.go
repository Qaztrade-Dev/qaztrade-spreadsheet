package adapters

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type AssignmentsRepositoryPostgresSuite struct {
	suite.Suite
	ctx  context.Context
	pg   *pgxpool.Pool
	repo *AssignmentsRepositoryPostgres
}

func (s *AssignmentsRepositoryPostgresSuite) SetupSuite() {
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

	if err := teardownAssignmentsRepositoryPostgres(ctx, s.pg); err != nil {
		fmt.Println(err)
	}

	s.ctx = ctx
	s.repo = NewAssignmentsRepositoryPostgres(s.pg)
}

func (s *AssignmentsRepositoryPostgresSuite) TearDownSuite() {
	s.pg.Close()
}

func (s *AssignmentsRepositoryPostgresSuite) TearDownTest() {
	if err := teardownAssignmentsRepositoryPostgres(s.ctx, s.pg); err != nil {
		fmt.Println(err)
	}
}

func TestAssignmentsRepositoryPostgres(t *testing.T) {
	suite.Run(t, new(AssignmentsRepositoryPostgresSuite))
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
	}

	sqlQuery := `
	insert into users 
		(id, email, hashed_password, attrs)
	values
		('%[1]s',	'1',	'',	'{"full_name": "John Doe"}'),
		('%[2]s',	'2',	'',	'{"full_name": "Jack Wolf"}')
	;

	insert into applications 
		(id, attrs, sign_link)
	values
		(
			'%[1]s',
			'{
				"application": {
					"from": "Facebook Inc.",
					"bin":	"012345678901"
				}
			}'::jsonb,
			'abc'
		),
		(
			'%[2]s',
			'{
				"application": {
					"from": "Google Ent.",
					"bin":	"098765432109"
				}
			}'::jsonb,
			'def'
		)
	;

	insert into "assignments" 
		(
			id, created_at,
			sheet_title, sheet_id, 
			total_rows, total_sum, 
			is_completed, completed_at,
			user_id, application_id
		)
	values
		(
			1, '2023-01-01T13:30:00',
			'Sheet1', 1,
			50, 150,
			false, null,
			'%[1]s', '%[1]s'
		),
		(
			2, '2023-01-01T13:31:00',
			'Sheet2', 2,
			1, 100,
			true, '2023-01-02T13:30:00',
			'%[1]s', '%[2]s'
		),
		(
			3, '2023-01-01T13:32:00',
			'Sheet1', 1,
			1, 49,
			false, null,
			'%[2]s', '%[1]s'
		)
	;

	insert into "assignment_results" 
		(assignment_id, total_completed)
	values
		(2, 100),
		(3, 14)
	;

	update "assignments" as ass
	set
		last_result_id = assres.id
	from "assignment_results" as assres
	where 
		assres.assignment_id = ass.id
	;
	`

	return sqlQuery, ids
}

func (s *AssignmentsRepositoryPostgresSuite) TestGetMany_All() {
	s.Require().Nil(applyFixture(s.ctx, s.pg, FixtureTestGetMany))

	var (
		expAssignmentsTotal  = 3
		expAssignmentObjects = []*domain.AssignmentView{
			{
				ID:             1,
				ApplicantName:  "Facebook Inc.",
				ApplicantBIN:   "012345678901",
				SheetTitle:     "Sheet1",
				SheetID:        1,
				AssignmentType: "",
				Link:           "#gid=1",
				SignLink:       "https://link.doodocs.kz/abc",
				AssigneeName:   "John Doe",
				TotalRows:      50,
				TotalSum:       150,
				RowsCompleted:  0,
				IsCompleted:    false,
				CompletedAt:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             2,
				ApplicantName:  "Google Ent.",
				ApplicantBIN:   "098765432109",
				SheetTitle:     "Sheet2",
				SheetID:        2,
				AssignmentType: "",
				Link:           "#gid=2",
				SignLink:       "https://link.doodocs.kz/def",
				AssigneeName:   "John Doe",
				TotalRows:      1,
				TotalSum:       100,
				RowsCompleted:  100,
				IsCompleted:    true,
				CompletedAt:    time.Date(2023, time.January, 2, 19, 30, 0, 0, time.Local),
			},
			{
				ID:             3,
				ApplicantName:  "Facebook Inc.",
				ApplicantBIN:   "012345678901",
				SheetTitle:     "Sheet1",
				SheetID:        1,
				AssignmentType: "",
				Link:           "#gid=1",
				SignLink:       "https://link.doodocs.kz/abc",
				AssigneeName:   "Jack Wolf",
				TotalRows:      1,
				TotalSum:       49,
				RowsCompleted:  14,
				IsCompleted:    false,
				CompletedAt:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		}
	)

	assignmentsList, err := s.repo.GetMany(s.ctx, &domain.GetManyInput{
		Limit:  10,
		Offset: 0,
	})
	s.Require().Nil(err)
	s.Require().Equal(assignmentsList.Total, expAssignmentsTotal)
	s.Require().Equal(assignmentsList.Objects, expAssignmentObjects)
}

func (s *AssignmentsRepositoryPostgresSuite) TestGetMany_User() {
	s.Require().Nil(applyFixture(s.ctx, s.pg, FixtureTestGetMany))

	var (
		userID               = "85925046-290b-4159-9ebb-89c6cdb00fde"
		expAssignmentsTotal  = 2
		expAssignmentObjects = []*domain.AssignmentView{
			{
				ID:             1,
				ApplicantName:  "Facebook Inc.",
				ApplicantBIN:   "012345678901",
				SheetTitle:     "Sheet1",
				SheetID:        1,
				AssignmentType: "",
				Link:           "#gid=1",
				SignLink:       "https://link.doodocs.kz/abc",
				AssigneeName:   "John Doe",
				TotalRows:      50,
				TotalSum:       150,
				RowsCompleted:  0,
				IsCompleted:    false,
				CompletedAt:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:             2,
				ApplicantName:  "Google Ent.",
				ApplicantBIN:   "098765432109",
				SheetTitle:     "Sheet2",
				SheetID:        2,
				AssignmentType: "",
				Link:           "#gid=2",
				SignLink:       "https://link.doodocs.kz/def",
				AssigneeName:   "John Doe",
				TotalRows:      1,
				TotalSum:       100,
				RowsCompleted:  100,
				IsCompleted:    true,
				CompletedAt:    time.Date(2023, time.January, 2, 19, 30, 0, 0, time.Local),
			},
		}
	)

	assignmentsList, err := s.repo.GetMany(s.ctx, &domain.GetManyInput{
		UserID: &userID,
		Limit:  10,
		Offset: 0,
	})
	s.Require().Nil(err)
	s.Require().Equal(assignmentsList.Total, expAssignmentsTotal)
	s.Require().Equal(assignmentsList.Objects, expAssignmentObjects)
}

// func (s *AssignmentsRepositoryPostgresSuite) TestGetSheets() {
// 	sheets, err := s.repo.GetSheets(s.ctx)
// 	s.Require().Nil(err)

// 	managersCount := 10
// 	pq := domain.DistributeAdvanced(managersCount, sheets)

// 	p := message.NewPrinter(language.English)
// 	// Print out manager assignments
// 	for i, manager := range pq.Managers {
// 		println("Manager ID: ", i)
// 		p.Printf("Total rows: %d\n", manager.TotalRows)
// 		p.Printf("Total sum: %f\n", manager.TotalSum)
// 		println("------------")
// 	}
// }

func (s *AssignmentsRepositoryPostgresSuite) TestGetInfo_All() {
	s.Require().Nil(applyFixture(s.ctx, s.pg, FixtureTestGetMany))

	var (
		expInfo = &domain.AssignmentsInfo{
			Total:     3,
			Completed: 1,
		}
	)

	assignmentsInfo, err := s.repo.GetInfo(s.ctx, &domain.GetInfoInput{})
	s.Require().Nil(err)
	s.Require().Equal(expInfo, assignmentsInfo)
}

func (s *AssignmentsRepositoryPostgresSuite) TestGetInfo_User() {
	s.Require().Nil(applyFixture(s.ctx, s.pg, FixtureTestGetMany))

	var (
		userID  = "85925046-290b-4159-9ebb-89c6cdb00fde"
		expInfo = &domain.AssignmentsInfo{
			Total:     2,
			Completed: 1,
		}
	)

	assignmentsInfo, err := s.repo.GetInfo(s.ctx, &domain.GetInfoInput{
		UserID: &userID,
	})
	s.Require().Nil(err)
	s.Require().Equal(expInfo, assignmentsInfo)
}

func (s *AssignmentsRepositoryPostgresSuite) TestChangeAssignee() {
	s.Require().Nil(applyFixture(s.ctx, s.pg, FixtureTestGetMany))

	var (
		newUserID           = "85925046-290b-4159-9ebb-89c6cdb00fde"
		assignmentID uint64 = 3

		// expAssignmentObjects = []*domain.AssignmentView{
		// 	{
		// 		ID:             1,
		// 		ApplicantName:  "Facebook Inc.",
		// 		ApplicantBIN:   "012345678901",
		// 		SheetTitle:     "Sheet1",
		// 		SheetID:        1,
		// 		AssignmentType: "",
		// 		Link:           "#gid=1",
		// 		SignLink:       "https://link.doodocs.kz/abc",
		// 		AssigneeName:   "John Doe",
		// 		TotalRows:      50,
		// 		TotalSum:       150,
		// 		RowsCompleted:  0,
		// 		IsCompleted:    false,
		// 		CompletedAt:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		// 	},
		// }
	)

	err := s.repo.ChangeAssignee(s.ctx, &domain.ChangeAssigneeInput{
		UserID:       newUserID,
		AssignmentID: assignmentID,
	})
	s.Require().Nil(err)

	assignmentsList, err := s.repo.GetMany(s.ctx, &domain.GetManyInput{
		UserID: &newUserID,
		Limit:  10,
		Offset: 0,
	})
	s.Require().Nil(err)
	s.Require().Equal(assignmentsList.Total, 3)
	s.Require().Equal(assignmentsList.Objects[2].AssigneeName, "John Doe")
}

func teardownAssignmentsRepositoryPostgres(ctx context.Context, pg *pgxpool.Pool) error {
	sqlQuery := `
		truncate users cascade;
		truncate applications cascade;
		truncate assignments cascade;
		truncate assignment_results cascade;
	`

	return performSQL(
		ctx,
		pg,
		sqlQuery,
	)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
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
