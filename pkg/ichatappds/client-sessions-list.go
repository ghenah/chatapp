package ichatappds

type ClientSessionsList interface {
	AddSession(userID uint, username string, outCh chan interface{}) (uint64, error)
	GetOutChannels([]uint) ([]ClSessChannelWithID, error)
	RemoveSessionsByID([]uint64)
}
