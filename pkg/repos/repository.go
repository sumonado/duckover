package repos

type Repository interface {
	Download(serial string) error
}
