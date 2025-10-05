package auth

type IUserRepository interface {
	IUserRepositoryReader
	IUserRepositoryWriter
}

type IUserRepositoryReader interface {
	GetUserByName()
	GetUserById()
}

type IUserRepositoryWriter interface {
	CreateUser()
}
