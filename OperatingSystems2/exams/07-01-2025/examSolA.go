package main

import (
    "fmt"
    "math/rand"
    "time"
)

// CONSTANTS:
const MAXBUFF = 100     // Maximum buffer size for channels
const MAXPROC = 10      // Number of user processes (goroutines)
const NT = 4            // Number of available physiotherapists
const MAX = 8           // Maximum number of people allowed inside (FUN + PHYSIO) at the same time
const FUN = 0           // Index representing the FUN area
const PHYSIO = 1        // Index representing the PHYSIO area

type Request struct {
    id  int        // ID of the requester (User or Lifeguard)
    ack chan int   // A channel used to receive acknowledgments
}

// Names for the two areas, used only for printing/logging
var Area [2]string = [2]string{"FUN", "PHYSIO"}

// CHANNELS:
var userEntry [2]chan Request     // userEntry[FUN] and userEntry[PHYSIO] for entering the respective areas
var lifeguardEntry = make(chan Request, MAXBUFF) // Lifeguards entering FUN area
var userExit [2]chan Request      // userExit[FUN] and userExit[PHYSIO] for exiting the respective areas
var lifeguardExit = make(chan Request, MAXBUFF)  // Lifeguards exiting FUN area

// CHANNELS for termination:
var done = make(chan bool)           // Signals that a goroutine (User or Lifeguard) has finished
var terminate = make(chan bool)      // Signals to the server that the entire center should terminate
var closeCenter = make(chan bool)    // Signals that all users are done, and we can proceed to close the center

// when is a helper function often used in Go concurrency examples to conditionally enable or disable a case in a select statement.
// If 'b' is false, the returned channel is nil. Reading from a nil channel blocks forever, effectively disabling that case.
func when(b bool, c chan Request) chan Request {
    if !b {
        return nil
    }
    return c
}

// sleepRandTime sleeps for a random amount of time between 1 and timeLimit seconds
func sleepRandTime(timeLimit int) {
    if timeLimit > 0 {
        time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
    }
}

// GOROUTINE: User
func User(id int) {
    fmt.Printf("[User %d] Starting\n", id)
    rounds := rand.Intn(4) // Each user will enter and exit a random number of times (0 to 3)

    for i := 0; i < rounds; i++ {
        fmt.Printf("[User %d] Round #%d\n", id, i+1)

        areaType := rand.Intn(2) // Randomly choose which area (FUN or PHYSIO) the user wants to enter
        r := Request{id, make(chan int)}

        fmt.Printf("[User %d] Wants to enter area %s\n", id, Area[areaType])
        userEntry[areaType] <- r     // Send request to enter
        <-r.ack                      // Wait for acknowledgment from the server

        // User spends some time in the chosen area
        sleepRandTime(7)

        // Now the user exits
        userExit[areaType] <- r
        <-r.ack
        fmt.Printf("[User %d] Has exited area %s\n", id, Area[areaType])

        // Optional break before possibly entering again
        sleepRandTime(2)
    }

    // Signal that the user has finished all their entries/exits
    done <- true
}

// GOROUTINE: Lifeguard
// The lifeguard terminates only when it receives a special signal from the server (ack = -1).
func lifeguard(id int) {
    fmt.Printf("[Lifeguard %d] Starting shift\n", id)
    for {
        r := Request{id, make(chan int)}

        // Attempt to enter (only relevant to the FUN area in this code)
        lifeguardEntry <- r
        res := <-r.ack // Wait for server's response

        if res == -1 {
            // Server is telling the lifeguard to terminate
            fmt.Printf("[Lifeguard %d] Must terminate now\n", id)
            done <- true
            return
        } else {
            fmt.Printf("[Lifeguard %d] Entered\n", id)
        }

        // Lifeguard stays for a random time inside the FUN area
        sleepRandTime(10)

        // Lifeguard exits
        lifeguardExit <- r
        <-r.ack
        fmt.Printf("[Lifeguard %d] Exited\n", id)

        // Wait briefly before potentially entering again
        sleepRandTime(2)
    }
}

// SERVER GOROUTINE:
// Manages the logic of who can enter/exit the FUN and PHYSIO areas based on:
// - The limit of total people inside (MAX).
// - The availability of lifeguards for the FUN area (must have at least 1 lifeguard present for FUN).
// - The availability of physiotherapists for PHYSIO area (NT physiotherapists).
// - A closing condition where eventually no new lifeguards are allowed in.
func server() {
    // State variables
    nFUN := 0                // Number of users in FUN area
    nPHYSIO := 0             // Number of users in PHYSIO area
    freePhysiotherapists := NT    // How many physiotherapists are available
    nLifeguards := 0         // Number of lifeguards currently in the FUN area
    isClosing := false       // Whether the center is in the closing phase

    fmt.Printf("[SERVER] Starting\n")

    for {
        // Print current state (for logging / debugging)
        fmt.Printf("\nCURRENT STATE:\n")
        fmt.Printf("  Users in FUN: %d, Users in PHYSIO: %d\n", nFUN, nPHYSIO)
        fmt.Printf("  Free physiotherapists: %d, Lifeguards present: %d\n", freePhysiotherapists, nLifeguards)
        fmt.Printf("  Queue lengths: Lifeguard entry: %d, User entry FUN: %d, User entry PHYSIO: %d, Lifeguard exit: %d\n\n",
            len(lifeguardEntry), len(userEntry[FUN]), len(userEntry[PHYSIO]), len(lifeguardExit),
        )

        select {
        // 1) Lifeguard entering (if the center is not closing yet)
        case r := <-when(!isClosing, lifeguardEntry):
            nLifeguards++
            r.ack <- 1

        // 2) User entering FUN area
        //    Conditions: total inside < MAX, at least 1 lifeguard present, and no lifeguards waiting in queue
        case r := <-when((nFUN+nPHYSIO < MAX) && (nLifeguards > 0) && (len(lifeguardEntry) == 0), userEntry[FUN]):
            nFUN++
            r.ack <- 1

        // 3) Lifeguard entry attempt when the center is closing (isClosing == true)
        //    This means the server will respond with -1 (terminate signal)
        case r := <-when(isClosing, lifeguardEntry):
            r.ack <- -1

        // 4) User entering PHYSIO area
        //    Conditions: total inside < MAX, there is at least one free physiotherapist,
        //    and no FUN users waiting in queue (trying to avoid blocking)
        case r := <-when((nFUN+nPHYSIO < MAX) && (freePhysiotherapists > 0) && (len(userEntry[FUN]) == 0), userEntry[PHYSIO]):
            nPHYSIO++
            freePhysiotherapists--
            r.ack <- 1

        // 5) Lifeguard exiting FUN area
        //    Condition: can exit if there's more than one lifeguard OR there are no users in FUN
        case r := <-when((nLifeguards > 1) || (nFUN == 0), lifeguardExit):
            nLifeguards--
            r.ack <- 1

        // 6) User exiting FUN area
        case r := <-userExit[FUN]:
            nFUN--
            r.ack <- 1

        // 7) User exiting PHYSIO area
        case r := <-userExit[PHYSIO]:
            nPHYSIO--
            freePhysiotherapists++
            r.ack <- 1

        // 8) Signal that all users have completely finished (close the center).
        //    This means we should not allow new lifeguards to enter and will let existing ones exit.
        case <-closeCenter:
            isClosing = true
            fmt.Printf("The Center is about to close...\n")

        // 9) Main termination signal: the center is officially closed.
        //    When we receive this, we end the server goroutine.
        case <-terminate:
            fmt.Println("The Center is now closed!")
            done <- true
            return
        }
    }
}

func main() {
    fmt.Printf("[MAIN] Starting\n\n")
    rand.Seed(time.Now().UnixNano())

    // Channel initialization for 2 areas: FUN and PHYSIO
    for i := 0; i < 2; i++ {
        userEntry[i] = make(chan Request, MAXBUFF)
        userExit[i] = make(chan Request, MAXBUFF)
    }

    // Launch server goroutine
    go server()

    // Launch user goroutines
    for i := 0; i < MAXPROC; i++ {
        go User(i)
    }

    // Launch lifeguard goroutines
    for i := 0; i < MAXPROC/2; i++ {
        go lifeguard(i)
    }

    // Wait for all users to finish
    for i := 0; i < MAXPROC; i++ {
        <-done
    }
    fmt.Printf("\nAll users have finished!\n\n")

    // Signal the server that all users are done
    closeCenter <- true

    // Wait for all lifeguards to finish
    for i := 0; i < MAXPROC/2; i++ {
        <-done
    }
    fmt.Printf("\nAll lifeguards have finished!\n\n")

    // Finally, tell the server to shut down the center
    terminate <- true
    <-done

    fmt.Printf("\n[MAIN] End\n")
}
