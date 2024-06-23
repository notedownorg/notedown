package persistor

type Persistor struct {
	root string
}

func NewPersistor() *Persistor {
	return &Persistor{}
}
