// -----------------------------------------------------------------------------------
// ENGLISH TRANSLATION AND EXPLANATION
//
// This program simulates a bike rental system with the following characteristics:
//
// - There are two types of bikes: 
//     * Traditional bikes (BT)
//     * Electric bikes (EB)
// - A client can request a specific type of bike (BT or EB), or "FLEX" which means
//   "prefer electric (EB) if available, otherwise take a traditional (BT)."
// - There is a limited number of electric bikes (N_EB) and a limited number of 
//   traditional bikes (N_BT).
// - Each client goroutine asks for a bike, receives it from the server goroutine 
//   if available, uses it for some time (simulated with a sleep), then returns it. 
// - If no bikes of the requested type are available, the client is put on hold
//   (the server keeps the request in a 'waiting' state) until a bike is returned.
// - The program ends when all clients have finished.
//
// Concurrency Elements:
//  - `richiesta`: A channel on which clients send requests of type `req`. 
//                 Each request has an `id` (the client’s ID) and a `tipo` 
//                 (BT, EB, or FLEX).
//  - `rilascio`:  A channel on which clients return a `bici` value (the type of bike).
//  - `risorsa[i]`: Each client i receives the allocated bike on this channel. 
//  - `done`:      Used to synchronize completion of clients and to signal the server's end.
//  - `termina`:   A signal sent from `main` to the server indicating that all clients
//                 have completed and the server can terminate.
//
// The `server` goroutine tracks how many electric bikes (dispEB) and how many
// traditional bikes (dispBT) are available, plus arrays `sospEB` and `sospBT` to 
// keep track of waiting clients. If a bike is returned (via `rilascio`), the server 
// checks if there are any waiting requests for that type and satisfies one immediately, 
// otherwise it increases the count of available bikes.
//
// In `main`, we create a certain number of client goroutines (each running `client()`)
// and one `server` goroutine. Once all clients have signaled they are done (reading 
// from `done`), we send a termination message to the server (on `termina`), 
// then wait for the final `done` from the server.
//
// -----------------------------------------------------------------------------------

package main

import (
    "fmt"
    "time"
)

// MAXPROC is the maximum number of client processes we can spawn
const MAXPROC = 100

// N_EB is the number of electric bikes (E-Bikes) available
// N_BT is the number of traditional bikes available
const N_EB = 1
const N_BT = 30

// Constants to represent bike types or request types
const BT = 0   // Traditional bike
const EB = 1   // Electric bike
const FLEX = 2 // "Flexible" request: prefers EB if available, otherwise BT

// bici is a simple alias type representing either BT or EB
type bici int

// 'req' is the structure of a request sent by a client:
//   - id:   the client's ID
//   - tipo: one of BT, EB, or FLEX
type req struct {
    id   int
    tipo int
}

// CHANNELS:

// 'richiesta': used by clients to request a bike (of type req).
var richiesta = make(chan req)

// 'rilascio': used by clients to return their allocated bike (of type bici).
var rilascio = make(chan bici)

// 'risorsa[i]': used by the server to send an allocated bike to client i.
var risorsa [MAXPROC]chan bici

// 'done': used to synchronize the completion of both clients and the server.
var done = make(chan int)

// 'termina': used by main to tell the server it can shut down 
//            once all clients have finished.
var termina = make(chan int)

// We keep track of whether each bike (EB or BT) is free or not using boolean arrays.
// (In this particular solution, we also have counters dispEB and dispBT in the server,
// so these arrays are somewhat redundant, but included as per the original code.)
var liberaEB [N_EB]bool
var liberaBT [N_BT]bool

// GOROUTINE: client
// 1) Prints its request (BT, EB, or FLEX).
// 2) Sends a request on 'richiesta'.
// 3) Waits to receive a 'bici' on risorsa[r.id] from the server.
// 4) Sleeps 2 seconds to simulate usage.
// 5) Returns the bike via 'rilascio'.
// 6) Signals 'done' to indicate the client is finished.
func client(r req) {
    // Announce request type
    if r.tipo == BT {
        fmt.Printf("[client %d] requesting a traditional bike (BT)...\n", r.id)
    } else if r.tipo == EB {
        fmt.Printf("[client %d] requesting an electric bike (EB)...\n", r.id)
    } else {
        fmt.Printf("[client %d] requesting FLEX (EB if possible, else BT)...\n", r.id)
    }

    // Send the request to the server
    richiesta <- r

    // Wait for the assigned bike from the server
    b := <-risorsa[r.id]
    if b == BT {
        fmt.Printf("[client %d] received a traditional bike\n", r.id)
    } else {
        fmt.Printf("[client %d] received an electric bike\n", r.id)
    }

    // Use the bike for 2 seconds
    time.Sleep(2 * time.Second)

    // Return the bike
    rilascio <- b

    // Signal completion
    done <- r.id
}

// GOROUTINE: server
// Manages the allocation and deallocation of bikes. Maintains counts of available
// electric and traditional bikes (dispEB, dispBT), and keeps track of waiting 
// requests in arrays 'sospEB' and 'sospBT'. If a bike is returned, it is allocated 
// to the first waiting request (if any). Otherwise, the count of available bikes 
// is incremented.
func server() {
    var dispEB int = N_EB
    var dispBT int = N_BT
    var b bici
    var r req

    // 'sospEB[i]': indicates if client i is waiting for an electric bike
    // 'sospBT[i]': indicates if client i is waiting for a traditional bike
    var sospEB [MAXPROC]bool
    var sospBT [MAXPROC]bool

    var nsospEB int = 0 // number of clients waiting for EB
    var nsospBT int = 0 // number of clients waiting for BT

    // Initialize the arrays for availability
    for i := 0; i < N_EB; i++ {
        liberaEB[i] = true
    }
    for i := 0; i < N_BT; i++ {
        liberaBT[i] = true
    }

    // Initialize waiting arrays
    for i := 0; i < MAXPROC; i++ {
        sospBT[i] = false
        sospEB[i] = false
    }

    for {
        // Sleep here just to slow down the server loop for demonstration
        time.Sleep(1 * time.Second)

        select {
        case b = <-rilascio:
            // A bike is being returned
            if b == EB {
                // Electric bike returned
                if nsospEB == 0 {
                    // No one is waiting for EB, increment available EB
                    dispEB++
                    fmt.Printf("[server] an electric bike was returned.\n")
                } else {
                    // Someone is waiting for an EB, assign it immediately
                    for i := 0; i < MAXPROC; i++ {
                        if sospEB[i] == true {
                            risorsa[i] <- b
                            nsospEB--
                            sospEB[i] = false
                            break
                        }
                    }
                }
            } else {
                // Traditional bike returned
                if nsospBT == 0 {
                    dispBT++
                    fmt.Printf("[server] a traditional bike was returned.\n")
                } else {
                    // Someone is waiting for a BT, assign it immediately
                    for i := 0; i < MAXPROC; i++ {
                        if sospBT[i] == true {
                            risorsa[i] <- b
                            nsospBT--
                            sospBT[i] = false
                            break
                        }
                    }
                }

        case r = <-richiesta:
            // A new client request arrived
            switch r.tipo {
            case FLEX:
                // FLEX: if an electric bike is available, give EB, else try BT,
                // otherwise client must wait for an EB
                if dispEB > 0 {
                    dispEB--
                    b = EB
                    fmt.Printf("[server] allocated electric bike to FLEX client %d\n", r.id)
                    risorsa[r.id] <- b
                } else if dispBT > 0 {
                    dispBT--
                    b = BT
                    fmt.Printf("[server] allocated traditional bike to FLEX client %d\n", r.id)
                    risorsa[r.id] <- b
                } else {
                    // No bikes available, must wait for EB
                    nsospEB++
                    sospEB[r.id] = true
                }

            case BT:
                // Client specifically requested a traditional bike
                if dispBT > 0 {
                    dispBT--
                    b = BT
                    fmt.Printf("[server] allocated traditional bike to client %d\n", r.id)
                    risorsa[r.id] <- b
                } else {
                    // Must wait
                    nsospBT++
                    sospBT[r.id] = true
                }

            case EB:
                // Client requested an electric bike
                if dispEB > 0 {
                    dispEB--
                    b = EB
                    fmt.Printf("[server] allocated electric bike to client %d\n", r.id)
                    risorsa[r.id] <- b
                } else {
                    // Must wait
                    nsospEB++
                    sospEB[r.id] = true
                }
            }

        case <-termina:
            // All clients have finished, time to end
            fmt.Println("END OF SERVER!")
            // Signal back to main that server is done
            done <- 1
            return
        }
    }
}

func main() {
    var cli int
    var r req

    // Ask how many clients to create (up to MAXPROC)
    fmt.Printf("\n How many clients (max %d)? ", MAXPROC)
    fmt.Scanf("%d", &cli)
    fmt.Println("Number of clients:", cli)

    // Initialize risorsa channels: one channel per client
    for i := 0; i < MAXPROC; i++ {
        risorsa[i] = make(chan bici)
    }

    // Create client goroutines
    // We'll assign them types (BT, EB, FLEX) by taking i % 3
    //   0 => BT, 1 => EB, 2 => FLEX
    for i := 0; i < cli; i++ {
        r.id = i
        r.tipo = i % 3
        go client(r)
    }

    // Create the server goroutine
    go server()

    // Wait for all clients to finish
    for i := 0; i < cli; i++ {
        <-done
    }

    // Send a termination signal to the server
    termina <- 1

    // Wait for the server's final 'done'
    <-done
}
