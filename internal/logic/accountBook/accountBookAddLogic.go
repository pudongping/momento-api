package accountBook

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
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
	userId := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

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

	isDefault := model.AccountBookMemberIsDefaultNo
	if count == 0 {
		isDefault = model.AccountBookMemberIsCreatorYes
	}

	var bookId int64
	now := time.Now().Unix()

	err = l.svcCtx.AccountBooksModel.Transaction(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. Create Account Book
		insertBookMap := map[string]interface{}{
			"name":            req.Name,
			"creator_user_id": userId,
			"member_count":    1,
			"created_at":      now,
			"updated_at":      now,
		}

		res, err := l.svcCtx.AccountBooksModel.InsertWithSession(ctx, session, insertBookMap)
		if err != nil {
			return errors.Wrap(err, "创建账本失败")
		}

		bookId, err := res.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "获取账本ID失败")
		}

		// 2. Add Creator as Member
		insertMemberMap := map[string]interface{}{
			"book_id":    bookId,
			"user_id":    userId,
			"is_creator": model.AccountBookMemberIsCreatorYes,
			"is_default": isDefault,
			"status":     model.AccountBookMemberStatusJoined,
			"joined_at":  now,
			"created_at": now,
			"updated_at": now,
		}

		_, err = l.svcCtx.AccountBookMembersModel.InsertWithSession(ctx, session, insertMemberMap)
		if err != nil {
			return errors.Wrap(err, "添加账本成员失败")
		}
		return nil
	})

	if err != nil {
		return nil, errcode.DBError.WithError(err).Msgr("创建账本失败")
	}

	return &types.AccountBookAddResp{
		BookId:        bookId,
		Name:          req.Name,
		CreatorUserId: cast.ToString(userId),
		MemberCount:   1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
