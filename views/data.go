package views

import "log"

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

// PublicError is an interface that allows the definition of
// methods for printing public messages
type PublicError interface {
	error // embeding the error interface
	Public() string
}

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

// SetAlert is a method that takes in an error and sets the alert for that data
func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// AlertError is a method on data that takes a string and puts in
// in the Alert -> Message field
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
