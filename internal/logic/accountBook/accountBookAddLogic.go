package accountBook

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AccountBookAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建账本
func NewAccountBookAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookAddLogic {
	return &AccountBookAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookAddLogic) AccountBookAdd(req *types.AccountBookAddReq) (resp *types.AccountBookAddResp, err error) {
	userId := ctxData.GetUIDFromCtx(l.ctx)
	userIdUint := cast.ToUint64(userId)

	// Check book count limit (max 12)
	countBuilder := l.svcCtx.AccountBookMembersModel.CountBuilder("book_id").
		Where(squirrel.Eq{"user_id": userId})
	count, err := l.svcCtx.AccountBookMembersModel.FindCount(l.ctx, countBuilder)
	if err != nil {
		return nil, errcode.DBError.WithError(err).Msgr("查询账本数量失败")
	}

	if count >= 12 {
		return nil, errcode.Fail.Msgr("每个用户最多只能创建或加入12个账本")
	}

	isDefault := int64(2)
	if count == 0 {
		isDefault = 1
	}

	var bookId int64
	now := time.Now().Unix()

	err = l.svcCtx.AccountBooksModel.Transaction(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. Create Account Book
		book := &model.AccountBooks{
			Name:          req.Name,
			CreatorUserId: userIdUint,
			MemberCount:   1,
			CreatedAt:     uint64(now),
			UpdatedAt:     uint64(now),
		}

		res, err := l.svcCtx.AccountBooksModel.WithSession(session).Insert(ctx, book)
		if err != nil {
			return err
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		bookId = lastId

		// 2. Add Creator as Member
		member := &model.AccountBookMembers{
			BookId:    uint64(bookId),
			UserId:    userIdUint,
			IsCreator: 1, // Yes
			IsDefault: isDefault,
			Status:    "joined",
			JoinedAt:  uint64(now),
			CreatedAt: uint64(now),
			UpdatedAt: uint64(now),
		}

		_, err = l.svcCtx.AccountBookMembersModel.WithSession(session).Insert(ctx, member)
		return err
	})

	if err != nil {
		return nil, errcode.DBError.WithError(err).Msgr("创建账本失败")
	}

	return &types.AccountBookAddResp{
		BookId:        bookId,
		Name:          req.Name,
		CreatorUserId: userId,
		MemberCount:   1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
