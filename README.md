#服务启动

#node1
./bin/etcd -name p1 -debug \
-heartbeat-interval  200 \
-initial-advertise-peer-urls http://0.0.0.0:2380 \
-listen-peer-urls http://0.0.0.0:2380 \
-listen-client-urls http://0.0.0.0:2379 \
-advertise-client-urls http://0.0.0.0:2379 \
-initial-cluster-token etcd-cluster-1 \
-initial-cluster p1=http://0.0.0.0:2380,p2=http://0.0.0.0:4380,p3=http://0.0.0.0:6380 \
-initial-cluster-state new

#node2
./bin/etcd -name p2 -debug \
-heartbeat-interval  200 \
-initial-advertise-peer-urls http://0.0.0.0:4380 \
-listen-peer-urls http://0.0.0.0:4380 \
-listen-client-urls http://0.0.0.0:4379 \
-advertise-client-urls http://0.0.0.0:4379 \
-initial-cluster-token etcd-cluster-1 \
-initial-cluster p1=http://0.0.0.0:2380,p2=http://0.0.0.0:4380,p3=http://0.0.0.0:6380 \
-initial-cluster-state new

#node3
./bin/etcd -name p3 -debug \
-heartbeat-interval  200 \
-initial-advertise-peer-urls http://0.0.0.0:6380 \
-listen-peer-urls http://0.0.0.0:6380 \
-listen-client-urls http://0.0.0.0:6379 \
-advertise-client-urls http://0.0.0.0:6379 \
-initial-cluster-token etcd-cluster-1 \
-initial-cluster p1=http://0.0.0.0:2380,p2=http://0.0.0.0:4380,p3=http://0.0.0.0:6380 \
-initial-cluster-state new