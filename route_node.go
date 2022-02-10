package wego

type nodeType int
const (
	nt_spot nodeType = iota
	nt_line
	nt_wild
)

type TreeNode struct {
	name     	string
	part     	string
	line  		*TreeNode
	wild  		*TreeNode
	children 	map[string]*TreeNode
	h_get  		*RouteInfo
	h_post 		*RouteInfo
	h_put 		*RouteInfo
	h_patch 	*RouteInfo
	t_delete 	*RouteInfo
	h_route		*RouteInfo
}

func (n *TreeNode) addChild(ntype nodeType, child *TreeNode) {
	part := child.part
	if ntype == nt_line {
		n.line = child
	} else if ntype == nt_wild {
		n.wild = child
	} else {
		if n.children == nil {
			n.children = map[string]*TreeNode{part: child}
			return
		}
		n.children[part] = child
	}
}

func (n *TreeNode) getChild(part string, ntype nodeType) *TreeNode {
	if ntype == nt_line && n.line != nil {
		return n.line
	} else if ntype == nt_wild && n.wild != nil {
		return n.wild
	} else if ntype == nt_spot {
		if n.children == nil {
			return nil
		} else {
			return n.children[part]
		}
	}
	return nil
}

func (n *TreeNode) hasChildren() bool {
	if n.line != nil {
		return true
	}
	if n.wild != nil {
		return true
	}
	if len(n.children) > 0 {
		return true
	}
	return false
}

func (n *TreeNode) isEndNode() bool {
	if n.line != nil {
		return false
	}
	if n.wild != nil {
		return false
	}
	if len(n.children) > 0 {
		return false
	}
	return true
}

func getPartType(part string) nodeType {
	if len(part) < 1 {
		return nt_spot
	}
	if part[0] == ':' {
		return nt_line
	} else if part[0] == '*' {
		return nt_wild
	} else {
		return nt_spot
	}
}

func (n *TreeNode) setHander(method string, pval *RouteInfo) {
	if method == MethodGet {
		n.h_get = pval
	} else if method == MethodPost {
		n.h_post = pval
	} else if method == MethodPut {
		n.h_put = pval
	} else if method == MethodPatch {
		n.h_patch = pval
	} else if method == MethodDelete {
		n.t_delete = pval
	} else if method == MethodPath {
		n.h_route = pval
	} else {
		panic("invalid method")
	}
}

func getUrlPart(path string) (string, string) {
	slen := len(path)
	if slen < 1 {
		return "", ""
	}

	if path[0] == '/' {
		path = path[1:]
		slen -= 1
	}

	for i:=0; i < slen; i++ {
		if path[i] == '/' {
			return path[:i], path[i:]
		}
	}

	return path, ""
}

func getPartName(part string) string {
	if part == "" {
		return part
	}

	if part[0] == ':' || part[0] == '*' {
		part = part[1:]
	}
	return part
}

func (n *TreeNode) AddRoute(method string, url_path string, pval *RouteInfo) {
	//pattern="" 时对应根目录
	if url_path == "" {
		n.part = ""
		n.setHander(method, pval)
		return
	}

	part, url_path := getUrlPart(url_path)
	for ntmp := n; ntmp != nil; {
		ntype := getPartType(part)
		child := ntmp.getChild(part, ntype)
		if child == nil {
			//创建并添加子结点
			name := getPartName(part)
			child = &TreeNode{part: part, name:name}
			ntmp.addChild(ntype, child)
		}

		//若pp == num -1，意味着这是最后的结点了
		//若子结点类型=nt_wild，意味着这是最后的结点了
		if url_path == "" || ntype == nt_wild {
			child.setHander(method, pval)
			return
		}

		ntmp = child
		part, url_path = getUrlPart(url_path)
	}
}

func (n *TreeNode) getHander(method string) *RouteInfo {
	var hinfo  *RouteInfo = nil
	if method == MethodGet {
		hinfo = n.h_get
	} else if method == MethodPost {
		hinfo = n.h_post
	} else if method == MethodPut {
		hinfo = n.h_put
	} else if method == MethodPatch {
		hinfo = n.h_patch
	} else if method == MethodDelete {
		hinfo = n.t_delete
	}
	if hinfo == nil {
		hinfo = n.h_route
	}
	return hinfo
}

func (n *TreeNode) GetRoute(method string, url_path string, params *PathParam) *RouteInfo {
	//url_path="" 时对应根目录
	if url_path == "" {
		return n.getHander(method)
	}

	url_part, url_path := getUrlPart(url_path)
	for ntmp := n; ntmp != nil; {
		//首先进行固定匹配
		//若存在，并且是最后结点，则返回
		//若存在，并且不是最后结点，则把继续child作为当前结点继续比较
		child := ntmp.getChild(url_part, nt_spot)
		if child != nil && url_path == "" {
			return child.getHander(method)
		} else if child != nil {
			ntmp = child
			url_part, url_path = getUrlPart(url_path)
			continue
		}

		//然后进行:模式匹配
		//若存在，将url_part作为参数加入params
		//若存在，并且是最后结点，则返回结果
		//若存在，并且不是最后结点，则把继续child作为当前结点继续比较
		child = ntmp.getChild(url_part, nt_line)
		if child != nil {
			params.SetValue(child.name, url_part)
		}
		if child != nil && url_path == "" {
			return child.getHander(method)
		} else if child != nil {
			ntmp = child
			url_part, url_path = getUrlPart(url_path)
			continue
		}

		//然后进行*模式匹配
		//若存在，将url_part作为参数加入params，然后返回结果
		child = ntmp.getChild(url_part, nt_wild)
		if child != nil {
			params.SetValue(child.name, url_part + url_path)
			return child.getHander(method)
		}

		ntmp = child
		url_part, url_path = getUrlPart(url_path)
	}

	return nil
}
