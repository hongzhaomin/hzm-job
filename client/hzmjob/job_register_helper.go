package hzmjob

import (
	"github.com/hongzhaomin/hzm-job/client/annotation"
	"github.com/hongzhaomin/hzm-job/client/internal"
	"github.com/hongzhaomin/hzm-job/client/internal/anno"
)

func AddJob(jobName string, job annotation.JobFunc) {
	internal.DefaultJobRegister().AddJob(jobName, job)
}

func AddJobs(jobs ...anno.Job) {
	internal.DefaultJobRegister().Registry(jobs...)
}
