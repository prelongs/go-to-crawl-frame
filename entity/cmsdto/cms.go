package cmsdto

type PageParam interface {
	InitPageParam()
}

type CmsBasePageQry struct {
	Id        int `p:"id"`
	PageIndex int `p:"pageIndex"`
	PageSize  int `p:"pageSize"`
}

func (qry *CmsBasePageQry) InitPageParam() {
	if qry.PageIndex < 0 {
		qry.PageIndex = 0
	}

	if qry.PageSize < 10 {
		qry.PageSize = 10
	}
}
