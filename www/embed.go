package www

import _ "embed"

//go:embed slide.html
var Slide string

//go:embed summaryStyle.css
var StyleSummary string

//go:embed pdfStyle.css
var StylePDF string

//go:embed devMode.html
var DevMode string

//go:embed presenter.html
var Presenter string

//go:embed zoom.js
var ZoomJS string
