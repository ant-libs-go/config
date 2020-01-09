/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:49:07
# File Name: parser.go
# Description:
####################################################################### */

package parser

import (
	"github.com/ant-libs-go/config/options"
)

type Parser interface {
	Unmarshal(cfg interface{}, opts *options.Options) error
	GetLastModTime(opts *options.Options) (int64, error)
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
