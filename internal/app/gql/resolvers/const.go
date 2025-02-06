package resolvers

const (
	cSymbolType             = "crypto"
	cSymbolSession          = "24x7"
	cSymbolTimezone         = "Etc/UTC"
	cSymbolMinMov           = 1
	cSymbolPricescale       = 10000
	cSymbolVisiblePlotsSets = "ohlcv"
	cSymbolVolumePrecision  = 6
	cSymbolDataStatus       = "streaming"
)

var supportedResolutions = []string{
	"1S",
	"2S",
	"5S",
	"1",
	"2",
	"5",
	"60",
	"120",
	"300",
	"1D",
	"2D",
	"1W",
	"2W",
	"1M",
	"2M",
	"3M",
}
