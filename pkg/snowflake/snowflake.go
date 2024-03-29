package mySnowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

func GetID() int64 {
	return node.Generate().Int64()
}

//func main() {
//	err := Init("2006-01-02", 1)
//	if err != nil {
//		fmt.Printf(", err: %v\n", err)
//	}
//	fmt.Println(GetID())
//}
