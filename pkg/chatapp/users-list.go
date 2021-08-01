package chatapp

import (
	"fmt"
	"sync"
)

type ProfilePictures struct {
	sync.Mutex
	pictures map[uint]string
}

func (pp *ProfilePictures) GetPicturesList(picturesList map[uint]string) {
	pp.Lock()
	defer pp.Unlock()
	for user := range picturesList {
		if picture, ok := pp.pictures[user]; ok {
			picturesList[user] = picture
		}
	}
}

func (pp *ProfilePictures) UpdatePicture(userID uint, picture string) {
	pp.Lock()
	defer pp.Unlock()
	fmt.Println("SESSION ADDED")
	pp.pictures[userID] = picture
}
