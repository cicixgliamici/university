package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ============================================================
//                    CONSTANTS / PARAMETERS
// ============================================================
const (
	MAXBUFFER   = 100  // Maximum buffer size for channels
	MAX_CLIENTS = 20   // Maximum number of Clients/Workers that can be managed

	TYPE_A   = 0       // First type of resource
	TYPE_B   = 1       // Second type of resource
	TYPE_MIX = 2       // "Mixed" type

	MAX_A = 4000       // Max capacity for resource type A
	MAX_B = 3000       // Max capacity for resource type B

	LOT_A   = 700      // Lot for resource type A
	LOT_B   = 300      // Lot for resource type B
	LOT_MIX = 500      // Lot for the "mixed" resource
)

// ============================================================
//                      DATA STRUCTURE
// ============================================================
// Request represents the request that a Client/Worker or a Supplier
// can make to the system. `ack` is the acknowledgment channel.
type Request struct {
	id   int      // Identifier of who is making the request
	tipo int      // Indicates the type of resource involved
	ack  chan int // Acknowledgment channel to signal completion of events,
	              // possibly can be int if returning information
}

// ============================================================
//                         CHANNELS
// ============================================================
// requestChan: channels used by Clients/Workers to request resources.
//   - requestChan[TYPE_A]
//   - requestChan[TYPE_B]
//   - requestChan[TYPE_MIX]
// It's an array of channels: each one manages the requests for a specific type.
//
// Usage example:
//   r := Request{id: 123, tipo: TYPE_A, ack: make(chan int)}
//   requestChan[TYPE_A] <- r // Sends the request on the channel for TYPE_A.
//
var requestChan [3]chan Request

// restockChan: channels used by Suppliers to restock the warehouse.
//   - restockChan[TYPE_A]
//   - restockChan[TYPE_B]
var restockChan [2]chan Request

// endRequest and endRestock: channels to signal the conclusion
// of a retrieval/restock operation.
var endRequest = make(chan Request, MAXBUFFER)
var endRestock = make(chan Request)

// Channels for process termination.
var done = make(chan bool)
var stopWarehouse = make(chan bool)
var stopSupplier = make(chan bool)

// ============================================================
//                     SUPPORT FUNCTIONS
// ============================================================

// Implements a logical guard: returns channel c if condition b is true, otherwise nil.
// Useful for conditional switch-select.
func when(b bool, c chan Request) chan Request {
	if !b {
		return nil
	}
	return c
}

// Waits a random amount of time (in seconds) in the range [0, max).
func sleepRandTime(max int) {
	if max > 0 {
		time.Sleep(time.Duration(rand.Intn(max)+1) * time.Second)
	}
}

// Waits a random amount of time (in seconds) in the range [min, max).
func sleepRandTimeRange(min, max int) {
	if min >= 0 && max > 0 && min < max {
		time.Sleep(time.Duration(rand.Intn(max-min)+min) * time.Second)
	}
}

// Returns a string based on the resource type (case constants are defined above).
func getResourceName(t int) string {
	switch t {
	case TYPE_A:
		return "type A"
	case TYPE_B:
		return "type B"
	case TYPE_MIX:
		return "MIXED type"
	default:
		return "unknown"
	}
}

// Debug print for the state of active retrievals, queues, etc.
func debug(activeA, activeB, lenA, lenB, lenMix int) {
	fmt.Printf("State: ongoing retrievals (A: %d, B: %d) | Queues(A: %d, B: %d, Mix: %d)\n",
		activeA, activeB, lenA, lenB, lenMix)
}

// ============================================================
//                         GOROUTINES
// ============================================================

// client simulates a generic “Worker” or “Operator” that cyclically
// requests and retrieves resources from the warehouse.
func client(id int) {
	tipo := -1
	r := Request{id: id, tipo: tipo, ack: make(chan int)}

	fmt.Printf("[CLIENT %d] Started\n", id)
	for i := 0; i < 5; i++ {
		// Random choice of resource type (TYPE_A, TYPE_B, or TYPE_MIX).
		tipoRand := rand.Intn(100)
		if tipoRand >= 80 {
			r.tipo = TYPE_MIX
		} else {
			r.tipo = tipoRand % 2 // 0 or 1
		}

		fmt.Printf("[CLIENT %d] Requesting resource %s\n", id, strings.ToUpper(getResourceName(r.tipo)))
		requestChan[r.tipo] <- r // send request
		<-r.ack                  // wait for start-ack

		fmt.Printf("[CLIENT %d] Retrieving resource %s...\n", id, strings.ToUpper(getResourceName(r.tipo)))
		sleepRandTime(3) // simulate retrieval

		endRequest <- r   // signal completion of retrieval
		<-r.ack           // wait for the warehouse to finish the operation
	}

	// After finishing all iterations, send a signal on the `done` channel
	// to notify that this client has finished its work.
	done <- true
	fmt.Printf("[CLIENT %d] Terminating\n", id)
}

// supplier simulates a generic supplier that cyclically restocks
// the warehouse with a certain type of resource.
func supplier(resourceType int) {
	// ID not necessary here as it’s not used in the example
	r := Request{tipo: resourceType, ack: make(chan int)}

	fmt.Printf("[SUPPLIER %s] Started\n", strings.ToUpper(getResourceName(resourceType)))
	for {
		sleepRandTimeRange(5, 10)

		fmt.Printf("[SUPPLIER %s] I want to restock the warehouse\n", strings.ToUpper(getResourceName(resourceType)))
		restockChan[resourceType] <- r // send restock request
		<-r.ack                        // wait for start-ack

		fmt.Printf("[SUPPLIER %s] Restocking in progress...\n", strings.ToUpper(getResourceName(resourceType)))
		sleepRandTimeRange(3, 5) // simulate restocking

		endRestock <- r    // signal completion
		<-r.ack            // wait for the warehouse to complete the operation
		fmt.Printf("[SUPPLIER %s] Restocking completed\n", strings.ToUpper(getResourceName(resourceType)))

		// This select checks whether a stop signal was sent to terminate the supplier.
		// If we don't receive anything from stopSupplier, continue the loop.
		// Otherwise, exit the loop and terminate.
		select {
		case <-stopSupplier:
			fmt.Printf("[SUPPLIER %s] Terminating\n", strings.ToUpper(getResourceName(resourceType)))
			done <- true
			return
		default:
			continue
		}
	}
}

// warehouse manages the logic of accessing and updating
// a warehouse with two types of resources (TYPE_A and TYPE_B)
// plus a “mixed” type (TYPE_MIX).
// Insert any precedence/priority conditions required by your exam or scenario here.
func warehouse() {
	// How many resources are available initially (you can change the initialization logic)
	resources := [2]int{MAX_A, MAX_B}

	// To track how many retrievals and restocks are in progress
	//
	// activePrel tracks how many clients are currently retrieving each resource type (A and B).
	// activeRestock indicates whether a restock is in progress for each resource type.
	activePrel := [2]int{0, 0}
	activeRestock := [2]bool{false, false}

	fmt.Printf("[WAREHOUSE] Started. Initial state: A: %d/%d, B: %d/%d\n",
		resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)

	for {
		select {
		//---------------------------------------------------
		//             RETRIEVAL (START)
		//---------------------------------------------------
		case req := <-when(
			// Conditions to allow retrieval of TYPE_A
			((LOT_A * (activePrel[TYPE_A] + 1)) <= resources[TYPE_A]) &&
				(!activeRestock[TYPE_A]) &&
				(len(requestChan[TYPE_MIX]) == 0), // e.g., give priority to TYPE_MIX
			requestChan[TYPE_A]):

			activePrel[TYPE_A]++
			fmt.Printf("[WAREHOUSE] Client %d begins retrieval of %d (type A)\n",
				req.id, LOT_A)
			req.ack <- 1 // unblock the client

		case req := <-when(
			// Conditions to allow retrieval of TYPE_B
			((LOT_B * (activePrel[TYPE_B] + 1)) <= resources[TYPE_B]) &&
				(!activeRestock[TYPE_B]) &&
				(len(requestChan[TYPE_MIX]) == 0 && len(requestChan[TYPE_A]) == 0),
			requestChan[TYPE_B]):

			activePrel[TYPE_B]++
			fmt.Printf("[WAREHOUSE] Client %d begins retrieval of %d (type B)\n",
				req.id, LOT_B)
			req.ack <- 1 // unblock the client

		case req := <-when(
			// Conditions to allow retrieval of TYPE_MIX (both A and B)
			((LOT_MIX * (activePrel[TYPE_A] + 1)) <= resources[TYPE_A] &&
				(LOT_MIX * (activePrel[TYPE_B] + 1)) <= resources[TYPE_B]) &&
				(!activeRestock[TYPE_A] && !activeRestock[TYPE_B]),
			requestChan[TYPE_MIX]):

			activePrel[TYPE_A]++
			activePrel[TYPE_B]++
			fmt.Printf("[WAREHOUSE] Client %d begins MIXED retrieval of %d (A) and %d (B)\n",
				req.id, LOT_MIX, LOT_MIX)
			req.ack <- 1 // unblock the client

		//---------------------------------------------------
		//             RETRIEVAL (END)
		//---------------------------------------------------
		case req := <-endRequest:
			switch req.tipo {
			case TYPE_A:
				resources[TYPE_A] -= LOT_A
				activePrel[TYPE_A]--
			case TYPE_B:
				resources[TYPE_B] -= LOT_B
				activePrel[TYPE_B]--
			case TYPE_MIX:
				resources[TYPE_A] -= LOT_MIX
				resources[TYPE_B] -= LOT_MIX
				activePrel[TYPE_A]--
				activePrel[TYPE_B]--
			default:
				fmt.Println("[WAREHOUSE] ERROR: invalid resource type.")
			}
			fmt.Printf("[WAREHOUSE] Client %d has finished. State: A: %d/%d, B: %d/%d\n",
				req.id, resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)
			// Unblock the client if needed:
			req.ack <- 1

		//---------------------------------------------------
		//           RESTOCK (START)
		//---------------------------------------------------
		case req := <-when(
			// Condition to restock TYPE_A
			(activePrel[TYPE_A] == 0) &&
				(resources[TYPE_A] <= resources[TYPE_B] || len(restockChan[TYPE_B]) == 0),
			restockChan[TYPE_A]):
			activeRestock[TYPE_A] = true
			fmt.Printf("[WAREHOUSE] Starting restock of A...\n")
			req.ack <- 1

		case req := <-when(
			// Condition to restock TYPE_B
			(activePrel[TYPE_B] == 0) &&
				(resources[TYPE_B] < resources[TYPE_A] || len(restockChan[TYPE_A]) == 0),
			restockChan[TYPE_B]):
			activeRestock[TYPE_B] = true
			fmt.Printf("[WAREHOUSE] Starting restock of B...\n")
			req.ack <- 1

		//---------------------------------------------------
		//           RESTOCK (END)
		//---------------------------------------------------
		case req := <-endRestock:
			switch req.tipo {
			case TYPE_A:
				resources[TYPE_A] = MAX_A
				activeRestock[TYPE_A] = false
				fmt.Printf("[WAREHOUSE] Finished restocking A. A: %d/%d, B: %d/%d\n",
					resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)
				req.ack <- 1
			case TYPE_B:
				resources[TYPE_B] = MAX_B
				activeRestock[TYPE_B] = false
				fmt.Printf("[WAREHOUSE] Finished restocking B. A: %d/%d, B: %d/%d\n",
					resources[TYPE_A], MAX_A, resources[TYPE_B], MAX_B)
				req.ack <- 1
			default:
				fmt.Println("[WAREHOUSE] ERROR: invalid resource type.")
				req.ack <- -1
			}

		//---------------------------------------------------
		//             TERMINATION
		//---------------------------------------------------
		case <-stopWarehouse:
			fmt.Printf("[WAREHOUSE] Terminating\n")
			done <- true
			return
		}
	}
}

// ============================================================
//                          MAIN
// ============================================================
func main() {
	fmt.Println("[MAIN] Start")
	rand.Seed(time.Now().UnixNano())

	// Add user input option for selecting the simulation mode
	var mode int
	fmt.Print("Enter the simulation mode (1=basic, 2=advanced): ")
	fmt.Scanf("%d\n", &mode)
	fmt.Printf("You chose mode %d. Customize behavior based on this value...\n", mode)

	// Example: specify how many clients and suppliers to start
	nClients := 5
	nSuppliers := 2

	fmt.Printf("[MAIN] How many Clients do you want to start? (max %d): ", MAX_CLIENTS)
	fmt.Scanf("%d\n", &nClients)
	if nClients < 2 {
		fmt.Printf("[MAIN] Too few clients. Using default value: 4.\n")
		nClients = 4
	}

	// Initialize main channels
	for i := 0; i < len(requestChan); i++ {
		requestChan[i] = make(chan Request, MAXBUFFER)
	}
	for i := 0; i < len(restockChan); i++ {
		restockChan[i] = make(chan Request, MAXBUFFER)
	}

	// Start goroutines
	go warehouse() // resource manager

	// Start suppliers (one for TYPE_A and one for TYPE_B, if nSuppliers=2)
	for i := 0; i < nSuppliers; i++ {
		go supplier(i) // i corresponds to TYPE_A or TYPE_B
	}

	// Start clients
	for i := 0; i < nClients; i++ {
		go client(i)
	}

	// The `done` channel is used as a completion signal by clients and suppliers.
	// When a client finishes, it sends `true` on the `done` channel.
	// The main waits for these signals to ensure orderly termination.

	// Wait for all clients to finish
	for i := 0; i < nClients; i++ {
		<-done
	}

	// Signal suppliers to terminate
	for i := 0; i < nSuppliers; i++ {
		stopSupplier <- true
	}
	// Wait for all suppliers to finish
	for i := 0; i < nSuppliers; i++ {
		<-done
	}

	// Signal warehouse to terminate
	stopWarehouse <- true
	<-done

	fmt.Println("[MAIN] End")
}

// ============================================================
//                  GO SEMANTICS (BASICS)
// ============================================================
// Below is a brief review of some key concepts:
//
// 1. Variable Declaration:
//    - Common forms:
//      var x int           // declare x of type int
//      x := 10             // declaration and assignment, type inferred
//      var y, z = 10, 20   // multiple declaration
//      x := rand.Intn(10)  // rand.Intn(n) generates a random integer in [0, n).
//                          // if you want range [0, n], you could do rand.Intn(n+1)
//
// 2. Channels:
//    - Declaration and creation:
//      var c chan int      // declare a channel of int, not yet created
//      c = make(chan int)  // creation of the channel
//    - Sending a value (send):
//      c <- value
//    - Receiving a value (receive):
//      <-c          // if the content isn't needed (e.g., for ack)
//      v := <-c     // if you need the content
//      or: v, ok := <-c   // where ok indicates if the channel is still open
//
// 3. Loops (for):
//    - Go doesn't have `while`; use `for`:
//      for i := 0; i < 10; i++ {
//          ...
//      }
//    - Conditional `for`:
//      for condition {
//          ...
//      }
//    - Infinite `for`:
//      for {
//          ...
//      }
//
