package worm

/*
func TestStmtCache (t *testing.T) {
	err := demoOpenDb()
	if err != nil {
		t.Error(err)
		return
	}

	engine, err := NewEngine("mysql", dbcon_test)
	if err != nil {
		t.Error(err)
		return
	}
	engine.ShowLog(true)
	engine.SetLogger(NewSimpleLogger())
	engine.stmt_cache.max_size = 5

	dbs := engine.NewSession()
	dbs.UsePrepare(true)

	var user DB_User
	dbs.Model(&user).Where("id=?", 1).Get()
	clist_num :=engine.stmt_cache.cacheList.Len()
	wlist_num :=engine.stmt_cache.wasteList.Len()
	//cache_list:a, waste_list:
	t.Logf("cache_list:%d, waste_list:%d\n", clist_num, wlist_num)

	i:=0
	for {
		i+=1
		time.Sleep(2*time.Second)
		dbs.Model(&user).Where("id=?", 1).Get()
		clist_num :=engine.stmt_cache.cacheList.Len()
		wlist_num :=engine.stmt_cache.wasteList.Len()
		t.Logf("cache_list:%d, waste_list:%d\n", clist_num, wlist_num)

		time.Sleep(8*time.Second)
		str := fmt.Sprintf("id=? and %d=%d", i, i)
		dbs.Model(&user).Where(str, 2).Get()
		clist_num =engine.stmt_cache.cacheList.Len()
		wlist_num =engine.stmt_cache.wasteList.Len()
		t.Logf("cache_list:%d, waste_list:%d\n", clist_num, wlist_num)
	}
	time.Sleep(10*time.Second)
}
*/