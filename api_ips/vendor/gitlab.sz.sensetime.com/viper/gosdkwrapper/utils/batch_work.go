package utils

import (
	"time"

	"github.com/jeffail/tunny"
	log "github.com/sirupsen/logrus"
)

type WorkerType int

const (
	GpuWorker WorkerType = iota
	// nolint: golint
	CpuWorker
)

type dataWithIndex struct {
	idx  int
	data interface{}
}

type Worker interface {
	TunnyJob(interface{}) interface{}
	TunnyReady() bool
	Release()
}

type WorkerPool struct {
	// nolint: golint
	CpuWorkers []tunny.TunnyWorker
	GpuWorkers []tunny.TunnyWorker
	cpuPool    *tunny.WorkPool
	gpuPool    *tunny.WorkPool
}

type WorkerPoolConfig struct {
	CpuConcurrency    int // nolint
	GpuConcurrency    int // nolint
	CpuWorkerInitFunc func() Worker
	GpuWorkerInitFunc func() Worker
}

func initWorker(workerType WorkerType, num int, initFunc func() Worker) []tunny.TunnyWorker {
	if num == 0 || initFunc == nil {
		log.Fatal("invalid param while initializing worker")
	}

	tag := "cpu"
	if workerType == GpuWorker {
		tag = "gpu"
	}
	log.Infof("init %d %v workers", num, tag)

	workers := make([]tunny.TunnyWorker, num)
	for i := range workers {
		workers[i] = initFunc()
	}
	return workers
}

func NewWorkerPool(config WorkerPoolConfig) *WorkerPool {
	cpuWorkers := initWorker(CpuWorker, config.CpuConcurrency, config.CpuWorkerInitFunc)
	gpuWorkers := initWorker(GpuWorker, config.GpuConcurrency, config.GpuWorkerInitFunc)
	cpuPool, err := tunny.CreateCustomPool(cpuWorkers).Open()
	if err != nil {
		panic(err)
	}
	gpuPool, err := tunny.CreateCustomPool(gpuWorkers).Open()
	if err != nil {
		panic(err)
	}

	return &WorkerPool{
		CpuWorkers: cpuWorkers,
		GpuWorkers: gpuWorkers,
		cpuPool:    cpuPool,
		gpuPool:    gpuPool,
	}
}

func (p *WorkerPool) RunBatch(workerType WorkerType, jobs []interface{}, timeout time.Duration) []interface{} {
	count := len(jobs)
	if count == 0 {
		return nil
	}
	pool := p.cpuPool
	if workerType == GpuWorker {
		pool = p.gpuPool
	}
	ch := make(chan dataWithIndex, count)
	for i, v := range jobs {
		idx := i
		pool.SendWorkTimedAsync(timeout, v, func(data interface{}, err error) {
			if err != nil {
				ch <- dataWithIndex{idx, err}
			} else {
				ch <- dataWithIndex{idx, data}
			}
		})
	}
	outputs := make([]interface{}, count)
	for i := 0; i < count; i++ {
		data := <-ch
		outputs[data.idx] = data.data
	}
	return outputs
}

func (p *WorkerPool) SendWorkTimed(workerType WorkerType, milliTimeout time.Duration, jobData interface{}) (interface{}, error) {
	if workerType == GpuWorker {
		return p.gpuPool.SendWorkTimed(milliTimeout, jobData)
	}
	return p.cpuPool.SendWorkTimed(milliTimeout, jobData)
}

func (p *WorkerPool) Release() {
	if err := p.cpuPool.Close(); err != nil {
		panic(err)
	}
	if err := p.gpuPool.Close(); err != nil {
		panic(err)
	}
	for _, v := range p.GpuWorkers {
		v.(Worker).Release()
	}
	for _, v := range p.CpuWorkers {
		v.(Worker).Release()
	}
}
