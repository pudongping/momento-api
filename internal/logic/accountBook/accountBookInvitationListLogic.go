package accountBook

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountBookInvitationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取邀请列表
func NewAccountBookInvitationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountBookInvitationListLogic {
	return &AccountBookInvitationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountBookInvitationListLogic) AccountBookInvitationList(req *types.AccountBookInvitationListReq) (resp []types.AccountBookInvitationListResp, err error) {
	userID := cast.ToUint64(ctxData.GetUIDFromCtx(l.ctx))

	// 1. 查询当前用户的邀请记录
	rowBuilder := l.svcCtx.AccountBookInvitationsModel.SelectBuilder().
		Where(squirrel.Eq{"target_uid": userID}).
		OrderBy("created_at DESC")

	invitations, err := l.svcCtx.AccountBookInvitationsModel.FindAll(l.ctx, rowBuilder)
	if err != nil {
		return nil, errcode.DBError.Msgr("获取邀请列表失败").WithError(errors.Wrap(err, "获取邀请列表失败"))
	}

	if len(invitations) == 0 {
		return nil, nil
	}

	// 2. 收集去重后的账本ID和邀请人ID
	bookIDs := make([]uint64, 0, len(invitations))
	inviterUIDs := make([]uint64, 0, len(invitations))
	bookIDMap := make(map[uint64]struct{})
	inviterUIDMap := make(map[uint64]struct{})

	for _, inv := range invitations {
		if _, ok := bookIDMap[inv.BookId]; !ok {
			bookIDs = append(bookIDs, inv.BookId)
			bookIDMap[inv.BookId] = struct{}{}
		}
		if _, ok := inviterUIDMap[inv.InviterUid]; !ok {
			inviterUIDs = append(inviterUIDs, inv.InviterUid)
			inviterUIDMap[inv.InviterUid] = struct{}{}
		}
	}

	// 3. 并行查询账本信息和邀请人信息
	var (
		booksMap = make(map[uint64]string)
		usersMap = make(map[uint64]string)
	)

	err = mr.Finish(func() error {
		booksMap = l.getBooksMap(bookIDs)
		return nil
	}, func() error {
		usersMap = l.getUsersMap(inviterUIDs)
		return nil
	})

	if err != nil {
		l.Logger.Errorf("AccountBookInvitationList mr.Finish error: %v", err)
		// 并发执行出错不应阻断主流程，因为是辅助信息
	}

	// 4. 组装数据
	resp = make([]types.AccountBookInvitationListResp, 0, len(invitations))
	for _, inv := range invitations {
		bookName := booksMap[inv.BookId]
		if bookName == "" {
			bookName = "未知账本" // 账本可能已被删除
		}

		inviterNickname := usersMap[inv.InviterUid]
		if inviterNickname == "" {
			inviterNickname = "未知用户" // 用户可能注销
		}

		resp = append(resp, types.AccountBookInvitationListResp{
			InvitationId:    int64(inv.InvitationId),
			BookId:          int64(inv.BookId),
			BookName:        bookName,
			InviterUid:      cast.ToString(inv.InviterUid),
			InviterNickname: inviterNickname,
			TargetUid:       cast.ToString(inv.TargetUid),
			Status:          inv.Status,
			CreatedAt:       int64(inv.CreatedAt),
		})
	}

	return resp, nil
}

func (l *AccountBookInvitationListLogic) getBooksMap(bookIDs []uint64) map[uint64]string {
	booksMap := make(map[uint64]string)
	if len(bookIDs) == 0 {
		return booksMap
	}

	bookBuilder := l.svcCtx.AccountBooksModel.SelectBuilder().Where(squirrel.Eq{"book_id": bookIDs})
	books, err := l.svcCtx.AccountBooksModel.FindAll(l.ctx, bookBuilder)
	if err != nil {
		l.Logger.Errorf("AccountBookInvitationList FindAll Books error: %v", err)
		return booksMap
	}

	for _, book := range books {
		booksMap[book.BookId] = book.Name
	}
	return booksMap
}

func (l *AccountBookInvitationListLogic) getUsersMap(userIDs []uint64) map[uint64]string {
	usersMap := make(map[uint64]string)
	if len(userIDs) == 0 {
		return usersMap
	}

	userBuilder := l.svcCtx.UserModel.SelectBuilder().Where(squirrel.Eq{"user_id": userIDs})
	users, err := l.svcCtx.UserModel.FindAll(l.ctx, userBuilder)
	if err != nil {
		l.Logger.Errorf("AccountBookInvitationList FindAll Users error: %v", err)
		return usersMap
	}

	for _, user := range users {
		usersMap[user.UserId] = user.Nickname
	}
	return usersMap
}
