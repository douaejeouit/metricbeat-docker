package docker

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/ingensi/metricbeat-docker/calculator"
	"github.com/fsouza/go-dockerclient"
	"strings"
	"time"
)
type DataGenerator struct {
	Socket            *string
	CalculatorFactory calculator.CalculatorFactory
	Period            time.Duration
}
func (d *DataGenerator) GetCpuData(container *docker.APIContainers, stats *docker.Stats) common.MapStr {

	calculator := d.CalculatorFactory.NewCPUCalculator(
		calculator.CPUData{
			PerCpuUsage:       stats.PreCPUStats.CPUUsage.PercpuUsage,
			TotalUsage:        stats.PreCPUStats.CPUUsage.TotalUsage,
			UsageInKernelmode: stats.PreCPUStats.CPUUsage.UsageInKernelmode,
			UsageInUsermode:   stats.PreCPUStats.CPUUsage.UsageInUsermode,
		},
		calculator.CPUData{
			PerCpuUsage:       stats.CPUStats.CPUUsage.PercpuUsage,
			TotalUsage:        stats.CPUStats.CPUUsage.TotalUsage,
			UsageInKernelmode: stats.CPUStats.CPUUsage.UsageInKernelmode,
			UsageInUsermode:   stats.CPUStats.CPUUsage.UsageInUsermode,
		},
	)

	event := common.MapStr{
		"@timestamp":      common.Time(stats.Read),
		"type":            "cpu",
		"containerID":     container.ID,
		"containerName":   d.extractContainerName(container.Names),
		"containerLabels": d.buildLabelArray(container.Labels),
		"dockerSocket":    d.Socket,
		"cpu": common.MapStr{
			"percpuUsage":       calculator.PerCpuUsage(),
			"totalUsage":        calculator.TotalUsage(),
			"usageInKernelmode": calculator.UsageInKernelmode(),
			"usageInUsermode":   calculator.UsageInUsermode(),
		},
	}
	return event
}


func (d *DataGenerator) extractContainerName(names []string) string {
	output := names[0]

	if cap(names) > 1 {
		for _, name := range names {
			if strings.Count(output, "/") > strings.Count(name, "/") {
				output = name
			}
		}
	}
	return strings.Trim(output, "/")
}

func (d *DataGenerator) buildLabelArray(labels map[string]string) []common.MapStr {

	output_labels := make([]common.MapStr, len(labels))

	i := 0
	for k, v := range labels {
		label := strings.Replace(k, ".", "_", -1)
		output_labels[i] = common.MapStr{
			"key":   label,
			"value": v,
		}
		i++
	}
	return output_labels
}
