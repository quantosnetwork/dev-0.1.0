package config

import (
	"encoding/json"
	"errors"
	"github.com/quantosnetwork/v0.1.0-dev/core/module"
	"github.com/quantosnetwork/v0.1.0-dev/fs"
	"path/filepath"
)

type ModulesConfig struct {
	repo           map[module.ModuleType]*module.ModuleRepository
	systemRequired []*module.Module
}

type modulesConfigFile struct {
	SystemRequired map[string]module.Module
}

func ReadModulesConfig(fUtils fs.FileUtils, mt string, path string) (*ModulesConfig, error) {
	mc := new(ModulesConfig)
	num, files, err := fUtils.GetFilesFromDir(filepath.Dir(path + "/" + mt))
	if err != nil {
		return nil, err
	}
	num2, files2 := fUtils.GetFilesLargerThan(int64(10), files)
	num3, withExt := fUtils.GetFilesWithExt(".config", files2)

	if num > 0 && num2 > 0 && num3 > 0 {

		for _, file := range withExt {
			s, _ := file.Stat()
			buf := make([]byte, s.Size())
			_, err2 := file.Read(buf)
			if err2 != nil {
				return nil, err2
			}
			var mod module.Module
			_ = json.Unmarshal(buf, &mod)
			mc.repo[mod.Type] = &module.ModuleRepository{Modules: map[string]*module.ModuleRepositoryItem{}}
			mc.repo[mod.Type].Modules[mod.Name] = &module.ModuleRepositoryItem{Module: mod}
			if mod.Type == module.SYSTEM {
				mc.systemRequired = append(mc.systemRequired, &mod)
			}
		}

		return mc, nil

	}
	return nil, errors.New("no config files exists")

}
