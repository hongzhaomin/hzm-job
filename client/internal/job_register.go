package internal

import (
	"github.com/hongzhaomin/hzm-job/client/internal/anno"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/core/tools"
	"reflect"
	"sync"
	"sync/atomic"
)

var defaultJobRegister atomic.Pointer[JobRegister]

func init() {
	defaultJobRegister.Store(&JobRegister{
		jobMap: make(map[string]anno.Job, 16),
		lock:   &sync.Mutex{},
	})
}

func DefaultJobRegister() *JobRegister {
	return defaultJobRegister.Load()
}

type JobRegister struct {
	jobMap map[string]anno.Job
	lock   *sync.Mutex
}

func (my *JobRegister) GetJob(jobName string) anno.Job {
	my.lock.Lock()
	defer my.lock.Unlock()
	job, ok := my.jobMap[jobName]
	if !ok {
		global.SingletonPool().Log.Error("job is not exist", "jobName", jobName)
		return nil
	}
	return job
}

func (my *JobRegister) AddJob(jobName string, job anno.Job) {
	if jobName == "" {
		global.SingletonPool().Log.Error("job name is empty, please check before registry")
		return
	}
	if job == nil {
		global.SingletonPool().Log.Error("<nil> job is not allowed registry")
		return
	}
	my.addJob(jobName, job)
}

func (my *JobRegister) Registry(jobs ...anno.Job) {
	for _, job := range jobs {
		my.registry(job)
	}
}

func (my *JobRegister) registry(job anno.Job) {
	if job == nil {
		global.SingletonPool().Log.Error("<nil> job is not allowed registry")
		return
	}
	jobName, ok := tools.FindAnnotationValueByType(job, (*anno.Job)(nil), "name")
	if !ok {
		rt := reflect.TypeOf(job)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		jobName = rt.Name()
	}
	if jobName == "" {
		global.SingletonPool().Log.Error("job name is empty, please check before registry")
		return
	}
	my.addJob(jobName, job)
}

func (my *JobRegister) addJob(jobName string, job anno.Job) {
	my.lock.Lock()
	defer my.lock.Unlock()

	if _, ok := my.jobMap[jobName]; ok {
		global.SingletonPool().Log.Error("job is already registered", "jobName", jobName)
		return
	}
	my.jobMap[jobName] = job
}
