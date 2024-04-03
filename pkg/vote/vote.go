package vote

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

var beginTime = "2024-03-15 07:46:43"

func GetTimeUnix(t time.Time) int64 {
	layout := "2006-01-02 15:04:05"
	parsedTime, _ := time.Parse(layout, beginTime)
	return t.Unix() - parsedTime.Unix()
}

/* 改进后的美国网站Reddit的帖子热度算法
满足以下几个条件
1.越新的帖子分数越高 （占大头）
2.反对票越多的帖子分数越低
3.讨论度越高的帖子分数越高
*/

func hot(up, down int64, ti time.Time) (score float64) {

	// t表示beginTime到发帖为止的时间戳   t = 发贴时间 - beginTime
	//t越大，得分越高，即新帖子的得分会高于老帖子。它起到自动将老帖子的排名往下拉的作用。
	t := GetTimeUnix(ti)

	// z表示赞成票与反对票之间差额的绝对值。如果对某个帖子的评价，越是一边倒，z就越大。如果赞成票等于反对票，z就等于1。
	// old: z = abs(up - down)
	z := float64(up + down)
	if z == 0 {
		z = 1
	}
	// 第一项：表示对z求一次以10为底的对数,这里可以表示z越大对分数的影响越小 也就是投票越多 投票对帖子热度的影响越小
	part1 := math.Log2(z)
	// 第二项：表示赞成多的帖子分数高，赞成低的帖子分数低，时间越靠前的帖子分数越低 t/45000表示一天的分数大约为两分
	part2 := float64(t / 45000)
	// 第三项：表示反对票如果比赞成票多，那么帖子的热度会下降 反对票每比赞成票多50 分数减1
	part3 := float64(max(0, (down-up)/50))

	score, _ = strconv.ParseFloat(fmt.Sprintf("%.7f", part1+part2+part3), 64)
	return score
}
