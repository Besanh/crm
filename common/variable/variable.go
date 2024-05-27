package variable

// Supplier
var SUPPLIER map[string]string = map[string]string{
	"fchat":     "fchat",
	"microsoft": "microsoft",
	"gmail":     "gmail",
}

// Cdr
var CDRSTATUS []string = []string{
	"answered",
	"no-answered",
	"busy-line",
	"busy",
	"not-available",
	"failed",
	"ivr",
	"cancel",
	"voicemail",
	"invalid-number",
	"phone-block",
	"congestion",
	"drop",
	"telco-block",
	"fail-system",
}
