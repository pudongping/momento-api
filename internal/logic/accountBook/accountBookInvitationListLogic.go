package accountBook

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/spf13/cast"

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
		l.Logger.Errorf("AccountBookInvitationList FindAll error: %v", err)
		return nil, errcode.DBError.Msgr("获取邀请列表失败")
	}

	resp = make([]types.AccountBookInvitationListResp, 0)
	if len(invitations) == 0 {
		return resp, nil
	}

	// 2. 收集账本ID和邀请人ID
	bookIDs := make([]uint64, 0, len(invitations))
	inviterUIDs := make([]uint64, 0, len(invitations))
	for _, inv := range invitations {
		bookIDs = append(bookIDs, inv.BookId)
		inviterUIDs = append(inviterUIDs, inv.InviterUid)
	}

	// 3. 批量查询账本信息
	// 注意：go-zero生成的model通常没有批量查询，这里使用SelectBuilder构建 IN 查询
	var booksMap = make(map[uint64]string)
	if len(bookIDs) > 0 {
		bookBuilder := l.svcCtx.AccountBooksModel.SelectBuilder().Where(squirrel.Eq{"book_id": bookIDs})
		books, err := l.svcCtx.AccountBooksModel.FindAll(l.ctx, bookBuilder)
		if err != nil {
			l.Logger.Errorf("AccountBookInvitationList FindAll Books error: %v", err)
			// 不阻断，继续
		} else {
			for _, book := range books {
				booksMap[book.BookId] = book.Name
			}
		}
	}

	// 4. 批量查询邀请人信息
	var usersMap = make(map[uint64]string)
	if len(inviterUIDs) > 0 {
		userBuilder := l.svcCtx.UserModel.SelectBuilder().Where(squirrel.Eq{"user_id": inviterUIDs})
		users, err := l.svcCtx.UserModel.FindAll(l.ctx, userBuilder)
		if err != nil {
			l.Logger.Errorf("AccountBookInvitationList FindAll Users error: %v", err)
			// 不阻断，继续
		} else {
			for _, user := range users {
				// 优先使用昵称，没有则使用默认
				usersMap[user.UserId] = user.Nickname
			}
		}
	}

	// 5. 组装数据
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
