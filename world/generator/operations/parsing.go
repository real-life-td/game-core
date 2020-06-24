package operations

import (
	"errors"
	"github.com/real-life-td/game-core/world/generator/parsing"
	"go/ast"
	"strings"
	"unicode"
)

type stage int
type action int

const (
	initStage stage = iota // Keep at start
	deltaStage
	numStages // Keep at end
)

const (
	// Order based on precedence (which actions should be applied first)
	setAction action = iota
	removeAction
	addAction
	deleteAction
	putAction
	putMultipleAction
)

var tokenToStage = map[string]stage{
	"init": initStage,
	"delta": deltaStage,
}

var tokenToAction = map[string]action{
	"SET":          setAction,
	"ADD":          addAction,
	"REMOVE":       removeAction,
	"PUT":          putAction,
	"PUT_MULTIPLE": putMultipleAction,
	"DELETE":       deleteAction,
}

type operation struct {
	field     string
	fieldType parsing.GoType
	stage     stage
	action    action
}

type StageOperations map[stage][]*operation

func (s StageOperations) add(o *operation) {
	curSlice, ok := s[o.stage]
	if ok {
		curSlice = append(curSlice, o)
		s[o.stage] = curSlice
	} else {
		s[o.stage] = []*operation{o}
	}
}

func (s StageOperations) get(stage stage) []*operation {
	return s[stage]
}

func (s StageOperations) iterable() map[stage][]*operation {
	return s
}

func findStructureOperations(structType *ast.StructType) StageOperations {
	operations := make(StageOperations)

	for _, field := range structType.Fields.List {
		if field.Comment != nil && strings.HasPrefix(field.Comment.Text(), "GEN:") {
			// Remove all spaces (including newlines) and the "GEN:" prefix from the string
			trimmedComment := removeWhitespace(field.Comment.Text())
			trimmedComment = strings.Replace(trimmedComment, "GEN:", "", 1)

			// Different operations stages will be separated by a semi-colon
			commands := strings.Split(trimmedComment, ";")
			for _, command := range commands {
				stage, err := parseStage(command)
				if err != nil {
					panic(err)
				}

				actions, err := parseActions(command)
				if err != nil {
					panic(err)
				}

				for _, action := range actions {
					// one line can have many comma separated fields
					for _, fieldName := range field.Names {
						operations.add(&operation{
							field:     fieldName.String(),
							fieldType: parsing.GoTypeFromExpr(field.Type),
							stage:     stage,
							action:    action,
						})
					}
				}
			}
		}
	}

	return operations
}

// Attempts to parse what stage the command is for and validates that after the stage token there is an open and closed
// parentheses
func parseStage(command string) (stage stage, err error) {
	for potentialStageToken, potentialStage := range tokenToStage {
		if strings.HasPrefix(command, potentialStageToken) {
			stage = potentialStage

			// Check that the stage string is followed by an open parentheses
			if command[len(potentialStageToken):len(potentialStageToken)+1] != "(" {
				return -1, errors.New("parseStage: command '" + command + "' missing open parentheses")
			}

			// Check that the stage string ends in a closing parentheses
			if command[len(command)-1:] != ")" {
				return -1, errors.New("parseStage: command: '" + command + "' missing closing parentheses")
			}

			return stage, nil
		}
	}

	return -1, errors.New("parseStage: stage from command: '" + command + "' is not valid")
}

func parseActions(command string) (actions []action, err error) {
	actions = make([]action, 0)

	// all text between the two parentheses of a command
	actionsString := command[strings.Index(command, "(")+1 : len(command)-1]

	for _, actionToken := range strings.Split(actionsString, ",") {
		validFound := false
		for potentialActionToken, potentialAction := range tokenToAction {
			if actionToken == potentialActionToken {
				actions = append(actions, potentialAction)
				validFound = true
				break
			}
		}

		if !validFound {
			return nil, errors.New("parseActions: action '" + actionToken + "' is not valid")
		}
	}

	return
}

// fastest implementation from https://stackoverflow.com/a/32081891
func removeWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
