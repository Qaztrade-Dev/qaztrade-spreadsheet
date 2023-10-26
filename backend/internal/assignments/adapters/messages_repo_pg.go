package adapters

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/pkg/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mattermost/squirrel"
)

type MessagesRepositoryPostgres struct {
	pg *pgxpool.Pool
}

var _ domain.MessagesRepository = (*MessagesRepositoryPostgres)(nil)

func NewMessagesRepositoryPostgres(pg *pgxpool.Pool) *MessagesRepositoryPostgres {
	return &MessagesRepositoryPostgres{
		pg: pg,
	}
}

func (r *MessagesRepositoryPostgres) CreateMessage(ctx context.Context, input *domain.CreateMessageInput) error {
	setMap := squirrel.Eq{
		"user_id":       input.UserID,
		"assignment_id": input.AssignmentID,
		"attrs":         input.Attrs,
	}

	if input.DoodocsDocumentID != "" {
		setMap["doodocs_document_id"] = input.DoodocsDocumentID
	}

	sql, args, err := psql.Insert("assignment_messages").SetMap(setMap).ToSql()
	if err != nil {
		return err
	}

	if _, err := r.pg.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func getMessagesQueryStatement(input *domain.GetMessageInput) squirrel.SelectBuilder {
	mainStmt := psql.
		Select(
			"msg.id",
			"msg.assignment_id",
		).
		From("assignment_messages msg")

	if input.DoodocsDocumentID != "" {
		mainStmt = mainStmt.Where("msg.doodocs_document_id = ?", input.DoodocsDocumentID)
	}

	return mainStmt
}

func queryMessages(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.Message, error) {
	var (
		objects = make([]*domain.Message, 0)

		// scans
		tmpMessageID    *string
		tmpAssignmentID *uint64
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&tmpMessageID,
		&tmpAssignmentID,
	}, func(pgx.QueryFuncRow) error {
		objects = append(objects, &domain.Message{
			MessageID:    postgres.Value(tmpMessageID),
			AssignmentID: postgres.Value(tmpAssignmentID),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, err
}

func (r *MessagesRepositoryPostgres) GetOne(ctx context.Context, input *domain.GetMessageInput) (*domain.Message, error) {
	stmt := getMessagesQueryStatement(input)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	objects, err := queryMessages(ctx, r.pg, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, domain.ErrMessageNotFound
	}

	return objects[0], nil
}

func (r *MessagesRepositoryPostgres) UpdateMessage(ctx context.Context, input *domain.UpdateMessageInput) error {
	stmt := psql.Update("assignment_messages").SetMap(squirrel.Eq{
		"doodocs_signed_at": input.DoodocsSignedAt,
		"doodocs_is_signed": input.DoodocsIsSigned,
	}).Where("id = ?", input.MessageID)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := r.pg.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
