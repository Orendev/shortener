package shortlinks

import "errors"

type Storage interface {
	Get(code string) (ShortLink, error)
	Add(shortLink ShortLink) error
}

type storage struct {
	data map[string]ShortLink
}

func (s *storage) Get(code string) (shortLink ShortLink, err error) {
	shortLink, ok := s.data[code]
	if !ok {
		err = errors.New("not found")
	}
	return
}

func (s *storage) Add(shortLink ShortLink) error {
	s.data[shortLink.Code] = shortLink
	return nil
}

func NewStorage(data map[string]ShortLink) (*storage, error) {
	return &storage{
		data: data,
	}, nil
}
