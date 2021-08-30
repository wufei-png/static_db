package lib_worker

import (
	"fmt"
	"math"
	"math/rand"
)

type Request struct {
	Vectors      Vectors //向量组
	Topk         int
	Request_type int
	I_delete     float64
	J_delete     float64
}
type Vectors struct {
	Vector   [][]float64 //向量
	Distance []float64
}
type Result struct {
	TopkDistance [][]float64
	VectorGroup  [][][]float64
}

type Node struct {
	Data     []float64
	Next     *Node
	Distance float64
}
type Node_first struct {
	Data       []float64
	Next       *Node
	Next_first *Node_first
	Distance   float64
}
type StatusDB struct {
	Head *Node_first
}

var DB StatusDB

// func Newnode_zz(dimension int) Node { //带指针的
// 	var newnode Node
// 	for i := 0; i < dimension; i++ {
// 		newnode.Data = append(newnode.Data, 10*rand.Float64()) //在这里用图片的向量
// 	}

// 	return newnode
// }
func Newnode(dimension int, index int) Node { //带指针的
	var newnode Node
	for i := 0; i < dimension; i++ {
		newnode.Data = append(newnode.Data, float64(index)+rand.Float64()) //在这里用图片的向量
	}

	return newnode
}
func Newnode_first(dimension int, value float64) Node_first {
	var newnode_first Node_first
	for i := 0; i < dimension; i++ {
		newnode_first.Data = append(newnode_first.Data, value+0.5) //在这里用图片的向量
	}
	return newnode_first
}
func Newnode_first_result(dimension int) Node_first {
	var newnode_first Node_first
	for i := 0; i < dimension; i++ {
		newnode_first.Data = append(newnode_first.Data, 0) //在这里用图片的向量
	}
	return newnode_first
}
func result_init(topk int, numOfVector int, dimension int) Node_first {
	var head_nodefirst Node_first
	var index_nodefirst *Node_first
	index_nodefirst = &head_nodefirst
	if numOfVector == 1 {
		if topk == 1 {
			return head_nodefirst
		}
		var index_newnode Node
		index_nodefirst.Next = &index_newnode
		var node *Node = &index_newnode
		for j := 0; j < topk-2; j++ {
			var newnode Node
			node.Next = &newnode
			node = &newnode
		}
		return head_nodefirst
	} else {
		for i := 0; i < numOfVector; i++ {
			if topk == 1 {
				if i == 0 {
					continue
				}
				var newnodefirst = Newnode_first_result(dimension)
				index_nodefirst.Next_first = &newnodefirst
				index_nodefirst = &newnodefirst
			} else {
				if i == 0 {
					var index_newnode Node
					index_nodefirst.Next = &index_newnode
					var node *Node = &index_newnode
					for j := 0; j < topk-2; j++ {
						var newnode Node
						node.Next = &newnode
						node = &newnode
					}
					continue
				}
				var newnodefirst = Newnode_first_result(dimension)
				index_nodefirst.Next_first = &newnodefirst
				index_nodefirst = &newnodefirst
				var index_newnode Node
				index_nodefirst.Next = &index_newnode
				var node *Node = &index_newnode
				for j := 0; j < topk-2; j++ {
					var newnode Node
					node.Next = &newnode
					node = &newnode
				}
			}

		}
		return head_nodefirst
	}
}
func Display_DBbyDistance() {
	for indexi := DB.Head; indexi != nil; indexi = indexi.Next_first {
		fmt.Print(indexi.Distance, "->")
		for indexj := indexi.Next; indexj != nil; indexj = indexj.Next {

			if indexj.Next == nil {
				fmt.Print(indexj.Distance)
				fmt.Println()
			} else {
				fmt.Print(indexj.Distance, "->")
			}
		}
	}
}
func Display_DBbydata() {
	for indexi := DB.Head; indexi != nil; indexi = indexi.Next_first {
		fmt.Print(indexi.Data, "->")
		for indexj := indexi.Next; indexj != nil; indexj = indexj.Next {
			if indexj.Next == nil {
				fmt.Print(indexj.Data)
				fmt.Println()
			} else {
				fmt.Print(indexj.Data, "->")
			}
		}
	}
}
func Display_allfirst() {
	for indexi := DB.Head; indexi != nil; indexi = indexi.Next_first {
		fmt.Print(indexi.Data)
	}
}
func findmin_index(start *Node_first, node *Node) *Node_first { //找到给定的一个node距离给出的nodefirst以及之后的距离他最近的nodefirst\
	var min float64
	var index int
	var j int
	for start_copy := start; start_copy != nil; start_copy = start_copy.Next_first {
		if start_copy == start {
			//fmt.Print(newnode.Data, index_nodefirst.Data, 1)
			min = Distance(start_copy.Data, node.Data)
			j++
			index = j
			continue
		} else {
			//fmt.Print(newnode.Data, index_nodefirst.Data, 2)
			if min > Distance(start_copy.Data, node.Data) {
				min = Distance(start_copy.Data, node.Data)
				index = j + 1
			}
		}
		j++
	}
	node.Distance = min
	return start.find_nodefirst(index)
}
func dbinit_train(i int) bool { //终止条件聚类中心不再改变
	// var db
	// nodefirst和db一样(firstnode已经变了，遍历dbnode,最近的firstnode调用add
	// 	db.=var的.
	fmt.Println("迭代第", i, "次")
	var stop bool = true
	for index_first := DB.Head; index_first != nil; index_first = index_first.Next_first {

		var sum = make([]float64, 4, 4)
		for i := range index_first.Data {
			//fmt.Print(index_first.D)
			sum[i] += index_first.Data[i]
		}

		for index_node := index_first.Next; index_node != nil; index_node = index_node.Next {
			for i := range index_node.Data {
				sum[i] += index_node.Data[i]
			}
		}
		for i := range sum {
			sum[i] = sum[i] / float64(index_first.length())
		} //至此算出平均值
		var min float64
		var index int
		var index_node = index_first.Next
		for j := 0; j < index_first.length(); j++ {
			if j == 0 {
				//fmt.Print(newnode.Data, index_nodefirst.Data, 1)
				min = Distance(sum, index_first.Data)
				continue
			} else {
				//fmt.Print(newnode.Data, index_nodefirst.Data, 2)
				if min > Distance(sum, index_node.Data) {
					min = Distance(sum, index_node.Data)
					index = j
				}
			}
			index_node = index_node.Next
		} //得到最小的索引index

		if index != 0 { //聚类中心变了
			stop = false //只要有一个就继续迭代
		} else {
			continue
		}
		var node = index_first.Next
		for i := 0; i < index-1; i++ {
			node = node.Next
		}
		for i := 0; i < len(index_first.Data); i++ {
			index_first.Data[i], node.Data[i] = node.Data[i], index_first.Data[i]
		} //交换向量，接下来Distance将重置，所以不需要交换Distance
	}

	var db_copy StatusDB
	var start Node_first
	db_copy.Head = &start
	start.Data = DB.Head.Data
	var pointer *Node_first = &start

	for index_first := DB.Head.Next_first; index_first != nil; index_first = index_first.Next_first {
		var newnode_first Node_first
		newnode_first.Data = index_first.Data
		pointer.Next_first = &newnode_first
		pointer = &newnode_first
	} //复制first

	for index_first := DB.Head; index_first != nil; index_first = index_first.Next_first {

		for index_node := index_first.Next; index_node != nil; index_node = index_node.Next { //遍历所有node
			//fmt.Print(1)
			var firstnode = findmin_index(db_copy.Head, index_node) //firstnode是dbdb的
			//fmt.Print(2)
			firstnode.add_tr(index_node)
			//fmt.Println(3)
		}
	}

	DB.Head = db_copy.Head //更新完毕
	if stop {
		fmt.Println("迭代结束")
	}
	return stop
}
func Dbinit_train() {
	var i int = 1
	for {
		if dbinit_train(i) {
			break
		}
		i++
	}
}
func Lib_worker_DBinit(num_node int, num_nodefirst int, dimension int) StatusDB {
	var newnodefirst Node_first
	head_nodefirst := &newnodefirst
	var i float64
	for ; int(i) < dimension; i++ {
		head_nodefirst.Data = append(head_nodefirst.Data, 0.5) //在这里用图片的向量
	}
	//第一个，不能调用new函数
	var index_nodefirst *Node_first
	index_nodefirst = head_nodefirst
	for i := 0; i < num_nodefirst-1; i++ {
		newnode_first := Newnode_first(dimension, float64(i+1)) //按1间隔，初始0.5
		index_nodefirst.Next_first = &newnode_first
		index_nodefirst = &newnode_first
	} //构建first
	// for i := 0; i < num_nodefirst-1; i++ {这里报错了，第三个first为nil
	// 	newnode_first := Newnode_first(dimension, index_nodefirst)
	// 	index_nodefirst = &newnode_first

	// } //构建first

	index_nodefirst = head_nodefirst //指针回去
	// for ; index_nodefirst != nil; index_nodefirst = index_nodefirst.Next_first {

	// }
	//fmt.Print(num_node)
	for i := 0; i < num_node; i++ {
		var newnode = Newnode(dimension, i*num_nodefirst/num_node)
		index_nodefirst = findmin_index(head_nodefirst, &newnode)
		// fmt.Print(nodefirst_location(head_nodefirst, index_nodefirst))
		// fmt.Println(i*num_nodefirst/num_node + 1)
		index_nodefirst.add(&newnode)
	}
	return StatusDB{Head: head_nodefirst}
}

// func nodefirst_location(start *Node_first, end *Node_first) int {
// 	var i int
// 	for index := start; index != end; index = index.Next_first {
// 		i++
// 	}
// 	return i + 1
// }
func Distance(a []float64, b []float64) float64 { //欧式
	if len(a) != len(b) {
		panic(fmt.Errorf("比较向量长度不等"))
	}
	var sum float64

	for i := range a {
		sum = float64(sum) + math.Pow(float64(a[i]-b[i]), 2)
	}
	sum = math.Pow(float64(sum), 0.5)
	return sum
}
func (head *Node_first) find_nodefirst(i int) *Node_first {
	var index_node *Node_first = head
	for index := 0; index < i-1; index++ {
		if index_node != nil {
			index_node = index_node.Next_first
		} else {
			fmt.Errorf("find方法错误")
		}
	}
	return index_node
}
func (nodefirst *Node_first) length() int {
	var i int = 1
	for index := nodefirst.Next; index != nil; index = index.Next {
		i++
	}
	return i
}
func (nodefirst *Node_first) add_tr(newnode *Node) { //??这里引用改变原来的没，在某个分片里加
	var node Node
	node.Data = newnode.Data
	nodefirst.add(&node)
}

func (nodefirst *Node_first) add(newnode *Node) { //??这里引用改变原来的没，在某个分片里加

	newnode.Distance = Distance(newnode.Data, nodefirst.Data)
	if nodefirst.Next == nil {
		nodefirst.Next = newnode
		return
	} else {
		var index *Node = nodefirst.Next
		if newnode.Distance < index.Distance { //第一个node
			nodefirst.Next = newnode
			newnode.Next = index

			return
		} else {
			var preindex = index
			index = index.Next
			for index != nil {
				if newnode.Distance < index.Distance {
					preindex.Next = newnode
					newnode.Next = index

					return
				} else {
					preindex = index
					index = index.Next

				}

			}
			preindex.Next = newnode

			return
		}
	}
}

// func swap(a *float64, b *float64) {
// 	buf := *a
// 	*a = *b
// 	*b = buf
// }

func (start *Node_first) add_s(blank []bool, Distance float64, data []float64) { //??这里引用改变原来的没，在某个分片里加
	if blank[0] == false { //对某个first第一次调用
		start.Distance = Distance
		//fmt.Println("startdis添加:", start.Distance)
		blank[0] = true //以后都是1
		start.Data = data
		// fmt.Print("startdta:")
		// fmt.Print(start.Data)
		return
	} else {
		if Distance < start.Distance { //小于头节点
			var buf1 = start.Distance
			var buf2 = start.Data
			//fmt.Println("startDistance更换:", start.Distance, "->", Distance)
			start.Distance = Distance

			start.Data = data
			// fmt.Print("startdta:")
			// fmt.Print(start.Data)
			start.add_s(blank, buf1, buf2) //递归
			return
		}
		var index *Node = start.Next
		var i = 0
		for ; index != nil; index = index.Next {
			i++
			if blank[i] == false { //节点没满
				index.Distance = Distance
				//fmt.Println("indexDistance添加:", index.Distance)
				// fmt.Print("indexdta:")
				// fmt.Print(index.Data)
				blank[i] = true
				index.Data = data
				return
			} else {
				if Distance < index.Distance {
					var buf1 = index.Distance
					var buf2 = index.Data
					//fmt.Println("indexDistance更换:", index.Distance, "->", Distance)
					index.Distance = Distance

					index.Data = data
					// fmt.Print("indexdta:")
					// fmt.Print(index.Data)
					start.add_s(blank, buf1, buf2)
					return
				} else {
					continue
				}
			}
		}
	}
}

// 搜索，添加删除，训练
// ！=nil试验一下
// 创建结果空数组
// 遍历每个向量
// 	暴力或者faiss:
// 		遍历每个数据库向量，通过Distance都用nodefirst.add一次（这里把Distance覆盖了，
// 		因为之前Distance代表距离类中心距离，现在我查询后不需要这个，
// 		就让他等于我第n个向量和db的距离然后用这个距离比较实现add核心（没满，满了）
// 	faiss:
// 	v1:对每个要找的向量，搜索最近的nodefirst,记录下index，通过根节点find找到这个first,start add所有的
// 	v2 1.不add所有的，add距离中心近的 2.（对距中心点比较近的点也搜索一些）

// 添加，找最近的first,计算距离并逐个比较插入
// 删除 find函数找到第i个first然后找到第j个向量
// 训练 递归结束条件 ：
// 计算平均值，找到最近的node，交换node和nodefirst的值，这个时候只剩下nodefist了，用train改变DB的排序
func Delete(req Request) string {
	var nodefirst = DB.Head.find_nodefirst(int(req.I_delete))
	return nodefirst.delete_node(int(req.J_delete))
}
func (start *Node_first) delete_node(i int) string { //找到第i个节点,中间末尾都分为上一个是first还是node，
	if i == 0 {
		var st string = "你要的是nodefirst诶"
		return st
	} else {
		if i == 1 {
			if start.Next.Next == nil {
				start.Next = nil
			} else {
				start.Next = start.Next.Next
			}
			//可以是nil
			var st string = "删除成功"
			return st
		}
		var j int
		var index *Node
		var preindex = index
		for index = start.Next; j < i-1; index = index.Next {
			if index == nil {
				var st string = "很明显超范围了first的那一堆没这么多向量"
				return st
			}
			j++
			preindex = index
		}
		if preindex.Next.Next == nil {
			preindex.Next = nil
		} else {
			preindex.Next = preindex.Next.Next
		}
	}
	var st string = "删除成功"
	return st
}
func Add(req Request) string {

	for i := 0; i < len(req.Vectors.Vector); i++ {
		if DB.Head == nil {
			var head Node_first
			head.Data = req.Vectors.Vector[0]
			DB.Head = &head
			continue
		} else if Length_first() < 3 { //10规定了一开始空库的时候add的时候聚类中心的数量
			var newnodefirst Node_first
			newnodefirst.Data = req.Vectors.Vector[i]
			DB.Head.find_nodefirst(Length_first()).Next_first = &newnodefirst
			continue
		} else {
			var min_index int
			var Distance1 float64
			var j int
			for indexi := DB.Head; indexi != nil; indexi = indexi.Next_first { //对每个向量遍历所有数据库
				if indexi == DB.Head { //初始化
					Distance1 = Distance(indexi.Data, req.Vectors.Vector[i])
					j++
					min_index = j
					continue
				} else if Distance(req.Vectors.Vector[i], indexi.Data) < Distance1 {
					Distance1 = Distance(req.Vectors.Vector[i], indexi.Data)
					min_index = j + 1
				}
				j++
			}
			var nodefirst = DB.Head.find_nodefirst(min_index)
			var node Node
			node.Data = req.Vectors.Vector[i]
			node.Distance = Distance1
			nodefirst.add(&node)
		}

	}
	var st string = "添加成功"
	return st
}
func Length_first() int {
	var i int
	for index := DB.Head; index != nil; index = index.Next_first {
		i++
	}
	return i
}

func Search(req Request) Result { //第几个分片
	var result Result
	var start Node_first = result_init(req.Topk, len(req.Vectors.Vector), len(req.Vectors.Vector[0])) //创建结果
	var start_copy = &start
	var Distance1 float64

	// for index6 := start_copy; index6 != nil; index6 = index6.Next_first {
	// 	fmt.Print(1)
	// }
	for i := 0; i < len(req.Vectors.Vector); i++ { //遍历每个向量
		var blank = make([]bool, start.length())
		if req.Request_type == 1 { //暴力

			for indexi := DB.Head; indexi != nil; indexi = indexi.Next_first { //对每个向量遍历所有数据库
				//fmt.Print(Distance(req.Vectors.Vector[i], indexi.Data))
				start_copy.add_s(blank, Distance(req.Vectors.Vector[i], indexi.Data), indexi.Data)
				for indexj := indexi.Next; indexj != nil; indexj = indexj.Next {
					Distance1 = Distance(req.Vectors.Vector[i], indexj.Data)
					//fmt.Print(Distance1)
					start_copy.add_s(blank, Distance1, indexj.Data)
				}
				// if Distance(req.Vectors.Vector[i], indexi.Data) > max {
				// 	max = Distance(req.Vectors.Vector[i], indexi.Data)
				// 	result.VectorGroup[i][j] = indexi.Data //第一个数据特殊处理
				// 	result.TopkDistance[i]
				// }
				// for indexj := indexi.Next; indexj.Next != nil; indexj = indexj.Next {
				// 	//Distance(req.Vectors.Vector[i],

				// }
			}

		} else if req.Request_type == 2 { //用kmean// 	v1:对每个要找的向量，
			//搜索最近的nodefirst,记录下index，通过根节点find找到这个first,start add所有的
			var min_index int
			var Distance1 float64
			var j int
			for indexi := DB.Head; indexi != nil; indexi = indexi.Next_first { //对每个向量遍历所有数据库
				if indexi == DB.Head { //初始化
					Distance1 = Distance(indexi.Data, req.Vectors.Vector[i])
					j++
					min_index = j
					continue
				} else if Distance(req.Vectors.Vector[i], indexi.Data) < Distance1 {
					Distance1 = Distance(req.Vectors.Vector[i], indexi.Data)
					min_index = j + 1
				}
				j++
			}
			var index_nf = DB.Head.find_nodefirst(min_index)
			start_copy.add_s(blank, Distance(index_nf.Data, req.Vectors.Vector[i]), index_nf.Data)
			for index_node := index_nf.Next; index_node != nil; index_node = index_node.Next {
				start_copy.add_s(blank[:], Distance(index_node.Data, req.Vectors.Vector[i]), index_node.Data)
			}

		}
		start_copy = start_copy.Next_first
	}

	// for indexi := &start; indexi != nil; indexi = indexi.Next_first {
	// 	fmt.Print(indexi.Distance)
	// 	for indexj := indexi.Next; indexj != nil; indexj = indexj.Next {
	// 		fmt.Print(indexj.Distance)
	// 	}
	// }
	ResuToint(&start, &result) //结果返回整形
	return result
}

func ResuToint(start *Node_first, result *Result) {
	var i int //第i+1个向量,行
	var j int = 1
	for indexi := start; indexi != nil; indexi = indexi.Next_first {
		//fmt.Print(indexi.Distance)
		//fmt.Print("一列")
		j = 1
		var topk_ls []float64
		var group1 [][]float64
		var group2 []float64
		topk_ls = append(topk_ls, indexi.Distance) //first特殊处理
		for k := 0; k < len(indexi.Data); k++ {
			group2 = append(group2, indexi.Data[k])
		}
		group1 = append(group1, group2)
		group2 = nil
		//fmt.Print(group1)

		i++

		for indexj := indexi.Next; indexj != nil; indexj = indexj.Next {
			//fmt.Print(indexj.Distance)
			// fmt.Print("一行")
			topk_ls = append(topk_ls, indexj.Distance)
			for k := 0; k < len(indexj.Data); k++ {
				group2 = append(group2, indexj.Data[k])
			}
			group1 = append(group1, group2)
			group2 = nil
			j++
		}
		result.TopkDistance = append(result.TopkDistance, topk_ls)
		topk_ls = nil
		result.VectorGroup = append(result.VectorGroup, group1)
		group1 = nil
	}

}

// for i := 0; nodefirst_fuben.Next_first != nil; i++ {
// 	var min int = 10000 //一个较大的值，否则需要一个bool每次都得判断
// 	if min > Distance(nodefirst.Data, a) {
// 		min = Distance(nodefirst.Data, a)
// 		lie = i
// 	}
// 	nodefirst_fuben = nodefirst_fuben.Next_first
// } //可以考虑在这一列附近的列，但是可能要保存列之间的距离
// index_nodefirst = head_nodefirst.find_nodefirst(index + 1)
// newnode.Distance = min
// index_nodefirst.add(newnode)

// start_fuben = *start_fuben.Next_first

// func check(a [][]int, i int, j int) {

// }
// func (startnode Newnode_first)xianruman(start *Node_first, head_db Node_first, req Request) {//结果，数据库，输入向量组

// 	// var buf [][]int
// 	// buf[0] = head_db.Data
// 	// for i := range buf {
// 	// 	if index.Next != nil {
// 	// 		dis
// 	// 		buf[i+1] = index.Data
// 	// 		index = index.Next
// 	// 	}
// 	// 	if i < topk-1 && index.Next == nil {
// 	// 		head_db = head_db.Next_first
// 	// 		i = i + 1
// 	// 		buf[i] = head_db.Data
// 	// 		index = head_db.Next
// 	// 	}
// 	// }
// 	for ; start.Next_first != nil; start = start.Next_first {
// 		start.Data = head_db.Data
// 		var buf *Node
// 		if head_db.Next != nil {
// 			buf = head_db.Next
// 		}
// 		for i := 0; i < req.Topk-1; i++ {
// 			start.add(*buf,i int,req Request)
// 			if buf.Next != nil {
// 				buf = buf.Next
// 			}
// 			if buf.Next == nil && i < req.Topk {
// 				head_db = *head_db.Next_first
// 				var newnode Node
// 				newnode.Data = head_db.Data
// 				start.add(newnode)
// 				buf = head_db.Next

// 			}

// 		}

// 	}

// }

// func transpose(A [][]int) [][]int {
// 	B := make([][]int, len(A[0]))
// 	for i := 0; i < len(A[0]); i++ {
// 		B[i] = make([]int, len(A))
// 		for j := 0; j < len(A); j++ {
// 			B[i][j] = A[j][i]
// 		}
// 	}
// 	return B
// }
