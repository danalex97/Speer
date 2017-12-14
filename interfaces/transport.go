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

type ControlTransport interface {
	ControlPing(string) bool
	ControlSend(string, interface{})
	ControlRecv() <-chan interface{}
}

type DataTransport interface {
	Node
	Connect(string) Link
}

type Transport interface {
	ControlTransport
	DataTransport
}
