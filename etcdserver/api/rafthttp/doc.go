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
//  |  |  |-etcdhttp.NewPeerHandler/peer.go                   创建http-handler，对外接收处理peer http请求
//  |  |-rafthttp                                             rafthttp包
//  |  |  |-stream                                   	      HTTP长连接，主要负责传输数据量较小、发送比较频繁的消息，例如，MsgApp消息、MsgHeartbeat消息、MsgVote消息等。
//  |  |  |-pipeline                                          Pipeline的消息通道在传输数据完成后立即关闭连接，主要负责传输数据量较大、发送频率较低的消息，例如，MsgSnap消息等。
