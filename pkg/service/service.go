package service

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

type Service struct {
	LeaseId    clientv3.LeaseID
	Addr       string
	Prefix     string
	Name       string
	Key        string
	etcdClient *clientv3.Client
	stop       chan error
}

func NewService(config clientv3.Config) (*Service, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		log.Fatalf("err new etcd client: %v", err)
		return nil, err
	}
	//defer cli.Close()

	service := &Service{
		etcdClient: cli,
	}

	return service, nil
}

func (s *Service) Register(Prefix string, serviceName string, addr string) error {
	s.Prefix = Prefix
	s.Name = serviceName
	s.Addr = addr
	s.Key = getRandString(10)

	lease, err := s.etcdClient.Grant(context.TODO(), 10)
	if err != nil {
		return err
	}
	s.LeaseId = lease.ID

	_, err = s.etcdClient.Put(context.TODO(), s.Prefix+s.Name+s.Key, addr, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	go s.KeepAlive()

	return nil
}

func (s *Service) UnRegister() error {
	_, err := s.etcdClient.Revoke(context.TODO(), s.LeaseId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(path string) ([]string, error) {
	nodes := []string{}
	resp, err := s.etcdClient.Get(context.TODO(), path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		nodes = append(nodes, string(kv.Value))
	}
	return nodes, nil
}

func (s *Service) Watch(path string, callback func()) {
	watchChan := s.etcdClient.Watch(context.Background(), path, clientv3.WithPrefix())
	for wresp := range watchChan {
		for range wresp.Events {
			callback()
		}
	}
}

func (s *Service) KeepAlive() {
	for {
		select {
		case <-s.stop:
			return
		case <-s.etcdClient.Ctx().Done():
			return
		default:
			if _, err := s.etcdClient.KeepAliveOnce(context.TODO(), s.LeaseId); err == rpctypes.ErrLeaseNotFound {
				log.Println("error when KeepAliveOnce")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func getRandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
