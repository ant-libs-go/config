/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:49:07
# File Name: parser.go
# Description:
####################################################################### */

package parser

import (
	"fmt"
	"os"

	"github.com/ant-libs-go/config/options"
)

type Parser interface {
	Unmarshal(cfg interface{}, opts *options.Options) error
	GetLastModTime(opts *options.Options) (int64, error)
}

func ParseFileLastModTime(file string) (r int64, err error) {
	fd, err := os.Stat(file)
	if err != nil {
		err = fmt.Errorf("parse file last modified time fail, %s", err)
		return
	}
	r = fd.ModTime().Unix()
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
