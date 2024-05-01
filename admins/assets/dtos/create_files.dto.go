package dtos

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/dangduoc08/ecommerce-api/admins/assets/commons"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/ctx"
	"github.com/dangduoc08/gogo/exception"
)

type CREATE_files_Body_DTO struct {
	CommonProvider commons.CommonProvider
	File           *ctx.DataFile `bind:"file"`
	Dir            string
}

func (instance CREATE_files_Body_DTO) Transform(c gogo.Context, medata common.ArgumentMetadata) any {
	currentDir, _ := os.Getwd()
	dir := instance.CommonProvider.CleanDir(c.Form().Get("dir"))
	instance.Dir = filepath.Join(currentDir, "publics", dir)

	return c.File().Bind(instance)
}

func (instance CREATE_files_Body_DTO) IsValid(dataFile *ctx.DataFile) bool {
	return strings.HasPrefix(dataFile.Type, "image")
}

func (instance CREATE_files_Body_DTO) Store(dataFile *ctx.DataFile, src multipart.File) {
	uploadPath := filepath.Join(instance.Dir, dataFile.Filename)
	uploadPath = instance.CommonProvider.GeneratePath(uploadPath, 1)

	dst, err := os.Create(uploadPath)
	if err != nil {
		dst.Close()
		panic(exception.BadRequestException(err.Error()))
	}
	_, err = io.Copy(dst, src)

	if err != nil {
		dst.Close()
		panic(exception.BadRequestException(err.Error()))
	}

	dst.Close()
	dataFile.Dest = dst.Name()
}
