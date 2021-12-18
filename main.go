package main

import (
	"flag"
	"log"
	"os"

	"github.com/contribsys/faktory/client"
	"github.com/robfig/cron/v3"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Faktory FaktoryConfig  `yaml:"faktory"`
	Jobs    []ScheduledJob `yaml:"jobs"`
}

type FaktoryConfig struct {
	URL string `yaml:"url"`
}

type ScheduledJob struct {
	Name  string `yaml:"name"`
	Every string `yaml:"every"`

	JobType string                 `yaml:"job_type"`
	Queue   *string                `yaml:"queue"`
	Args    []interface{}          `yaml:"args"`
	Custom  map[string]interface{} `yaml:"custom"`
}

func main() {
	logger := log.Default()

	configPath := flag.String("config", "config.yaml", "Path to configuration file")

	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		logger.Fatalf("could not open configuration: %v", err)
	}

	decoder := yaml.NewDecoder(file)

	var config Config
	if err = decoder.Decode(&config); err != nil {
		logger.Fatalf("could not decode configuration: %v", err)
	}

	if len(config.Faktory.URL) > 0 {
		os.Setenv("FAKTORY_URL", config.Faktory.URL)
	}

	faktoryPool, err := client.NewPool(2)
	if err != nil {
		logger.Fatalf("could not connect to faktory: %v", err)
	}

	c := cron.New()

	for _, job := range config.Jobs {
		logger.Printf("Registering job %s to run with pattern %s\n", job.Name, job.Every)

		currentJob := job
		c.AddFunc(job.Every, func() {
			logger.Printf("Running job %s\n", currentJob.Name)

			newJob := client.NewJob(currentJob.JobType, currentJob.Args...)
			if currentJob.Queue != nil {
				newJob.Queue = *currentJob.Queue
			}
			if currentJob.Custom != nil {
				newJob.Custom = currentJob.Custom
			}

			faktory, err := faktoryPool.Get()
			if err != nil {
				log.Fatalf("could not get faktory connection: %v", err)
			}
			defer faktoryPool.Put(faktory)

			if err = faktory.Push(newJob); err != nil {
				log.Fatalf("could not enqueue faktory job: %v", err)
			}
		})
	}

	c.Start()
	select {}
}
