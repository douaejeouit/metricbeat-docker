package memory

/*import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/ingensi/metricbeat-docker/calculator"
	"github.com/fsouza/go-dockerclient"
        "github.com/ingensi/metricbeat-docker/module/docker"
	"strings"
	"time"
)

func (d *DataGenerator) GetMemoryData(container *docker.APIContainers, stats *docker.Stats) common.MapStr {
	event := common.MapStr{
		"@timestamp":      common.Time(stats.Read),
		"type":            "memory",
		"containerID":     container.ID,
		"containerName":   d.extractContainerName(container.Names),
		"containerLabels": d.buildLabelArray(container.Labels),
		"dockerSocket":    d.Socket,
		"memory": common.MapStr{
			"failcnt":    stats.MemoryStats.Failcnt,
			"limit":      stats.MemoryStats.Limit,
			"maxUsage":   stats.MemoryStats.MaxUsage,
			"totalRss":   stats.MemoryStats.Stats.TotalRss,
			"totalRss_p": float64(stats.MemoryStats.Stats.TotalRss) / float64(stats.MemoryStats.Limit),
			"usage":      stats.MemoryStats.Usage,
			"usage_p":    float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit),
		},
	}

	return event
}*/