/**
* Created by GoLand.
* User: luffy
* Date: 2019-07-03
* Time: 11:35
 */
/**
* Nginx基于权重的轮询算法的实现 Upstream: smooth weighted round-robin balancing.
* https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
* https://blog.csdn.net/zhangskd/article/details/50194069
* https://lihaoquan.me/2018/4/24/go-reverse-proxy.html
 */

package weighted

type SWRR struct {
	nodes []*smoothWRR
	size int
}
type smoothWRR struct {
	Item interface{}
	Weight int //自身权重
	EffectiveWeight int //有效权重，初始值为weight，如果发现和后端通信失败时，就减小effectiveWeight,此后再有新的请求过来，选取后端的过程时再逐步增加effectiveWeight,最终恢复到weight
	CurrentWeight int //当前权重，初始化为0
}

func (s *SWRR) Add (item interface{}, weight int) {
	smooth := &smoothWRR{Item:item, Weight:weight, EffectiveWeight:weight}
	s.nodes = append(s.nodes, smooth)
	s.size ++
}

func (s *SWRR) Next() (item interface{}) {
	if s.size == 0 {
		return nil
	}
	if s.size == 1 {
		return s.nodes[0].Item
	}
	return nextSWRR(s.nodes).Item
}

func nextSWRR(nodes []*smoothWRR) (best *smoothWRR) {
	total := 0
	size := len(nodes)
	for i :=0 ; i < size; i ++ {
		node := nodes[i]
		if node == nil {
			continue
		}
		node.CurrentWeight += node.EffectiveWeight
		total += node.EffectiveWeight
		if node.EffectiveWeight < node.Weight {
			node.EffectiveWeight ++
		}
		if best == nil || node.CurrentWeight > best.CurrentWeight {
			best = node
		}
	}
	if best == nil {
		return nil
	}
	best.CurrentWeight -= total
	return best
}

func(s *SWRR) Reset() {
	for _, node := range s.nodes {
		node.EffectiveWeight = node.Weight
		node.CurrentWeight = 0
	}
}
func(s *SWRR) RemoveAll() {
	s.nodes = s.nodes[:0]
	s.size = 0
}