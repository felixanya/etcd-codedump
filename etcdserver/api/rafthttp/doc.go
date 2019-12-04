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

// Package rafthttp implements HTTP transportation layer for etcd/raft pkg.
package rafthttp

// rafthttp
// rafthttp外层封装了etcdhttp，用做server端的服务监听
// rafthttp本身作为client端使用，链接其他peers
//  |-api
//  |  |-etcdhttp
//  |  |  |-etcdhttp.NewPeerHandler/peer.go                   创建http-handler，对外接收处理peer 请求
//  |  |-rafthttp                                             rafthttp包
//  |  |  |-stream                                   	      HTTP长连接，主要负责传输数据量较小、发送比较频繁的消息，例如，MsgApp消息、MsgHeartbeat消息、MsgVote消息等。
//  |  |  |-pipeline                                          Pipeline的消息通道在传输数据完成后立即关闭连接，主要负责传输数据量较大、发送频率较低的消息，例如，MsgSnap消息等。

// raft-http几个重要strcut
// stream: stream是raft节点之间信息交流的结构体，通过ServeHTTP方法拿到http conection,然后attach到outgoingcon，不断write resp.Body（或者 read resp.Body），提高通信性能
// pipeline: stream不可用的时候，会使用pipeline

// raft节点之间的消息传递并不是简单的request-response模型，而是读写分离模型，
// 即每两个server之间会建立两条链路，对于每一个server来说，一条链路专门用来发送数据，另一条链路专门用来接收数据.
// 在代码实现中，通过streamWriter发送数据，通过streamReader接收数据。即通过streamReader接收数据接收到数据后会直接响应，在处理完数据后通过streamWriter将响应发送到对端
// 对于每个server来说，不管是leader、candicate还是follower，都会维持一个peers数组，每个peer对应集群中的一个server，负责处理server之间的一些数据交互。
// 当server需要向其他server发送数据时，只需要找到其他server对应的peer，然后向peer的streamWriter的msgc通道发送数据即可，streamWriter会监听msgc通道的数据并发送到对端server；
// 而streamReader会在一个goroutine中循环读取对端发送来的数据，一旦接收到数据，就发送到peer的p.propc或p.recvc通道，而peer会监听这两个通道的事件，写入到node的n.propc或n.recvc通道，node只需要监听这两个通道的数据并处理即可。这就是在etcd的raft实现中server间数据交互的流程。
