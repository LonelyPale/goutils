package ipfscluster

import (
	"context"
	"fmt"

	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs-cluster/ipfs-cluster/api/rest/client"
	ma "github.com/multiformats/go-multiaddr"
)

type Config struct {
	ID            string `value:"${ipfs-cluster.id:=}"`
	Host          string `value:"${ipfs-cluster.host:=localhost}"`
	Port          string `value:"${ipfs-cluster.port:=9094}"`
	ProxyAddrHost string `value:"${ipfs-cluster.proxy_addr_host:=localhost}"`
	ProxyAddrPort string `value:"${ipfs-cluster.proxy_addr_port:=5001}"`
}

type Client struct {
	client.Client
}

func NewClient(conf *Config) (*Client, error) {
	nodeMAddr := ma.StringCast(fmt.Sprintf("/ip4/%s/tcp/%s", conf.ProxyAddrHost, conf.ProxyAddrPort))
	cfg := &client.Config{
		Host:      conf.Host,
		Port:      conf.Port,
		ProxyAddr: nodeMAddr,
	}

	cli, err := client.NewDefaultClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{cli}, nil
}

func (c *Client) Pin(cid string) (api.Pin, error) {
	ctx := context.Background()

	ci, err := api.DecodeCid(cid)
	if err != nil {
		return api.Pin{}, err
	}

	return c.Client.Pin(ctx, ci, api.DefaultAddParams().PinOptions)
}

// MustPin todo
//func (c *Client) MustPin(cid string) (_ api.Pin, err error) {
//	ctx := context.Background()
//
//	ci, err := api.DecodeCid(cid)
//	if err != nil {
//		return api.Pin{}, err
//	}
//
//	pin, err := c.Client.Pin(ctx, ci, api.DefaultAddParams().PinOptions)
//	if err != nil {
//		return api.Pin{}, err
//	}
//
//	pinInfo, err := c.Status(ctx, ci, false)
//	if err != nil {
//		return api.Pin{}, err
//	}
//
//	pinnedNum := 0
//	peerNum := len(pinInfo.PeerMap)
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		done := make(chan struct{})
//		for {
//			select {
//			case <-done:
//				return
//			default:
//				var info api.GlobalPinInfo
//				info, err = c.Status(ctx, ci, false)
//				if err != nil {
//					return
//				}
//
//				for _, peer := range info.PeerMap {
//					switch peer.Status {
//					case api.TrackerStatusPinning:
//						continue
//					case api.TrackerStatusPinned:
//					default:
//						wg.Done()
//					}
//				}
//			}
//		}
//
//	}()
//}
