package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ========================== CONSTANTS & TYPES ==========================
// Parking spot types
const (
	MAXI     = 0  // Large parking spot
	STANDARD = 1  // Standard parking spot
)

// Vehicle types
const (
	CAR      = 0
	CAMPER   = 1
	SNOWPLOW = 2
)

// System capacities
const (
	STANDARD_SPOTS = 10  // Standard parking spots
	MAXI_SPOTS     = 5   // Large parking spots
	NUM_TOURISTS   = 25  // Total tourists (cars + campers)
	MAXBUFF        = 100 // Max channel buffer size
)

// Traffic directions
const (
	UPHILL   = 0
	DOWNHILL = 1
)

// ========================== CHANNELS ==========================
var (
	// Uphill traffic channels (vehicle type -> channel)
	startUphill [3]chan int      // Request to enter uphill
	endUphill   [3]chan int      // Notify end of uphill journey
	
	// Downhill traffic channels (vehicle type -> channel)
	startDownhill [3]chan Parking // Request to enter downhill (with parking info)
	endDownhill   [3]chan int     // Notify end of downhill journey
	
	// Acknowledgment channels
	ACK_tourist  [NUM_TOURISTS]chan int // Per-tourist ACK channels
	ACK_snowplow = make(chan int, MAXBUFF) // Snowplow ACK channel
	
	// Termination channels
	done               = make(chan bool)          // Unbuffered for sync
	terminate          = make(chan bool)          // Castle termination
	terminateSnowplow  = make(chan bool)          // Snowplow termination
)

// Parking struct carries parking spot information
type Parking struct {
	index       int // Vehicle ID
	parkingType int // Spot type (only relevant for cars)
}

// ========================== HELPER FUNCTIONS ==========================
// Conditional channel selector for int channels
func when(b bool, c chan int) chan int {
	if !b {
		return nil // Channel disabled
	}
	return c
}

// Conditional channel selector for Parking channels
func whenParking(b bool, c chan Parking) chan Parking {
	if !b {
		return nil
	}
	return c
}

// Random sleep to simulate real-world delays
func sleepRandTime(timeLimit int) {
	if timeLimit > 0 {
		time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
	}
}

// ========================== GOROUTINES ==========================
// Tourist (car/camper) behavior
func tourist(index int, vehicleType int) {
	var parkingType int
	
	// Request uphill access
	startUphill[vehicleType] <- index
	parkingType = <-ACK_tourist[index] // Wait for parking assignment
	
	// Simulate uphill journey
	sleepRandTime(3)
	
	// Notify uphill completion
	endUphill[vehicleType] <- index
	<-ACK_tourist[index] // Wait for confirmation
	
	// Visit the castle
	sleepRandTime(4)
	
	// Request downhill access
	startDownhill[vehicleType] <- Parking{index, parkingType}
	<-ACK_tourist[index] // Wait for confirmation
	
	// Simulate downhill journey
	sleepRandTime(2)
	
	// Notify downhill completion
	endDownhill[vehicleType] <- index
	<-ACK_tourist[index]
	done <- true // Signal completion
}

// Snowplow maintenance vehicle
func snowplow() {
	var res int
	sleepRandTime(4) // Initial delay
	
	for {
		// Request downhill access
		startDownhill[SNOWPLOW] <- Parking{-1, -1}
		res = <-ACK_snowplow
		
		if res == -1 { // Termination signal
			fmt.Printf("[snowplow] terminating...\n")
			done <- true
			return
		}
		
		// Downhill journey
		fmt.Printf("[snowplow] entered downhill direction\n")
		sleepRandTime(2)
		endDownhill[SNOWPLOW] <- 1
		res = <-ACK_snowplow
		
		// Request uphill return
		sleepRandTime(8)
		startUphill[SNOWPLOW] <- 1
		res = <-ACK_snowplow
		fmt.Printf("[snowplow] entered uphill direction\n")
		
		// Uphill journey
		sleepRandTime(2)
		endUphill[SNOWPLOW] <- 1
		res = <-ACK_snowplow
		fmt.Printf("[snowplow] entered the castle successfully!\n")
		sleepRandTime(8)
	}
}

// Castle (central coordinator)
func castle() {
	var (
		index           int
		p               Parking
		stop            = false
		numCampersOnRoad= [2]int{0, 0} // [UPHILL, DOWNHILL]
		numCarsOnRoad   = [2]int{0, 0}
		snowplowActive  = false
		freeStandardSpots = STANDARD_SPOTS
		freeMaxiSpots     = MAXI_SPOTS
	)
	
	fmt.Printf("[castle] The road is open!\n")
	
	for {
		select {
		// === UPHILL REQUESTS ===
		case index = <-when(
			freeMaxiSpots > 0 && 
			numCampersOnRoad[DOWNHILL]+numCarsOnRoad[DOWNHILL] == 0 &&
			!snowplowActive &&
			len(startDownhill[CAMPER])+len(startDownhill[CAR])+len(startDownhill[SNOWPLOW]) == 0, 
			startUphill[CAMPER]):
			// Camper entering uphill
			freeMaxiSpots--
			numCampersOnRoad[UPHILL]++
			fmt.Printf("[castle] CAMPER %d entered uphill\n", index)
			ACK_tourist[index] <- MAXI

		case index = <-when(
			(freeStandardSpots+freeMaxiSpots > 0) &&
			numCampersOnRoad[DOWNHILL] == 0 &&
			!snowplowActive &&
			len(startUphill[CAMPER]) == 0 &&
			len(startDownhill[CAMPER])+len(startDownhill[CAR])+len(startDownhill[SNOWPLOW]) == 0, 
			startUphill[CAR]):
			// Car entering uphill
			parkingType := STANDARD
			if freeStandardSpots > 0 {
				freeStandardSpots--
			} else {
				freeMaxiSpots--
				parkingType = MAXI
			}
			numCarsOnRoad[UPHILL]++
			fmt.Printf("[castle] CAR %d entered uphill\n", index)
			ACK_tourist[index] <- parkingType

		case <-when(
			(numCampersOnRoad[DOWNHILL]+numCarsOnRoad[DOWNHILL]+numCampersOnRoad[UPHILL]+numCarsOnRoad[UPHILL] == 0) &&
			(len(startUphill[CAMPER])+len(startUphill[CAR]) == 0) && 
			(len(startDownhill[CAMPER])+len(startDownhill[CAR]) == 0), 
			startUphill[SNOWPLOW]):
			// Snowplow entering uphill
			snowplowActive = true
			fmt.Printf("[castle] SNOWPLOW entered uphill\n")
			ACK_snowplow <- 1

		// === UPHILL COMPLETIONS ===
		case index = <-endUphill[CAMPER]:
			numCampersOnRoad[UPHILL]--
			fmt.Printf("[castle] CAMPER %d arrived\n", index)
			ACK_tourist[index] <- 1

		case index = <-endUphill[CAR]:
			numCarsOnRoad[UPHILL]--
			fmt.Printf("[castle] CAR %d arrived\n", index)
			ACK_tourist[index] <- 1

		case <-endUphill[SNOWPLOW]:
			snowplowActive = false
			fmt.Printf("[castle] SNOWPLOW arrived\n")
			ACK_snowplow <- 1

		// === DOWNHILL REQUESTS ===
		case p = <-whenParking(
			(numCampersOnRoad[UPHILL]+numCarsOnRoad[UPHILL] == 0) &&
			!snowplowActive &&
			len(startDownhill[SNOWPLOW]) == 0, 
			startDownhill[CAMPER]):
			// Camper leaving
			numCampersOnRoad[DOWNHILL]++
			freeMaxiSpots++
			fmt.Printf("[castle] CAMPER %d exiting\n", p.index)
			ACK_tourist[p.index] <- 1

		case p = <-whenParking(
			(numCampersOnRoad[UPHILL] == 0) &&
			!snowplowActive &&
			len(startDownhill[SNOWPLOW])+len(startDownhill[CAMPER]) == 0, 
			startDownhill[CAR]):
			// Car leaving
			numCarsOnRoad[DOWNHILL]++
			if p.parkingType == MAXI {
				freeMaxiSpots++
			} else {
				freeStandardSpots++
			}
			fmt.Printf("[castle] CAR %d exiting\n", p.index)
			ACK_tourist[p.index] <- 1

		case <-whenParking(
			!stop &&
			(numCampersOnRoad[DOWNHILL]+numCarsOnRoad[DOWNHILL]+numCampersOnRoad[UPHILL]+numCarsOnRoad[UPHILL] == 0), 
			startDownhill[SNOWPLOW]):
			// Snowplow exiting
			snowplowActive = true
			fmt.Printf("[castle] SNOWPLOW exiting\n")
			ACK_snowplow <- 1

		// === DOWNHILL COMPLETIONS ===
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

		// === TERMINATION HANDLING ===
		case <-terminateSnowplow:
			stop = true
			fmt.Printf("[castle] Stopping snowplow...\n")

		case <-whenParking(stop, startDownhill[SNOWPLOW]):
			ACK_snowplow <- -1

		case <-terminate:
			fmt.Printf("[castle] Terminating...\n")
			done <- true
			return
		}
	}
}

// ========================== MAIN ==========================
func main() {
	rand.Seed(time.Now().UnixNano())
	
	// Channel initialization
	for i := 0; i < 3; i++ {
		startUphill[i] = make(chan int, MAXBUFF)
		endUphill[i] = make(chan int, MAXBUFF)
		startDownhill[i] = make(chan Parking, MAXBUFF)
		endDownhill[i] = make(chan int, MAXBUFF)
	}
	
	// ACK channels initialization
	for i := 0; i < NUM_TOURISTS; i++ {
		ACK_tourist[i] = make(chan int, MAXBUFF)
	}
	
	// Start system components
	go castle()
	go snowplow()
	for i := 0; i < NUM_TOURISTS; i++ {
		vehicleType := rand.Intn(2) // 0=car, 1=camper
		go tourist(i, vehicleType)
	}
	
	// Wait for tourists to finish
	for i := 0; i < NUM_TOURISTS; i++ {
		<-done
	}
	
	// Shutdown sequence
	terminateSnowplow <- true
	<-done            // Wait for snowplow
	terminate <- true // Signal castle
	<-done            // Wait for castle
	fmt.Println("[main] All goroutines terminated")
}

/*
KEY SYSTEM LOGIC:

1. Concurrency Model:
- Uses Go's channel-based concurrency
- Castle goroutine acts as central traffic controller
- 3 types of actors: Tourists (cars/campers), Snowplow, Castle

2. Channel Usage:
- Buffered channels (MAXBUFF=100) for most operations to handle bursts
- Unbuffered channels (done, terminate) for precise synchronization
- Separate channels for different vehicle types and directions

3. Traffic Rules:
- Uphill priority: No downhill traffic when vehicles are going up
- Snowplow has exclusive access when operating
- Campers require MAXI spots, cars can use STANDARD or MAXI
- Strict vehicle ordering to prevent collisions

4. Parking Management:
- Tracks available STANDARD/MAXI spots
- Dynamically assigns parking types to cars based on availability
- Releases spots when vehicles exit

5. Snowplow Behavior:
- Operates in cycles: downhill -> uphill -> repeat
- Can only run when no other vehicles are present
- Termination handled via dedicated channel

6. Termination Sequence:
1. Wait for all tourists to complete
2. Signal snowplow termination
3. Wait for snowplow confirmation
4. Terminate castle
5. Confirm full shutdown

TYPICAL EXAM SCENARIO:
Question: "Add emergency vehicle priority that can interrupt all traffic"
Solution:
- Add new vehicle type constant EMERGENCY
- Modify when() conditions to allow emergency vehicles through
- Add priority channel handling in castle's select
- Implement preemption logic for existing vehicles
*/
