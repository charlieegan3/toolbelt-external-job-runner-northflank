package main

import (
	"log"

	"github.com/charlieegan3/toolbelt-external-job-runner-northflank/pkg/runner"

	"github.com/spf13/viper"
)

type exampleExternalJob struct {
	cfg map[string]any
}

func (e *exampleExternalJob) Name() string {
	return "example-external-job"
}

func (e *exampleExternalJob) RunnerName() string {
	return "nop"
}

func (e *exampleExternalJob) Config() map[string]any {
	return e.cfg
}

func main() {
	var err error
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}

	r := runner.Northflank{}
	err = r.Configure(viper.GetStringMap("northflank"))
	if err != nil {
		log.Fatalf("failed to configure runner: %v", err)
	}

	e := exampleExternalJob{
		cfg: viper.GetStringMap("job"),
	}

	err = r.RunJob(&e)
	if err != nil {
		log.Fatalf("failed to run job: %v", err)
	}
}
