package views

const (
	AlertLvlError   = "danger"
	AlertLvlWaring  = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random eror
	// is encoutered by our backedn.
	AlertMsgGeneric = "Something whent wrong. Please try " +
		"again, and contact us if the problem persists."
)

// Data is the top level structure that views expect data
// to come in.
type Data struct {
	Alert *Alert
	Yield interface{}
}

// Alert is used to render Bootstrap Alert Messages in templates
type Alert struct {
	Level   string
	Message string
}
