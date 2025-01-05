package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Parking spot types:
const MAXI = 0
const STANDARD = 1

// Vehicle types:
const CAR = 0
const CAMPER = 1
const SNOWPLOW = 2

const STANDARD_SPOTS = 10 // Standard parking spots for cars
const MAXI_SPOTS = 5      // Maxi parking spots for cars or campers

const NUM_TOURISTS = 25 // Number of tourists (campers and cars)
const MAXBUFF = 100

// Traffic direction
const UPHILL = 0
const DOWNHILL = 1

// Uphill/downhill road channels (car = 0, camper = 1, snowplow = 2)
var startUphill [3]chan int
var endUphill [3]chan int
var startDownhill [3]chan Parking
var endDownhill [3]chan int

// Acknowledgment channels
var ACK_tourist [NUM_TOURISTS]chan int
var ACK_snowplow = make(chan int, MAXBUFF)

// Termination channels
var done = make(chan bool)
var terminate = make(chan bool)
var terminateSnowplow = make(chan bool)

type Parking struct {
	index        int // ID of the vehicle
	parkingType  int // Type of parking spot occupied (standard/maxi), relevant only for cars
}

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func whenParking(b bool, c chan Parking) chan Parking {
	if !b {
		return nil
	}
	return c
}

func sleepRandTime(timeLimit int) {
	if timeLimit > 0 {
		time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
	}
}

// Goroutines:
func tourist(index int, vehicleType int) {
	var parkingType int
	sleepRandTime(4)
	startUphill[vehicleType] <- index
	parkingType = <-ACK_tourist[index]
	sleepRandTime(3) // Time to go uphill
	endUphill[vehicleType] <- index
	<-ACK_tourist[index]
	// Visit the castle
	sleepRandTime(4)
	p := Parking{index, parkingType}
	startDownhill[vehicleType] <- p
	<-ACK_tourist[index]
	sleepRandTime(2)
	endDownhill[vehicleType] <- index
	<-ACK_tourist[index]
	done <- true
	return
}

func snowplow() {
	var res int
	sleepRandTime(4)
	for {
		startDownhill[SNOWPLOW] <- Parking{-1, -1}
		res = <-ACK_snowplow
		if res == -1 {
			fmt.Printf("[snowplow] terminating...\n")
			done <- true
			return
		}
		fmt.Printf("[snowplow] entered downhill direction\n")
		sleepRandTime(2)
		endDownhill[SNOWPLOW] <- 1
		res = <-ACK_snowplow
		sleepRandTime(8)
		startUphill[SNOWPLOW] <- 1
		res = <-ACK_snowplow
		fmt.Printf("[snowplow] entered uphill direction\n")
		sleepRandTime(2)
		endUphill[SNOWPLOW] <- 1
		res = <-ACK_snowplow
		fmt.Printf("[snowplow] entered the castle successfully!\n")
		sleepRandTime(8)
	}
}

func castle() {
	var index int
	var p Parking
	var stop = false
	numCampersOnRoad := [2]int{0, 0}
	numCarsOnRoad := [2]int{0, 0}
	var snowplowActive = false
	var freeStandardSpots, freeMaxiSpots int = STANDARD_SPOTS, MAXI_SPOTS
	fmt.Printf("[castle] The road is open!\n")
	for {
		select {
		case index = <-when((freeMaxiSpots > 0) && (numCampersOnRoad[DOWNHILL]+numCarsOnRoad[DOWNHILL] == 0) && (!snowplowActive) &&
			(len(startDownhill[CAMPER])+len(startDownhill[CAR])+len(startDownhill[SNOWPLOW]) == 0), startUphill[CAMPER]):
			freeMaxiSpots--
			numCampersOnRoad[UPHILL]++
			fmt.Printf("[castle] CAMPER %d entered uphill direction\n", index)
			ACK_tourist[index] <- MAXI
		case index = <-when((freeStandardSpots+freeMaxiSpots > 0) && (numCampersOnRoad[DOWNHILL] == 0) && (!snowplowActive) && (len(startUphill[CAMPER]) == 0) &&
			(len(startDownhill[CAMPER])+len(startDownhill[CAR])+len(startDownhill[SNOWPLOW]) == 0), startUphill[CAR]):
			parkingType := -1
			if freeStandardSpots > 0 {
				freeStandardSpots--
				parkingType = STANDARD
			} else {
				freeMaxiSpots--
				parkingType = MAXI
			}
			numCarsOnRoad[UPHILL]++
			fmt.Printf("[castle] CAR %d entered uphill direction\n", index)
			ACK_tourist[index] <- parkingType
		case <-when((numCampersOnRoad[DOWNHILL]+numCarsOnRoad[DOWNHILL]+numCampersOnRoad[UPHILL]+numCarsOnRoad[UPHILL] == 0) &&
			(len(startUphill[CAMPER])+len(startUphill[CAR]) == 0) && (len(startDownhill[CAMPER])+len(startDownhill[CAR]) == 0), startUphill[SNOWPLOW]):
			snowplowActive = true
			fmt.Printf("[castle] SNOWPLOW entered uphill direction\n")
			ACK_snowplow <- 1
		case index = <-endUphill[CAMPER]:
			numCampersOnRoad[UPHILL]--
			fmt.Printf("[castle] CAMPER %d entered the castle\n", index)
			ACK_tourist[index] <- 1
		// Car
		case index = <-endUphill[CAR]:
			numCarsOnRoad[UPHILL]--
			fmt.Printf("[castle] CAR %d entered the castle\n", index)
			ACK_tourist[index] <- 1
		// Snowplow
		case <-endUphill[SNOWPLOW]:
			snowplowActive = false
			fmt.Printf("[castle] SNOWPLOW entered the castle\n")
			ACK_snowplow <- 1
		case p = <-whenParking((numCampersOnRoad[UPHILL]+numCarsOnRoad[UPHILL] == 0) && (!snowplowActive) && (len(startDownhill[SNOWPLOW]) == 0), startDownhill[CAMPER]):
			numCampersOnRoad[DOWNHILL]++
			freeMaxiSpots++
			fmt.Printf("[castle] CAMPER %d entered downhill direction\n", index)
			ACK_tourist[p.index] <- 1
		case p = <-whenParking((numCampersOnRoad[UPHILL] == 0) && (!snowplowActive) && (len(startDownhill[SNOWPLOW])+len(startDownhill[CAMPER]) == 0), startDownhill[CAR]):
			numCarsOnRoad[DOWNHILL]++
			if p.parkingType == MAXI {
				freeMaxiSpots++
			} else {
				freeStandardSpots++
			}
			fmt.Printf("[castle] CAR %d entered downhill direction\n", index)
			ACK_tourist[p.index] <- 1
		case <-whenParking((stop == false) && (numCampersOnRoad[DOWNHILL]+numCarsOnRoad[DOWNHILL]+numCampersOnRoad[UPHILL]+numCarsOnRoad[UPHILL] == 0), startDownhill[SNOWPLOW]):
			snowplowActive = true
			fmt.Printf("[castle] SNOWPLOW entered downhill direction\n")
			ACK_snowplow <- 1
		case index = <-endDownhill[CAMPER]:
			numCampersOnRoad[DOWNHILL]--
			fmt.Printf("[castle] CAMPER %d exited\n", index)
			ACK_tourist[index] <- 1
		case index = <-endDownhill[CAR]:
			numCarsOnRoad[DOWNHILL]--
			fmt.Printf("[castle] CAR %d exited\n", index)
			ACK_tourist[index] <- 1
		case <-endDownhill[SNOWPLOW]:
			snowplowActive = false
			fmt.Printf("[castle] SNOWPLOW exited\n")
			ACK_snowplow <- 1
		// Termination:
		case <-terminateSnowplow:
			stop = true
			fmt.Printf("[castle] tourists terminated, stopping snowplow...\n")
		case <-whenParking((stop == true), startDownhill[SNOWPLOW]):
			ACK_snowplow <- -1
		case <-terminate:
			fmt.Printf("[castle] Terminating...\n")
			done <- true
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed for random generation
	// Initialize channels
	for i := 0; i < 3; i++ {
		startUphill[i] = make(chan int, MAXBUFF)
		endUphill[i] = make(chan int, MAXBUFF)
		startDownhill[i] = make(chan Parking, MAXBUFF)
		endDownhill[i] = make(chan int, MAXBUFF)
	}
	for i := 0; i < NUM_TOURISTS; i++ {
		ACK_tourist[i] = make(chan int, MAXBUFF)
	}
	// Start goroutines
	go castle()
	go snowplow()
	for i := 0; i < NUM_TOURISTS; i++ {
		vehicleType := rand.Intn(2) // 0 = CAR, 1 = CAMPER
		go tourist(i, vehicleType)
	}
	// Wait for all tourist goroutines to finish
	for i := 0; i < NUM_TOURISTS; i++ {
		<-done
	}
	// Signal snowplow termination
	terminateSnowplow <- true
	// Wait for snowplow to terminate
	<-done
	// Terminate castle
	terminate <- true
	<-done
	fmt.Println("[main] All goroutines have terminated. Program finished.")
}
