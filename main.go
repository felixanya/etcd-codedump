// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is a simple wrapper of the real etcd entrypoint package
// (located at go.etcd.io/etcd/etcdmain) to ensure that etcd is still
// "go getable"; e.g. `go get go.etcd.io/etcd` works as expected and
// builds a binary in $GOBIN/etcd
//
// This package should NOT be extended or modified in any way; to modify the
// etcd binary, work in the `go.etcd.io/etcd/etcdmain` package.
//
package main

import (
	"go.etcd.io/etcd/etcdmain"
)

// etcd服务入口：
// 1. etcdman加载服务启动配置
// 2. etcdman/etcd.go: startEtcd   ==>   embed/etcd.go: StartEtcd
// 3. etcd的proxy在etcdman中实现，对下层服务无感；StartEtcd 是服务真正启动的入口
func main() {
	etcdmain.Main()
}

// ***** Etcd服务结构

// main()                                etcdmain/main.go
//  |-checkSupportArch()
//  |-startEtcdOrProxyV2()               etcdmain/etcd.go
//    |-newConfig()                      etcd配置参数
//    |-startEtcd() —>
//    | |-embed.StartEtcd() —>           embed/etcd.go
//    |   |-configurePeerListeners()     listenners
//    |   |-configureClientListeners()   listenners
//    |   |-EtcdServer.ServerConfig()    生成新的配置
//    |   |
//    |   |-EtcdServer.NewServer()       etcdserver/server.go 正式启动RAFT服务<<<1>>>
//    |   |
//    |   |-EtcdServer.Start() ->        etcdserver/server.go 开始启动服务
//    |   | |-EtcdServer.start() ->
//    |   |   |-wait.New()               新建WaitGroup组以及一些管道服务
//    |   |
//    |   |   |-EtcdServer.run()         etcdserver/raft.go 启动应用层的处理协程<<<2>>>
//    |   |
//    |   |-Etcd.servePeers()            启动集群内部通讯
//    |   | |-etcdhttp.NewPeerHandler()  启动http服务
//    |   | |-v3rpc.Server()             启动gRPC服务 api/v3rpc/grpc.go，这里真正监听了frpc请求
//    |   |   |-grpc.NewServer()         调用gRPC的接口创建
//    |   |   |-pb.RegisterKVServer()    注册各种的服务，这里包含了多个
//    |   |   |-pb.RegisterWatchServer()
//    |   |
//    |   |-Etcd.serveClients()          启动协程处理客户请求
//    |   |
//    |   |-Etcd.serveMetrics()
//    |-notifySystemd()
//    |-select()                         等待stopped信号，及时return
//    |-osutil.Exit()

// 注:
// etcd/etcdserver/etcdserverpb/rpc/proto 声明了etcd所有的service
// 基础服务的入口
// service KV
// service Watch
// service Lease
// service Cluster
// service Maintenance
// service Auth

// 代码结构树
// etcd
//  |-bin                                  编译后的可执行文件
//  |-client                               客户端
//  |-clientv3                             客户端v3
//  |-contrib                              模块代码范例:raft
//  |-embed                                插件和扩展
//  |-etcdtl                               命令行工具
//  |-etcdman                              服务启动入口，加载服务参数
//  |-etcdserver                           核心，从逻辑上代表了一个完整的Etcd服务
//  |  |-v3server                          入口文件
//  |  |-etcdserverpb                      grpc描述文件
//  |  |   |-rpc.proto                     etcdservice描述主文件
//  |  |
//  |  |-api                               server服务接口/rafthttp接口

//  |-lease                                租约，计时器
//  |-mvcc                                 存储
//  |-raft                                 raft模块，只处理raft核心算法，不具备HTTP能力
//  |-wal                                  预写式日志

// etcd启动参数
// name 节点名称
// data-dir 指定节点的数据存储目录
// listen-peer-urls 监听URL，用于与其他节点通讯
// initial-advertise-peer-urls 该节点同伴监听地址，这个值会告诉集群中其他节点
// listen-client-urls 对外提供服务的地址：比如 http://ip:2379,http://127.0.0.1:2379 ，客户端会连接到这里和 etcd 交互
// advertise-client-urls 对外公告的该节点客户端监听地址，这个值会告诉集群中其他节点
// initial-cluster 集群中所有节点的信息，格式为 node1=http://ip1:2380,node2=http://ip2:2380,… 。注意：这里的 node1 是节点的 --name 指定的名字；后面的 ip1:2380 是 --initial-advertise-peer-urls 指定的值
// initial-cluster-state 新建集群的时候，这个值为 new ；假如已经存在的集群，这个值为 existing
// initial-cluster-token 创建集群的 token，这个值每个集群保持唯一。这样的话，如果你要重新创建集群，即使配置和之前一样，也会再次生成新的集群和节点 uuid；否则会导致多个集群之间的冲突，造成未知的错误
