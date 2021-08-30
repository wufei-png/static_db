package client

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/localcache"
)

// const ()

var backoffTimes = []time.Duration{
	50 * time.Millisecond,
	1 * time.Second,
	2 * time.Second,
	5 * time.Second,
	10 * time.Second,
	20 * time.Second,
}

// BackoffTime returns: 50ms, 1s, 2s, 5s, 10s, 20s, 20s, 20s...
func BackoffTime(i int) time.Duration {
	if i >= len(backoffTimes) {
		i = len(backoffTimes) - 1
	}
	return backoffTimes[i] + time.Duration(rand.Intn(100))*time.Millisecond
}

func getDefaultCache() *Cache {
	root := os.Getenv("MODEL_CACHE_ROOT")
	if root == "" {
		root = "/tmp/client-model-cache"
	}
	endpoint := os.Getenv("MODEL_CACHE_ADDRESS")
	c, err := NewCache(root, endpoint, DefaultTimeout)
	if err != nil {
		log.Warn("WARN: failed to create default model cache: ", err)
		return nil
	}
	return c
}

var defaultModelCache struct {
	sync.Mutex
	c *Cache
}

func GetDefaultModelCache() *Cache {
	defaultModelCache.Lock()
	defer defaultModelCache.Unlock()
	if defaultModelCache.c == nil {
		defaultModelCache.c = getDefaultCache()
	}
	return defaultModelCache.c
}

type Cache struct {
	c *localcache.LocalCache
	m *Manager
}

func NewCache(rootDir string, endpoint string, timeout time.Duration) (*Cache, error) {
	c, err := localcache.NewLocalCache(rootDir)
	if err != nil {
		return nil, err
	}
	var m *Manager
	if endpoint != "" {
		m = NewManager(endpoint, timeout)
	} else {
		log.Info("creating client cache with offline mode")
	}
	return &Cache{
		c: c,
		m: m,
	}, nil
}

func (c *Cache) IsOffline() bool {
	return c.m == nil
}

type ModelFetchOpts struct {
	LocalPath    string
	DisableProbe bool
	ForceOffline bool
	NoRetry      bool
}

func isDefaultName(s string) bool {
	return s == "" || s == localcache.DefaultName
}

func (c *Cache) FindRefInCache(mp api.ModelPath, probe bool) (*localcache.MetaFile, string, error) {
	mf, path, err := c.c.Get(&mp)
	if err == nil {
		return mf, path, err
	}

	if !probe {
		return mf, path, err
	}

	// first try default hw
	if !isDefaultName(mp.Hardware) {
		mp.Hardware = localcache.DefaultName
		mf, path, err = c.c.Get(&mp)
		if err == nil {
			return mf, path, err
		}
	}

	// then try default rt
	if !isDefaultName(mp.Runtime) {
		mp.Runtime = localcache.DefaultName
		mf, path, err = c.c.Get(&mp)
		if err == nil {
			return mf, path, err
		}
	}

	return mf, path, err
}

func (c *Cache) FindRefInRemote(mp api.ModelPath, probe bool) (*localcache.MetaFile, error) {
	if c.IsOffline() {
		return nil, ErrModelNotFound
	}
	mf, err := c.m.GetModelMeta(&mp)
	if err == nil {
		return mf, err
	}

	if !probe {
		return mf, err
	}

	// first try default hw
	if !isDefaultName(mp.Hardware) {
		mp.Hardware = localcache.DefaultName
		mf, err = c.m.GetModelMeta(&mp)
		if err == nil {
			return mf, err
		}
	}

	// then try default rt
	if !isDefaultName(mp.Runtime) {
		mp.Runtime = localcache.DefaultName
		mf, err = c.m.GetModelMeta(&mp)
		if err == nil {
			return mf, err
		}
	}

	return mf, err
}

// GetModelLocalPath set retries <= 0 for local only
func (c *Cache) GetModelLocalPath(ctx context.Context, modelRef string, opts ModelFetchOpts) (*api.Model, string, error) {
	entry := log.WithField("model_ref", modelRef)
	if opts.LocalPath != "" {
		_, err := os.Stat(opts.LocalPath)
		if err == nil {
			entry.Info("using local model file: ", opts.LocalPath)
			return nil, opts.LocalPath, nil
		}
	}

	if modelRef == "" {
		return nil, "", ErrModelNotFound
	}
	mp, err := ParseModelRef(modelRef)
	if err != nil {
		return nil, "", err
	}

	offline := c.IsOffline() || opts.ForceOffline

	if offline {
		mf, path, err := c.FindRefInCache(*mp, !opts.DisableProbe)
		if err == nil {
			entry.Info("using model: ", mf.Meta, ", blob: ", path)
			return mf.Meta, path, nil
		}
		entry.Error("failed to fetch offline model: ", err)
		return nil, "", ErrModelNotFound
	}

	retries := len(backoffTimes)
	if opts.NoRetry {
		retries = 1
	}
	// in online mode, always check for
	for r := 0; r < retries; r++ {
		// backoff
		if r > 0 {
			entry.Info("try get model meta for ", modelRef, ", retry: ", r)
			t := BackoffTime(r)
			select {
			case <-ctx.Done():
				return nil, "", ctx.Err()
			case <-time.After(t):
			}
		}

		var mf *localcache.MetaFile
		var path string

		mf, err = c.FindRefInRemote(*mp, !opts.DisableProbe)
		if err == ErrModelNotFound {
			continue
		}
		if err != nil {
			entry.Warn("failed to fetch model in remote: ", err, ", try using local cache")
			mf, path, err = c.FindRefInCache(*mp, !opts.DisableProbe)
			if err != nil {
				entry.Warn("find ref in cache failed: ", err)
				continue
			}
			entry.Info("using cached model: ", mf.Meta, ", blob: ", path)
			return mf.Meta, path, nil
		}

		path, err = c.downloadModelIfNeeded(mf)
		if err != nil {
			entry.Error("failed to fetch model: ", err)
			continue
		}

		entry.Info("using fetched model: ", mf.Meta, ", blob: ", path)
		return mf.Meta, path, nil
	}
	return nil, "", err
}

// for test
var downloadSleepTime = time.Duration(0)

func (c *Cache) downloadModelIfNeeded(mf *localcache.MetaFile) (string, error) {
	path, err := c.c.ValidateBlob(mf.Meta)
	if err == nil {
		return path, nil
	}

	writer, err := c.c.LockForWriter(mf)
	if err != nil {
		return "", err
	}

	log.Info("downloading blob for: ", mf.Meta)
	startTime := time.Now()
	reader, err := c.m.GetModelBlob(mf.Meta)
	if err != nil {
		writer.Abort() // nolint
		return "", err
	}
	defer reader.Close() // nolint

	err = localcache.CopyWithChecksum(writer, reader, mf.Meta.GetSize(), mf.Meta.GetChecksum())
	if err != nil {
		writer.Abort() // nolint
		return "", err
	}
	if downloadSleepTime > 0 {
		time.Sleep(downloadSleepTime)
	}
	endTime := time.Now()
	log.Info("downloaded blob for: ", mf.Meta, ", time: ", endTime.Sub(startTime))
	if err := writer.Commit(); err != nil {
		return "", err
	}

	return c.c.ValidateBlob(mf.Meta)
}
