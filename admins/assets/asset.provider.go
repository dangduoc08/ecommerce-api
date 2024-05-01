package assets

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dangduoc08/ecommerce-api/admins/assets/commons"
	"github.com/dangduoc08/gogo/core"
)

type AssetProvider struct {
	MaxDepth       int
	PublicPath     string
	CommonProvider commons.CommonProvider
}

func (instance AssetProvider) NewProvider() core.Provider {
	instance.MaxDepth = 1

	return instance
}

func (instance AssetProvider) List(dirPath string) ([]*AssetModel, error) {
	assets := []*AssetModel{}
	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Count(path, string(os.PathSeparator))-strings.Count(dirPath, string(os.PathSeparator)) > instance.MaxDepth {
			return fs.SkipDir
		}

		if dirPath != path {
			fileInfo, _ := d.Info()
			ext := filepath.Ext(path)
			if ext != "" {
				ext = ext[1:]
			}

			asset := AssetModel{
				Name:      d.Name(),
				Size:      fileInfo.Size(),
				UpdatedAt: fileInfo.ModTime(),
				IsDir:     d.IsDir(),
				Extension: ext,
			}
			assets = append(assets, &asset)
		}

		return nil
	})

	return assets, nil
}

func (instance AssetProvider) Mkdir(dirPath, dirName string) (string, error) {
	path := instance.CommonProvider.GeneratePath(dirPath+dirName, 1)

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return "", err
	}

	path = strings.Replace(path, instance.PublicPath, "", 1)
	if path != "" {
		path = path[1:]
	}

	return path, nil
}
