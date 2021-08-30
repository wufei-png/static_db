package leaderelect

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

type Role int

const (
	Follower Role = 0
	Leader        = 1
)

type Status struct {
	LeaderID      string
	LeaderAddress string
	Role          Role
}

const (
	statusAcquiring = 0
	statusMaster    = 1
	statusResigned  = 2
)

const dbTimeout = 2 * time.Second

type Config struct {
	Hosts                  []string
	Keyspace               string
	TableName              string
	HeartbeatTimeoutSecond int
	AdvertiseAddress       string
	NodeID                 string
	ResourceName           string
}

func NewConfig(node string, resource string) *Config {
	return &Config{
		TableName:              "leader_elect",
		HeartbeatTimeoutSecond: 30,
		NodeID:                 node,
		ResourceName:           resource,
	}
}

type LeaderElector struct {
	config  *Config
	session *gocql.Session
	status  int32

	ev   chan Status
	dead chan bool
}

func NewLeaderElector(config *Config) (*LeaderElector, error) {
	if config.Keyspace == "" {
		panic("no keyspace")
	}
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.ConnectTimeout = dbTimeout
	cluster.Timeout = dbTimeout
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &LeaderElector{
		config:  config,
		session: session,
		ev:      make(chan Status, 1000),
		dead:    make(chan bool),
	}, nil
}

func NewLeaderElectorWithSession(session *gocql.Session, config *Config) *LeaderElector {
	if config.Keyspace == "" {
		panic("no keyspace")
	}
	return &LeaderElector{
		config:  config,
		session: session,
		ev:      make(chan Status, 1000),
		dead:    make(chan bool),
	}
}

func (l *LeaderElector) CurrentNodeID() string {
	return l.config.NodeID
}

func (l *LeaderElector) Start() {
	go l.elect()
}

func (l *LeaderElector) sleep(leader bool) {
	sec := l.config.HeartbeatTimeoutSecond
	if leader {
		sec = (l.config.HeartbeatTimeoutSecond - 1) / 2
	}
	deadline := time.Now().Add(time.Duration(sec) * time.Second)
	for {
		cur := atomic.LoadInt32(&l.status)
		if cur == statusResigned {
			break
		}
		now := time.Now()
		if deadline.Before(now) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (l *LeaderElector) elect() {
Loop:
	for {
		cur := atomic.LoadInt32(&l.status)
		switch cur {
		case statusAcquiring:
			status, err := l.createLease()
			if err != nil {
				log.Warn("failed to acquire lease: ", err)
				l.sleep(false)
				continue
			}
			l.ev <- status
			if status.Role == Leader {
				if !atomic.CompareAndSwapInt32(&l.status, cur, statusMaster) {
					continue
				}
			}
			l.sleep(status.Role == Leader)
		case statusMaster:
			status, err := l.updateLease()
			if err != nil {
				log.Warn("master failed to acquire lease, back to follower: ", err)
				status = Status{
					Role: Follower,
				}
			}
			l.ev <- status
			if status.Role != Leader {
				if !atomic.CompareAndSwapInt32(&l.status, cur, statusAcquiring) {
					continue
				}
			}
			l.sleep(status.Role == Leader)
		case statusResigned:
			if err := l.removeLease(); err != nil {
				log.Warn("failed to remove lease: ", err)
			}
			break Loop
		default:
			panic("unreachable")
		}
	}
	log.Info("election resigned: ", l.config.NodeID)
	close(l.ev)
	close(l.dead)
}

func (l *LeaderElector) Resign() {
	atomic.StoreInt32(&l.status, statusResigned)
	<-l.dead
}

func (l *LeaderElector) Status() <-chan Status {
	return l.ev
}

func (l *LeaderElector) statusFromDB(applied bool, id string, val string) Status {
	if applied {
		return Status{
			LeaderID:      l.config.NodeID,
			LeaderAddress: l.config.AdvertiseAddress,
			Role:          Leader,
		}
	}
	if id == l.CurrentNodeID() {
		return Status{
			LeaderID:      l.config.NodeID,
			LeaderAddress: l.config.AdvertiseAddress,
			Role:          Leader,
		}
	}
	return Status{
		LeaderID:      id,
		LeaderAddress: val,
		Role:          Follower,
	}
}

func (l *LeaderElector) updateLease() (Status, error) {
	cql := fmt.Sprintf(`UPDATE %s.%s SET leader_id = ?, value = ? WHERE resource_name = ? IF leader_id = ?`, l.config.Keyspace, l.config.TableName)
	q := l.session.Query(cql, l.config.NodeID, l.config.AdvertiseAddress,
		l.config.ResourceName, l.config.NodeID)
	defer q.Release()

	var id, val string
	applied, err := q.ScanCAS(&id, &val)
	if err != nil {
		return Status{}, err
	}
	return l.statusFromDB(applied, id, val), nil
}

func (l *LeaderElector) createLease() (Status, error) {
	cql := fmt.Sprintf(`INSERT INTO %s.%s (resource_name, leader_id, value) VALUES (?,?,?) IF NOT EXISTS`, l.config.Keyspace, l.config.TableName)
	q := l.session.Query(cql, l.config.ResourceName, l.config.NodeID, l.config.AdvertiseAddress)
	defer q.Release()

	var rn, id, val string
	applied, err := q.ScanCAS(&rn, &id, &val)
	if err != nil {
		return Status{}, err
	}
	return l.statusFromDB(applied, id, val), nil
}

func (l *LeaderElector) removeLease() error {
	cql := fmt.Sprintf(`DELETE FROM %s.%s WHERE resource_name=? IF leader_id=?`,
		l.config.Keyspace, l.config.TableName)
	q := l.session.Query(cql, l.config.ResourceName, l.config.NodeID)
	defer q.Release()
	var id string
	_, err := q.ScanCAS(&id)
	return err
}

func (l *LeaderElector) CreateTable(defaultTTLSec int) error {
	cql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
		resource_name text PRIMARY KEY,
		leader_id text,
		value text
	) with default_time_to_live = %d`, l.config.Keyspace, l.config.TableName, defaultTTLSec)
	return l.session.Query(cql).Exec()
}
