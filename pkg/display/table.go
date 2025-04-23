package display

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/leekchan/accounting"
	"github.com/rgoomar/aws-cct/pkg/models"
)

// TableDisplay handles the display of cost data in table format
type TableDisplay struct {
	ac *accounting.Accounting
}

// NewTableDisplay creates a new table display
func NewTableDisplay() *TableDisplay {
	return &TableDisplay{
		ac: accounting.DefaultAccounting("$", 2),
	}
}

// Render renders the cost comparison as a table
func (t *TableDisplay) Render(comparison *models.CostComparison, outputFormat string) string {
	tw := table.NewWriter()
	var secondMonthHeader = comparison.SecondMonthStart
	if comparison.IsProjection {
		secondMonthHeader += " (PROJECTION)"
	}

	// Remove commas to make output compatible for CSVs
	if outputFormat == "csv" {
		t.ac.SetThousandSeparator("")
	}

	tw.AppendHeader(table.Row{"Service", comparison.FirstMonthStart, secondMonthHeader, "Delta", "Delta Percent"})
	for _, serviceCosts := range comparison.ServiceCosts {
		tw.AppendRow(table.Row{
			serviceCosts.ServiceName,
			t.ac.FormatMoney(serviceCosts.Amount),
			t.ac.FormatMoney(serviceCosts.SecondAmount),
			t.ac.FormatMoney(serviceCosts.Delta),
			fmt.Sprintf("%s%%", accounting.FormatNumber(serviceCosts.DeltaPercent, 1, "", ".")),
		})
	}

	tw.AppendFooter(table.Row{
		"TOTAL",
		t.ac.FormatMoney(comparison.TotalAmount),
		t.ac.FormatMoney(comparison.TotalSecondAmount),
		t.ac.FormatMoney(comparison.TotalDelta),
		fmt.Sprintf("%s%%", accounting.FormatNumber(comparison.TotalDeltaPercent, 1, "", ".")),
	})

	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: comparison.FirstMonthStart, Align: text.AlignRight, AlignFooter: text.AlignRight},
		{Name: secondMonthHeader, Align: text.AlignRight, AlignFooter: text.AlignRight},
		{Name: "Delta", Align: text.AlignRight, AlignFooter: text.AlignRight},
		{Name: "Delta Percent", Align: text.AlignRight, AlignFooter: text.AlignRight},
	})

	switch outputFormat {
	case "csv":
		return tw.RenderCSV()
	default:
		return "\n" + tw.Render() + "\n"
	}
}
