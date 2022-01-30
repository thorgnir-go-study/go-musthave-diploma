package repository

type UserEntity struct {
	ID    string `db:"id"`
	Login string `db:"login"`
}

type UserWithPassword struct {
	UserEntity
	Password string `db:"password"`
}
