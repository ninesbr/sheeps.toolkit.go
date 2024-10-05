package jsonstorage

type JsonStorage[T any] struct {
	client JsonStorageInterface
}

func NewJsonStorage[T any](client JsonStorageInterface) JsonStorage[T] {
	return JsonStorage[T]{
		client: client,
	}
}

func (s JsonStorage[T]) Push(docs ...T) error {
	var documents [][]byte
	for _, doc := range docs {
		d, err := structToDoc(doc)
		if err != nil {
			return err
		}
		buf, err := d.Marshal()
		if err != nil {
			return err
		}
		documents = append(documents, buf)
	}
	return s.client.pushDocuments(documents)
}

func (s JsonStorage[T]) Patch(docs ...T) error {
	var documents [][]byte
	for _, doc := range docs {
		d, err := structToDoc(doc)
		if err != nil {
			return err
		}
		buf, err := d.Marshal()
		if err != nil {
			return err
		}
		documents = append(documents, buf)
	}
	return s.client.patchDocuments(documents)
}

func (s JsonStorage[T]) Drop(ids ...string) error {
	return s.client.deleteDocuments(ids...)
}

func (s JsonStorage[T]) Get(id string) (*T, error) {
	var document T
	err := s.client.getDocument(id, &document)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (s JsonStorage[T]) Find(query map[string]string) ([]T, error) {
	var documents []T
	err := s.client.getDocuments(query, &documents)
	if err != nil {
		return nil, err
	}
	return documents, nil
}
