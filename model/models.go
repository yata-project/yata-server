package model

type YataList struct {
	UserID UserID
	ListID ListID
	Title  string
}

type YataItem struct {
	UserID  UserID
	ListID  ListID
	ItemID  ItemID
	Content string
}
