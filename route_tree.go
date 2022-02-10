package wego

type RouteTree struct {
	tree  	TreeNodeX
	routes 	[]*RouteInfo
}

func trimSpace(path string) string {
	//trim left
	if len(path) > 0 {
		i:=0
		for ; i < len(path); i++ {
			if path[i] != ' ' {
				break
			}
		}
		if i > 0 {
			path = path[i:]
		}
	}

	//trim right
	if len(path) > 0 {
		i:=len(path)-1
		for ; i>=0; i-- {
			if path[i] != ' ' {
				break
			}
		}
		if i < len(path)-1 {
			path = path[0:i+1]
		}
	}

	return path
}

func (r *RouteTree) addRoute(method string, path string, rinfo *RouteInfo) {
	tree := &r.tree
	path = trimSpace(path)
	tree.AddRoute(method, path, rinfo)
	r.routes = append(r.routes, rinfo)
}

func (r *RouteTree) getRoute(method string, path string, params *PathParam) *RouteInfo {
	tree := &r.tree
	path = trimSpace(path)
	return tree.GetRoute(method, path, params)
}

func (r *RouteTree) initRouteFilter() {
	for _, item := range r.routes {
		if len(item.filters) < 1 {
			item.filters = item.group.getHandlerFilter()
		}
		if item.before_func == nil && item.before_mthd == nil {
			item.before_func = item.group.getBeforExec()
		}
		if item.after_func == nil && item.after_mthd == nil {
			item.after_func = item.group.getAfterExec()
		}
	}
}
