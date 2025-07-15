package errs

import "errors"

var (
	ErrInternal             = errors.New("internal error")
	ErrNotFound             = errors.New("not found")
	ErrAlreadyExists        = errors.New("already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUnauthenticated      = errors.New("unauthenticated")
	ErrCommentsNotAllowed   = errors.New("comments are not allowed on this post")
	ErrPostAndReplyMismatch = errors.New("reply id's post id doesn't match provided post id")
	ErrUnauthorized         = errors.New("you are not authorized to do that")
)
