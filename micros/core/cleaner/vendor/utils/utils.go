package utils

import (
	"constants"

	"log"

	//http2curl "moul.io/http2curl"
	"context"

	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func CleanerControl(key string) string {
	result := "OK"
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{constants.DOCKER_GW_REGISTER_IP + ":2379"},
		DialTimeout: 5 * time.Second,
	})

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := cli.Delete(ctx, key)

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

func CleanerLogin(input string) string {

	return CleanerControl("public" + input)
}
