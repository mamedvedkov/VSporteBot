package cron

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type Job struct {
	schedule string
	f        cron.FuncJob
}

// NewJob adds new cron job with shedule in format of cron '* * * * * *'
// where 'second minute hour dayOfMonth month dayOfWeek'
func NewJob(shedule string, f func()) *Job {
	return &Job{
		schedule: shedule,
		f:        f,
	}
}

type Cron struct {
	parser cron.Parser
	base   *cron.Cron
}

// New cron job controller
func New() *Cron {
	return &Cron{
		parser: cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
		base:   cron.New(),
	}
}

// MustAddJobs same as AddJobs but panic if error
func (c *Cron) MustAddJobs(jobs ...*Job) {
	if err := c.AddJobs(jobs...); err != nil {
		panic(fmt.Errorf("failed add cron jobs: %w", err))
	}
}

// AddJobs adds job to cron controller, error if controller can't parse schedule
func (c *Cron) AddJobs(jobs ...*Job) error {
	for _, job := range jobs {
		schedule, err := c.parser.Parse(job.schedule)
		if err != nil {
			return fmt.Errorf("parse schedule %s: %w", job.schedule, err)
		}

		c.base.Schedule(schedule, job.f)
	}

	return nil
}

// Run starts all added jobs and wait for context
func (c *Cron) Run(ctx context.Context) error {
	c.base.Start()
	<-ctx.Done()

	<-c.base.Stop().Done()

	return ctx.Err()
}
