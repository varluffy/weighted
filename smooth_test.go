/**
* Created by GoLand.
* User: luffy
* Date: 2019-07-03
* Time: 19:04
 */
package weighted

import (
	"fmt"
	"testing"
)

func Test_SWRR(t *testing.T) {
	s := &SWRR{}
	s.Add(1, 30)
	s.Add(2,40)
	s.Add(3,30)
	r := make(map[int]int)
	for i:=0; i < 100; i ++ {
		ss := s.Next().(int)
		r[ss] ++
	}
	fmt.Printf("SWRR = %#v", r)
}