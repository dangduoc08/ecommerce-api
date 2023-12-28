package dtos

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/ctx"
	"github.com/dangduoc08/gooh/exception"
)

type CREATE_files_Body struct {
	providers.HandleAsset
	File *ctx.DataFile `bind:"file"`
	Dir  string
}

func (self CREATE_files_Body) Transform(c gooh.Context, medata common.ArgumentMetadata) any {
	currentDir, _ := os.Getwd()
	dir := self.CleanDir(c.Form().Get("dir"))
	self.Dir = filepath.Join(currentDir, "public", dir)

	return c.File().Bind(self)
}

func (self CREATE_files_Body) IsValid(dataFile *ctx.DataFile) bool {
	return strings.HasPrefix(dataFile.Type, "image")
}

func (self CREATE_files_Body) Store(dataFile *ctx.DataFile, src multipart.File) {
	uploadPath := filepath.Join(self.Dir, dataFile.Filename)
	uploadPath = self.GeneratePath(uploadPath, 1)

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
