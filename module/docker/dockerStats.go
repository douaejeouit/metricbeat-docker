package docker

import (
	"errors"
	"time"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/ingensi/metricbeat-docker/calculator"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)


type SocketConfig struct {
	socket    string
}

// check docker config ??
type DockerStats struct {
	period               time.Duration
	socketConfig         SocketConfig
	dockerClient         *docker.Client
	dataGenerator         DataGenerator
}
// TO DO : Add the configuration


 func New() *DockerStats {
	 return &DockerStats{}
 }
func (bt *DockerStats) GetDockerClient() (*docker.Client, error) {
	var client *docker.Client
	var err error
	 // TO ADD TO THE CONFIG METHOD
	/*bt.socketConfig = SocketConfig{
		socket:    "unix:///var/run/docker.sock",
	}
	if bt.socketConfig.enableTls {
		client, err = docker.NewTLSClient(
			bt.socketConfig.socket,
			bt.socketConfig.certPath,
			bt.socketConfig.keyPath,
			bt.socketConfig.caPath,
		)
	} else {*/
	client, err = docker.NewClient(bt.socketConfig.socket)

	return client, err
}

func (bt *DockerStats) InitDockerCLient() error{
	var clientErr error
	var err error
	bt.period = 10
	bt.socketConfig = SocketConfig{
		socket:    "unix:///var/run/docker.sock",
	}
	bt.dockerClient, clientErr = bt.GetDockerClient()
	bt.dataGenerator = DataGenerator{
		Socket: & bt.socketConfig.socket,
		CalculatorFactory: calculator.CalculatorFactoryImpl{},
		Period:    bt.period,
	}


	 if clientErr != nil {
		 err = errors.New(fmt.Sprintf(" Unable to create dockerCLient"))
	 }
	return err
}
func (d *DockerStats) GetDockerStats() ([]common.MapStr) {
	d.InitDockerCLient()
	logp.Info("DockerSTat is running")

	timerStart := time.Now()
	myStats ,_:= d.FetchSTats()
	timerEnd := time.Now()

	if myStats != nil {
		logp.Info(" Great, stats are available! \n")
		logp.Info(" Data: %v", myStats)
		return myStats
	}
	duration := timerEnd.Sub(timerStart)
	logp.Info(" Duration is : %d",duration)
	logp.Info("Oups, No stats available ")
	return nil
}
func (d *DockerStats) FetchSTats() ([]common.MapStr, error ){
	containers, err := d.dockerClient.ListContainers(docker.ListContainersOptions{})

	myEvents := []common.MapStr{}
	if err == nil {
		//export stats for each container
		for _, container := range containers {
			myEvents = append(myEvents, d.ExportContainerStats(container))
		}
	} else {
		logp.Err("Cannot get container list: %v", err)
	}

	//d.dataGenerator.CleanOldStats(containers)

	return myEvents, err
}
func (d *DockerStats) ExportContainerStats(container docker.APIContainers) common.MapStr  {
	// statsOptions creation
	statsC := make(chan *docker.Stats)
	done := make(chan bool)
	errC := make(chan error, 1)
	events := common.MapStr{}
	// the stream bool is set to false to only listen the first stats
	statsOptions := docker.StatsOptions{
		ID:      container.ID,
		Stats:   statsC,
		Stream:  false,
		Done:    done,
		Timeout: -1,
	}
	// goroutine to listen to the stats
	go func() {
		errC <- d.dockerClient.Stats(statsOptions)
		close(errC)
	}()
	// goroutine to get the stats & publish it
	//go func() {
		stats := <-statsC
		err := <-errC

		if err == nil && stats != nil {
			events = d.dataGenerator.GetCpuData(&container, stats)
			//d.events.PublishEvents(events)
		} else if err == nil && stats == nil {
			logp.Warn("Container was existing at listing but not when getting statistics: %v", container.ID)
			//d.publishLogEvent(WARN, fmt.Sprintf("Container was existing at listing but not when getting statistics: %v", container.ID))
		} else {
			logp.Err("An error occurred while getting docker stats: %v", err)
			//d.publishLogEvent(ERROR, fmt.Sprintf("An error occurred while getting docker stats: %v", err))
		}
	//}()

	return events
}
