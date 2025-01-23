// -----------------------------------------------------------------------------------
// ENGLISH TRANSLATION AND EXPLANATION (IN-CODE)
//
// This Go program is another variation of a "bike rental" system with the following logic:
//
// 1) There are two types of bikes: traditional (BT) and electric (EB). 
//    A FLEX request means the client will prefer EB if available, otherwise take BT.
// 2) We have separate buffered channels for each kind of request:
//      - richiestaBT:   client requests a traditional bike
//      - richiestaEB:   client requests an electric bike
//      - richiestaFLEX: a "flexible" request, i.e., prefer electric bike if available
// 3) The server goroutine handles these requests using "guarded" select cases with
//    the helper function "when(cond, channel)" which returns nil if cond is false,
//    effectively disabling that case.
// 4) If no bikes are available for a particular request, the server either queues
//    the request or re-sends it to the appropriate channel (for example, a FLEX
//    request might get re-sent to richiestaEB if only EB is needed later).
// 5) Clients each request a bike, wait to receive it (through a dedicated channel
//    risorsa[clientID]), then hold onto the bike for some time before returning
//    it on rilascio. 
// 6) The system ends when all clients have finished, and main sends a termination
//    signal to the server.
//
// Below is the code with inline commentary in English.
//
// -----------------------------------------------------------------------------------

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// MAXPROC: maximum number of clients we can handle
const MAXPROC = 100

// N_EB, N_BT: number of electric (EB) and traditional (BT) bikes available
const N_EB = 3
const N_BT = 10

// Constants identifying a type of bike or type of request
const BT = 0   // traditional bike
const EB = 1   // electric bike
const FLEX = 2 // "flexible" request: prefer electric, otherwise accept traditional

// DIMBUF: size for buffered channels
const DIMBUF = 300

// bici is a custom type (int) used to represent either a traditional or electric bike
type bici int

// req is the request structure:
//  - id:   unique identifier of the client
//  - tipo: BT, EB, or FLEX
type req struct {
    id   int
    tipo int
}

// We have separate channels for each request type, plus one for releasing bikes.
//  - richiestaBT:   requests for a traditional bike
//  - richiestaEB:   requests for an electric bike
//  - richiestaFLEX: flexible requests (prefer EB, else BT)
//  - rilascio:      used by clients to return a bike
var richiestaBT = make(chan req, DIMBUF)
var richiestaEB = make(chan req, DIMBUF)
var richiestaFLEX = make(chan req, DIMBUF)
var rilascio = make(chan bici, DIMBUF)

// Each client has its own 'risorsa[clientID]' channel to receive the allocated bike
var risorsa [MAXPROC]chan bici

// done is used for waiting for client completion and also for the server to signal
var done = make(chan int)

// termina is sent from main to the server indicating it can shut down
var termina = make(chan int)

// when is a helper function used for "guarded" select statements:
// it returns 'c' if b == true, or nil if b == false, effectively blocking that case.
func when(b bool, c chan req) chan req {
    if !b {
        return nil
    }
    return c
}

// client simulates a user who requests a bike, receives it, uses it, then releases it.
func client(r req) {
    var b bici

    // Print the request according to the type
    if r.tipo == BT {
        fmt.Printf("[client %d] requesting a traditional bike (BT)...\n", r.id)
        // Send the request on richiestaBT
        richiestaBT <- r
    } else if r.tipo == EB {
        fmt.Printf("[client %d] requesting an electric bike (EB)...\n", r.id)
        richiestaEB <- r
    } else {
        fmt.Printf("[client %d] making a FLEX request...\n", r.id)
        richiestaFLEX <- r
    }

    // Wait for the server to send the allocated bike on risorsa[r.id]
    b = <-risorsa[r.id]

    // Announce which bike type was assigned
    if b == BT {
        fmt.Printf("[client %d] received a traditional bike (BT)\n", r.id)
    } else {
        fmt.Printf("[client %d] received an electric bike (EB)\n", r.id)
    }

    // Simulate using the bike for 2 seconds
    time.Sleep(time.Second * 2)

    // Release the bike
    rilascio <- b

    // Signal that this client has finished
    done <- r.id
}

// server is the goroutine that manages the available bikes, receiving requests
// and returning bikes. It uses a select statement with guarded channels via 'when'.
func server() {
    // dispEB, dispBT track how many EB or BT bikes are currently available
    var dispEB int = N_EB
    var dispBT int = N_BT

    // used for receiving a returned bike or an incoming request
    var b bici
    var r req

    for {
        // Slow down the loop a bit for demonstration
        time.Sleep(time.Second * 1)

        select {
        case b = <-rilascio:
            // A bike is being returned
            switch b {
            case EB:
                dispEB++
                fmt.Printf("[server] an electric bike was returned.\n")
            case BT:
                dispBT++
                fmt.Printf("[server] a traditional bike was returned.\n")
            }

        // A request for a traditional bike (BT)
        case r = <-when(dispBT > 0, richiestaBT):
            dispBT--
            b = BT
            fmt.Printf("[server] assigned a traditional bike to client %d\n", r.id)
            risorsa[r.id] <- b

        // A request for an electric bike (EB)
        case r = <-when(dispEB > 0, richiestaEB):
            dispEB--
            b = EB
            fmt.Printf("[server] assigned an electric bike to client %d\n", r.id)
            risorsa[r.id] <- b

        // A FLEX request: if there's an EB available, assign EB first
        case r = <-when(dispEB > 0, richiestaFLEX):
            dispEB--
            b = EB
            fmt.Printf("[server] assigned an electric bike to FLEX client %d\n", r.id)
            risorsa[r.id] <- b

        // Another FLEX case: if no EB is left but there's a BT, assign BT
        case r = <-when(dispEB == 0 && dispBT > 0, richiestaFLEX):
            dispBT--
            b = BT
            fmt.Printf("[server] assigned a traditional bike to FLEX client %d\n", r.id)
            risorsa[r.id] <- b

        // If both EB and BT are 0, we queue the FLEX request as an EB request,
        // effectively waiting for an electric bike. 
        case r = <-when(dispEB == 0 && dispBT == 0, richiestaFLEX):
            fmt.Printf("[server] FLEX client %d is queued for an electric bike...\n", r.id)
            // re-send the request on richiestaEB, so the client is effectively waiting
            richiestaEB <- r

        // Termination case: all clients done
        case <-termina:
            fmt.Println("END OF SERVER!")
            done <- 1
            return
        }
    }
}

func main() {
    var cli int
    var r req

    // Seed for random type assignments
    rand.Seed(time.Now().Unix())

    // Ask how many clients (max 100)
    fmt.Printf("\nHow many clients (max %d)? ", MAXPROC)
    fmt.Scanf("%d", &cli)
    fmt.Println("Number of clients:", cli)

    // Initialize the channels for each client to receive a bike
    for i := 0; i < MAXPROC; i++ {
        risorsa[i] = make(chan bici, DIMBUF)
    }

    // Create client goroutines
    // We randomly decide if each one is BT, EB, or FLEX
    for i := 0; i < cli; i++ {
        r.id = i
        r.tipo = rand.Intn(3) // 0=BT, 1=EB, 2=FLEX
        go client(r)
    }

    // Create the server goroutine
    go server()

    // Wait until all clients have finished (each sends 'done' once)
    for i := 0; i < cli; i++ {
        <-done
    }

    // Signal the server to terminate
    termina <- 1

    // Wait for the server to confirm termination
    <-done
}
