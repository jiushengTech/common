package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	Watcher clientv3.Watcher
}

var ETCD Etcd

func NewEtcdUtil(e Etcd) {
	ETCD = e
}

func GetKeyValue(key string) error {
	getResp, err := ETCD.Kv.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
	}

	// 遍历所有任务, 进行反序列化
	for _, kvPair := range getResp.Kvs {
		fmt.Println(kvPair)
	}
	return err
}

func PutKeyValue(key, value string) error {
	// 保存到etcd
	putResp, err := ETCD.Kv.Put(context.TODO(), key, value, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	// 如果是更新, 那么返回旧值
	if putResp.PrevKv != nil {
		fmt.Println(putResp.PrevKv.Value)
	}
	return err
}

func DeleteKey(key string) error {
	// 从etcd中删除它
	delResp, err := ETCD.Kv.Delete(context.TODO(), key, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	// 返回被删除的值
	if len(delResp.PrevKvs) != 0 {
		fmt.Println(delResp.PrevKvs[0].Value)
	}
	return err
}

func PutKeyWithLease(key string, timeout int64) (err error) {

	leaseGrantResp, err := ETCD.Lease.Grant(context.TODO(), timeout)
	if err != nil {
		return
	}

	// 租约ID
	leaseId := leaseGrantResp.ID
	_, err = ETCD.Kv.Put(context.TODO(), key, "", clientv3.WithLease(leaseId))
	if err != nil {
		return
	}
	return
}

func GetValueByKey(key string) (string, error) {
	getResp, err := ETCD.Kv.Get(context.TODO(), key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 如果找到匹配的键，返回其值
	if len(getResp.Kvs) > 0 {
		return string(getResp.Kvs[0].Value), nil
	}

	return "", nil // 如果未找到匹配的键，返回空字符串
}
