package main

import (
	"fmt"
	"os"

	"k8s.io/klog/v2"

	"github.com/cfanbo/kubectr/app"
)

func main() {
	if err := app.NewCmd().Execute(); err != nil {
		if klog.V(1).Enabled() {
			klog.Fatalf("%+v", err) // with stack trace
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
