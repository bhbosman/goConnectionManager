all: DeleteAllGeneratedFiles IConnectionManager genmock

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go

IConnectionManager:
	mockgen -package goConnectionManager -generateWhat ddd -destination IConnectionManagerInterfaceMethods.go . IConnectionManager

genmock:
	mockgen -package goConnectionManager -generateWhat mockgen -destination ConnectionServiceManagerMock.go . IService



