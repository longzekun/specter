package core

import "sync"

type Job struct {
	ID          int
	Name        string
	Description string
	Protocol    string
	Port        uint16
	JobCtl      chan bool
}

var (
	Jobs = &jobs{
		active: &sync.Map{},
	}
	jobID = 0
)

type jobs struct {
	active *sync.Map
}

func (j *jobs) All() []*Job {
	all := []*Job{}
	j.active.Range(func(key, value interface{}) bool {
		all = append(all, value.(*Job))
		return true
	})
	return all
}

func (j *jobs) Add(job *Job) {
	j.active.Store(job.ID, job)
}

func (j *jobs) Remove(job *Job) {
	j.active.Delete(job.ID)
}

func (j *jobs) Get(id int) *Job {
	if id <= 0 {
		return nil
	}

	val, ok := j.active.Load(id)
	if ok {
		return val.(*Job)
	}
	return nil
}

func NextJobID() int {
	newID := jobID + 1
	jobID++
	return newID
}
