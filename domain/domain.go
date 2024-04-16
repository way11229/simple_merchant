package domain

import "time"

// time function type
type FuncTimeType func() time.Time

type ServiceManager struct {
	UserService    UserService
	AuthService    AuthService
	ProductService ProductService
}
