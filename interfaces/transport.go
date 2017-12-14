package interfaces

type Data struct {
	Id   string
	Size int
}

type NodeCapacity interface {
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
	From() NodeCapacity
	To() NodeCapacity
}

type ControlTransport interface {
	ControlPing(string) bool
	ControlSend(string, interface{})
	ControlRecv() <-chan interface{}
}

type DataTransport interface {
	NodeCapacity
	Connect(string) Link
}

type Transport interface {
	ControlTransport
	DataTransport
}
