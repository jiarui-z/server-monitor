package terminal

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"runtime"
)

func byteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

type controller struct {
	Grid *ui.Grid

	HeapObjectsSparkline     *widgets.Sparkline
	HeapObjectSparklineGroup *widgets.SparklineGroup
	HeapObjectsData          *StatRing

	SysText       *widgets.Paragraph
	GCCPUFraction *widgets.Gauge

	HeapAllocBarChart     *widgets.BarChart
	HeapAllocBarChartData *StatRing

	HeapPie *widgets.PieChart
}

func (p *controller) Render(data *runtime.MemStats) {
	p.HeapObjectsData.Push(data.HeapObjects)
	p.HeapObjectsSparkline.Data = p.HeapObjectsData.NormalizedData()
	p.HeapObjectSparklineGroup.Title = fmt.Sprintf("Live heap object count: %d", data.HeapObjects)

	p.SysText.Text = fmt.Sprint(byteCountBinary(data.Sys))

	f := data.GCCPUFraction
	if f < 0.01 && f > 0 {
		for f < 1 {
			f = f * 10.0
		}
	}
	p.GCCPUFraction.Percent = int(f)

	p.GCCPUFraction.Label = fmt.Sprintf("%.2f%%", data.GCCPUFraction*100)

	p.HeapAllocBarChartData.Push(data.HeapAlloc)
	p.HeapAllocBarChart.Data = p.HeapAllocBarChartData.ToFloat()
	p.HeapAllocBarChart.Labels = nil
	for _, v := range p.HeapAllocBarChart.Data {
		p.HeapAllocBarChart.Labels = append(p.HeapAllocBarChart.Labels, byteCountBinary(uint64(v)))
	}

	p.HeapPie.Data = []float64{float64(data.HeapIdle), float64(data.HeapInuse)}

	ui.Render(p.Grid)
}

func (p *controller) resize() {
	w, h := ui.TerminalDimensions()
	p.Grid.SetRect(0, 0, w, h)
}

func (p *controller) Resize() {
	p.resize()
	ui.Render(p.Grid)
}

func (p *controller) initUI() {
	p.resize()

	p.HeapObjectsSparkline.LineColor = ui.Color(89) // xterm color DeepPink4
	p.HeapObjectSparklineGroup = widgets.NewSparklineGroup(p.HeapObjectsSparkline)

	p.SysText.Title = "The total bytes of memory obtained from the OS"
	p.SysText.PaddingLeft = 30
	p.SysText.PaddingTop = 2

	p.GCCPUFraction.Title = "GCCPUFraction %"
	p.GCCPUFraction.BarColor = ui.Color(50) // xterm color Cyan2

	p.HeapAllocBarChart.BarGap = 2
	p.HeapAllocBarChart.BarWidth = 8
	p.HeapAllocBarChart.Title = "Bytes of allocated heap objects"
	p.HeapAllocBarChart.NumFormatter = func(f float64) string { return "" }

	p.HeapPie.Title = "HeapInuse vs HeapIdle"
	p.HeapPie.LabelFormatter = func(idx int, _ float64) string { return []string{"Idle", "Inuse"}[idx] }

	p.Grid.Set(
		ui.NewRow(.2, p.HeapObjectSparklineGroup),
		ui.NewRow(.8,
			ui.NewCol(.5,
				ui.NewRow(.2, p.SysText),
				ui.NewRow(.2, p.GCCPUFraction),
				ui.NewRow(.6, p.HeapAllocBarChart),
			),
			ui.NewCol(.5, p.HeapPie),
		),
	)

}

func NewController() *controller {

	terminalWidth, _ := ui.TerminalDimensions()
	ctl := &controller{
		Grid: ui.NewGrid(),

		HeapObjectsSparkline: widgets.NewSparkline(),
		HeapObjectsData:      NewStatRing(terminalWidth),

		SysText: widgets.NewParagraph(),

		GCCPUFraction: widgets.NewGauge(),

		HeapAllocBarChart:     widgets.NewBarChart(),
		HeapAllocBarChartData: NewStatRing(6),

		HeapPie: widgets.NewPieChart(),
	}

	ctl.initUI()

	return ctl
}
