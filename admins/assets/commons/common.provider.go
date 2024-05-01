package commons

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dangduoc08/gogo/core"
)

type CommonProvider struct{}

func (instance CommonProvider) NewProvider() core.Provider {
	return instance
}

func (instance CommonProvider) CleanDir(dir string) string {
	if dir != "" {
		if dir[0] != filepath.Separator {
			dir = "/" + dir
		}

		dir = filepath.Clean(dir)
	}
	return dir
}

func (instance CommonProvider) GeneratePath(path string, i int) string {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ext := filepath.Ext(path)
		pathWithoutExt := strings.TrimSuffix(path, ext)
		if i > 1 {
			pathWithoutExt = strings.TrimSuffix(pathWithoutExt, fmt.Sprintf(" (%v)", i))
		}

		i++
		newPath := fmt.Sprintf("%v (%v)%v", pathWithoutExt, i, ext)
		return instance.GeneratePath(newPath, i)
	}

	return path
}
