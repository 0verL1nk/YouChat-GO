package utils

import (
	"github.com/bwmarrin/snowflake"
)

// 雪花算法生成uint64的随机ID
func GenNumId() (res int64, err error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0, err
	}
	return node.Generate().Int64(), nil
}
