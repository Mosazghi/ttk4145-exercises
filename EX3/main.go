package main

import (
	"SingleElevator/elevio"
	"flag"
	"fmt"
	"time"
)

type Behavior int

const (
	IDLE Behavior = iota
	MOVING
	DOOROPEN
)

func ReadInitialButtons() [4][3]bool {
	var orders [4][3]bool
	for f := range orders {
		for b := range orders[f] {
			if elevio.GetButton(elevio.ButtonType(b), f) {
				orders[f][b] = true
			}
		}
	}
	return orders
}

var numFloors = 4

type Order struct {
	floor int
	rType elevio.ButtonType
}

type ElevState struct {
	target    Order
	currFloor int
	dir       elevio.MotorDirection
	behavior  Behavior
	orders    [4][3]bool
}

func main() {
	portNum := flag.String("port", "15657", "specify port number")
	flag.Parse()
	elevio.Init("localhost:"+*portNum, numFloors)

	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)

	stateMachine(drvButtons, drvFloors, drvObstr, drvStop)
}

func stateMachine(drvButtons chan elevio.ButtonEvent, drvFloors chan int, drvObst chan bool, drvStop chan bool) {
	s := ElevState{
		Order{-1, elevio.BT_Cab},
		elevio.GetFloor(),
		elevio.MD_Stop,
		IDLE,
		ReadInitialButtons(),
	}

	_ = elevio.MD_Stop
	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()

	printState := func() {
		fmt.Printf("State: %+v\n", s)
	}

	for {
		select {
		case a := <-drvButtons:
			fmt.Printf("[BTNS] %+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			s.orders[a.Floor][a.Button] = true
			printState()

		case a := <-drvFloors:
			fmt.Printf("[FLOORS] %+v\n", a)
			if a == numFloors-1 {
				s.dir = elevio.MD_Down
			} else if a == 0 {
				s.dir = elevio.MD_Up
			}
			s.currFloor = a
			elevio.SetFloorIndicator(s.currFloor)
			printState()

		case a := <-drvObst:
			fmt.Printf("[OBSTR] %+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(s.dir)
			}
			printState()

		case a := <-drvStop:
			fmt.Printf("[STOP] %+v\n", a)
			for f := range numFloors {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			printState()

		case <-ticker.C:
			switch s.behavior {
			case IDLE:
				if s.target.floor == -1 {
					// Check for any active orders
					found := false
					for f := 0; f < numFloors && !found; f++ {
						for b := elevio.ButtonType(0); b < 3 && !found; b++ {
							if s.orders[f][b] {
								s.target.floor = f
								s.target.rType = b
								found = true
							}
						}
					}

					if found {
						fmt.Println("Target floor set to:", s.target.floor)
						elevio.SetDoorOpenLamp(false)
						s.behavior = MOVING
					}
				}
			case MOVING:
				tF := s.target.floor
				cF := s.currFloor
				// fmt.Println("Moving. Current floor:", cF, "Target floor:", tF)

				// only decide direction when arrived at NEW floor
				// if cF != prevFloor {
				// prevFloor = cF
				if tF == cF {
					s.behavior = DOOROPEN
					fmt.Printf("Arrived at floor: %v, curr=%v\n", tF, cF)
					continue
				} else if tF > cF {
					s.dir = elevio.MD_Up
					fmt.Println("Dir=UP")
				} else {
					s.dir = elevio.MD_Down
					fmt.Println("Dir=DOWN")
				}

				// if prevDir != s.dir {
				// fmt.Printf("Changing dir=%v, from %v\n", s.dir, prevDir)
				elevio.SetMotorDirection(s.dir)
				// prevDir = s.dir
				// }
				// }

			case DOOROPEN:
				fmt.Println("Door open at floor:", s.currFloor)
				elevio.SetStopLamp(true)
				for b := range s.orders[s.currFloor] {
					s.orders[s.currFloor][b] = false
					elevio.SetButtonLamp(elevio.ButtonType(b), s.currFloor, false)
				}
				s.target.floor = -1
				s.target.rType = elevio.BT_Cab
				s.dir = elevio.MD_Stop
				elevio.SetMotorDirection(s.dir)
				s.behavior = IDLE
				elevio.SetDoorOpenLamp(true)
				time.Sleep(time.Second)
				elevio.SetDoorOpenLamp(false)
				printState()
			}
		}
	}
}
