// -----------------------------------------------------------------------------------
// ENGLISH TRANSLATION AND EXPLANATION OF THE CODE
// This Go program simulates a shop where:
//   - Clients (either "regular" or "occasional") enter, need a mask, and must be
//     supervised by shop assistants (commessi).
//   - Shop assistants (commessi) can supervise up to 3 clients each.
//   - A supplier (fornitore) delivers batches of masks (NM each time).
//   - The shop has a maximum capacity (MAX) that counts both clients and assistants.
//   - There is a limited number of masks in the shop (mascherine). Each client
//     consumes exactly 1 mask when entering.
//   - Each assistant can come in and out multiple times but can only exit if they
//     have no currently assigned clients.
//   - Clients and assistants signal termination after finishing their tasks.
//   - The shop (negozio) server goroutine manages all the constraints by
//     reading requests from channels.
//
// The main concurrency elements are:
//   - Channels for clients entering (separated for regular vs. occasional clients),
//     for shop assistants entering and exiting, for clients exiting, and for
//     delivering mask batches (supplier).
//   - A single "shop" goroutine that enforces capacity constraints, mask availability,
//     assignment of clients to assistants, and assistant exit conditions.
//   - Separate goroutines for each client, each assistant, and the supplier.
//
// Data Structures:
//   - Richiesta: A struct used by clients/assistants to request entering/exiting,
//     carrying an 'id' and an 'ack' (acknowledgment) channel of type bool.
//   - Commesso: The state of a shop assistant, including whether they are
//     inside the shop, want to exit, how many clients they are assigned, etc.
//
// Flow of a client goroutine (cliente):
//   1) Sleeps a random time, then tries to enter the shop (sending a request on
//      either `entraClienteAbituale` or `entraClienteOccasionale`).
//   2) Waits for acknowledgment from the shop (the 'negozio' goroutine).
//   3) Simulates staying inside (sleep), then sends an exit request to `esciCliente`.
//   4) Signals its own termination.
//
// Flow of an assistant goroutine (commesso):
//   1) In a loop, sleeps randomly, requests to enter the shop by sending on `entraCommesso`.
//   2) Waits for ack. Stays inside a random time, then requests to exit on `esciCommesso`.
//   3) If the shop server allows them to leave (no assigned clients), they exit. Otherwise
//      they wait until all assigned clients have exited.
//   4) They can be signaled to terminate (via 'terminaCommesso'), in which case they finish.
//
// Flow of the supplier goroutine (fornitore):
//   1) Repeatedly sleeps a random time, then attempts to deliver a batch of masks
//      by sending `true` on 'deposita'.
//   2) Waits for the shop to consume that message and respond on the same channel
//      (so effectively 'deposita' is used for both signal and ack).
//   3) If the main program sends a termination signal (through 'terminaFornitore'),
//      the supplier finishes.
//
// The shop goroutine (negozio) uses a select block with conditions (via the helper
// function 'whenRichiesta') to manage concurrency. It tracks the state of all
// assistants in a 'commessi' array, how many masks are currently available, how
// many are inside the shop, and so forth. The shop ends when the main function
// sends a termination signal ('terminaNegozio').
//
// -----------------------------------------------------------------------------------

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// BUFFER AND CAPACITY CONSTANTS
const MAXBUFF int = 100  // general buffer size for channels
const MAX int = 18       // maximum capacity of the shop (clients + assistants)
const N_COMMESSI int = 8 // number of shop assistants
const N_CLIENTI int = 70 // total number of clients
const NM = 10            // each batch of masks delivered by the supplier

// CLIENT TYPES
const ABITUALE int = 0
const OCCASIONALE int = 1

// An array to print the client's type in a human-readable form.
var tipoClienteStr [2]string = [2]string{"ABITUALE", "OCCASIONALE"}

// Richiesta is used by both clients and assistants to request entry/exit.
// 'id' is the ID (unique to each goroutine).
// 'ack' is a channel on which the shop server (negozio) sends a boolean ack.
type Richiesta struct {
    id  int
    ack chan bool
}

// Commesso represents the state of a shop assistant:
//  - dentro:        whether the assistant is currently inside the shop
//  - vuoleUscire:   whether the assistant wants to exit but is waiting for
//                   assigned clients to finish
//  - clientiAssegnati: an array of up to 3 client IDs
//  - numeroClientiAssegnati: how many clients the assistant is currently supervising
//  - ackUscita:     a channel used to signal the assistant can exit
type Commesso struct {
    dentro                 bool
    vuoleUscire            bool
    clientiAssegnati       [3]int
    numeroClientiAssegnati int
    ackUscita              chan bool
}

// Helper function for conditional select on a channel of Richiesta.
func whenRichiesta(b bool, c chan Richiesta) chan Richiesta {
    if !b {
        return nil
    }
    return c
}

// Helper function for conditional select on a channel of int.
func whenInt(b bool, c chan int) chan int {
    if !b {
        return nil
    }
    return c
}

// Utility: sleeps a random time between 1 and timeLimit seconds
func sleepRandTime(timeLimit int) {
    if timeLimit > 0 {
        time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
    }
}

// GOROUTINE: Client (either ABITUALE or OCCASIONALE)
func cliente(id int, tipo int, entra chan Richiesta, esci chan int, termina chan bool) {
    // Prepare a request
    var ric Richiesta
    ric.id = id
    ric.ack = make(chan bool, MAXBUFF)

    // Simulate a random initialization time
    sleepRandTime(5)

    fmt.Printf("[CLIENT %s %d] I want to enter the shop...\n", tipoClienteStr[tipo], id)

    // Send a request to enter
    entra <- ric
    // Wait for acknowledgment
    <-ric.ack

    fmt.Printf("[CLIENT %s %d] I have entered the shop...\n", tipoClienteStr[tipo], id)

    // Simulate shopping / being inside
    sleepRandTime(7)

    // Now exit
    esci <- id
    fmt.Printf("[CLIENT %s %d] I have left the shop...\n", tipoClienteStr[tipo], id)

    fmt.Printf("[CLIENT %s %d] Terminating...\n", tipoClienteStr[tipo], id)

    // Signal that this client has finished
    termina <- true
    return
}

// GOROUTINE: Shop assistant (commesso)
func commesso(id int, entra chan Richiesta, esci chan Richiesta, termina chan bool, done chan bool) {
    // We'll reuse the same Richiesta structure every time they enter/exit
    var ric Richiesta
    ric.id = id
    ric.ack = make(chan bool, MAXBUFF)

    for {
        sleepRandTime(5)

        fmt.Printf("[ASSISTANT %d] I want to enter the shop...\n", id)

        // Request to enter
        entra <- ric
        // Wait for ack
        <-ric.ack

        fmt.Printf("[ASSISTANT %d] I have entered the shop...\n", id)

        // Simulate working inside
        sleepRandTime(9)

        // Request to exit
        esci <- ric
        <-ric.ack

        fmt.Printf("[ASSISTANT %d] I have left the shop...\n", id)

        // Check if we should terminate
        select {
        case <-termina:
            {
                fmt.Printf("[ASSISTANT %d] Terminating...\n", id)
                done <- true
                return
            }
        default:
            {
                // Not terminating yet, wait a bit
                sleepRandTime(2)
            }
        }
    }
}

// GOROUTINE: Supplier (fornitore)
// Delivers NM masks every time it can, repeatedly, until it is asked to terminate.
func fornitore(deposita chan bool, termina chan bool) {
    for {
        sleepRandTime(5)

        fmt.Printf("[SUPPLIER] I want to deliver a batch of masks...\n")
        // Send a signal that we have a batch to deposit
        deposita <- true
        // Wait for the shop to confirm
        <-deposita
        fmt.Printf("[SUPPLIER] Delivery completed...\n")

        // Check if we should terminate
        select {
        case <-termina:
            {
                fmt.Printf("[SUPPLIER] Terminating...\n")
                termina <- true
                return
            }
        default:
            {
                sleepRandTime(2)
            }
        }
    }
}

// GOROUTINE: The "shop" (negozio) server
// This goroutine manages:
//   - The maximum capacity inside (clients + assistants <= MAX).
//   - How many assistants are inside, how many are free (supervising < 3 clients).
//   - How many masks are available (mascherine).
//   - The assignment of clients to assistants, so each assistant can supervise up to 3.
//   - Whether an assistant can exit (only if they have 0 assigned clients).
//   - The supplier's deliveries of masks.
func negozio(
    entraClienteAbituale chan Richiesta,
    entraClienteOccasionale chan Richiesta,
    entraCommesso chan Richiesta,
    esciCliente chan int,
    esciCommesso chan Richiesta,
    deposita chan bool,
    termina chan bool,
) {

    // Track how many clients and assistants are inside
    clientiDentro := 0
    commessiDentro := 0

    // How many assistants are currently free (supervising < 3 clients)
    commessiLiberi := 0

    // Array that holds the state of each assistant
    commessi := make([]Commesso, N_COMMESSI)

    // Initialize each assistant's data
    for i := 0; i < N_COMMESSI; i++ {
        commessi[i].dentro = false
        commessi[i].vuoleUscire = false
        commessi[i].numeroClientiAssegnati = 0
        commessi[i].ackUscita = nil
        for j := 0; j < 3; j++ {
            commessi[i].clientiAssegnati[j] = -1
        }
    }

    // Number of masks currently available
    mascherine := 0

    fmt.Printf("MAX: %d, NM: %d, N_CLIENTI: %d, N_COMMESSI: %d...\n", MAX, NM, N_CLIENTI, N_COMMESSI)

    // Main loop of the shop server
    for {
        fmt.Printf("[SHOP] ClientsInside: %d, AssistantsInside: %d, FreeAssistants: %d, Masks: %d...\n",
            clientiDentro, commessiDentro, commessiLiberi, mascherine)

        select {
        // 1) Supplier deposit
        case <-deposita:
            {
                mascherine += NM
                fmt.Printf("[SHOP] The supplier delivered %d masks...\n", NM)
                // Send ack back to the supplier on the same channel
                deposita <- true
            }

        // 2) An assistant wants to enter the shop
        case ric := <-whenRichiesta(clientiDentro+commessiDentro < MAX, entraCommesso):
            {
                commessiDentro++
                commessiLiberi++
                commessi[ric.id].dentro = true
                commessi[ric.id].vuoleUscire = false
                commessi[ric.id].numeroClientiAssegnati = 0
                for i := 0; i < 3; i++ {
                    commessi[ric.id].clientiAssegnati[i] = -1
                }
                fmt.Printf("[SHOP] Assistant %d enters the shop...\n", ric.id)
                ric.ack <- true
            }

        // 3) An assistant requests to exit the shop
        case ric := <-esciCommesso:
            {
                if commessi[ric.id].numeroClientiAssegnati == 0 {
                    // If the assistant has no assigned clients, they can exit immediately
                    fmt.Printf("[SHOP] Assistant %d leaves the shop...\n", ric.id)
                    commessi[ric.id].dentro = false
                    commessi[ric.id].vuoleUscire = false
                    commessi[ric.id].ackUscita = nil
                    ric.ack <- true
                    commessiLiberi--
                    commessiDentro--
                } else {
                    // The assistant must wait until all clients are done
                    fmt.Printf("[SHOP] Assistant %d wants to exit but is waiting (%d assigned clients)...\n",
                        ric.id, commessi[ric.id].numeroClientiAssegnati)
                    commessi[ric.id].vuoleUscire = true
                    commessi[ric.id].ackUscita = ric.ack
                }
            }

        // 4) A REGULAR client (ABITUALE) wants to enter
        //    Conditions:
        //      - There is at least 1 assistant inside and free
        //      - At least 1 mask available
        //      - The shop is not full
        //      - No one is queued in entraCommesso
        //      - The queue for occasional clients does not have priority
        case ric := <-whenRichiesta(
            commessiDentro > 0 && commessiLiberi > 0 && mascherine >= 1 &&
                len(entraCommesso) == 0 &&
                (clientiDentro+commessiDentro < MAX),
            entraClienteAbituale,
        ):
            {
                found := false
                // Search for a free assistant with <3 assigned clients
                for i := 0; i < N_COMMESSI && !found; i++ {
                    if commessi[i].dentro && commessi[i].numeroClientiAssegnati < 3 {
                        for j := 0; j < 3 && !found; j++ {
                            if commessi[i].clientiAssegnati[j] < 0 {
                                // Assign this client to the assistant
                                commessi[i].clientiAssegnati[j] = ric.id
                                commessi[i].numeroClientiAssegnati++
                                if commessi[i].numeroClientiAssegnati == 3 {
                                    // This assistant is now fully occupied
                                    commessiLiberi--
                                }
                                clientiDentro++
                                mascherine--
                                found = true
                                ric.ack <- true
                                fmt.Printf("[SHOP] Regular client %d enters the shop...\n", ric.id)
                                fmt.Printf("[SHOP] Assigning assistant %d to regular client %d...\n", i, ric.id)
                            }
                        }
                    }
                }
                if !found {
                    fmt.Printf("[DEBUG SHOP] Unable to find a free assistant for a REGULAR client...\n")
                }
            }

        // 5) An OCCASIONAL client wants to enter
        //    Conditions:
        //      - No one is queued in entraClienteAbituale (regular clients have priority)
        //      - At least 1 free assistant
        //      - At least 1 mask available
        //      - The shop is not full
        //      - No one is queued in entraCommesso
        case ric := <-whenRichiesta(
            len(entraClienteAbituale) == 0 && commessiDentro > 0 && commessiLiberi > 0 && mascherine >= 1 &&
                len(entraCommesso) == 0 &&
                (clientiDentro+commessiDentro < MAX),
            entraClienteOccasionale,
        ):
            {
                found := false
                // Search for a free assistant
                for i := 0; i < N_COMMESSI && !found; i++ {
                    if commessi[i].dentro && commessi[i].numeroClientiAssegnati < 3 {
                        for j := 0; j < 3 && !found; j++ {
                            if commessi[i].clientiAssegnati[j] < 0 {
                                commessi[i].clientiAssegnati[j] = ric.id
                                commessi[i].numeroClientiAssegnati++
                                if commessi[i].numeroClientiAssegnati == 3 {
                                    commessiLiberi--
                                }
                                clientiDentro++
                                mascherine--
                                found = true
                                ric.ack <- true
                                fmt.Printf("[SHOP] Occasional client %d enters the shop...\n", ric.id)
                                fmt.Printf("[SHOP] Assigning assistant %d to occasional client %d...\n", i, ric.id)
                            }
                        }
                    }
                }
                if !found {
                    fmt.Printf("[DEBUG SHOP] Unable to find a free assistant for an OCCASIONAL client...\n")
                }
            }

        // 6) A client exits the shop (esciCliente)
        case id := <-esciCliente:
            {
                found := false
                // Find which assistant was assigned to this client
                for i := 0; i < N_COMMESSI && !found; i++ {
                    if commessi[i].dentro {
                        for j := 0; j < 3 && !found; j++ {
                            if commessi[i].clientiAssegnati[j] == id {
                                // Free that slot
                                commessi[i].clientiAssegnati[j] = -1

                                if commessi[i].numeroClientiAssegnati == 3 {
                                    // This assistant was at full capacity, now has space
                                    commessiLiberi++
                                }
                                commessi[i].numeroClientiAssegnati--
                                clientiDentro--
                                found = true
                                fmt.Printf("[SHOP] Client %d leaves the shop...\n", id)
                                fmt.Printf("[SHOP] Freeing assistant %d from supervising client %d...\n", i, id)

                                // Check if the assistant was waiting to exit
                                if commessi[i].dentro && commessi[i].vuoleUscire && commessi[i].numeroClientiAssegnati == 0 {
                                    // The assistant can now exit
                                    fmt.Printf("[SHOP] Assistant %d leaves the shop...\n", i)
                                    commessi[i].dentro = false
                                    commessi[i].vuoleUscire = false
                                    commessi[i].ackUscita <- true
                                    commessi[i].ackUscita = nil
                                    // Reset the assigned clients array
                                    for j := 0; j < 3; j++ {
                                        commessi[i].clientiAssegnati[j] = -1
                                    }
                                    commessiLiberi--
                                    commessiDentro--
                                }
                            }
                        }
                    }
                }
            }

        // 7) The shop receives a termination signal
        case <-termina:
            {
                fmt.Printf("[SHOP] Terminating...\n")
                termina <- true
                return
            }
        }
    }
}

func main() {
    // Channels for clients:
    //   - separate channels for regular (abituale) and occasional (occasionale) entry
    //   - a shared channel for exiting
    entraClienteAbituale := make(chan Richiesta, MAXBUFF)
    entraClienteOccasionale := make(chan Richiesta, MAXBUFF)

    // Channels for assistants entering and exiting
    entraCommesso := make(chan Richiesta, MAXBUFF)
    esciCliente := make(chan int)
    esciCommesso := make(chan Richiesta)

    // Channel used by the supplier to deposit mask batches
    deposita := make(chan bool)

    // Termination signals
    terminaCliente := make(chan bool)             // used by clients
    terminaCommesso := make([]chan bool, N_COMMESSI) // one channel per assistant
    done := make(chan bool)                       // used by assistants to confirm they've terminated
    terminaFornitore := make(chan bool)           // used by the supplier
    terminaNegozio := make(chan bool)             // used by the shop

    // Seed random generator
    rand.Seed(time.Now().Unix())

    // Create client goroutines
    for i := 0; i < N_CLIENTI; i++ {
        var tipo int
        var entraCliente chan Richiesta

        // 30% chance to be regular (ABITUALE), 70% to be occasional (OCCASIONALE),
        // based on a random threshold
        if rand.Intn(100) > 70 {
            tipo = ABITUALE
            entraCliente = entraClienteAbituale
        } else {
            tipo = OCCASIONALE
            entraCliente = entraClienteOccasionale
        }
        go cliente(i, tipo, entraCliente, esciCliente, terminaCliente)
    }

    // Create assistant goroutines
    for i := 0; i < N_COMMESSI; i++ {
        terminaCommesso[i] = make(chan bool, MAXBUFF)
        go commesso(i, entraCommesso, esciCommesso, terminaCommesso[i], done)
    }

    // Create supplier goroutine
    go fornitore(deposita, terminaFornitore)

    // Create the shop server goroutine
    go negozio(
        entraClienteAbituale,
        entraClienteOccasionale,
        entraCommesso,
        esciCliente,
        esciCommesso,
        deposita,
        terminaNegozio,
    )

    // Wait for all clients to terminate
    for i := 0; i < N_CLIENTI; i++ {
        <-terminaCliente
    }

    // Terminate the supplier
    terminaFornitore <- true
    <-terminaFornitore

    // Signal each assistant to terminate
    for i := 0; i < N_COMMESSI; i++ {
        terminaCommesso[i] <- true
    }

    // Wait for each assistant to confirm termination
    for i := 0; i < N_COMMESSI; i++ {
        <-done
    }

    // Finally, terminate the shop
    terminaNegozio <- true
    <-terminaNegozio
}

// -----------------------------------------------------------------------------------
// END OF CODE
//
// Explanation recap (in brief, repeated from above):
//  - "negozio" is the central server goroutine controlling shop capacity, assistant
//    availability (each can supervise up to 3 clients), and mask counts.
//  - "cliente" goroutines represent clients entering, shopping, and exiting,
//    each needing 1 mask. They cannot enter unless there's at least 1 free
//    assistant and a mask available.
//  - "commesso" goroutines represent shop assistants who repeatedly enter and
//    exit. They cannot exit if they still have assigned clients, so they must
//    wait until all assigned clients have left.
//  - "fornitore" repeatedly delivers batches of masks (NM each time).
//  - The main function spawns all these goroutines, waits for the clients to
//    finish, then signals the supplier and assistants to finish, and finally
//    signals the shop to terminate.
//
// This pattern (one server goroutine + multiple request goroutines) is a common
// approach in Go to manage shared state safely without explicit locks, relying
// on channels and 'select' to serialize state changes.
// -----------------------------------------------------------------------------------
