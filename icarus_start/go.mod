module liam-tool/icarus_start

go 1.18

replace liam-tool/modules => ../modules

require (
	git.ace/icarus/icarusclient/v5 v5.4.4
	liam-tool/modules v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
)
