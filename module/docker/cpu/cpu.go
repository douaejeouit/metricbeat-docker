package cpu

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/ingensi/metricbeat-docker/module/docker"
	"time"
	"fmt"

)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("docker", "cpu", New); err != nil {
		panic(err)
	}
}

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	ds *docker.DockerStats
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	fmt.Printf("HI there : new method")

	config := struct {
		Period    time.Duration  `config:"period"`
		Socket    string  `config:"socket"`
		enableTls bool  `config:"enable"`
		CaPath    string  `config:"ca_path"`
		CertPath  string   `config:"ca_path"`
		KeyPath   string  `config:"key_path"`
	}{
		Period: 1,
		Socket: "unix:///var/run/docker.sock",
		enableTls: false,

	}

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}
	return &MetricSet{
		BaseMetricSet: base,
		ds : docker.CreateDS(config.Period, config.Socket, config.enableTls),
	}, nil
	/*if config.enableTls == false {
		return &MetricSet{
			BaseMetricSet: base,
			ds : docker.CreateDS(config.Period, config.Socket, config.enableTls),
		}, nil
	} else {
		return &MetricSet{
			BaseMetricSet: base,
			ds: docker.CreateDSE(config.Period, config.Socket, config.enableTls, config.CaPath,config.CertPath,config.KeyPath),
		}, nil
	}*/
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {

	events := m.ds.GetDockerStats()
	return events, nil
}
