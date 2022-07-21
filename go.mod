module github.com/bhbosman/goConnectionManager

go 1.18

require (
	github.com/bhbosman/gocommon v0.0.0-20220627073905-4951fb81c325
	github.com/cskr/pubsub v1.0.2
	github.com/golang/mock v1.6.0
	go.uber.org/fx v1.17.1
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20220318055525-2edf467146b5 // indirect
)

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20220617134815-f277ff266f47

replace github.com/bhbosman/gocommon => ../gocommon
replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions
