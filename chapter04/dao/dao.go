package dao

type Dao interface {
	Get(string) (string, error)
	Set(string, string) error
}


