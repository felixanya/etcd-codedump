// Copyright 2016 The etcd Authors
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

// Package api manages the capabilities and features that are exposed to clients by the etcd cluster.
package api

// etcd的网络服务有grpc和http两类，其中 grpc是etcd的server，对外接收客户端的put、del等请求
// http是raft的服务，符合raft节点的信息交换
// http: (1)etcdhttp属于server端
//     	 (2)rafthttp属于client端

// grpc                                                             api/v3rpc/grpc.go 服务在这里注册
//  |-etcdserver.Server()                                           返回grpcserver,执行网络监听
//  |  |-pb.RegisterKVServer()                                      注册KV服务
//  |  |-pb.RegisterWatchServer()                                   注册Watch服务
//  |  |-pb.RegisterLeaseServer()                                   注册LEASE服务
//  |  |-pb.RegisterClusterServer()                                 注册Cluster服务
//  |  |-pb.RegisterAuthServer()                                    注册Auth服务
//  |  |-pb.RegisterMaintenanceServer()                             注册Maintenance服务

// rafthttp
// rafthttp外层封装了etcdhttp，用做server端的服务监听

//  |-api
//  |  |-etcdhttp
//  |  |  |-etcdhttp.NewPeerHandler/peer.go                   创建http-handler，对外接收处理peer http请求
//  |  |-rafthttp                                             rafthttp包，负责发送请求
//  |  |  |-stream                                   	      HTTP长连接，主要负责传输数据量较小、发送比较频繁的消息，例如，MsgApp消息、MsgHeartbeat消息、MsgVote消息等。
//  |  |  |-pipeline                                          Pipeline的消息通道在传输数据完成后立即关闭连接，主要负责传输数据量较大、发送频率较低的消息，例如，MsgSnap消息等。
// rafthttp/doc.go写了一部分注释，具体分析看那里

// 注：
// v3_server.go是raftkv inteface，在EtcdServer上添加的func，对server.go的补充

// grpc以put请求为例
// quotaKVServer.Put()                             api/v3rpc/quota.go 首先检查是否满足需求
//  |-quotoAlarm.check()                           检查
//  |-KVServer.Put()                               api/v3rpc/key.go 真正的处理请求
//    |-checkPutRequest()                          校验请求参数是否合法
//    |-RaftKV.Put()                               etcdserver/v3_server.go 处理请求
//    |-EtcdServer.Put()                           实际调用的是该函数
//    | |-raftRequest()
//    |   |-raftRequestOnce()
//    |     |-processInternalRaftRequestOnce()     真正开始处理请求
//    |       |-context.WithTimeout()              创建超时的上下文信息
//    |       |-raftNode.Propose()                 raft/node.go
//    |         |-raftNode.step()                  对于类型为MsgProp类型消息，向propc通道中传入数据
//    |-header.fill()                              etcdserver/api/v3rpc/header.go填充响应的头部信息

// 网络服务的启动是在 embed.servePeers()中进行的
//
