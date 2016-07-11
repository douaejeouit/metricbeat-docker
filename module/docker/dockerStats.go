package docker

import (
	"errors"
	"time"
	"fmt"
	"sync"
	"github.com/fsouza/go-dockerclient"
	"github.com/ingensi/metricbeat-docker/calculator"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

)
//var BeatConfig *config.Config

type SocketConfig struct {
	socket    string
	enableTls bool
	caPath    string
	certPath  string
	keyPath   string
}

// check docker config ??
type DockerStats struct {
	period               time.Duration
	socketConfig         SocketConfig
	dockerClient         *docker.Client
	dataGenerator         DataGenerator
}

// if tls diseable
 func  CreateDS(pPeriod time.Duration, pSocket string, pEnable bool) *DockerStats {
	 return &DockerStats{
		 period: pPeriod,
		 socketConfig: SocketConfig{
			 socket: pSocket,
			 enableTls: pEnable,
		 },
	 }
 }

//if tls enable
/*func  CreateDSE(pPeriod time.Duration, pSocket string, pEnable bool,
 				pCapath string, pCertpath string, pKeypath string){

	tmpDS := CreateDS(pPeriod,pSocket,pEnable)
	return tmpDS{
		tmpDS.socketConfig.enableTls: pEnable,
		tmpDS.socketConfig.caPath: pCapath,
		tmpDS.socketConfig.certPath: pCertpath,
		tmpDS.socketConfig.keyPath: pKeypath,
	}
}
*/
func (bt *DockerStats) InitDockerCLient() error{
	logp.Info("Je suis à : InitDockerCLient ")
	var clientErr error
	var err error

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
func (bt *DockerStats) GetDockerClient() (*docker.Client, error) {
	logp.Info("Je suis à :GetDockerClient ")
	var client *docker.Client
	var err error
	if bt.socketConfig.enableTls ==true{
		client, err = docker.NewTLSClient(
			bt.socketConfig.socket,
			bt.socketConfig.certPath,
			bt.socketConfig.keyPath,
			bt.socketConfig.caPath,
		)
	}else {
		client, err = docker.NewClient(bt.socketConfig.socket)

	}
	return client, err
}
func (d *DockerStats) GetDockerStats() ([]common.MapStr) {
	logp.Info("Je suis à : GetDockerStats")
	/*fmt.Printf(" ",d.period)
	fmt.Printf("socket : ",d.socketConfig.socket)
	fmt.Printf("enable : ",d.socketConfig.enableTls)
	*/
	d.InitDockerCLient()
	logp.Info("DockerSTat is running")
	//ticker := time.NewTicker(d.period)
	//defer ticker.Stop()

	//for {
		/*select {
		//case <-bt.done:
		//	return nil
		case <-ticker.C:
		}*/
		//timerStart := time.Now()
		myStats, _ := d.FetchSTats()
		//
	// timerEnd := time.Now()
		if myStats != nil {
			logp.Info(" Great, stats are available! \n")
			logp.Info(" Data: %v", myStats)
			//duration := (timerEnd.Sub(timerStart) * time.Second)

			return myStats
		}

		logp.Info("Oups, No stats available ")
		return nil
	//}

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
	//done := make(chan bool)
	errC := make(chan error, 1)
	var wg sync.WaitGroup
	events := common.MapStr{}
	// the stream bool is set to false to only listen the first stats
	statsOptions := docker.StatsOptions{
		ID:      container.ID,
		Stats:   statsC,
		Stream:  false,
		Timeout: -1,
	}
	wg.Add(2)
	// goroutine to listen to the stats
	go func() {
		defer wg.Done()
		errC <- d.dockerClient.Stats(statsOptions)
		close(errC)
	}()
	// goroutine to get the stats & publish it
	go func() {
		defer wg.Done()
		stats := <-statsC
		err := <-errC

		if err == nil && stats != nil {
			events = d.dataGenerator.GetCpuData(&container, stats)
		} else if err == nil && stats == nil {
			logp.Warn("Container was existing at listing but not when getting statistics: %v", container.ID)

		} else {
			logp.Err("An error occurred while getting docker stats: %v", err)

		}
	}()
	wg.Wait()
	return events
}
