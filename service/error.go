package service

import "errors"

var (
	UsError               error = errors.New("Call Us Error")
	NotFoundError         error = errors.New("Not found in DB")
	MongoError            error = errors.New("Query Mongo Error")
	ScmError              error = errors.New("Call Scm Error")
	EncryptError          error = errors.New("Encrypt token Error")
	CreateConflictError   error = errors.New("Cannot create the same identity")
	NotSupportedAuthError error = errors.New("Auth not supported")
	UserNotFoundError     error = errors.New("User Not found in DB")
	ScmNotFoundError      error = errors.New("Scm Not found in DB")
	DataTransformError    error = errors.New("Data transform failed")
)
