package elevator

import (
	"fmt"

	"SingleElevator/elevio"
)

func HasOrders(e *ElevState) bool {
	for f := range e.Orders {
		for b := range e.Orders[f] {
			if e.Orders[f][b] {
				return true
			}
		}
	}
	return false
}

func HasOrdersAbove(e *ElevState) bool {
	for f := e.CurrFloor + 1; f < len(e.Orders); f++ {
		for b := range e.Orders[f] {
			if e.Orders[f][b] {
				return true
			}
		}
	}
	return false
}

func HasOrdersBelow(e *ElevState) bool {
	for f := 0; f < e.CurrFloor; f++ {
		for b := range e.Orders[f] {
			if e.Orders[f][b] {
				return true
			}
		}
	}
	return false
}

func ShouldStop(e *ElevState) bool {
	switch e.Dir {
	case elevio.Down:
		return e.Orders[e.CurrFloor][elevio.HallDown] ||
			e.Orders[e.CurrFloor][elevio.Cab] ||
			!HasOrdersBelow(e) // Reached "dead-end"
	case elevio.Up:

		return e.Orders[e.CurrFloor][elevio.HallUp] ||
			e.Orders[e.CurrFloor][elevio.Cab] ||
			!HasOrdersAbove(e) // Reached "dead-end"
	case elevio.Stop:
		return true
	}

	return true
}

func ChooseDirection(e *ElevState) (elevio.MotorDirection, Behavior) {
	if e.CurrFloor < e.Target.Floor {
		return elevio.Up, Moving
	}

	if e.CurrFloor > e.Target.Floor {
		return elevio.Down, Moving
	}

	return elevio.Stop, DoorOpen
}

func ClearAtCurrentFloor(e *ElevState) {
	e.Orders[e.CurrFloor][elevio.Cab] = false
}

func PrintOrders(e *ElevState) {
	for f := range e.Orders {
		for b := range e.Orders[f] {
			if e.Orders[f][b] {
				fmt.Printf("Order at floor %d, button %v\n", f, elevio.ButtonType(b))
			}
		}
	}
}
