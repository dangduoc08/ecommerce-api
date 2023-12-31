package providers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/models"
	"github.com/dangduoc08/gooh/core"
)

type HandleAsset struct {
	MaxDepth   int
	PublicPath string
}

func (self HandleAsset) NewProvider() core.Provider {
	self.MaxDepth = 1

	return self
}

func (self HandleAsset) List(dirPath string) ([]*models.Asset, error) {
	assets := []*models.Asset{}
	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Count(path, string(os.PathSeparator))-strings.Count(dirPath, string(os.PathSeparator)) > self.MaxDepth {
			return fs.SkipDir
		}

		if dirPath != path {
			fileInfo, _ := d.Info()
			ext := filepath.Ext(path)
			if ext != "" {
				ext = ext[1:]
			}

			asset := *&models.Asset{
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

func (self HandleAsset) Mkdir(dirPath, dirName string) (string, error) {
	path := self.GeneratePath(dirPath+dirName, 1)

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return "", err
	}

	path = strings.Replace(path, self.PublicPath, "", 1)
	if path != "" {
		path = path[1:]
	}

	return path, nil
}

func (self HandleAsset) CleanDir(dir string) string {
	if dir != "" {
		if dir[0] != filepath.Separator {
			dir = "/" + dir
		}

		dir = filepath.Clean(dir)
	}
	return dir
}

func (self HandleAsset) GeneratePath(path string, i int) string {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ext := filepath.Ext(path)
		pathWithoutExt := strings.TrimSuffix(path, ext)
		if i > 1 {
			pathWithoutExt = strings.TrimSuffix(pathWithoutExt, fmt.Sprintf(" (%v)", i))
		}

		i++
		newPath := fmt.Sprintf("%v (%v)%v", pathWithoutExt, i, ext)
		return self.GeneratePath(newPath, i)
	}

	return path
}
