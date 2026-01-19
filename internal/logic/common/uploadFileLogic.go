// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/coreKit/helpers/fileHelper"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件上传
func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, r *http.Request) (*types.UploadFileResp, error) {
	file, fileHeader, err := r.FormFile("file") // 读取入参 file 字段的上传文件信息
	if err != nil {
		return nil, errcode.ErrorUploadFileFail.WithError(errors.Wrap(err, "获取上传文件失败"))
	}

	userId := ctxData.GetUIDFromCtx(l.ctx)
	if userId <= 0 {
		return nil, errcode.Unauthorized.Msgr("用户未登录")
	}

	savePath := filepath.Join(l.svcCtx.Config.UploadFile.SavePath, req.FileType)
	localStorageIns := fileHelper.NewLocalStorage(
		savePath,
		file,
		fileHeader,
		fileHelper.WithMaxSize(l.svcCtx.Config.UploadFile.MaxSize),
	)
	// 检查文件类型
	if !localStorageIns.CheckFileExts(fileHeader, l.svcCtx.Config.UploadFile.AllowExts) {
		return nil, errcode.BadRequest.Msgr("不支持的文件类型")
	}
	// 保存文件
	dst, err := localStorageIns.Save()
	if err != nil {
		return nil, errcode.ErrorUploadFileFail.Msgr("文件保存失败").WithError(errors.Wrap(err, "保存文件失败"))
	}

	fileSize, err := fileHelper.FileContentSize(file)
	if err != nil {
		return nil, errcode.InternalServerError.Msgr("获取文件大小失败").WithError(errors.Wrap(err, "获取文件大小失败"))
	}

	// 网址路径
	dst = filepath.Join("/", dst)
	absoluteUrl := l.svcCtx.Config.AppService.StaticFSRelativePath + dst

	now := time.Now().Unix()
	uploadFile := &model.UploadFiles{
		UserId:       cast.ToUint64(userId),
		RelativePath: dst,
		AbsoluteUrl:  absoluteUrl,
		FileSize:     cast.ToUint64(fileSize),
		FileType:     req.FileType,
		BusinessType: req.BusinessType,
		UploadTime:   cast.ToUint64(now),
		CreatedAt:    cast.ToUint64(now),
	}

	result, err := l.svcCtx.UploadFilesModel.Insert(l.ctx, uploadFile)
	if err != nil {
		return nil, errcode.DBError.Msgr("保存文件记录失败").WithError(errors.Wrap(err, "保存文件记录失败"))
	}

	fileId, err := result.LastInsertId()
	if err != nil {
		return nil, errcode.DBError.Msgr("获取文件ID失败").WithError(errors.Wrap(err, "获取文件ID失败"))
	}

	var resp types.UploadFileResp
	if err := copier.Copy(&resp, uploadFile); err != nil {
		return nil, errcode.Fail.Msgr("数据转换失败").WithError(errors.Wrap(err, "数据转换失败"))
	}
	resp.FileId = cast.ToInt64(fileId)
	resp.UserId = cast.ToString(userId)
	resp.FileSize = cast.ToInt64(fileSize)

	return &resp, nil
}
