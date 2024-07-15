package spi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

var (
	ErrDirNotFound        = errors.New("ekit: 目录不存在")
	ErrSymbolNameIsEmpty  = errors.New("ekit: 结构体名不能为空")
	ErrOpenPluginFailed   = errors.New("ekit: 打开插件失败")
	ErrSymbolNameNotFound = errors.New("ekit: 从插件中查找对象失败")
	ErrInvalidSo          = errors.New("ekit: 插件非该接口类型")
)

func LoadService[T any](dir string, symName string) ([]T, error) {
	var services []T
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w", ErrDirNotFound)
	}
	if symName == "" {
		return nil, fmt.Errorf("%w", ErrSymbolNameIsEmpty)
	}
	// 遍历目录下的所有 .so 文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			// 打开插件
			p, err := plugin.Open(path)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrOpenPluginFailed, err)
			}
			// 查找变量
			sym, err := p.Lookup(symName)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrSymbolNameNotFound, err)
			}

			// 尝试将符号断言为接口类型
			service, ok := sym.(T)
			if !ok {
				return fmt.Errorf("%w", ErrInvalidSo)
			}
			// 收集服务
			services = append(services, service)
		}
		return nil
	})
	return services, err
}
