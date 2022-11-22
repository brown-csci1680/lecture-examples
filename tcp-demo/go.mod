module ip-demo

go 1.19

//replace gvisor.dev/gvisor v0.0.0-20221121220602-9ff1c425909e => github.com/google/gvisor v0.0.0-20221121220602-9ff1c425909e

require (
	golang.org/x/net v0.2.0
	gvisor.dev/gvisor v0.0.0-20221122005124-1af676b175b3
)

require (
	github.com/google/btree v1.1.2 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
)
