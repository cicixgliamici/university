package main

import (
    "fmt"
    "math/rand"
    "time"
)

// CONSTANTS
const scolari = 25     // Number of people in a school group
const N = 40           // Max number of people allowed in the hall (including supervisors)
const NC = 30          // Max number of people allowed in the corridor (combined IN + OUT directions)
const MaxS = 4         // Max number of supervisors allowed in the hall
const MAXBUFF = 15
const MAXPROC = 5

// For corridor directions:
const IN int = 0
const OUT int = 1

// Types of visitors:
const SING int = 0  // single visitor
const SCOL int = 1  // school group of 25 people
const SORV int = 2  // supervisor

// The request structure is sent on channels when a process (visitor or supervisor) wants to move:
//   - id:   ID of the request (or goroutine)
//   - tipo: which type of entity (single, school group, or supervisor)
//   - ack:  a channel to receive acknowledgment (server replies with an int)
type richiesta struct {
    id   int
    tipo int
    ack  chan int
}

// Channels to enter the corridor in direction IN.
var entrataC_IN [3]chan richiesta

// Channels to enter the corridor in direction OUT.
var entrataC_OUT [3]chan richiesta

// Channels to exit the corridor (either from IN or OUT direction).
var uscitaC_IN = make(chan richiesta, MAXBUFF)
var uscitaC_OUT = make(chan richiesta, MAXBUFF)

// Channels for synchronization and termination
var done = make(chan bool, MAXBUFF)    // Signifies that a goroutine has completed
var termina = make(chan bool, MAXBUFF) // Tells the server to stop

// Utility function: prints the type of visitor/supervisor
func printTipo(typ int) string {
    switch typ {
    case SING:
        return "single visitor"
    case SCOL:
        return "school group"
    case SORV:
        return "supervisor"
    }
    return ""
}

// Utility function: prints direction (IN or OUT)
func printDirezione(typ int) string {
    switch typ {
    case IN:
        return "entrance"
    case OUT:
        return "exit"
    }
    return ""
}

// SERVER GOROUTINE
// Manages corridor usage (IN and OUT directions) and checks constraints:
//  - The corridor can hold at most NC people overall (IN + OUT).
//  - The hall can hold at most N people (including supervisors).
//  - At most MaxS supervisors in the hall at once.
//  - A school group has 25 members (they enter/exit as a block).
//  - Supervisors must be present for single visitors or school groups to enter, etc.
func server() {
    scolaresche_in_C := [2]int{0, 0} // number of school groups in the corridor, indexed by direction [IN, OUT]
    persone_in_C := [2]int{0, 0}     // number of people in the corridor, indexed by direction [IN, OUT]
    var persone_in_sala = 0          // how many people are currently in the hall
    var sorveglianti_in_sala = 0     // how many supervisors are currently in the hall

    for {
        select {
        // -----------------------------
        // ENTRANCE: corridor direction IN
        // 1) A SUPERVISOR enters the corridor IN
        // Conditions:
        //   - No school groups in the OUT corridor ( scolaresche_in_C[OUT] == 0 )
        //   - Total corridor usage < NC
        //   - People in the hall < N
        //   - Supervisors in the hall < MaxS
        //   - No pending school groups or supervisors or single visitors waiting to enter the corridor OUT
        case x := <-when(
            scolaresche_in_C[OUT] == 0 &&
                (persone_in_C[IN]+persone_in_C[OUT]) < NC &&
                persone_in_sala < N &&
                sorveglianti_in_sala < MaxS &&
                (len(entrataC_OUT[SCOL])+len(entrataC_OUT[SORV])+len(entrataC_OUT[SING]) == 0),
            entrataC_IN[SORV],
        ):
            persone_in_C[IN]++
            persone_in_sala++
            sorveglianti_in_sala++
            x.ack <- 1

        // 2) A SINGLE VISITOR enters the corridor IN
        // Conditions:
        //   - No school groups in OUT corridor
        //   - Corridor usage < NC
        //   - Hall usage < N
        //   - At least 1 supervisor in hall ( sorveglianti_in_sala > 0 )
        //   - No supervisors waiting in the IN corridor or anything else in the OUT corridor with higher priority
        case x := <-when(
            scolaresche_in_C[OUT] == 0 &&
                (persone_in_C[IN]+persone_in_C[OUT]) < NC &&
                persone_in_sala < N &&
                sorveglianti_in_sala > 0 &&
                (len(entrataC_IN[SORV])+len(entrataC_OUT[SCOL])+len(entrataC_OUT[SORV])+len(entrataC_OUT[SING]) == 0),
            entrataC_IN[SING],
        ):
            persone_in_C[IN]++
            persone_in_sala++
            x.ack <- 1

        // 3) A SCHOOL GROUP enters the corridor IN
        // Conditions:
        //   - No people in OUT corridor
        //   - Enough space in the corridor for 25 ( scolari )
        //   - Enough space in the hall for 25
        //   - At least 1 supervisor present
        //   - No one is queued in the IN corridor for supervisor/single visitor or anything in the OUT corridor
        case x := <-when(
            persone_in_C[OUT] == 0 &&
                (persone_in_C[IN]+persone_in_C[OUT])+scolari <= NC &&
                persone_in_sala+scolari <= N &&
                sorveglianti_in_sala > 0 &&
                (len(entrataC_IN[SORV])+len(entrataC_IN[SING])+len(entrataC_OUT[SCOL])+len(entrataC_OUT[SORV])+len(entrataC_OUT[SING]) == 0),
            entrataC_IN[SCOL],
        ):
            persone_in_C[IN] += scolari
            scolaresche_in_C[IN]++
            persone_in_sala += scolari
            x.ack <- 1

        // -----------------------------
        // ENTRANCE: corridor direction OUT
        // 4) A SUPERVISOR enters the corridor OUT
        // Conditions:
        //   - No school groups in the IN corridor
        //   - Corridor usage < NC
        //   - (sorveglianti_in_sala > 1) OR (persone_in_sala == 1) so there's still at least 1 supervisor or 0 people left
        //   - No school groups or single visitors in OUT corridor
        case x := <-when(
            scolaresche_in_C[IN] == 0 &&
                (persone_in_C[IN]+persone_in_C[OUT]) < NC &&
                (sorveglianti_in_sala > 1 || persone_in_sala == 1) &&
                (len(entrataC_OUT[SCOL])+len(entrataC_OUT[SING]) == 0),
            entrataC_OUT[SORV],
        ):
            persone_in_C[OUT]++
            persone_in_sala--
            sorveglianti_in_sala--
            x.ack <- 1

        // 5) A SINGLE VISITOR enters the corridor OUT
        // Conditions:
        //   - No school groups in the IN corridor
        //   - Corridor usage < NC
        //   - No school groups waiting in the OUT corridor
        case x := <-when(
            scolaresche_in_C[IN] == 0 &&
                (persone_in_C[IN]+persone_in_C[OUT]) < NC &&
                (len(entrataC_OUT[SCOL]) == 0),
            entrataC_OUT[SING],
        ):
            persone_in_C[OUT]++
            persone_in_sala--
            x.ack <- 1

        // 6) A SCHOOL GROUP enters the corridor OUT
        // Conditions:
        //   - Nobody in the IN corridor
        //   - Enough space in the corridor for 25 people
        case x := <-when(
            persone_in_C[IN] == 0 &&
                (persone_in_C[IN]+persone_in_C[OUT])+scolari <= NC,
            entrataC_OUT[SCOL],
        ):
            persone_in_C[OUT] += scolari
            scolaresche_in_C[OUT]++
            persone_in_sala -= scolari
            x.ack <- 1

        // -----------------------------
        // EXIT from the corridor (IN or OUT direction)
        // 7) A visitor or supervisor exiting from the IN corridor
        case x := <-uscitaC_IN:
            if x.tipo == SCOL {
                persone_in_C[IN] -= scolari
                scolaresche_in_C[IN]--
            } else {
                // single visitor or supervisor
                persone_in_C[IN]--
            }
            x.ack <- 1

        // 8) A visitor or supervisor exiting from the OUT corridor
        case x := <-uscitaC_OUT:
            if x.tipo == SCOL {
                persone_in_C[OUT] -= scolari
                scolaresche_in_C[OUT]--
            } else {
                // single visitor or supervisor
                persone_in_C[OUT]--
            }
            x.ack <- 1

        // -----------------------------
        // SERVER TERMINATION
        case <-termina: // all processes have finished
            fmt.Println("\nEND!!!")
            done <- true
            return
        }
    }
}

// GOROUTINE: Visitor (single or school group)
func visitatore(id int, tipo int) {
    // 'tt' is a random time used to simulate delays
    var tt int
    var r richiesta

    // Random initialization delay
    tt = rand.Intn(2) + 1
    fmt.Printf("\nInitializing visitor %d of type %s in %d seconds\n", id, printTipo(tipo), tt)
    time.Sleep(time.Duration(tt) * time.Second)

    // Prepare the request
    r = richiesta{id, tipo, make(chan int, MAXBUFF)}

    // 1) Enter corridor IN
    entrataC_IN[tipo] <- r
    <-r.ack
    fmt.Printf("\n[Visitor %d, type %s] entering corridor in direction IN\n", id, printTipo(tipo))

    // 2) Exit corridor IN
    tt = rand.Intn(2) + 1
    time.Sleep(time.Duration(tt) * time.Second)
    uscitaC_IN <- r
    <-r.ack
    fmt.Printf("\n[Visitor %d, type %s] entered the hall\n", id, printTipo(tipo))

    // 3) Visit/stay inside the hall
    tt = rand.Intn(5) + 1
    time.Sleep(time.Duration(tt) * time.Second)

    // 4) Enter corridor OUT
    entrataC_OUT[tipo] <- r
    <-r.ack
    fmt.Printf("\n[Visitor %d, type %s] entering corridor in direction OUT\n", id, printTipo(tipo))

    // 5) Exit corridor OUT
    tt = rand.Intn(2) + 1
    time.Sleep(time.Duration(tt) * time.Second)
    uscitaC_OUT <- r
    <-r.ack
    fmt.Printf("\n[Visitor %d, type %s] left the corridor in direction OUT and is going home...\n", id, printTipo(tipo))

    // Signal that this goroutine is done
    done <- true
}

// GOROUTINE: Supervisor
func sorvegliante(id int) {
    var tt int
    var r richiesta
    var tipo = SORV

    // The supervisor will enter and exit multiple times
    // so there is always a chance for at least one supervisor present in the hall.
    // 'volte' sets how many times they loop.
    volte := 2 * MAXPROC

    tt = rand.Intn(2) + 1
    fmt.Printf("\nInitializing supervisor %d in %d seconds...\n", id, tt)
    time.Sleep(time.Duration(tt) * time.Second)

    r = richiesta{id, tipo, make(chan int, MAXBUFF)}

    for i := 0; i < volte; i++ {
        // 1) Enter corridor IN
        entrataC_IN[tipo] <- r
        <-r.ack
        fmt.Printf("\n[Supervisor %d] entered corridor IN\n", id)

        tt = rand.Intn(2) + 1
        time.Sleep(time.Duration(tt) * time.Second)

        // 2) Exit corridor IN
        uscitaC_IN <- r
        <-r.ack
        fmt.Printf("\n[Supervisor %d] is now in the hall\n", id)

        // 3) Supervision time in the hall
        tt = rand.Intn(5) + 1
        time.Sleep(time.Duration(tt) * time.Second)

        // 4) Enter corridor OUT
        entrataC_OUT[tipo] <- r
        <-r.ack
        fmt.Printf("\n[Supervisor %d] entered corridor OUT\n", id)

        tt = rand.Intn(2) + 1
        time.Sleep(time.Duration(tt) * time.Second)

        // 5) Exit corridor OUT
        uscitaC_OUT <- r
        <-r.ack
        fmt.Printf("\n[Supervisor %d] left the corridor OUT\n", id)

        tt = rand.Intn(1) + 1
        time.Sleep(time.Duration(tt) * time.Second)
    }

    fmt.Printf("\n[Supervisor %d] done and going home...\n", id)
    done <- true
}

// Helper function for conditional select statements
func when(b bool, c chan richiesta) chan richiesta {
    if !b {
        return nil
    }
    return c
}

func main() {
    var scolaresche int
    var singoli int
    var sorveglianti int
    rand.Seed(time.Now().Unix())

    // Ask user how many school groups, singles, and supervisors to launch
    fmt.Printf("\n[main] How many school groups? (max %d)\n", MAXPROC)
    fmt.Scanf("%d", &scolaresche)
    fmt.Printf("\n[main] How many single visitors? (max %d)\n", MAXPROC)
    fmt.Scanf("%d", &singoli)
    fmt.Printf("\n[main] How many supervisors? (max %d)\n", MAXPROC)
    fmt.Scanf("%d", &sorveglianti)

    // Initialize the corridor IN/OUT channels for each of the 3 entity types (SING, SCOL, SORV)
    for i := 0; i < 3; i++ {
        entrataC_IN[i] = make(chan richiesta, MAXBUFF)
        entrataC_OUT[i] = make(chan richiesta, MAXBUFF)
    }

    // Start the server goroutine
    go server()

    // Create supervisor goroutines
    for i := 0; i < sorveglianti; i++ {
        go sorvegliante(i)
    }

    // Create single visitor goroutines
    for i := 0; i < singoli; i++ {
        go visitatore(i, SING)
    }

    // Create school group goroutines
    for i := 0; i < scolaresche; i++ {
        go visitatore(i, SCOL)
    }

    // Wait for all visitors and supervisors to finish
    for i := 0; i < (sorveglianti + singoli + scolaresche); i++ {
        <-done
        fmt.Printf("\n[main] The %d-th process has completed...\n\n", i+1)
    }

    fmt.Printf("\n[main] All user processes are done.\n\n")

    // Signal the server to terminate
    termina <- true
    <-done  // wait for the server to confirm it has ended

    fmt.Println()
}
