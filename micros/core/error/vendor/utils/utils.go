package utils

import (
	"constants"

	"log"

	//http2curl "moul.io/http2curl"
	"context"
	"encoding/json"

	"time"

	"interfaces"

	"strings"

	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func ErrorMngmnt(input string) string {

	var ret string
	resGeneralErrorKeyValuePayload := interfaces.ResGeneralErrorKeyValuePayload{}
	if !strings.Contains(input, "audit") {

		err := json.Unmarshal([]byte(input), &resGeneralErrorKeyValuePayload)
		if err != nil {
			ret = searchEtcd("error" + input)
		} else {
			if resGeneralErrorKeyValuePayload.Value.Code != constants.EMPTY {
				initEtcdClient("error"+resGeneralErrorKeyValuePayload.Key, input)
				ret = "OK"
			}
		}
	} else {
		res := strings.Split(input, "audit")
		SetAuditEtcdClient(res[1])
	}

	return ret
}

func searchEtcd(sUuid string) string {
	sUuid = strings.ReplaceAll(sUuid, "\n", "")

	ret := ""

	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	var k string
	var v string

	//Search with prefix
	resp, _ := cli.Get(
		ctx,
		sUuid,
		clientv3.WithPrefix(),
		clientv3.WithSort(
			clientv3.SortByKey,
			clientv3.SortDescend))

	for _, ev := range resp.Kvs {
		k = fmt.Sprintf("%s", ev.Key)
		log.Println("encontrado: ", k)
		v = fmt.Sprintf("%s", ev.Value)
		if !strings.Contains(k, sUuid) {
			ret = ""
		} else {
			ret = v
		}
	}
	return ret
}

func initEtcdClient(key string, value string) {
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := cli.Put(ctx, key, value)

	if err != nil {
		switch err {
		case context.Canceled:
			log.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Println("ctx is attached with a deadline is exceeded: %v", err)

		default:
			log.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	/*
		response, err := cli.Get(ctx, key)

		for _, ev := range response.Kvs {
			log.Printf("%s : %s\n", ev.Key, ev.Value)
		}
	*/

}

func SetAuditEtcdClient(key string) {

	value := searchEtcd("error" + key)
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := cli.Put(ctx, "audit"+key, value)

	if err != nil {
		switch err {
		case context.Canceled:
			log.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Println("ctx is attached with a deadline is exceeded: %v", err)

		default:
			log.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	_, err = cli.Delete(ctx, "error"+key)
	if err != nil {
		switch err {
		case context.Canceled:
			log.Println("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Println("ctx is attached with a deadline is exceeded: %v", err)

		default:
			log.Println("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}

}
