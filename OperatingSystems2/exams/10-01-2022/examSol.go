package main

import (
	"fmt"
	"math/rand"
	"time"
)

//General constants
const NUM_OFFICES = 5         //Number of offices (consultants)
const MAX_WAITING_ROOM = 10   //Capacity of the waiting room
const NUM_USERS = 100         //Total number of users
const MAX_BUFFER = 50         //Buffer size for channels

//Constants for user priority in the waiting room
const USER_TYPES = 3
const ADMIN = 0          // Administrator
const PRIVATE_SINGLE = 1 // Private owner, without accompanying person
const PRIVATE_WITH = 2   // Private owner, with accompanying person

//Constants for service priority in offices
const FINANCE_TYPES = 2
const SUPERBONUS = 0     // Superbonus-related services
const OTHER = 1          // Other financial services

//Struct to represent user data
type User struct {
	id          int      // User ID
	userType    int      // Type of user (administrator, private, etc.)
	serviceType int      // Type of financial service (superbonus, other)
	reply       chan int // Channel for user replies
}

//General communication channels
var terminate chan bool // Channel to signal termination
var done chan bool      // Channel to signal task completion

//Specific communication channels
var enterWaitingRoom [USER_TYPES]chan User      //Channels for entering the waiting room by user type
var enterOffice      [FINANCE_TYPES]chan User   //Channels for entering offices based on service type
var exitOffice chan int                         //Channel for exiting the office

//Utility function: simulate random sleep
func sleepRandom() {
	time.Sleep(time.Duration(1e9 * ((rand.Intn(30)) + 1)))     //Random sleep between 1-30 seconds, you need the 1e9 because it's in nanosecond
}

//Utility function: conditional channel activation (Logic guard)
func when(condition bool, ch chan User) chan User {
	if !condition {
		return nil // Return nil channel if condition is false
	}
	return ch // Return the original channel if condition is true
}

func server() {
	var i int
	waitingRoomCount := 0 // Number of people in the waiting room
	officesOccupied := 0  // Number of occupied offices
	var officeOccupied [NUM_OFFICES]bool // Tracks whether each office is occupied
	for i = 0; i < NUM_OFFICES; i++ {
		officeOccupied[i] = false // Initialize all offices as unoccupied
	}
	fmt.Printf("The consulting service is open.\n\n")
	for {
		select {
		//Case 1: An administrator enters the waiting room
		case request := <-when(waitingRoomCount < MAX_WAITING_ROOM, enterWaitingRoom[ADMIN]):
			waitingRoomCount += 1
			fmt.Printf("SERVER: Administrator %d entered the waiting room.\n", request.id)
			request.reply <- 1 // Notify the client that they entered successfully

		//Case 2: A private individual without an accompanist enters the waiting room
		case request := <-when(waitingRoomCount < MAX_WAITING_ROOM && len(enterWaitingRoom[ADMIN]) == 0, enterWaitingRoom[PRIVATE_SINGLE]):
			waitingRoomCount += 1
			fmt.Printf("SERVER: Private individual (alone) %d entered the waiting room.\n", request.id)
			request.reply <- 1

		//Case 3: A private individual with an accompanist enters the waiting room
		case request := <-when(waitingRoomCount+2 <= MAX_WAITING_ROOM && len(enterWaitingRoom[ADMIN]) == 0 && len(enterWaitingRoom[PRIVATE_SINGLE]) == 0, enterWaitingRoom[PRIVATE_WITH]):
			waitingRoomCount += 2
			fmt.Printf("SERVER: Private individual with accompanist %d entered the waiting room.\n", request.id)
			request.reply <- 1

		//Case 4: A client enters an office for a Superbonus service
		case request := <-when(officesOccupied < NUM_OFFICES, enterOffice[SUPERBONUS]):
			for i = 0; i < NUM_OFFICES; i++ { // Find the first available office
				if !officeOccupied[i] {
					break
				}
			}
			officeOccupied[i] = true
			officesOccupied++
			if request.userType == PRIVATE_WITH {
				waitingRoomCount -= 2 // Free up 2 spots in the waiting room
				fmt.Printf("SERVER: Private individual with accompanist for Superbonus %d entered office %d.\n", request.id, i)
			} else {
				waitingRoomCount -= 1 // Free up 1 spot in the waiting room
				if request.userType == ADMIN {
					fmt.Printf("SERVER: Administrator for Superbonus %d entered office %d.\n", request.id, i)
				} else {
					fmt.Printf("SERVER: Private individual (alone) for Superbonus %d entered office %d.\n", request.id, i)
				}
			}
			request.reply <- i // Send the office number to the client

		//Case 5: A client enters an office for a different service
		case request := <-when(officesOccupied < NUM_OFFICES && len(enterOffice[SUPERBONUS]) == 0, enterOffice[OTHER]):
			for i = 0; i < NUM_OFFICES; i++ { // Find the first available office
				if !officeOccupied[i] {
					break
				}
			}
			officeOccupied[i] = true
			officesOccupied++
			if request.userType == PRIVATE_WITH {
				waitingRoomCount -= 2 // Free up 2 spots in the waiting room
				fmt.Printf("SERVER: Private individual with accompanist for Other service %d entered office %d.\n", request.id, i)
			} else {
				waitingRoomCount -= 1 // Free up 1 spot in the waiting room
				if request.userType == ADMIN {
					fmt.Printf("SERVER: Administrator for Other service %d entered office %d.\n", request.id, i)
				} else {
					fmt.Printf("SERVER: Private individual (alone) for Other service %d entered office %d.\n", request.id, i)
				}
			}
			request.reply <- i // Send the office number to the client

		//Case 6: A client exits an office
		case release := <-exitOffice:
			officeOccupied[release] = false // Mark the office as unoccupied
			officesOccupied--

		//Case 7: Terminate the service
		case <-terminate:
			fmt.Printf("The consulting service is closing.\n")
			done <- true
			return
		}
	}
}

func user(id int) {
	userType := rand.Intn(USER_TYPES)   // Administrator, individual, or accompanied
	serviceType := rand.Intn(FINANCE_TYPES) // Type of financing (Superbonus or Other)
	var ack = make(chan int)
	var request User
	request := User{id, userType, serviceType, ack}
	sleepRandom()
	enterWaitingRoom[userType] <- request
	<-request.reply
	enterOffice[serviceType] <- request
	officeAssigned := <-request.reply
	sleepRandom()
	exitOffice <- officeAssigned
	fmt.Printf("User [%d]: I have exited office %d. Terminating.\n", id, officeAssigned)
	done <- true
}

func main() {
	rand.Seed(time.Now().UnixNano())
	terminate = make(chan bool)
	done = make(chan bool)
	for i := 0; i < USER_TYPES; i++ {
		enterWaitingRoom[i] = make(chan User, MAX_BUFFER)
	}
	for i := 0; i < FINANCE_TYPES; i++ {
		enterOffice[i] = make(chan User, MAX_BUFFER)
	}
	exitOffice = make(chan int, MAX_BUFFER)
	go server()
	for id := 0; id < NUM_USERS; id++ {
		go user(id)
	}
	for id := 0; id < NUM_USERS; id++ {
		<-done
	}
	terminate <- true
	<-done
}
