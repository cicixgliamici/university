// -----------------------------------------------------------------------------------
// ENGLISH TRANSLATION AND IN-CODE EXPLANATION
//
// This Go program simulates a one-lane bridge that can be used by both pedestrians
// (PED) and cars (AUT), coming from either direction (N = North, S = South). 
// The bridge has a capacity limit (MAX) of 35 "units". A pedestrian counts as 1 unit, 
// and a car counts as 10 units.
//
// The logic is enforced by a single "server" goroutine that manages the current usage
// of the bridge. Pedestrians and cars send requests to enter from either North or 
// South and then exit accordingly. The server checks constraints before allowing 
// them to enter. 
//
// Key points in this simulation:
// - Pedestrians from each side share separate channels: 
//       entrataN_P (pedestrian from North) and entrataS_P (pedestrian from South).
// - Cars from each side also have separate channels: 
//       entrataN_A (car from North) and entrataS_A (car from South).
// - Each client (user) is randomly assigned to be a pedestrian or a car, and also 
//   assigned a direction (N or S). 
// - The server uses a select statement with conditions (guarded channels via `when()`).
//   For instance, a pedestrian can enter from the South if `tot < MAX` and there are
//   no cars from the North occupying the bridge (`contN[AUT] == 0`), etc.
// - Exiting is done by sending an object `msg_out` with a `tipo` (PED or AUT) and 
//   an `id`, to either `uscitaN` or `uscitaS`. The server updates counters based 
//   on whether it was a pedestrian or a car (subtract 1 or 10 units).
// - The program ends when all "users" are done, and `main` sends a termination signal
//   to the server.
//
// Below is the code with inline commentary.
//
// -----------------------------------------------------------------------------------

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// BUFFER AND CONCURRENCY CONSTANTS
const MAXBUFF = 100
const MAXPROC = 100

// The bridge capacity (in "units"). 
// A pedestrian uses 1 unit, a car uses 10 units.
const MAX = 35

// Directions
const N int = 0 // North
const S int = 1 // South

// Types of users
const PED int = 0 // pedestrian
const AUT int = 1 // car

// msg_out is used by a user to indicate exit, carrying two fields:
//   - tipo: PED or AUT
//   - id:   the user's ID
type msg_out struct{ tipo, id int }

// done channel signals a goroutine is finished (client or server)
var done = make(chan bool)

// termina channel signals the server to shut down once everyone is finished
var termina = make(chan bool)

// Channels for entry from N or S, separated by type (pedestrian or car)
var entrataN_A = make(chan int, MAXBUFF) // North, car
var entrataN_P = make(chan int, MAXBUFF) // North, pedestrian
var entrataS_A = make(chan int, MAXBUFF) // South, car
var entrataS_P = make(chan int, MAXBUFF) // South, pedestrian

// Channels for exiting from North or South (send a msg_out object)
var uscitaN = make(chan msg_out)
var uscitaS = make(chan msg_out)

// Each user has an ACK channel to synchronize entry. ACK[i] is for user i.
var ACK [MAXPROC]chan int

// Helper function for conditional select statement: if b is false, return nil
// so that case is effectively disabled.
func when(b bool, c chan int) chan int {
    if !b {
        return nil
    }
    return c
}

// GOROUTINE: utente (user)
// Each user has an ID, a direction (N or S), and randomly becomes either a
// pedestrian (PED) or a car (AUT). Once assigned a direction and type, it:
//  1) Waits a random initialization time.
//  2) Sends a request to the appropriate channel (e.g., entrataN_P for
//     a pedestrian from North).
//  3) Waits for ACK[myid] to get permission from the server.
//  4) Sleeps a random time to simulate being on the bridge.
//  5) Sends a msg_out to the appropriate exit channel (uscitaN or uscitaS).
//  6) Signals done.
func utente(myid int, dir int) {
    var m_out msg_out
    var tt int

    // Randomly decide if this user is a pedestrian or a car
    tipo := rand.Intn(2) // 0 => PED, 1 => AUT

    // Prepare a random initialization sleep
    tt = rand.Intn(5) + 1
    fmt.Printf("Initialization user %d direction %d type %d, waiting %d seconds\n", myid, dir, tipo, tt)
    time.Sleep(time.Duration(tt) * time.Second)

    // Prepare the exit message
    m_out.tipo = tipo
    m_out.id = myid

    // Depending on direction (N or S) and type (PED or AUT), we choose the channel
    if dir == N {
        if tipo == PED {
            // Pedestrian from North
            entrataN_P <- myid
            <-ACK[myid] 
            fmt.Printf("[pedestrian %d] entered from NORTH\n", myid)
            tt = rand.Intn(5)
            time.Sleep(time.Duration(tt) * time.Second)
            uscitaN <- m_out
        } else {
            // Car from North
            entrataN_A <- myid
            <-ACK[myid]
            fmt.Printf("[car %d] entered from NORTH\n", myid)
            tt = rand.Intn(5)
            time.Sleep(time.Duration(tt) * time.Second)
            uscitaN <- m_out
        }
    } else { 
        // direction S
        if tipo == PED {
            // Pedestrian from South
            entrataS_P <- myid
            <-ACK[myid]
            fmt.Printf("[pedestrian %d] entered from SOUTH\n", myid)
            tt = rand.Intn(5)
            time.Sleep(time.Duration(tt) * time.Second)
            uscitaS <- m_out
        } else {
            // Car from South
            entrataS_A <- myid
            <-ACK[myid]
            fmt.Printf("[car %d] entered from SOUTH\n", myid)
            tt = rand.Intn(5)
            time.Sleep(time.Duration(tt) * time.Second)
            uscitaS <- m_out
        }
    }

    // Signal user completion
    done <- true
}

// GOROUTINE: server
// Maintains counts of pedestrians and cars for each direction, plus a 'tot' that
// represents the total usage of the bridge in "units". Pedestrians = 1 unit, cars = 10 units.
// The select statement uses "when(..., channel)" to enforce capacity constraints and 
// mutual exclusion rules. 
//
// Conditions for entering:
//   - Pedestrian from South can enter if tot < MAX and no cars are in the North (contN[AUT] == 0).
//   - Pedestrian from North can enter if tot < MAX and no cars in South (contS[AUT] == 0) and
//     no waiting queue of pedestrians from South (len(entrataS_P) == 0).
//   - A car from South can enter if (tot+10 <= MAX), no pedestrians/cars from North (contN[AUT] + contN[PED] == 0), 
//     and no pedestrians from either side are queued (the example checks len(entrataN_P) + len(entrataS_P) == 0).
//   - Similar for a car from North. 
//
// Exiting is handled by reading from uscitaN or uscitaS, subtracting 1 from tot if
// it's a pedestrian, or 10 if it's a car. 
func server() {
    // contN and contS track pedestrians and cars in each direction:
    //   contN[PED], contN[AUT], contS[PED], contS[AUT].
    var contN [2]int
    var contS [2]int

    // tot tracks the total usage of the bridge (sum in "units" - each ped = 1, car = 10).
    var tot int

    for {
        select {

        // 1) Pedestrian from South can enter if capacity is not exceeded (tot < MAX) 
        //    and there are no cars from North (contN[AUT] == 0).
        case x := <-when((tot < MAX) && (contN[AUT] == 0), entrataS_P):
            contS[PED]++
            tot++
            ACK[x] <- 1

        // 2) Pedestrian from North can enter if tot < MAX, no cars from South, 
        //    and there are no waiting pedestrians from South (len(entrataS_P) == 0).
        case x := <-when((tot < MAX) && (contS[AUT] == 0) && (len(entrataS_P) == 0), entrataN_P):
            contN[PED]++
            tot++
            ACK[x] <- 1

        // 3) Car from South can enter if tot + 10 <= MAX, no one from North is on the bridge, 
        //    and no pedestrians are waiting to enter (entrataN_P + entrataS_P == 0).
        case x := <-when((tot+10 <= MAX) && (contN[PED]+contN[AUT] == 0) && 
                         (len(entrataN_P)+len(entrataS_P) == 0), entrataS_A):
            contS[AUT]++
            tot += 10
            ACK[x] <- 1

        // 4) Car from North can enter if tot + 10 <= MAX, no one from South is on the bridge, 
        //    and no pedestrians from either side are waiting.
        case x := <-when((tot+10 <= MAX) && (contS[PED]+contS[AUT] == 0) && 
                         (len(entrataN_P)+len(entrataS_P)+len(entrataS_A) == 0), entrataN_A):
            contN[AUT]++
            tot += 10
            ACK[x] <- 1

        // 5) Exiting from North
        case x := <-uscitaN:
            contN[x.tipo]--
            if x.tipo == PED {
                tot--
            } else {
                tot -= 10
            }

        // 6) Exiting from South
        case x := <-uscitaS:
            contS[x.tipo]--
            if x.tipo == PED {
                tot--
            } else {
                tot -= 10
            }

        // 7) Termination signal
        case <-termina:
            fmt.Println("BRIDGE IS CLOSING!")
            done <- true
            return
        }
    }
}

func main() {
    var VN int // number of users from North
    var VS int // number of users from South

    fmt.Printf("\nHow many users from the NORTH (max %d)? ", MAXPROC/2)
    fmt.Scanf("%d", &VN)
    fmt.Printf("\nHow many users from the SOUTH (max %d)? ", MAXPROC/2)
    fmt.Scanf("%d", &VS)

    // Initialize ACK channels
    for i := 0; i < VN+VS; i++ {
        ACK[i] = make(chan int, MAXBUFF)
    }

    // Seed random
    rand.Seed(time.Now().Unix())
    // Start the server goroutine
    go server()

    // Create user goroutines for the South side first
    for i := 0; i < VS; i++ {
        go utente(i, S)
    }
    // Then create user goroutines for the North side
    for j := VS; j < VN+VS; j++ {
        go utente(j, N)
    }

    // Wait for all users to finish
    for i := 0; i < VN+VS; i++ {
        <-done
    }

    // Tell the server to terminate
    termina <- true
    <-done
    fmt.Printf("\nALL DONE\n")
}
