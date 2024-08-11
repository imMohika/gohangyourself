package sub

type Command interface {
	Handle(args []string)
}
