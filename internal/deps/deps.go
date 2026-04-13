package deps

import (
	"os"

	"github.com/pbm-org/pbm/internal/config"
)

// create cache proto dir
// git clone --depth 1 --branch $ref https://github.com/medianexapp/plugin_api plugin_api@$ref
// git clone https://github.com/medianexapp/plugin_api plugin_api@ref && git checkout $ref

func Clone(dep config.Dep) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	_, err = os.Stat(cacheDir)
	if err != nil {
		err = os.Mkdir(cacheDir, 0644)
		if err != nil {
			return err
		}
	}

	// dir := filepath.Base(dep.Path)
	// destDir := filepath.Join(cacheDir,dir)
	return nil

}
