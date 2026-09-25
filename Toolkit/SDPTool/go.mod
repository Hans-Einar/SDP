module github.com/Hans-Einar/SDP/Toolkit/SDPTool

go 1.26.0

require (
	github.com/Hans-Einar/SDP/SystemDesignLanguage/go v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/Hans-Einar/SDP/SystemDesignLanguage/go => ../../SDL/go

replace github.com/Hans-Einar/SDP/SDUI/go => ../../SDUI/go
