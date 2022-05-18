package bot

// var allowedOps = []string{"save", "update", "delete"}

type Op int

const (
	Save Op = iota
	Update
	Delete
)

func (o Op) String() string {
	switch o {
	case 0:
		return "save"
	case 1:
		return "update"
	case 2:
		return "delete"
	default:
		panic("invalid string")
	}
}

func IsValidOp(operation string) bool {
	for _, op := range []string{Save.String(), Update.String(), Delete.String()} {
		if op == operation {
			return true
		}
	}
	return false
}

func ToOp(op string) Op {
	switch op {
	case "save":
		return Save
	case "update":
		return Update
	case "delete":
		return Delete
	default:
		panic("invalid op")
	}
}
