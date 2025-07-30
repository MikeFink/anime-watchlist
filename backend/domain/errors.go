package domain

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

var (
	ErrInvalidTitle     = &ValidationError{Field: "title", Message: "title is required"}
	ErrInvalidAnilistID = &ValidationError{Field: "anilist_id", Message: "anilist_id must be positive"}
) 