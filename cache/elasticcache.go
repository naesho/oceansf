package cache

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ohsaean/oceansf/lib"
	log "github.com/sirupsen/logrus"
)

var (
	Mcache *Client
)

// Node is a single ElastiCache node
type Node struct {
	URL  string
	Host string
	IP   string
	Port int
}

// Item embeds the memcache client's type of the same name
type Item memcache.Item

// Client embeds the memcache client so we can hide those details away
type Client struct {
	*memcache.Client
}

func (c *Client) Get(key string) (*memcache.Item, error) {
	// TODO performance check
	return c.Client.Get(key)
}

func (c *Client) Delete(key string) error {
	return c.Client.Delete(key)
}

func (c *Client) Set(item *memcache.Item) error {
	return c.Client.Set(item)
}

func (c *Client) Add(item *memcache.Item) error {
	return c.Client.Add(item)
}

func (c *Client) Cas(item *memcache.Item) error {
	return c.Client.CompareAndSwap(item)
}

func (c *Client) CasDelayed(item *Item) error {
	// queue item --> when db transaction success --> flush and do CompareAndSwap operation
}

// New returns an instance of the memcache client
func NewMemcache(endpoint string) (*Client, error) {
	urls, err := clusterNodes(endpoint) // cluster endpoint ex) mycluster.fnjyzo.cfg.use1.cache.amazonaws.com
	if err != nil {
		return &Client{Client: memcache.New()}, err
	}

	return &Client{Client: memcache.New(urls...)}, nil
}

// Initialize instance for using main package
func InitMemcache(endpoint string) {
	var err error
	Mcache, err = NewMemcache(endpoint)
	log.Debug("memcached init")
	lib.CheckError(err)
}

func clusterNodes(endpoint string) ([]string, error) {

	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		log.Debugf("Socket Dial (%s): %s", endpoint, err.Error())
		return nil, err
	}
	defer conn.Close()

	command := "config get cluster\r\n"
	fmt.Fprintf(conn, command)

	response, err := parseNodes(conn)
	if err != nil {
		return nil, err
	}

	urls, err := parseURLs(response)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func parseNodes(conn io.Reader) (string, error) {
	var response string

	count := 0
	location := 3 // AWS docs suggest that nodes will always be listed on line 3

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		count++
		if count == location {
			response = scanner.Text()
		}
		if scanner.Text() == "END" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Scanner: ", err.Error())
		return "", err
	}

	log.Debug("ElastiCache nodes found: ", response)
	return response, nil
}

func parseURLs(response string) ([]string, error) {
	var urls []string
	var nodes []Node

	items := strings.Split(response, " ")

	for _, v := range items {
		fields := strings.Split(v, "|") // ["host", "ip", "port"]

		port, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Error("Integer conversion: ", err.Error())
			return nil, err
		}

		node := Node{fmt.Sprintf("%s:%d", fields[1], port), fields[0], fields[1], port}
		nodes = append(nodes, node)
		urls = append(urls, node.URL)

		log.Debugf("Host: %s, IP: %s, Port: %d, URL: %s", node.Host, node.IP, node.Port, node.URL)
	}

	return urls, nil
}
