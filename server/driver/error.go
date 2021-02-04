package driver

import "fmt"

type IncorrectActionError struct {
	seat Seat
	upperActionIndex int
}

func (e IncorrectActionError) Error() string {
	return fmt.Sprintf("an action is required for seat [%d] within range [0 to %d]", e.seat, e.upperActionIndex)
}

type TransitionLimitReachedError struct {
	transitionLimit int
	StateHistory []string
}

func (e TransitionLimitReachedError) Error() string {
	return fmt.Sprintf("transitioning to next action state took more than [%d] steps", e.transitionLimit)
}

type GameLogicError struct {
	Err error
}

func (e GameLogicError) Error() string {
	return fmt.Sprintf("game logic error: %s", e.Err.Error())
}