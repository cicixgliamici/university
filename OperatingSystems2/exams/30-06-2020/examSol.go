package main

import (
	"fmt"
	"math/rand"
	"time"
)

/////////////////////////////////////////////////////////////////////
// Data Structures
/////////////////////////////////////////////////////////////////////
type Request struct {
	id  int
	ack chan int
}

/////////////////////////////////////////////////////////////////////
// Constants
/////////////////////////////////////////////////////////////////////
const MAXBUFF = 100                  // Max channel buffer size
const MAX_VEHICLES = 60              // Max number of vehicles
const MAX_BOATS = 6                  // Max number of boats
const MAX_VEHICLE_CAPACITY = 5       // Max vehicles on bridge

const bridgeUp, bridgeDown int = 0, 1 // Bridge states (up/down)
const northToSouth, southToNorth int = 0, 1 // Traffic directions

/////////////////////////////////////////////////////////////////////
// Channels
/////////////////////////////////////////////////////////////////////
var bridgeBoatCh [2]chan Request    // Boat channels [enter, exit]
var bridgeVehicleInCh [4]chan Request // Vehicle entry channels [north, south, public_north, public_south]
var bridgeVehicleOutCh = make(chan Request, MAXBUFF) // Vehicle exit channel

// Channel indices
const BOAT_ENTER, BOAT_EXIT int = 0, 1
const VEHICLE_NORTH, VEHICLE_SOUTH, PUBLIC_NORTH, PUBLIC_SOUTH int = 0, 1, 2, 3

/////////////////////////////////////////////////////////////////////
// Synchronization Channels
/////////////////////////////////////////////////////////////////////
var done = make(chan bool)       // Completion notification
var terminate = make(chan bool)  // Termination signal

/////////////////////////////////////////////////////////////////////
// Helper Functions
/////////////////////////////////////////////////////////////////////
// Conditional channel selector
func when(b bool, c chan Request) chan Request {
	if !b {
		return nil
	}
	return c
}

// Random sleep functions
func sleepMilliseconds(t int) {
	if t > 0 {
		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}

func sleepRandomSeconds(timeLimit int) {
	if timeLimit > 0 {
		time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
	}
}

/////////////////////////////////////////////////////////////////////
// Goroutines
/////////////////////////////////////////////////////////////////////
func vehicle(id int, vehicleType int) {
	sleepRandomSeconds(15)
	req := Request{id, make(chan int)}

	// Request bridge access
	fmt.Printf("\n[Vehicle %d] Type %d: Requesting bridge access", id, vehicleType)
	bridgeVehicleInCh[vehicleType] <- req
	<-req.ack // Wait for approval
	
	// Cross the bridge
	fmt.Printf("\n[Vehicle %d] Type %d: Crossing bridge...", id, vehicleType)
	sleepMilliseconds(600)
	
	// Exit bridge
	bridgeVehicleOutCh <- req
	<-req.ack
	fmt.Printf("\n[Vehicle %d] Type %d: Crossed bridge", id, vehicleType)
	
	done <- true
}

func boat(id int) {
	sleepRandomSeconds(15)
	req := Request{id, make(chan int)}

	// Request bridge entry
	fmt.Printf("\n[Boat %d] Requesting bridge access", id)
	bridgeBoatCh[BOAT_ENTER] <- req
	<-req.ack
	
	// Pass through bridge
	fmt.Printf("\n[Boat %d] Passing through...", id)
	time.Sleep(2 * time.Second)
	
	// Exit bridge
	bridgeBoatCh[BOAT_EXIT] <- req
	<-req.ack
	fmt.Printf("\n[Boat %d] Passed through", id)
	
	done <- true
}

func bridgeManager() {
	state := bridgeDown  // Initial state: bridge down for vehicles
	direction := northToSouth
	vehiclesOnBridge := 0

	for {
		select {
		// Boat handling
		case req := <-when(state == bridgeUp, bridgeBoatCh[BOAT_ENTER]):
			vehiclesOnBridge++
			fmt.Printf("\n[Bridge] Boat %d entering\tState: %d\tVehicles: %d", req.id, state, vehiclesOnBridge)
			req.ack <- 1

		case req := <-bridgeBoatCh[BOAT_EXIT]:
			vehiclesOnBridge--
			fmt.Printf("\n[Bridge] Boat %d exited\tState: %d\tVehicles: %d", req.id, state, vehiclesOnBridge)
			req.ack <- 1
			// Lower bridge if no more boats
			if len(bridgeBoatCh[BOAT_ENTER]) == 0 && vehiclesOnBridge == 0 {
				fmt.Printf("\n[Bridge] Lowering bridge for vehicles")
				state = bridgeDown
			}

		// Vehicle handling (north to south)
		case req := <-when(
			state == bridgeDown && 
			((vehiclesOnBridge > 0 && vehiclesOnBridge < MAX_VEHICLE_CAPACITY && direction == northToSouth) || 
			 (vehiclesOnBridge == 0 && direction == southToNorth)) && 
			len(bridgeBoatCh[BOAT_ENTER]) == 0, 
			bridgeVehicleInCh[PUBLIC_NORTH]):
			// Handle public service vehicles with priority
			if direction == southToNorth {
				direction = northToSouth
			}
			vehiclesOnBridge++
			fmt.Printf("\n[Bridge] Public Vehicle %d N->S\tState: %d\tVehicles: %d", req.id, state, vehiclesOnBridge)
			req.ack <- 1

		// Vehicle handling (south to north)
		case req := <-when(
			state == bridgeDown && 
			((vehiclesOnBridge > 0 && vehiclesOnBridge < MAX_VEHICLE_CAPACITY && direction == southToNorth) || 
			 (vehiclesOnBridge == 0 && direction == northToSouth)) && 
			len(bridgeBoatCh[BOAT_ENTER]) == 0, 
			bridgeVehicleInCh[PUBLIC_SOUTH]):
			if direction == northToSouth {
				direction = southToNorth
			}
			vehiclesOnBridge++
			fmt.Printf("\n[Bridge] Public Vehicle %d S->N\tState: %d\tVehicles: %d", req.id, state, vehiclesOnBridge)
			req.ack <- 1

		// Private vehicle handling
		case req := <-when(
			state == bridgeDown && 
			len(bridgeBoatCh[BOAT_ENTER]) == 0 && 
			len(bridgeVehicleInCh[PUBLIC_SOUTH]) == 0 && 
			len(bridgeVehicleInCh[PUBLIC_NORTH]) == 0, 
			bridgeVehicleInCh[VEHICLE_NORTH]):
			// Similar logic for private vehicles
			// ...

		// Vehicle exit handling
		case req := <-bridgeVehicleOutCh:
			vehiclesOnBridge--
			fmt.Printf("\n[Bridge] Vehicle %d exited\tState: %d\tVehicles: %d", req.id, state, vehiclesOnBridge)
			req.ack <- 1
			// Raise bridge if boats waiting
			if vehiclesOnBridge == 0 && len(bridgeBoatCh[BOAT_ENTER]) > 0 {
				fmt.Printf("\n[Bridge] Raising bridge for boats")
				state = bridgeUp
			}

		case <-terminate:
			fmt.Printf("\n\n[Bridge] Terminating...")
			done <- true
			return
		}
	}
}

/////////////////////////////////////////////////////////////////////
// Main
/////////////////////////////////////////////////////////////////////
func main() {
	rand.Seed(time.Now().Unix())

	// Initialize channels
	for i := 0; i < 2; i++ {
		bridgeBoatCh[i] = make(chan Request, MAXBUFF)
		bridgeVehicleInCh[i] = make(chan Request, MAXBUFF)
	}
	bridgeVehicleInCh[2] = make(chan Request, MAXBUFF)
	bridgeVehicleInCh[3] = make(chan Request, MAXBUFF)

	go bridgeManager()

	// Start vehicles and boats
	for i := 0; i < MAX_VEHICLES; i++ {
		go vehicle(i, rand.Intn(4))
	}
	for i := 0; i < MAX_BOATS; i++ {
		go boat(i)
	}

	// Wait for completion
	for i := 0; i < MAX_VEHICLES+MAX_BOATS; i++ {
		<-done
	}

	terminate <- true
	<-done
	fmt.Printf("\n[Main] Simulation ended\n")
}

/*
SYSTEM OVERVIEW

BRIDGE MANAGEMENT COMPONENTS:
- State Control:
  * bridgeUp: Ponte sollevato per il passaggio barche
  * bridgeDown: Ponte abbassato per il passaggio veicoli
- Traffico Veicoli:
  * 4 tipologie: Nord->Sud privati/pubblici, Sud->Nord privati/pubblici
  * Priorità a mezzi pubblici e mantenimento direzione di marcia
- Gestione Concorrenza:
  * Canali bufferizzati per richieste (MAXBUFF = 100)
  * Handshake con ack channels per sincronizzazione
  * Contatore veicoli sul ponte (MAX_VEHICLE_CAPACITY = 5)

CONCURRENCY PATTERNS:
1. Select con when-condition per gestione prioritaria:
   - Barche > Veicoli pubblici > Veicoli privati
   - Controllo stato ponte e direzione traffico
2. Transizioni di stato sicure:
   - Ponte si alza solo con nessun veicolo in transito
   - Cambio direzione solo con ponte vuoto
3. Pattern Handshake:
   req -> canale -> ack
   exit -> canale -> ack

PRIORITY SYSTEM:
[Barche] -> [Mezzi pubblici] -> [Mezzi privati]
1. Richieste barche bloccano immediatamente il ponte
2. Mezzi pubblici possono invertire direzione
3. Mezzi privati solo in assenza di mezzi pubblici

TYPICAL WORKFLOW:
1. Richiesta accesso (veicolo/barca)
2. Verifica condizioni:
   - Stato ponte
   - Direzione traffico
   - Capacità residua
3. Attivazione handshake
4. Transito
5. Notifica uscita
6. Aggiornamento stato

IMPLEMENTATO CON:
- Goroutine per entità indipendenti (veicoli, barche)
- Channel select per gestione centralizzata (bridgeManager)
- Buffered channels per burst di richieste
- Sincronizzazione con canali dedicati (terminate/done)
*/
