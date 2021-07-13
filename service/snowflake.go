package service

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var flowNode *snowflake.Node
var flowNodeOnce sync.Once

func GetSnowflakeId() snowflake.ID {
	//
	flowNodeOnce.Do(func() {

		var nodeId int64 = 1
		var err error
		flowNode, err = snowflake.NewNode(nodeId)
		if err != nil {
			panic(err)
		}
	})

	return flowNode.Generate()
}
