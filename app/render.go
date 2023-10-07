package app

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func render(containers []Container) {
	data := make([][]string, 0, len(containers))

	for _, c := range containers {
		ready := "0"
		if c.Ready() {
			ready = "1"
		}
		data = append(data, []string{
			c.Name,
			ready,
			c.State(),
			c.RestartCol(),
			c.Age(),
			c.Ports,
			c.Image,
			c.ImagePullPolicy,
			string(c.Type),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "READY", "STATUS", "RESTARTS", "AGE", "PORTS", "IMAGE", "PULLPOLICY", "TYPE"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	//table.EnableBorder(false)
	table.SetTablePadding("\t  ") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
