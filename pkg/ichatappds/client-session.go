package ichatappds

type ClientSession struct {
	ID            uint64
	OwnerID       uint
	OwnerUsername string
	OutCh         chan interface{}
}

type ClSessChannelWithID struct {
	ID    uint64
	OutCh chan interface{}
}
