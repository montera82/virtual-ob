package mock

import "github.com/bwmarrin/snowflake"

type SnowFlakeNode struct {
	GenerateFunc func() snowflake.ID
}

func (sf *SnowFlakeNode) Generate() snowflake.ID {
	return sf.GenerateFunc()
}
