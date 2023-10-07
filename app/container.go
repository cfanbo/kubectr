package app

import (
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var (
	Waiting    ContainerState = "Waiting"
	Running    ContainerState = "Running"
	Terminated ContainerState = "Terminated"
)

type Container struct {
	Type            containerType
	Name            string
	Image           string
	ImagePullPolicy string
	Ports           string
	Status          corev1.ContainerStatus
}

func (c Container) State() string {
	var state ContainerState
	if c.Status.State.Waiting != nil {
		state = Waiting
	} else if c.Status.State.Running != nil {
		state = Running
	} else {
		state = Terminated
	}

	return string(state)
}

func (c Container) Ready() bool {
	return c.Status.Ready
}

func (c Container) RestartCount() int32 {
	return c.Status.RestartCount
}

func (c Container) RestartDurationString() string {
	if c.RestartCount() == 0 {
		return ""
	}

	if c.Status.LastTerminationState.Terminated != nil {
		dur := time.Since(c.Status.LastTerminationState.Terminated.FinishedAt.Time)
		return formatContainerUptime(dur, 1)
	}
	return ""
}

func (c Container) RestartCol() string {
	if c.RestartDurationString() != "" {
		return fmt.Sprintf("%d (%s ago)", c.RestartCount(), c.RestartDurationString())
	}

	return strconv.Itoa(int(c.RestartCount()))
}

func (c Container) Age() string {
	if c.Status.State.Running != nil {
		dur := time.Since(c.Status.State.Running.StartedAt.Time)
		return formatContainerUptime(dur, 2)
	}
	return "-"
}
