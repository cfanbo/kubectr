package app

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// join ports
func parsePort(ports []corev1.ContainerPort) string {
	if ports == nil {
		return "-"
	}

	result := make([]string, 0, len(ports))
	for _, cport := range ports {
		result = append(result, fmt.Sprintf("%d/%s", cport.ContainerPort, cport.Protocol))
	}
	return strings.Join(result, ",")
}

func formatContainerUptime(uptime time.Duration, segmentCount int) string {
	days := int(uptime.Hours() / 24)
	uptime -= time.Duration(days) * 24 * time.Hour

	hours := int(uptime.Hours())
	uptime -= time.Duration(hours) * time.Hour

	minutes := int(uptime.Minutes())
	uptime -= time.Duration(minutes) * time.Minute

	seconds := int(uptime.Seconds())
	uptime -= time.Duration(seconds) * time.Second

	result := ""
	segment := 0
	if days > 0 {
		result += fmt.Sprintf("%dd", days)
		segment++
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
		segment++
	}
	if minutes > 0 && segment < segmentCount {
		result += fmt.Sprintf("%dm", minutes)
		segment++
	}
	if seconds > 0 && segment < segmentCount {
		result += fmt.Sprintf("%ds", seconds)
	}

	return result
}
