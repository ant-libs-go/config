/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:59:22
# File Name: toml.go
# Description:
####################################################################### */

package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/ant-libs-go/config/options"
	stringutil "github.com/naoina/go-stringutil"
	"github.com/naoina/toml"
)

type TomlParser struct {
	toml    toml.Config
	impToml toml.Config
}

func NewTomlParser() *TomlParser {
	tomlNormFieldName := func(typ reflect.Type, s string) string {
		return strings.Replace(strings.ToLower(s), "_", "", -1)
	}
	tomlSnakeCase := func(typ reflect.Type, s string) string {
		return stringutil.ToSnakeCase(s)
	}
	tomlMissingField := func(typ reflect.Type, key string) error {
		return fmt.Errorf("field corresponding to `%s' is not defined in %v", key, typ)
	}

	defaultToml := toml.Config{
		NormFieldName: tomlNormFieldName,
		FieldToKey:    tomlSnakeCase,
		MissingField: func(typ reflect.Type, field string) error {
			if field == "import" {
				return nil
			}
			return tomlMissingField(typ, field)
		},
	}
	impToml := toml.Config{
		NormFieldName: tomlNormFieldName,
		FieldToKey:    tomlSnakeCase,
		MissingField: func(typ reflect.Type, field string) error {
			return nil
		},
	}

	return &TomlParser{
		toml:    defaultToml,
		impToml: impToml,
	}
}

type Import struct {
	Import []string
}

func (this *TomlParser) Unmarshal(cfg interface{}, opts *options.Options) (err error) {
	files := []string{opts.File}

	for len(files) > 0 {
		file := files[len(files)-1]
		files = files[:len(files)-1]

		imp := &Import{}
		if err = this.load(imp, file); err != nil {
			return
		}
		for i := len(imp.Import); i > 0; i-- {
			files = append(files, path.Join(path.Dir(file), imp.Import[i-1]+".toml"))
		}

		if err = this.load(cfg, file); err != nil {
			return
		}
	}
	return
}

func (this *TomlParser) load(cfg interface{}, file string) (err error) {
	if len(file) == 0 {
		return fmt.Errorf("Config file not specified")
	}
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Config file open fail, %s", err)
	}
	defer f.Close()

	buff, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("Config file read fail, %s", err)
	}

	if _, ok := cfg.(*Import); !ok {
		err = this.toml.Unmarshal(buff, cfg)
	} else {
		err = this.impToml.Unmarshal(buff, cfg)
	}
	if err != nil {
		return fmt.Errorf("Config file Unmarshal fail, %s", err)
	}
	return
}

func (this *TomlParser) GetLastModTime(opts *options.Options) (r int64, err error) {
	files := []string{opts.File}

	for len(files) > 0 {
		file := files[len(files)-1]
		files = files[:len(files)-1]

		imp := &Import{}
		if err = this.load(imp, file); err != nil {
			return
		}
		for i := len(imp.Import); i > 0; i-- {
			files = append(files, path.Join(path.Dir(file), imp.Import[i-1]+".toml"))
		}

		var modTime int64
		if modTime, err = this.getLastModTime(file); err != nil {
			return
		}
		if modTime > r {
			r = modTime
		}
	}
	return
}

func (this *TomlParser) getLastModTime(file string) (r int64, err error) {
	fd, err := os.Stat(file)
	if err != nil {
		err = fmt.Errorf("Get file last modified time fail, %s", err)
		return
	}
	r = fd.ModTime().Unix()
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
