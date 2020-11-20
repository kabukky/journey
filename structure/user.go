package structure

type User struct {
	Id       int64
	Name     []byte
	Slug     string
	Email    []byte
	Image    []byte
	Cover    []byte
	Bio      []byte
	Website  []byte
	Location []byte
	Role     int // 1 = Administrator, 2 = Editor, 3 = Author, 4 = Owner
}
