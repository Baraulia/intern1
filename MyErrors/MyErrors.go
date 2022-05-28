package MyErrors

type Error string

func (e Error) Error() string { return string(e) }

const DoesNotExist = Error("object with this id does not exist")
