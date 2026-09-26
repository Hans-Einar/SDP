module github.com/Hans-Einar/SDP/Toolkit/SDPTool

go 1.26.0

require (
	github.com/Hans-Einar/SDP/SDUI/go v0.0.0
	github.com/Hans-Einar/SDP/SystemDesignLanguage/go v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/text v0.42.0 // indirect

replace github.com/Hans-Einar/SDP/SystemDesignLanguage/go => ../../SDL/go

replace github.com/Hans-Einar/SDP/SDUI/go => ../../SDUI/go
