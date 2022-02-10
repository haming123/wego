package worm

import (
	"database/sql"
	"time"
)

type DbSlave struct {
	db_name string
	db_raw *sql.DB
	db_weight int
	weight_no int
}

func getRandInt(seed int64, w int) int {
	val := seed%int64(w)
	return int(val)
}

func getSlaveByWeight(slaves []*DbSlave) *DbSlave {
	if len(slaves) == 1 {
		return slaves[0]
	}

	tt := 0
	for i := 0; i < len(slaves); i++ {
		slaves[i].weight_no = tt
		tt += slaves[i].db_weight
	}

	//seed := time.Now().UnixNano() / 1000
	//var r = rand.New(rand.NewSource(seed))
	//sel := slaves[r.Intn(len(slaves))]
	//no := r.Intn(tt)
	seed := time.Now().Unix()
	//fmt.Println(seed)
	sel := slaves[getRandInt(seed, len(slaves))]
	no := getRandInt(seed, tt)
	for i := 0; i < len(slaves); i++ {
		if no >= slaves[i].weight_no && no < slaves[i].weight_no + slaves[i].db_weight {
			return slaves[i]
		}
	}
	return sel
}