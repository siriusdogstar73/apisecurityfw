package utils

import (
	"constants"
	"fmt"

	"log"

	"encoding/json"

	//http2curl "moul.io/http2curl"
	"context"

	"time"

	"interfaces"

	clientv3 "go.etcd.io/etcd/client/v3"
)

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

func initEtcdClientWirhControl(key string, value string) string {
	result := "OK"
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
	if err == nil {
		result = "OK"
	} else {
		result = "KO"
	}
	return result

}

func searchEtcdClientPublic(input string) string {
	var k string
	var v string
	key := "public" + input
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	resp, _ := cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))

	for _, ev := range resp.Kvs {
		k = fmt.Sprintf("%s", ev.Key)
		log.Println("encontrado: ", k)
		v = fmt.Sprintf("%s", ev.Value)
	}
	return v
}

func SearchProcess(sHash string) string {
	jsonString := SearchEtcd(sHash)

	return jsonString
}

func SearchEtcd(key string) string {
	log.Println("key: ", key)

	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	response, err := cli.Get(ctx, key)

	sum := 1
	for sum < 100 {
		sum += sum
		response, err = cli.Get(ctx, key)
	}
	log.Println("response.Count: ", response.Count)
	if err != nil {
		log.Println("err: ", err.Error())
	}
	var res string
	for _, ev := range response.Kvs {
		log.Println("ev.Value: ", ev.Value)
		res = fmt.Sprintf("%s\n", ev.Value)
		log.Println("res: ", res)
	}

	return ""

}

func SaveServerKeys(input string) string {

	sPublicKeyRequest := interfaces.PublicKeyRequest{}
	err := json.Unmarshal([]byte(input), &sPublicKeyRequest)
	if err != nil {
		log.Println(err)
	}
	return initEtcdClientWirhControl("public"+sPublicKeyRequest.Uuid, input)
}

func SearchServerKeys(input string) string {

	return searchEtcdClientPublic(input)
}
