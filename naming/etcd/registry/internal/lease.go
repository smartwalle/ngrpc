package internal

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Lease struct {
	client *clientv3.Client
	id     clientv3.LeaseID
	cancel context.CancelFunc
}

func NewLease(ctx context.Context, client *clientv3.Client, ttl int64) (*Lease, error) {
	var leaseID clientv3.LeaseID
	ctx, cancel := context.WithCancel(ctx)
	if ttl > 0 {
		grantRsp, err := client.Grant(ctx, ttl)
		if err != nil {
			cancel()
			return nil, err
		}
		leaseID = grantRsp.ID

		keepAliveRsp, err := client.KeepAlive(ctx, leaseID)
		if err != nil || keepAliveRsp == nil {
			cancel()
			return nil, err
		}

		go func() {
			for range keepAliveRsp {
			}
		}()
	}

	var s = &Lease{}
	s.client = client
	s.id = leaseID
	s.cancel = cancel
	return s, nil
}

func (s *Lease) ID() clientv3.LeaseID {
	return s.id
}

func (s *Lease) Revoke(ctx context.Context) error {
	s.cancel()
	_, err := s.client.Revoke(ctx, s.id)
	return err
}
