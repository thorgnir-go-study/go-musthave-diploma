package auth

type userEntity struct {
	ID    string `db:"id"`
	Login string `db:"login"`
}
