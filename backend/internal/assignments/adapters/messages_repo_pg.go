package adapters

import (
	"context"
	"time"

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
			"msg.created_at",
			"msg.attrs",
			"msg.user_id",
			"u.email",
			"u.attrs->>'full_name'",
			"msg.doodocs_is_signed",
			"msg.doodocs_signed_at",
		).
		From("assignment_messages msg").
		Join("users u on u.id = msg.user_id").
		OrderBy("msg.created_at asc")

	if input.DoodocsDocumentID != "" {
		mainStmt = mainStmt.Where("msg.doodocs_document_id = ?", input.DoodocsDocumentID)
	}

	if input.AssignmentID != 0 {
		mainStmt = mainStmt.Where("msg.assignment_id = ?", input.AssignmentID)
	}

	return mainStmt
}

func queryMessages(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.Message, error) {
	var (
		objects = make([]*domain.Message, 0)

		// scans
		tmpMessageID       *string
		tmpAssignmentID    *uint64
		tmpCreatedAt       *time.Time
		tmpAttrs           *map[string]interface{}
		tmpUserID          *string
		tmpEmail           *string
		tmpFullName        *string
		tmpDoodocsIsSigned *bool
		tmpDoodocsSignedAt *time.Time
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&tmpMessageID,
		&tmpAssignmentID,
		&tmpCreatedAt,
		&tmpAttrs,
		&tmpUserID,
		&tmpEmail,
		&tmpFullName,
		&tmpDoodocsIsSigned,
		&tmpDoodocsSignedAt,
	}, func(pgx.QueryFuncRow) error {
		objects = append(objects, &domain.Message{
			MessageID:       postgres.Value(tmpMessageID),
			AssignmentID:    postgres.Value(tmpAssignmentID),
			CreatedAt:       postgres.Value(tmpCreatedAt),
			Attrs:           postgres.Value(tmpAttrs),
			UserID:          postgres.Value(tmpUserID),
			Email:           postgres.Value(tmpEmail),
			FullName:        postgres.Value(tmpFullName),
			DoodocsIsSigned: postgres.Value(tmpDoodocsIsSigned),
			DoodocsSignedAt: postgres.Value(tmpDoodocsSignedAt),
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

func (r *MessagesRepositoryPostgres) GetMany(ctx context.Context, input *domain.GetMessageInput) ([]*domain.Message, error) {
	stmt := getMessagesQueryStatement(input)

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	return queryMessages(ctx, r.pg, sql, args...)
}
