package auth

type userEntity struct {
	ID    string `db:"id"`
	Login string `db:"login"`
}

type userWithPassword struct {
	userEntity
	Password string `db:"password"`
}
