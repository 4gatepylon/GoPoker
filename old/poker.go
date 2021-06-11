package main

import "errors"
import "math/rand"

////////// CONSTS //////////

const MAX_PLAYERS := 6
const MAX_ROOMS := 4
const ROOMNAME_LENGTH := 5

const CommandID (
	// once inside a room you can always check, fold or bet
	// but this will only be accepted if it's your turn and otherwise
	// the server should send back an informational error message
	fold = iota
	check
	bet

	// once in a room you can list the players in that room
	// or you can also leave the room (but you must be in a room and you cannot
	// leave if you are "playing" which just means you are in a hand)
	listPlayers
	leaveRoom

	// before joining a room you are able to list rooms or join rooms
	// and the message sent back should be informational (i.e. how many
	// people are in that room since six is the maximum)... also if you join
	// a room and pass in nil as a name a new room will be created
	listRooms
	joinRoom
	
	// if you are an admin in a room you are allowed at any time to issue chips
	// to someone who is not currently playing (i.e. not in the hand: this means that
	// if you run out you will probably have to sit out a round)
	issueChips
)

////////// STRUCTS //////////

type Card struct {
	int num // 1 = A, 11 = J, 12 = Q, 13 = K
	byte suit // D = diamonds, S = spades, H = hearts, C = clubs
}

// A Player represents a single player who has a hand and some chips
// and is an admin if they created the room they are in (admins can
// issue new chips arbitrarily, but it is impossible to take chips away)
// (names must be unique)
type Player struct {
	name string
	hand [2]*Card
	chips int
	admin bool
	playing bool
}

// A room represents and individual game and keeps track of who is playing
// and the state of the overall game
type Room struct {
	id string
	players [MAX_PLAYERS]*Player
	middle [5]*Card
}

type Server struct {
	rooms [MAX_ROOMS]*Room
	addr2player map[string]*Player
}

// Commands are recieved from TCP connections
// and can contain metadata which is int | string int | string, where
// if it's string and the id was joinRoom it's to join that room by name
// if it's int it's bet then that's the ammount you are betting
// and if it's string int you are an admin and issuing int chips to a named player
type Command struct {
	id CommandID
	addr string
	metadata string
}

////////// FUNCS //////////

// TODO
func genRandName() string, error {
	name := [ROOMNAME_LENGTH]bytes{}
	minByte := 'a'
	maxByte := 'z'
	for i := 0; i < ROOMNAME_LENGTH; i++ {
		name[i] = 1
	}
	return string(name)
}