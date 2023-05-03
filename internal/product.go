package internal

type SqliteProduct struct {
	Name        string
	Extension   string
	RelativeUrl string
	Version     string
	SizeInBytes int64
	Sha3Hash    string
}
