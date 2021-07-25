package chatapp

import (
	"sync"

	"github.com/ghenah/chatapp/pkg/ichatappds"
)

var TEMPORARYgenClientSessionID = makeIDGenerator64()

type ClientSessionsList struct {
	sync.Mutex
	// [sessionID]*session
	clientSessions map[uint64]*ichatappds.ClientSession
	// [ownerID]*sessionID
	ownerSessionsTable map[uint][]uint64
}

func (csl *ClientSessionsList) AddSession(userID uint, username string, outCh chan interface{}) (uint64, error) {

	session := &ichatappds.ClientSession{
		ID:            TEMPORARYgenClientSessionID(),
		OwnerID:       userID,
		OwnerUsername: username,
		OutCh:         outCh,
	}

	csl.Lock()
	defer csl.Unlock()
	csl.clientSessions[session.ID] = session
	csl.ownerSessionsTable[userID] = append(csl.ownerSessionsTable[userID], session.ID)

	return session.ID, nil
}

func (csl *ClientSessionsList) GetOutChannels(ownersIDs []uint) ([]ichatappds.ClSessChannelWithID, error) {
	csl.Lock()

	outChannels := []ichatappds.ClSessChannelWithID{}
	var s *ichatappds.ClientSession
	// Go through provided owners.
	for _, id := range ownersIDs {

		// Add all the sessions of an owner to the list.
		for _, sessionID := range csl.ownerSessionsTable[id] {
			s = csl.clientSessions[sessionID]
			outChannels = append(outChannels, ichatappds.ClSessChannelWithID{ID: s.ID, OutCh: s.OutCh})
		}

	}
	csl.Unlock()

	return outChannels, nil

}

func (csl *ClientSessionsList) RemoveSessionsByID(sessionIDs []uint64) {
	csl.Lock()
	for _, id := range sessionIDs {
		sessionOwnerID := csl.clientSessions[id].OwnerID
		csl.ownerSessionsTable[sessionOwnerID] = removeSessionReference(csl.ownerSessionsTable[sessionOwnerID], id)
		delete(csl.clientSessions, id)
	}
	csl.Unlock()
}

func removeSessionReference(s []uint64, e uint64) []uint64 {
	sOut := []uint64{}
	for _, el := range s {
		if e != el {
			sOut = append(sOut, el)
		}
	}

	return sOut
}
