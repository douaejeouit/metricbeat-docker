package cpu


/*
import (
	"github.com/elastic/beats/libbeat/common"
	h "github.com/elastic/beats/metricbeat/helper"
)
// Map data to MapStr
func eventMapping([]common.MapStr) []common.MapStr {

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

*/