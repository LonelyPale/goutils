package pagination

import "github.com/LonelyPale/goutils/types"

type Pagination struct {
	Current  int         `json:"current"`  //当前页数, 默认值 1
	PageSize int         `json:"pageSize"` //每页条数, 默认值 10
	Total    int         `json:"total"`    //数据总数
	Data     interface{} `json:"data"`     //切片指针
}

func (p *Pagination) Skip() int64 {
	if p.Current <= 0 {
		p.Current = 1
	}
	return int64((p.Current - 1) * p.PageSize)
}

func (p *Pagination) Limit() int64 {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return int64(p.PageSize)
}

func (p *Pagination) Result() interface{} {
	if p.Data == nil {
		p.Data = []types.M{}
	}
	return p.Data
}

func (p *Pagination) SetTotal(n int64) {
	p.Total = int(n)
}
