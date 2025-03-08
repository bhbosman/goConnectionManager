module github.com/bhbosman/goConnectionManager

go 1.23.0

toolchain go1.24.0

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20230730212737-00ad0cf16194
	github.com/bhbosman/gocommon v0.0.0-20250307235859-f370cb0a3bac
	github.com/cskr/pubsub v1.0.2
	github.com/golang/mock v1.6.0
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.23.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/icza/gox v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af
