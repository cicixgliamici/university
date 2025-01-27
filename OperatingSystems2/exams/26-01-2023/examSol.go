package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Constants defining system parameters
const MAX_BUFFER = 100        // Max buffer size for channels
const MAX_CLIENTS = 100       // Max number of clients

// Bottle types
const SmallBottle = 0         // 0.5 liters, costs 0.10
const LargeBottle = 1         // 1.5 liters, costs 0.20

// Bottle capacities
const CapacitySmall = 0.5     // Small bottle capacity
const CapacityLarge = 1.5     // Large bottle capacity

const TankCapacity = 50.0     // Total tank capacity in liters

// Coin boxes for different coin types
const SmallCoinsBox = 0       // Box for 10-cent coins (SmallBottle)
const LargeCoinsBox = 1       // Box for 20-cent coins (LargeBottle)

// Max coins before needing a refill
const MaxSmallCoins = 15      // Max 10-cent coins before refill
const MaxLargeCoins = 20      // Max 20-cent coins before refill

// Channels for client requests
var start_request [2]chan request // Buffered channels for starting requests (index 0: Small, 1: Large)
var end_request = make(chan request, MAX_BUFFER) // Buffered channel for ending requests

// Channels for operator actions
// start refill is buffered because the operator can request multiple refills.
// terminate Operator is not buffered to ensure synchronous handshake.
var start_refill = make(chan int, MAX_BUFFER)    // Operator starts refill
var end_refill = make(chan int, MAX_BUFFER)      // Operator ends refill
var ack_operator = make(chan int, MAX_BUFFER)    // Acknowledgment for operator

// Termination channels (unbuffered for synchronization)
var done = make(chan bool)              // Signals client completion
var terminate = make(chan bool)         // Signals waterStation to terminate
var terminateOperator = make(chan bool) // Signals operator to terminate

// Request structure for client requests
type request struct {
	index int     // Client ID
	kind  int     // Bottle type (SmallBottle/LargeBottle)
	ack   chan int // Acknowledgment channel for synchronization
}

// Helper function to conditionally enable a channel based on a boolean condition
func whenRequest(b bool, c chan request) chan request {
	if !b {
		return nil // Channel disabled if condition is false
	}
	return c // Channel enabled if condition is true
}

// Helper function to conditionally enable an int channel
func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

// Simulates random delays for realistic concurrency behavior
func sleepRandomTime(limit int) {
	if limit > 0 {
		time.Sleep(time.Duration(rand.Intn(limit)+1) * time.Second)
	}
}

// Client goroutine: Simulates client behavior
func client(index int) {
	kind := rand.Intn(2) // Randomly choose bottle type
	r := request{index, kind, make(chan int)}
	sleepRandomTime(2) // Simulate payment time

	// Send request to appropriate channel
	if r.kind == SmallBottle {
		fmt.Printf("[client %d] requested a small bottle\n", index)
	} else {
		fmt.Printf("[client %d] requested a large bottle\n", index)
	}
	start_request[kind] <- r // Send request to small or large channel
	<-r.ack                  // Wait for server acknowledgment

	sleepRandomTime(3)       // Simulate bottle filling time
	end_request <- r         // Notify server filling is done
	<-r.ack                  // Wait for final acknowledgment
	fmt.Printf("[client %d] finished filling my bottle, exiting!\n", index)
	done <- true // Signal completion to main
}

// Operator goroutine: Manages refilling the tank and coin boxes
func operator() {
	var response int
	sleepRandomTime(4) // Simulate initial delay
	for {
		start_refill <- 1           // Request to start refill
		response = <-ack_operator   // Wait for acknowledgment
		if response == -1 {
			fmt.Printf("[operator] exiting...\n")
			done <- true // Signal termination
			return
		}
		fmt.Printf("[operator] starting the refill process...\n")
		sleepRandomTime(3)          // Simulate refill time
		end_refill <- 1             // Notify refill completion
		<-ack_operator              // Wait for acknowledgment
		fmt.Printf("[operator] Refill complete, water station is operational again...\n")
		sleepRandomTime(5)          // Simulate downtime after refill
	}
}

// Server goroutine: Manages the water station's state and coordination
func waterStation() {
	var currentWater = TankCapacity // Track remaining water
	var smallCoinCount = 0          // 10-cent coins collected
	var largeCoinCount = 0          // 20-cent coins collected
	var busy = false                // Whether the station is busy
	var stop = false                // Termination flag

	fmt.Printf("[waterStation] Water station is operational!\n")
	for {
		select {
		// Handle SmallBottle request if:
		// - Not busy, enough water, small coin box not full, and 
		//   (large coin box isn't full OR no pending refill)
		case x := <-whenRequest(!busy && currentWater >= CapacitySmall && smallCoinCount < MaxSmallCoins &&
			(largeCoinCount < MaxLargeCoins || (largeCoinCount == MaxLargeCoins && len(start_refill) == 0)), start_request[SmallBottle]):
			busy = true
			smallCoinCount++                // Add 10-cent coin
			currentWater -= CapacitySmall   // Deduct water
			fmt.Printf("[waterStation] Client %d started filling a bottle of type %d\n", x.index, x.kind)
			x.ack <- 1                      // Acknowledge client

		// Handle LargeBottle request if:
		// - Not busy, enough water, large coin box not full, no pending small requests,
		//   and (small coin box isn't full OR no pending refill)
		case x := <-whenRequest(!busy && currentWater >= CapacityLarge && largeCoinCount < MaxLargeCoins && 
			len(start_request[SmallBottle]) == 0 &&
			(smallCoinCount < MaxSmallCoins || (smallCoinCount == MaxSmallCoins && len(start_refill) == 0)), start_request[LargeBottle]):
			busy = true
			largeCoinCount++                // Add 20-cent coin
			currentWater -= CapacityLarge   // Deduct water
			fmt.Printf("[waterStation] Client %d started filling a bottle of type %d\n", x.index, x.kind)
			x.ack <- 1                      // Acknowledge client

		// Handle refill request from operator if:
		// - Not stopped, not busy, and (coin boxes full/water empty OR no pending requests)
		case <-when(!stop && !busy &&
			((smallCoinCount == MaxSmallCoins || largeCoinCount == MaxLargeCoins || currentWater == 0) ||
				(len(start_request[SmallBottle])+len(start_request[LargeBottle]) == 0)), start_refill):
			busy = true
			currentWater = TankCapacity // Refill water
			smallCoinCount = 0          // Reset coin counters
			largeCoinCount = 0
			fmt.Printf("[waterStation] Operator started refilling the tank and emptying coin boxes\n")
			ack_operator <- 1           // Acknowledge operator

		// Handle end of client request (bottle filled)
		case x := <-end_request:
			busy = false               // Free the station
			x.ack <- 1                 // Final acknowledgment

		// Handle end of refill process
		case <-end_refill:
			busy = false              // Free the station
			ack_operator <- 1         // Acknowledge operator

		// Handle operator termination signal
		case <-terminateOperator:
			stop = true              // Stop further refills
			fmt.Printf("[waterStation] All clients served, notifying operator to terminate\n")

		// Handle termination of refill process
		case <-when(stop, start_refill):
			ack_operator <- -1       // Signal operator to exit

		// Handle general termination
		case <-terminate:
			fmt.Printf("[waterStation] Shutting down!\n")
			done <- true             // Signal main to exit
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix()) // Seed random generator

	// Initialize channels for small and large requests
	for i := 0; i < 2; i++ {
		start_request[i] = make(chan request, MAX_BUFFER)
	}

	// Start all client goroutines
	for i := 0; i < MAX_CLIENTS; i++ {
		go client(i)
	}

	// Start operator and waterStation goroutines
	go operator()
	go waterStation()

	fmt.Printf("\n[MAIN] Water station is open.\n")

	// Wait for all clients to finish
	for i := 0; i < MAX_CLIENTS; i++ {
		<-done
	}

	// Terminate operator and waterStation
	terminateOperator <- true
	<-done                // Wait for operator to exit
	terminate <- true     // Signal waterStation to exit
	<-done                // Wait for waterStation to exit
	fmt.Printf("\n[MAIN] Water station is closed.\n")
}

/*
KEY CONCEPTS AND LOGIC:

1. Concurrency Control:
   - Uses Go's channels for synchronization between clients, operator, and the waterStation.
   - Buffered channels (start_request, end_request) handle multiple simultaneous requests.
   - Unbuffered channels (done, terminate) ensure strict synchronization for termination.

2. Client Workflow:
   - Each client randomly selects a bottle type, sends a request, waits for ack, simulates filling, then exits.
   - start_request[kind] sends the request to the appropriate channel (Small/Large).
   - end_request signals completion, freeing the station for the next client.

3. Operator Workflow:
   - Periodically refills the water tank and resets coin counters when conditions are met.
   - Uses start_refill and end_refill to coordinate with the waterStation.

4. WaterStation Logic:
   - Manages state: current water level, coin counts, busy status.
   - Uses select with conditional channels (when/whenRequest) to prioritize actions:
     - Small bottles have lower priority if large coins are full and a refill is pending.
     - Refill is triggered when coin boxes are full, water is empty, or no pending requests.
   - Ensures only one operation (client serving/refill) happens at a time via the busy flag.

5. Termination Sequence:
   - After all clients finish (all <-done received), main signals terminateOperator.
   - Operator exits, then main signals terminate to waterStation.
   - Channels ensure all goroutines exit gracefully.

6. Condition Handling:
   - The when/whenRequest helper functions dynamically enable/disable channels in select statements based on current state.
   - This allows the waterStation to handle requests only when safe (e.g., enough water, coin boxes not full).

TIPS FOR EXAM:
- Focus on how channels coordinate goroutines (client requests, operator actions).
- Understand the select statement in waterStation: conditions determine which case runs.
- Note the use of ack channels for step-by-step synchronization (e.g., client waits after sending request).
- Trace the flow for a single client and the operator to see interactions.
*/
