package ichatappds

type ProfilePictures interface {
	GetPicturesList(picturesList map[uint]string)
	UpdatePicture(userID uint, picture string)
}
