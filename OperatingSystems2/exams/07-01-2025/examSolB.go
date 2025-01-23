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

        fmt.Printf("[TRAINER %d]
