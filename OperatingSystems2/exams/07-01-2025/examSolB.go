package main

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)

// CONSTANTS
const MAXBUFF = 100      // Maximum buffer size for channels
const MAXPROC = 10       // Not used directly in this example, but available if needed
const MAXCICLI = 4       // Maximum number of activity cycles per user

// Identifiers for different areas
const AREAPESI = 0
const AREACORSI = 1
const NumAree = 2        // Number of different resource types (two areas)

// Capacity constraints
const NP = 15   // Maximum number of people allowed in the weights area
const NT = 5    // Number of personal trainers
const MAX = 18  // Overall gym capacity (all users combined)

// Request is sent across channels when a user or trainer wants to enter/exit an area.
// 'id' is the ID of the requesting goroutine (user or trainer).
// 'tipo' indicates which area (AREAPESI or AREACORSI) for a user, or it could be used to store extra info.
// 'ack' is a channel where the server sends a boolean response (true/false).
type Request struct {
    id   int
    tipo int
    ack  chan bool
}

// Trainer state
// - dentro:       whether the trainer is currently inside the gym
// - vuoleUscire:  whether the trainer wants to exit but is currently busy
// - utenteAssegnato: which user is assigned to this trainer (-1 if none)
// - ackUscita:    a channel used to acknowledge trainer exit once they are free
type Trainer struct {
    dentro          bool
    vuoleUscire     bool
    utenteAssegnato int
    ackUscita       chan bool
}

// CHANNELS
// For users entering each area:
var IngressoArea [NumAree]chan Request

// Single channel for user exit (they indicate from which area they are exiting via 'tipo')
var Uscita = make(chan Request, MAXBUFF)

// For personal trainers entering (IngressoPT) and exiting (UscitaPT)
var IngressoPT = make(chan Request, MAXBUFF)
var UscitaPT = make(chan Request)

// CHANNELS for termination
var done = make(chan bool)           // Signals that a goroutine has finished
var termina = make(chan bool)        // Signals trainers that it’s time to stop
var terminaServer = make(chan bool)  // Signals the server that it should terminate

// Helper function that returns the given channel if b is true, or nil if b is false.
// Used in select statements to conditionally enable/disable a case.
func when(b bool, c chan Request) chan Request {
    if !b {
        return nil
    }
    return c
}

// Sleep for a random duration between 1 and timeLimit seconds
func sleepRandTime(timeLimit int) {
    if timeLimit > 0 {
        time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
    }
}

// Utility function to convert area type to string
func getTipo(t int) string {
    switch t {
    case AREAPESI:
        return "Area Pesi"
    case AREACORSI:
        return "Area Corsi"
    default:
        return ""
    }
}

// GOROUTINE: User
// A user will perform a random number of cycles (up to MAXCICLI).
// In each cycle, the user:
// 1) Chooses a random area (weights or courses).
// 2) Requests entry via IngressoArea[tipo], then waits for ack.
// 3) Sleeps to simulate training.
// 4) Requests exit by sending on Uscita, then waits for ack.
func utente(id int) {
    fmt.Printf("[USER %d] Start...\n", id)
    r := Request{id, -1, make(chan bool, MAXBUFF)}

    cycles := rand.Intn(MAXCICLI) + 1 // up to MAXCICLI times

    for i := 0; i < cycles; i++ {
        // Choose an area at random
        tipo := rand.Intn(NumAree)
        r.tipo = tipo

        fmt.Printf("[USER %d] requests to enter %s\n", id, strings.ToUpper(getTipo(tipo)))
        IngressoArea[tipo] <- r     // ask to enter
        <-r.ack                     // wait for server acknowledgment

        fmt.Printf("[USER %d] training in %s...\n", id, strings.ToUpper(getTipo(tipo)))
        sleepRandTime(5)

        fmt.Printf("[USER %d] leaving %s\n", id, strings.ToUpper(getTipo(tipo)))
        Uscita <- r                // request to exit
        <-r.ack                    // wait for server acknowledgment
    }

    fmt.Printf("[USER %d] finished and leaving the gym completely\n", id)
    done <- true
}

// GOROUTINE: Personal Trainer
// A trainer repeatedly:
// 1) Requests to enter "Area Corsi" (symbolically) via IngressoPT.
// 2) Sleeps to simulate being inside.
// 3) Requests to exit via UscitaPT.
// 4) Checks whether it’s time to stop (via 'termina' channel). If so, exits.
func trainer(id int) {
    var req Request
    req.id = id
    req.tipo = AREAPESI // Not really relevant, but we store a default
    req.ack = make(chan bool, MAXBUFF)

    for {
        // Some random idle time before asking to enter
        sleepRandTime(5)

        fmt.Printf("[TRAINER %d] wants to enter AREA CORSI...\n", id)
        IngressoPT <- req
        <-req.ack

        fmt.Printf("[TRAINER %d] is now inside...\n", id)
        sleepRandTime(15)

        UscitaPT <- req
        <-req.ack

        fmt.Printf("[TRAINER %d] has exited...\n", id)

        // Check if we should terminate
        select {
        case <-termina:
            fmt.Printf("[TRAINER %d] done!\n", id)
            done <- true
            return
        default:
            // Continue if no termination signal
            sleepRandTime(2)
        }
    }
}

// SERVER GOROUTINE: "palestra" (the gym)
// Manages all entries (users to weights area or courses area, and trainers) and exits.
// Maintains state variables about how many users and trainers are inside, and which user
// is assigned to which trainer.
func palestra() {
    utentiInPalestra := 0             // total users in the gym
    utentiInAP := 0                   // users in the weights area
    trainer := make([]Trainer, NT)    // state of each trainer

    // Initialize trainer state
    for i := 0; i < NT; i++ {
        trainer[i].dentro = false
        trainer[i].vuoleUscire = false
        trainer[i].utenteAssegnato = -1
        trainer[i].ackUscita = nil
    }

    trainerLiberi := 0   // how many trainers are free (not assigned to a user)
    trainerDentro := 0   // how many trainers are currently inside the gym

    fmt.Printf("[GYM] Opened!\n")

    for {
        select {
        // 1) User entering the WEIGHTS area (AREAPESI)
        //    Condition: total users < MAX, users in weights area < NP,
        //    and nobody is waiting to enter the courses area first
        case r := <-when(utentiInPalestra < MAX && utentiInAP < NP && len(IngressoArea[AREACORSI]) == 0, IngressoArea[AREAPESI]):
            utentiInPalestra++
            utentiInAP++
            fmt.Printf("[GYM] User %d entered the weights area.\n", r.id)
            r.ack <- true

        // 2) User entering the COURSES area (AREACORSI)
        //    Condition: total users < MAX, at least 1 free trainer,
        //    and no trainers waiting to enter (IngressoPT) at the moment
        case r := <-when(utentiInPalestra < MAX && trainerLiberi > 0 && len(IngressoPT) == 0, IngressoArea[AREACORSI]):
            utentiInPalestra++
            // Search for a free trainer
            found := false
            i := 0
            for i = 0; i < NT && !found; i++ {
                if trainer[i].utenteAssegnato == -1 && trainer[i].dentro {
                    found = true
                    trainer[i].utenteAssegnato = r.id // trainer i is assigned to user r.id
                }
            }
            trainerLiberi--
            fmt.Printf("[GYM] User %d is in the courses area, training with trainer %d.\n", r.id, i)
            r.ack <- true

        // 3) A trainer requests to enter
        case r := <-IngressoPT:
            fmt.Printf("[GYM] Trainer %d entered.\n", r.id)
            trainer[r.id].dentro = true
            trainer[r.id].vuoleUscire = false
            trainer[r.id].utenteAssegnato = -1
            trainer[r.id].ackUscita = nil
            trainerDentro++
            trainerLiberi++
            r.ack <- true

        // 4) A user requests to exit from either area
        case r := <-Uscita:
            utentiInPalestra--
            fmt.Printf("[GYM] User %d exiting from %s\n", r.id, strings.ToUpper(getTipo(r.tipo)))

            // If exiting from courses area, free up the trainer assigned to that user
            if r.tipo == AREACORSI {
                found := false
                for i := 0; i < NT && !found; i++ {
                    if trainer[i].utenteAssegnato == r.id {
                        found = true
                        trainer[i].utenteAssegnato = -1
                        trainerLiberi++
                        // If this trainer wanted to exit but was waiting for the user to finish:
                        if trainer[i].vuoleUscire && trainer[i].dentro {
                            fmt.Printf("[GYM] Trainer %d is now allowed to exit the gym...\n", i)
                            trainer[i].dentro = false
                            trainer[i].vuoleUscire = false
                            trainer[i].ackUscita <- true  // let the trainer exit
                            trainer[i].ackUscita = nil
                            trainerDentro--
                            trainerLiberi--
                        }
                    }
                }
            } else {
                // Exiting from the weights area
                utentiInAP--
            }
            r.ack <- true

        // 5) A trainer requests to exit
        case req := <-UscitaPT:
            fmt.Printf("[GYM] Trainer %d is asking to exit...\n", req.id)
            if trainer[req.id].utenteAssegnato == -1 {
                // Trainer is not assigned to any user, can exit right away
                fmt.Printf("[GYM] Trainer %d is free and is leaving the gym...\n", req.id)
                trainer[req.id].dentro = false
                trainer[req.id].vuoleUscire = false
                trainer[req.id].ackUscita = nil
                trainerLiberi--
                trainerDentro--
                req.ack <- true
            } else {
                // Trainer is busy with a user -> must wait
                fmt.Printf("[GYM] Trainer %d is busy and waits to exit.\n", req.id)
                trainer[req.id].vuoleUscire = true
                trainer[req.id].ackUscita = req.ack
            }

        // 6) The server receives a termination signal
        case <-terminaServer:
            fmt.Printf("[GYM] Closing.\n")
            done <- true
            return
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    var nUtenti int
    var nTrainer int

    nUtenti = 50
    nTrainer = NT

    // Initialize channels for user entry in both areas
    for i := 0; i < len(IngressoArea); i++ {
        IngressoArea[i] = make(chan Request, MAXBUFF)
    }

    // Start the server goroutine (the gym)
    go palestra()

    // Create trainer goroutines
    for i := 0; i < nTrainer; i++ {
        go trainer(i)
    }

    // Create user goroutines
    for i := 0; i < nUtenti; i++ {
        go utente(i)
    }

    // Wait for all users to finish
    for i := 0; i < nUtenti; i++ {
        <-done
    }

    // Signal all trainers to terminate and wait for them
    for i := 0; i < nTrainer; i++ {
        termina <- true
        <-done
    }

    // Finally, tell the server to terminate
    terminaServer <- true
    <-done

    fmt.Printf("\n\n[MAIN] The gym is closed!\n")
}
