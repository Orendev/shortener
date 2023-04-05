package shortlinks

type ShortLink struct {
	Code string
	Link string
}

type ShortLinks struct {
	storage storage
}

func (s *ShortLinks) Get(code string) (shortLink ShortLink, err error) {
	return s.storage.Get(code)
}

func (s *ShortLinks) Add(shortLink ShortLink) error {
	return s.storage.Add(shortLink)
}

func New(s storage) (*ShortLinks, error) {
	return &ShortLinks{
		storage: s,
	}, nil
}
