package main

import "context"

type JobResult struct {
	output string
	err    error
}

type Job struct {
	script        string
	cb            func(output string, err error)
	outputChannel chan JobResult
	ctx           context.Context
}

func NewJob(script string, ctx context.Context) *Job {
	j := &Job{
		script:        script,
		outputChannel: make(chan JobResult),
		ctx:           ctx,
	}
	j.cb = func(output string, err error) {
		j.outputChannel <- JobResult{
			output: output,
			err:    err,
		}
	}
	return j
}

func (j *Job) GetOutput() (string, error) {
	res := <-j.outputChannel
	return res.output, res.err
}

type Worker struct {
	queue chan *Job
}

func NewWorker() *Worker {
	return &Worker{
		queue: make(chan *Job),
	}
}

func (w *Worker) Push(j *Job) {
	w.queue <- j
}

func (w *Worker) Work() {
	for {
		select {
		case j := <-w.queue:
			output, err := EvaluateScript(j.script, j.ctx)
			j.cb(output, err)
		}
	}
}
