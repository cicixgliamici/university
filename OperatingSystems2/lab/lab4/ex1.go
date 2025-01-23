// -----------------------------------------------------------------------------------
// ENGLISH TRANSLATION WITH COMMENTARY (IN-CODE)
//
// This Go program simulates a single-lane bridge that can hold up to MAX vehicles at 
// once, traveling either from North (N) to South (S) or from South to North (N). 
// The bridge is "one-way" at a time: vehicles from one direction cannot enter if there 
// are vehicles currently crossing in the opposite direction. Moreover, there is a 
// priority rule favoring vehicles from the North (N).
//
// Key points:
//  - Up to MAX vehicles from one direction can be on the bridge simultaneously.
//  - If there are vehicles traveling from North, vehicles from the South can only 
//    enter if there are zero vehicles on the bridge from the North AND no waiting 
//    queue from the North (len(entrataN) == 0).
//  - Each vehicle goroutine (veicolo) picks a direction (N or S), waits a random 
//    initialization time, and then attempts to enter. Once it crosses, it signals 
//    exit. The server goroutine manages the capacity and direction constraints.
//  - The direction from the North is prioritized in the select statement: 
//    vehicles from the South can only enter when contN == 0 and len(entrataN) == 0.
//
// Channels:
//  - entrataN, entrataS: vehicles send their IDs here to request entry from North or 
//    South, respectively.
//  - uscitaN, uscitaS: vehicles send their IDs upon exiting.
//  - ACK_N[i], ACK_S[i]: each vehicle i from N or S receives acknowledgment (1) when 
//    it is allowed to enter.
//  - done: signals each goroutine has finished.
//  - termina: signals the server to shut down after all vehicles are done.
//
// Operation flow for a vehicle going North (N):
//  1) Sleep a random time (initialization).
//  2) Send its ID on entrataN.
//  3) Wait on ACK_N[myID] to be sure it can cross.
//  4) Cross (simulate crossing by sleeping).
//  5) Send ID on uscitaN.
//  6) Send a done signal.
//
// The server tracks how many vehicles from North (contN) and South (contS) are 
// currently on the bridge. The constraints allow a new North vehicle if contN < MAX 
// and contS == 0. A South vehicle can enter if contS < MAX, contN == 0, and nobody 
// from the North is waiting (len(entrataN) == 0). 
// 
// The program ends after all vehicles have finished crossing, and main sends a 
// termina signal, at which point the server shuts down.
//
// -----------------------------------------------------------------------------------

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// Buffer size and concurrency limits
const MAXBUFF = 100
const MAXPROC = 100

// MAX is the capacity of the bridge (number of vehicles allowed on it simultaneously)
const MAX = 5

// Directions
const N int = 0 // North
const S int = 1 // South

// Channels for synchronization
var done = make(chan bool)     // signals a goroutine (vehicle/server) has finished
var termina = make(chan bool)  // signals the server to terminate

// Channels for entering from North or South
var entrataN = make(chan int, MAXBUFF)
var entrataS = make(chan int, MAXBUFF)

// Channels for exiting from North or South
var uscitaN = make(chan int)
var uscitaS = make(chan int)

// Each vehicle from North or South has an acknowledgment channel to confirm 
// permission to enter the bridge.
var ACK_N [MAXPROC]chan int
var ACK_S [MAXPROC]chan int

// Helper function for "guarded" channels: returns c if b is true, or nil otherwise.
// This effectively enables/disables a select case based on the condition b.
func when(b bool, c chan int) chan int {
    if !b {
        return nil
    }
    return c
}

// GOROUTINE: veicolo
// A vehicle i traveling in direction dir. 
// 1) Sleeps a random init time.
// 2) Requests entry by sending ID on entrataN or entrataS.
// 3) Waits for ACK_N[i] or ACK_S[i] to confirm it can enter.
// 4) Sleeps to simulate crossing.
// 5) Signals exit via uscitaN or uscitaS.
// 6) Sends a done signal.
func veicolo(myid int, dir int) {
    // Random initialization delay
    var tt int
    tt = rand.Intn(5) + 1
    fmt.Printf("Initializing vehicle %d direction %d in %d seconds\n", myid, dir, tt)
    time.Sleep(time.Duration(tt) * time.Second)

    if dir == N {
        // Request to enter from North
        entrataN <- myid 
        // Wait for server acknowledgment
        <-ACK_N[myid]   
        fmt.Printf("[vehicle %d] entered the bridge heading NORTH\n", myid)

        // Cross the bridge (random time)
        tt = rand.Intn(5)
        time.Sleep(time.Duration(tt) * time.Second)

        // Signal exit
        uscitaN <- myid
        fmt.Printf("[vehicle %d] left the bridge heading NORTH\n", myid)

    } else {
        // Request to enter from South
        entrataS <- myid
        // Wait for server acknowledgment
        <-ACK_S[myid]
        fmt.Printf("[vehicle %d] entered the bridge heading SOUTH\n", myid)

        // Cross the bridge (random time)
        tt = rand.Intn(5)
        time.Sleep(time.Duration(tt) * time.Second)

        // Signal exit
        uscitaS <- myid
        fmt.Printf("[vehicle %d] left the bridge heading SOUTH\n", myid)
    }

    // Signal completion of this vehicle
    done <- true
}

// GOROUTINE: server
// Manages the number of vehicles on the bridge: contN for North, contS for South.
// The constraints are:
//   - A North vehicle can enter if contN < MAX and contS == 0.
//   - A South vehicle can enter if contS < MAX, contN == 0, and no one from North is waiting.
//
// Priority is given to vehicles from the North. The code in select ensures that 
// if there are vehicles from the North waiting, the South can't enter 
// (unless contN == 0 and len(entrataN) == 0).
func server() {
    var contN int = 0 // how many North vehicles are currently on the bridge
    var contS int = 0 // how many South vehicles are currently on the bridge

    for {
        select {
        // 1) A North vehicle tries to enter if contN < MAX and contS == 0
        case x := <-when((contN < MAX) && (contS == 0), entrataN):
            contN++
            ACK_N[x] <- 1 // allow the vehicle to enter

        // 2) A South vehicle tries to enter if contS < MAX, contN == 0, and no North waiting
        case x := <-when((contS < MAX) && (contN == 0) && (len(entrataN) == 0), entrataS):
            contS++
            ACK_S[x] <- 1 // allow the vehicle to enter

        // 3) A North vehicle leaves
        case <-uscitaN:
            contN--

        // 4) A South vehicle leaves
        case <-uscitaS:
            contS--

        // 5) Termination signal from main
        case <-termina:
            fmt.Println("END!!!")
            done <- true
            return
        }
    }
}

func main() {
    var VN int // number of vehicles from North
    var VS int // number of vehicles from South

    fmt.Printf("\nHow many NORTH vehicles (max %d)? ", MAXPROC)
    fmt.Scanf("%d", &VN)
    fmt.Printf("\nHow many SOUTH vehicles (max %d)? ", MAXPROC)
    fmt.Scanf("%d", &VS)

    // Initialize the acknowledgment channels for each North vehicle
    for i := 0; i < VN; i++ {
        ACK_N[i] = make(chan int, MAXBUFF)
    }

    // Initialize the acknowledgment channels for each South vehicle
    for i := 0; i < VS; i++ {
        ACK_S[i] = make(chan int, MAXBUFF)
    }

    // Seed random generator
    rand.Seed(time.Now().Unix())

    // Start the server goroutine
    go server()

    // Create SOUTH vehicle goroutines
    for i := 0; i < VS; i++ {
        go veicolo(i, S)
    }

    // Create NORTH vehicle goroutines
    for i := 0; i < VN; i++ {
        go veicolo(i, N)
    }

    // Wait until all vehicles are done
    for i := 0; i < VN+VS; i++ {
        <-done
    }

    // Signal the server to terminate
    termina <- true

    // Wait for the server's done
    <-done
    fmt.Printf("\nALL FINISHED\n")
}
