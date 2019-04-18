package interfaces

type Data struct {
	Id   string
	Size int
}

type Node interface {
	Up() int
	Down() int
}

type Link interface {
	// Functions provided for Uploading and Downloading data from a link.
	Upload(Data)
	Download() <-chan Data

	// Function used to clear the data to be uploaded to the channel.
	Clear()

	// Upload and download capacities(together with an identifier) for each
	// end of the connection.
	From() Node
	To() Node
}

type Transport interface {
	Node

	// Create a connection used for data transfer.
	Connect(string) Link

	// Control message interfaces.
	ControlPing(string) bool
	ControlSend(string, interface{})
	ControlRecv() <-chan interface{}
}
