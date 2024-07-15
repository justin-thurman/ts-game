package persistence

import "context"

type playerStore interface {
	GetAccountID(context.Context, string) (int32, error)
	CreateAccount(context.Context, string, string, string) (int32, error)
}

type PlayerPersistence struct {
	store playerStore
}

func NewPersitence(store playerStore) PlayerPersistence {
	return PlayerPersistence{store: store}
}

// TODO: Should this go here? Or in the player package? I'm thinking player, but not exactly certain.
// Maybe we don't need the intermediate interface. Just pass queryEngine to player struct. Player struct
// requires something that implements an interface with the save player SQL call and that's it.
