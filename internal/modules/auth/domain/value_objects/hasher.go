package valueobjects

type Hasher interface {
	Hash(string) (string, error)
	Compare(hashedPassword, password string) error
}