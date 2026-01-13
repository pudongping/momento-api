package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UploadFilesModel = (*customUploadFilesModel)(nil)

type (
	// UploadFilesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUploadFilesModel.
	UploadFilesModel interface {
		uploadFilesModel
		withSession(session sqlx.Session) UploadFilesModel
	}

	customUploadFilesModel struct {
		*defaultUploadFilesModel
	}
)

// NewUploadFilesModel returns a model for the database table.
func NewUploadFilesModel(conn sqlx.SqlConn) UploadFilesModel {
	return &customUploadFilesModel{
		defaultUploadFilesModel: newUploadFilesModel(conn),
	}
}

func (m *customUploadFilesModel) withSession(session sqlx.Session) UploadFilesModel {
	return NewUploadFilesModel(sqlx.NewSqlConnFromSession(session))
}
