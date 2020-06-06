package uuid

import "github.com/google/uuid"

type UUID = uuid.UUID

func New() UUID {
	//return uuid.New()
	return Must(NewUUID())
}

func NewUUID() (UUID, error) {
	return uuid.NewRandom()
}

// Must returns uuid if err is nil and panics otherwise.
func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}
