package bot

import (
	"fmt"
	"strconv"
	"strings"
)

type Op int

const (
	Save Op = iota
	Update
	Delete
	Get
)

func (o Op) String() string {
	switch o {
	case 0:
		return "save"
	case 1:
		return "update"
	case 2:
		return "delete"
	case 3:
		return "get"
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
	case "get":
		return Get
	default:
		panic("invalid op")
	}
}

type requirements struct {
	numberOfCommands int
}

var r map[Op]requirements = map[Op]requirements{
	Save:   {numberOfCommands: 2},
	Update: {numberOfCommands: 2},
	Get:    {numberOfCommands: 1},
	Delete: {numberOfCommands: 1},
}

func (op Op) Requirements() requirements { return r[op] }

// the text argument has the unformatted request string as an array
// eg ["save", "paschal", "31/12"]
func (op Op) Convert(text []string) (string, int, int, error) {
	switch op {
	case Save, Update:
		date := strings.Split(text[2], "/")
		if len(date) != 2 {
			return "", 0, 0, fmt.Errorf("invalid date %s", text[2])
		}
		day, err := strconv.Atoi(date[0])
		if err != nil {
			return "", 0, 0, err
		}
		if err := validateDay(day); err != nil {
			return "", 0, 0, err
		}
		month, err := strconv.Atoi(date[1])
		if err != nil {
			return "", 0, 0, err
		}
		if err := validateMonth(month); err != nil {
			return "", 0, 0, err
		}
		return text[1], day, month, nil
	case Get, Delete:
		month, err := strconv.Atoi(text[1])
		if err != nil {
			// this means that its not the month alone
			what := strings.Split(text[1], "/")
			if len(what) == 1 {
				// this means its a name, as a date will have two parts
				return what[0], 0, 0, nil
			}
			// else its an actual date
			day, err := strconv.Atoi(what[0])
			if err != nil {
				return "", 0, 0, err
			}
			if err := validateDay(day); err != nil {
				return "", 0, 0, err
			}
			month, err := strconv.Atoi(what[1])
			if err != nil {
				return "", 0, 0, err
			}
			if err := validateMonth(month); err != nil {
				return "", 0, 0, err
			}
			return "", day, month, nil
		}
		if err := validateMonth(month); err != nil {
			return "", 0, 0, err
		}
		return "", 0, month, nil
	default:
		panic("invalid op")
	}
}

func validateDay(day int) error {
	if day <= 0 || day > 31 {
		// invalid day, must be between 1 and 31
		return nil
	}
	return fmt.Errorf("invalid day %d", day)
}

func validateMonth(month int) error {
	if month <= 0 || month > 12 {
		// invalid month, month must be between 1 and 12
		return nil
	}
	return fmt.Errorf("invalid month %d", month)
}
