package wego

import "strings"

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
	MethodPath   = "PATH"
)

type TreeNodeX struct {
	name     string
	part     string
	items    [256]*TreeNodeX
	h_get    *RouteInfo
	h_post   *RouteInfo
	h_put    *RouteInfo
	h_patch  *RouteInfo
	t_delete *RouteInfo
	h_route  *RouteInfo
}

func (n *TreeNodeX) setHander(method string, pval *RouteInfo) {
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

func getSameIndex(str1 string, str2 string) int {
	index := -1
	for ss := 0; ss < len(str1) && ss < len(str2); ss++ {
		if str1[ss] != str2[ss] {
			break
		}
		index = ss
	}
	return index
}

func splitPath(path string) (string, string, string) {
	index := -1
	for c := 0; c < len(path); c++ {
		if path[c] == ':' || path[c] == '*' {
			index = c
			break
		}
	}

	if index < 0 {
		return path, "", ""
	}

	str1 := path[0:index]
	path = path[index:]
	index = -1
	for c := 0; c < len(path); c++ {
		if path[c] == '/' {
			index = c
			break
		}
	}

	if index < 0 {
		return str1, path, ""
	}

	name := path[0:index]
	str2 := path[index:]
	return str1, name, str2
}

func (n *TreeNodeX) AddRoute(method string, url_path string, pval *RouteInfo) {
	//url_path="" 时对应根目录
	if url_path == "" {
		n.part = ""
		n.setHander(method, pval)
		return
	}

	node_cur := n
	for {
		//若不存在对应槽位的子结点，则在该槽位新增一个结点
		child := node_cur.items[url_path[0]]
		if child == nil {
			//若存在":"或"*"，将path拆分位3个部分
			str1, ppp, str2 := splitPath(url_path)
			if ppp != "" {
				//若存在路由参数，则首先创建str1结点
				if str1 != "" {
					node_new := &TreeNodeX{part: str1}
					node_cur.items[str1[0]] = node_new
					node_cur = node_new
				}
				//然后创建参数结点
				node_new := &TreeNodeX{part: ppp, name: ppp[1:]}
				node_cur.items[ppp[0]] = node_new
				if str2 == "" {
					//例如:/hello/:name
					url_path = str2
					node_new.setHander(method, pval)
					return
				} else {
					//例如:/hello/:name/...
					url_path = str2
					node_cur = node_new
					continue
				}
			} else {
				//若不存在路由参数,则创建一个新的结点，并返回
				url_path = str2
				node_new := &TreeNodeX{part: str1}
				node_cur.items[str1[0]] = node_new
				node_new.setHander(method, pval)
				return
			}
		}

		//若存在对应槽位的子结点，则与现路由进行字符串比较
		str_child := child.part
		ss := getSameIndex(str_child, url_path)
		if ss == len(str_child)-1 && ss == len(url_path)-1 {
			//新增路由与现存路由完全相同，则覆盖
			child.setHander(method, pval)
			return
		} else if ss == len(str_child)-1 {
			//新增路由包含了现存路由，则将新增路由截断，
			//并将截断的部分作为现存路由的下级加入
			node_cur = child
			url_path = url_path[ss+1:]
			continue
		} else if ss == len(url_path)-1 {
			//现存路由包含了新增路由，则将现存路由截断
			//用截断的部分创建一个新的结点作为孙子结点
			node_new := &TreeNodeX{part: str_child[0 : ss+1]}
			node_cur.items[str_child[0]] = node_new
			//把原来的结点作为孙子结点
			grandson := child
			grandson.part = str_child[ss+1:]
			node_new.items[grandson.part[0]] = grandson

			node_new.setHander(method, pval)
			return
		} else {
			//现存路由与新增路由部分相同，则将现存路由截断
			//用截断的部分创建一个新的结点作为孙子结点
			//并将新增路由也截断，截断的部分作为孙子结点
			node_new := &TreeNodeX{part: str_child[0 : ss+1]}
			node_cur.items[str_child[0]] = node_new
			//用现存截断的部分创建一个新的结点作为孙子结点
			grandson := child
			grandson.part = str_child[ss+1:]
			node_new.items[grandson.part[0]] = grandson

			//将新增路由截断新的部分创建一个新的结点作为孙子结点
			node_cur = node_new
			url_path = url_path[ss+1:]
			continue
		}
	}
}

func (n *TreeNodeX) getHander(method string) *RouteInfo {
	var hinfo *RouteInfo = nil
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

func trimValue(url_path string) (string, string) {
	index := strings.Index(url_path, "/")
	if index < 0 {
		return url_path, ""
	}

	value := url_path[0:index]
	url_path = url_path[index:]
	return value, url_path
}

func (n *TreeNodeX) GetRoute(method string, url_path string, params *PathParam) *RouteInfo {
	//url_path="" 时对应根目录
	if url_path == "" {
		return n.getHander(method)
	}

	var m_node *TreeNodeX = nil
	var m_path string = ""
	var m_num int = 0
	var x_node *TreeNodeX = nil
	var x_path string = ""
	var x_num int = 0
	node_cur := n
	for {
		char0 := url_path[0]
		child := node_cur.items[char0]
		if node_cur.items[':'] != nil {
			m_node = node_cur.items[':']
			m_path = url_path
			m_num = len(params.items)
		}
		if node_cur.items['*'] != nil {
			x_node = node_cur.items['*']
			x_path = url_path
			x_num = len(params.items)
		}

		//若不存在对应槽位的子结点，则查询":"结点
		//若不匹配":"结点， 则查询"*"结点，
		//若不匹配"*"结点， 则返回nil
		if child != nil {
			str_child := child.part
			ss := getSameIndex(str_child, url_path)
			if ss == len(str_child)-1 && ss == len(url_path)-1 {
				//path与路由完全相同，则返回
				return child.getHander(method)
			} else if ss == len(str_child)-1 {
				//path包含了路由，则将新path截断，并继续匹配截断的部分
				node_cur = child
				url_path = url_path[ss+1:]
				continue
			}
		}

		//查询":"结点
		if m_node != nil {
			child = m_node
			m_node = nil
			url_path = m_path
			if len(params.items) > m_num {
				params.items = params.items[0:m_num]
			}

			value, str2 := trimValue(url_path)
			params.SetValue(child.name, value)
			if str2 == "" {
				url_path = str2
				return child.getHander(method)
			} else {
				url_path = str2
				node_cur = child
				continue
			}
		}

		//查询"*"结点
		if x_node != nil {
			child = x_node
			x_node = nil
			url_path = x_path
			if len(params.items) > x_num {
				params.items = params.items[0:x_num]
			}

			params.SetValue(child.name, url_path)
			return child.getHander(method)
		} else {
			return nil
		}
	}

	return nil
}
