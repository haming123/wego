package worm

type StringRow map[string]string

type StringTable struct {
	columns []string
	data []string
}

func (data *StringTable)GetColCount() int {
	return len(data.columns)
}

func (data *StringTable)GetRowCount() int {
	n_col := len(data.columns)
	n_cell := len(data.data)
	return n_cell/n_col
}

func (data *StringTable)GetColumns() []string {
	return data.columns
}

func (data *StringTable)GetRowData(r_no int) StringRow {
	row_data := make(StringRow)
	n_col := len(data.columns)
	n_cell := len(data.data)
	beg := n_col* r_no
	end := beg + n_col
	if end > n_cell {
		return row_data
	}

	for i:=0; i < n_col; i++{
		row_data[data.columns[i]] = data.data[beg+i]
	}
	return row_data
}