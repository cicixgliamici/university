// -----------------------------------------------------------------------------------
// ENGLISH TRANSLATION WITH IN-CODE EXPLANATION
//
// This Go program simulates a small factory scenario where two different models of cars
// (A and B) are assembled by two robots (RobotA for model A, RobotB for model B). Each
// car needs 4 wheels composed of two parts: a "wheel rim" (cerchio) and a "tire" 
// (pneumatico). The parts come in two variations:
//
//   - For model A: cerchio A (CA) and pneumatico A (PA)
//   - For model B: cerchio B (CB) and pneumatico B (PB)
//
// We have:
//   - 4 conveyor belts (nastri) continuously delivering parts: PA, PB, CA, CB. 
//   - A storage "deposit" that can hold up to maxP = 3 tires (pneumatici) 
//     and maxC = 3 wheel rims (cerchi). 
//   - 2 robots (one for model A, one for model B) that pick up the needed parts 
//     from the deposit and assemble 4 wheels for each car.
//
// Each conveyor belt runs in a loop, delivering parts of its assigned type. When
// it tries to deliver a part, it sends a message (e.g., consegnaPA) to the deposit. 
// The deposit either accepts (and stores it if there's room) or eventually signals
// termination (sending an ack of -1) if the production is finished. 
//
// Each robot runs in a loop, building cars. For each wheel:
//   - Robot A picks up one CA (rim A) and one PA (tire A).
//   - Robot B picks up one CB (rim B) and one PB (tire B).
// Once a robot finishes assembling 4 wheels, that means it finished 1 car of its model. 
// The deposit keeps a count of how many model A cars and model B cars are completed. 
// We want TOT = 10 total cars (any combination of A and B).
//
// The deposit goroutine manages inventories of CA, CB, PA, PB in separate counters,
// ensuring no more than maxP tires (combined? Actually separate for A or B) and 
// maxC rims (again, separate for A or B) are stored. It uses a select statement with 
// "when(...)" to allow or block certain actions. For example, a rim or tire can only
// be delivered if there is space. A robot can only pick up a part if there is a
// matching part available in the deposit. 
//
// Once the deposit detects that the sum of cars built (model A + model B) equals TOT,
// it sets a flag (fine = true) so that subsequent attempts to deliver or pick up
// parts are answered with -1 (termination). After all conveyors and robots exit, 
// 'main' tells the deposit to terminate as well.
//
// -----------------------------------------------------------------------------------

package main

import (
    "fmt"
    "math/rand"
    "time"
)

// Limits on how many tires (pneumatici) and rims (cerchi) can be stored
const maxP = 3 // max tires in the deposit
const maxC = 3 // max rims in the deposit

// Types of parts:
//   0 => pneumatico A (PA)
//   1 => pneumatico B (PB)
//   2 => cerchio A (CA)
//   3 => cerchio B (CB)
const tipoPA = 0
const tipoPB = 1
const tipoCA = 2
const tipoCB = 3

// Robots: 0 => Robot model A, 1 => Robot model B
const RobotA = 0
const RobotB = 1

// TOT is the total number of cars (model A or B) we want to build
const TOT = 10

// Human-readable names
var tipoRobot = [2]string{"Modello A", "Modello B"}
var tipoNastro = [4]string{"pneumatico A", "pneumatico B", "cerchio A", "cerchio B"}

// Channels for termination and synchronization
var done = make(chan bool)
var terminaDeposito = make(chan bool)

// Channels for ROBOTS to pick up parts from the deposit
var prelievoPA = make(chan int, 100)
var prelievoPB = make(chan int, 100)
var prelievoCA = make(chan int, 100)
var prelievoCB = make(chan int, 100)

// Channels for CONVEYOR BELTS to deliver parts to the deposit
var consegnaPA = make(chan int, 100)
var consegnaPB = make(chan int, 100)
var consegnaCA = make(chan int, 100)
var consegnaCB = make(chan int, 100)

// Acknowledgment channels
var ack_robotA = make(chan int)   // ack for robot A
var ack_robotB = make(chan int)   // ack for robot B
var ack_nastroPA = make(chan int) // ack for conveyor delivering PA
var ack_nastroPB = make(chan int) // ack for conveyor delivering PB
var ack_nastroCA = make(chan int) // ack for conveyor delivering CA
var ack_nastroCB = make(chan int) // ack for conveyor delivering CB

// Helper function for conditional case in select
func when(b bool, c chan int) chan int {
    if !b {
        return nil
    }
    return c
}

// Robot goroutine: each robot builds cars of a specific model (A or B).
// For each car, the robot assembles 4 wheels, each wheel requires a rim + a tire.
//   RobotA => cerchio A (CA) + pneumatico A (PA)
//   RobotB => cerchio B (CB) + pneumatico B (PB)
//
// The robot loops indefinitely, and after building 4 wheels (one car), it logs it.
// If it receives a -1 from the deposit (ack_robotA or ack_robotB), it terminates.
func Robot(tipo int) {
    fmt.Printf("[Robot %s]: starting up!\n", tipoRobot[tipo])
    var assemblyCount, ackVal int
    var tt int

    for {
        // For each car, we build 4 wheels:
        for i := 0; i < 4; i++ {
            if tipo == RobotA {
                // 1) Pick up rim CA
                prelievoCA <- tipo
                ackVal = <-ack_robotA
                if ackVal == -1 {
                    fmt.Printf("[Robot %s]: terminating now!\n", tipoRobot[tipo])
                    done <- true
                    return
                }
                fmt.Printf("[Robot %s]: picked up rim CA\n", tipoRobot[tipo])
                tt = rand.Intn(2) + 1
                time.Sleep(time.Duration(tt) * time.Second) // mounting time for the rim

                // 2) Pick up tire PA
                prelievoPA <- tipo
                ackVal = <-ack_robotA
                if ackVal == -1 {
                    fmt.Printf("[Robot %s]: terminating now!\n", tipoRobot[tipo])
                    done <- true
                    return
                }
                fmt.Printf("[Robot %s]: picked up tire PA\n", tipoRobot[tipo])
                tt = rand.Intn(2) + 1
                time.Sleep(time.Duration(tt) * time.Second) // mounting time for the tire

            } else { // RobotB
                // 1) Pick up rim CB
                prelievoCB <- tipo
                ackVal = <-ack_robotB
                if ackVal == -1 {
                    fmt.Printf("[Robot %s]: terminating now!\n", tipoRobot[tipo])
                    done <- true
                    return
                }
                fmt.Printf("[Robot %s]: picked up rim CB\n", tipoRobot[tipo])
                tt = rand.Intn(2) + 1
                time.Sleep(time.Duration(tt) * time.Second) // mounting time for the rim

                // 2) Pick up tire PB
                prelievoPB <- tipo
                ackVal = <-ack_robotB
                if ackVal == -1 {
                    fmt.Printf("[Robot %s]: terminating now!\n", tipoRobot[tipo])
                    done <- true
                    return
                }
                fmt.Printf("[Robot %s]: picked up tire PB\n", tipoRobot[tipo])
                tt = rand.Intn(2) + 1
                time.Sleep(time.Duration(tt) * time.Second) // mounting time for the tire
            }
        }
        // After 4 wheels, we finished one car
        assemblyCount++
        fmt.Printf("[Robot %s]: completed car #%d\n", tipoRobot[tipo], assemblyCount)
    }
}

// Conveyor belt goroutine for delivering a particular type of part (PA, PB, CA, CB).
// It loops infinitely, sleeping a random time (1-2 seconds) each iteration to simulate
// transport time, then sends a piece to the deposit (e.g. consegnaPA <- 1) and waits
// for an acknowledgment. If it receives -1, it terminates.
func nastro(myType int) {
    var tt, ackVal, countDelivered int

    for {
        tt = rand.Intn(2) + 1
        time.Sleep(time.Duration(tt) * time.Second) // simulating belt movement

        switch myType {
        case tipoPA:
            consegnaPA <- 1
            ackVal = <-ack_nastroPA
            if ackVal == -1 {
                fmt.Printf("[conveyor %s]: terminating!\n", tipoNastro[myType])
                done <- true
                return
            }
            fmt.Printf("[conveyor %s]: delivered %s\n", tipoNastro[myType], tipoNastro[myType])

        case tipoPB:
            consegnaPB <- 1
            ackVal = <-ack_nastroPB
            if ackVal == -1 {
                fmt.Printf("[conveyor %s]: terminating!\n", tipoNastro[myType])
                done <- true
                return
            }
            fmt.Printf("[conveyor %s]: delivered %s\n", tipoNastro[myType], tipoNastro[myType])

        case tipoCA:
            consegnaCA <- 1
            ackVal = <-ack_nastroCA
            if ackVal == -1 {
                fmt.Printf("[conveyor %s]: terminating!\n", tipoNastro[myType])
                done <- true
                return
            }
            fmt.Printf("[conveyor %s]: delivered %s\n", tipoNastro[myType], tipoNastro[myType])

        case tipoCB:
            consegnaCB <- 1
            ackVal = <-ack_nastroCB
            if ackVal == -1 {
                fmt.Printf("[conveyor %s]: terminating!\n", tipoNastro[myType])
                done <- true
                return
            }
            fmt.Printf("[conveyor %s]: delivered %s\n", tipoNastro[myType], tipoNastro[myType])

        default:
            fmt.Printf("[conveyor %d]: invalid initialization, terminating!\n", myType)
            done <- true
            return
        }
        countDelivered++
    }
}

// deposito goroutine: central "warehouse" or "stock" that handles storing parts
// (up to maxP tires and up to maxC rims) and letting robots pick them up.
//
// Variables like numCA, numCB, numPA, numPB track how many of each part is currently
// in the deposit, while numCAMontati, numCBMontati, etc. track how many have been
// taken by the robots for assembly. 
//
// Once the deposit sees that the total number of assembled cars (model A + model B)
// equals TOT, it sets 'fine = true' and from that point on, any conveyor or robot
// request is answered with -1 in the ack channel, forcing them to terminate.
func deposito() {
    // Counters for parts that have been "used" by the robots for each model
    var numCAMontati int
    var numCBMontati int
    var numPAMontati int
    var numPBMontati int

    // Counters for how many complete cars have been built
    var numAMontati int
    var numBMontati int

    // Current amount of each part in stock
    var numCA int
    var numCB int
    var numPA int
    var numPB int

    // totP = total tires, totC = total rims in storage
    var totP int
    var totC int

    var fine bool = false // becomes true when TOT cars are built

    for {
        select {
        // 1) Receiving cerchio A (CA)
        case <-when(
            !fine && (totC < maxC && numCA < maxC-1) &&
            (numAMontati < numBMontati ||
                (numAMontati >= numBMontati && len(consegnaCB) == 0)),
            consegnaCA,
        ):
            numCA++
            totC++
            ack_nastroCA <- 1
            fmt.Printf("[deposit] added rim A: now CA=%d, CB=%d, total rims=%d\n", numCA, numCB, totC)

        // 2) Receiving cerchio B (CB)
        case <-when(
            !fine && (totC < maxC && numCB < maxC-1) &&
            (numAMontati >= numBMontati ||
                (numAMontati < numBMontati && len(consegnaCA) == 0)),
            consegnaCB,
        ):
            numCB++
            totC++
            ack_nastroCB <- 1
            fmt.Printf("[deposit] added rim B: now CA=%d, CB=%d, total rims=%d\n", numCA, numCB, totC)

        // 3) Receiving pneumatico A (PA)
        case <-when(
            !fine && (totP < maxP && numPA < maxP-1) &&
            (numAMontati < numBMontati ||
                (numAMontati >= numBMontati && len(consegnaPB) == 0)),
            consegnaPA,
        ):
            numPA++
            totP++
            ack_nastroPA <- 1
            fmt.Printf("[deposit] added tire A: now PA=%d, PB=%d, total tires=%d\n", numPA, numPB, totP)

        // 4) Receiving pneumatico B (PB)
        case <-when(
            !fine && (totP < maxP && numPB < maxP-1) &&
            (numAMontati >= numBMontati ||
                (numAMontati < numBMontati && len(consegnaPA) == 0)),
            consegnaPB,
        ):
            numPB++
            totP++
            ack_nastroPB <- 1
            fmt.Printf("[deposit] added tire B: now PA=%d, PB=%d, total tires=%d\n", numPA, numPB, totP)

        // 5) Robot A picking up a cerchio A (CA)
        case <-when(
            !fine && numCA > 0 &&
                (numAMontati < numBMontati ||
                    (numAMontati >= numBMontati && len(prelievoCB) == 0)),
            prelievoCA,
        ):
            numCA--
            totC--
            numCAMontati++
            ack_robotA <- 1
            fmt.Printf("[deposit] robot A took rim A: total rims now=%d\n", totC)

        // 6) Robot B picking up a cerchio B (CB)
        case <-when(
            !fine && numCB > 0 &&
                (numAMontati >= numBMontati ||
                    (numAMontati < numBMontati && len(prelievoCA) == 0)),
            prelievoCB,
        ):
            numCB--
            totC--
            numCBMontati++
            ack_robotB <- 1
            fmt.Printf("[deposit] robot B took rim B: total rims now=%d\n", totC)

        // 7) Robot A picking up a pneumatico A (PA)
        case <-when(
            !fine && numPA > 0 &&
                (numAMontati < numBMontati ||
                    (numAMontati >= numBMontati && len(prelievoPB) == 0)),
            prelievoPA,
        ):
            numPA--
            totP--
            numPAMontati++
            ack_robotA <- 1
            fmt.Printf("[deposit] robot A took tire A: total tires now=%d\n", totC)

        // 8) Robot B picking up a pneumatico B (PB)
        case <-when(
            !fine && numPB > 0 &&
                (numAMontati >= numBMontati ||
                    (numAMontati < numBMontati && len(prelievoPA) == 0)),
            prelievoPB,
        ):
            numPB--
            totP--
            numPBMontati++
            ack_robotB <- 1
            fmt.Printf("[deposit] robot B took tire B: total tires now=%d\n", totC)

        // 9) Once 'fine' is set, any incoming requests get an ack of -1 (termination)
        case <-when(fine, consegnaCA):
            ack_nastroCA <- -1
        case <-when(fine, consegnaCB):
            ack_nastroCB <- -1
        case <-when(fine, consegnaPA):
            ack_nastroPA <- -1
        case <-when(fine, consegnaPB):
            ack_nastroPB <- -1
        case <-when(fine, prelievoCA):
            ack_robotA <- -1
        case <-when(fine, prelievoCB):
            ack_robotB <- -1
        case <-when(fine, prelievoPA):
            ack_robotA <- -1
        case <-when(fine, prelievoPB):
            ack_robotB <- -1

        // 10) The main function eventually sends terminaDeposito here
        case <-terminaDeposito:
            fmt.Printf("[deposit] Terminating now.\n")
            done <- true
            return
        }

        // Each time we finish a wheel pickup for Robot A, we increment numCAMontati or numPAMontati. 
        // If numCAMontati == 4 and numPAMontati == 4, that means 4 rims A + 4 tires A have been used, 
        // so we've built 1 model A car.
        if numCAMontati == 4 && numPAMontati == 4 {
            numAMontati++
            numCAMontati = 0
            numPAMontati = 0
        }
        // Similarly for Robot B
        if numCBMontati == 4 && numPBMontati == 4 {
            numBMontati++
            numCBMontati = 0
            numPBMontati = 0
        }
        fmt.Printf("[deposit] Model A cars built=%d, Model B cars built=%d\n", numAMontati, numBMontati)

        // If total cars built = TOT, set fine = true to stop further production
        if numAMontati+numBMontati == TOT {
            fine = true
        }
    }
}

func main() {
    rand.Seed(time.Now().Unix())

    fmt.Printf("[main] Starting 4 conveyor belts and 2 robots.\n")

    // Start the deposit goroutine
    go deposito()

    // Create 4 conveyor belt goroutines, one for each part type: PA, PB, CA, CB
    for i := 0; i < 4; i++ {
        go nastro(i)
    }

    // Create 2 robot goroutines
    for i := 0; i < 2; i++ {
        go Robot(i)
    }

    // Wait for the 4 conveyor belts to finish
    for i := 0; i < 4; i++ {
        <-done
    }

    // Wait for the 2 robots to finish
    for i := 0; i < 2; i++ {
        <-done
    }

    // Signal the deposit to terminate
    terminaDeposito <- true
    <-done

    fmt.Printf("[main] APPLICATION FINISHED\n")
}
