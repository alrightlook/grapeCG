// 2017-07-24
// 测试网络层基本收法行为
package main

import (
	gnet "github.com/koangel/grapeNet/Net"
)

func main() {
	vnet := gnet.NewEmptyTcp() // 创建一个空的对象

	vnet.Runnable()
}
