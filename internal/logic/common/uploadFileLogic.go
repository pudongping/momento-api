// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package common

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pudongping/momento-api/coreKit/ctxData"
	"github.com/pudongping/momento-api/coreKit/errcode"
	"github.com/pudongping/momento-api/internal/svc"
	"github.com/pudongping/momento-api/internal/types"
	"github.com/pudongping/momento-api/model"

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

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, r *http.Request) (resp *types.UploadFileResp, err error) {
	file, fileHeader, err := r.FormFile("file")
	// 获取用户ID
	userId := ctxData.GetUIDFromCtx(l.ctx)
	if userId == 0 {
		return nil, errcode.Unauthorized.Msgr("用户未登录")
	}

	// 验证文件大小（最大 10MB）
	maxFileSize := int64(10 * 1024 * 1024)
	if fileHeader.Size > maxFileSize {
		return nil, errcode.BadRequest.Msgf("文件大小不能超过 %dMB", maxFileSize/1024/1024)
	}

	// 验证文件类型
	allowedTypes := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"pdf":  true,
	}

	fileExt := filepath.Ext(fileHeader.Filename)
	if len(fileExt) > 0 {
		fileExt = fileExt[1:] // 移除 .
	}

	if !allowedTypes[fileExt] {
		return nil, errcode.BadRequest.Msgr("不支持的文件类型")
	}

	// 生成文件保存路径
	now := time.Now()
	uploadDir := fmt.Sprintf("public/uploads/%d/%d/%d", now.Year(), now.Month(), now.Day())
	os.MkdirAll(uploadDir, 0755)

	// 生成唯一文件名
	newFileName := fmt.Sprintf("%d_%s_%s", userId, strconv.FormatInt(time.Now().UnixNano(), 10), fileHeader.Filename)
	filePath := filepath.Join(uploadDir, newFileName)

	// 读取文件大小
	fileContent, err := io.ReadAll(file)
	if err != nil {
		l.Errorf("读取文件失败: %v", err)
		return nil, errcode.InternalServerError.Msgr("文件读取失败")
	}
	fileSize := int64(len(fileContent))

	// 保存文件到本地
	err = os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		l.Errorf("保存文件失败: %v", err)
		return nil, errcode.InternalServerError.Msgr("文件保存失败")
	}

	// 生成相对路径和绝对 URL
	relativePath := fmt.Sprintf("/%s/%s", uploadDir, newFileName)
	absoluteUrl := fmt.Sprintf("http://%s:%d%s", l.svcCtx.Config.Host, l.svcCtx.Config.Port, relativePath)

	// 保存文件记录到数据库
	uploadTime := time.Now().Unix()
	uploadFile := &model.UploadFiles{
		UserId:       uint64(userId),
		RelativePath: relativePath,
		AbsoluteUrl:  absoluteUrl,
		FileSize:     uint64(fileSize),
		FileType:     fileExt,
		BusinessType: req.BusinessType,
		UploadTime:   uint64(uploadTime),
		CreatedAt:    uint64(uploadTime),
	}

	result, err := l.svcCtx.UploadFilesModel.Insert(l.ctx, uploadFile)
	if err != nil {
		l.Errorf("保存文件记录失败: %v", err)
		return nil, errcode.DBError.Msgr("保存文件记录失败")
	}

	fileId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取文件ID失败: %v", err)
		return nil, errcode.InternalServerError.Msgr("获取文件ID失败")
	}

	resp = &types.UploadFileResp{
		FileId:       fileId,
		UserId:       strconv.FormatInt(userId, 10),
		RelativePath: relativePath,
		AbsoluteUrl:  absoluteUrl,
		FileSize:     fileSize,
		FileType:     fileExt,
		BusinessType: req.BusinessType,
		UploadTime:   uploadTime,
	}

	return resp, nil
}
