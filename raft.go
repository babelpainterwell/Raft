package raft

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type State int 

const (
	Follower State = iota 
	Candidate 
	Leader 
)

type Raft struct {
	mu 						sync.Mutex
	id 						int
	currentTerm 			int 
	state 					State 
	votedFor 				int 			// the candidate who receives the vote, -1 for nil
	log 					[]LogEntry
	electionTimer 			*time.Timer
	heartbeatInterval 		time.Duration
}

func GenerateRandom(lowerBound int64, upperBound int64) int64 {
	// to ensure different random values for each server 
	// rand.Seed(time.Now().UnixNano()) --> deprecared 
	// create a random number generator 
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rng.Int63n(upperBound - lowerBound) + lowerBound

}

// Set up a new Raft istance with a random election timeout between 150 - 300 ms
func NewRaft(id int) *Raft {

	randomNum := GenerateRandom(150, 300)
	minTimeout := 150 * time.Millisecond
	electionTimeout := minTimeout + time.Duration(randomNum)

	return &Raft{
		id: id,
		currentTerm: 0, 
		state: Follower,
		electionTimer: time.NewTimer(electionTimeout),
		heartbeatInterval: 50 * time.Millisecond,
	}
}

// reset a randomized timer in case of term resets, failures, or state changes 
func (node *Raft) ResetElectionTimer() {
	randomNum := GenerateRandom(150, 300)
	minTimeout := 150 * time.Millisecond
	electionTimeout := minTimeout + time.Duration(randomNum)
	node.electionTimer.Reset(electionTimeout)
	fmt.Printf("Server %d: Reset election timer to %v\n", node.id, electionTimeout)
}

// Methods for state transtions (ONLY), excluding operations such as reset timer 
func (node *Raft) FollowerToCandidate(){
	node.state = Candidate

	// increase its term for a new election
	node.currentTerm += 1
	
	fmt.Printf("Server %d becomes a candidate\n", node.id)
}

func (node *Raft) CandidateToLeader(){
	node.state = Leader 
	fmt.Printf("Server %d becomes a leader\n", node.id)
}

func (node *Raft) LeaderToFollower(){
	node.state = Follower

	fmt.Printf("Server %d becomes a follower\n", node.id)
}

func (node *Raft) CandidateToFollower() {
	node.state = Follower

	// reset votedFor 
	node.votedFor = -1

	fmt.Printf("Server %d becomes a follower\n", node.id)
}


//
// RPC HANDLERS
//

