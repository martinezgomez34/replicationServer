package main

import (
	servidorMain "C1E/ServerMain"
	servidorReplication "C1E/ServerReplication"
)

func main() {
	go servidorMain.Run()
	go servidorReplication.Run()
	select {}
}
