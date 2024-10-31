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
	mu sync.Mutex
	id int
	term int 
	state State 
	votes int 
	electionTimer *time.Timer
	heartbeatInterval time.Duration
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
		term: 0, 
		state: Follower,
		votes: 0,
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
// without any RPC communication. eg. new candidate sends RequestVote, new leader sends empty AppendEntry RPC (heartbeat)
func (node *Raft) FollowerToCandidate(){
	node.state = Candidate

	// increase its term 
	node.term += 1

	// candidate votes for itself first 
	node.votes = 1

	// restart the election timeout and it waits the timeout to eplase before starting a new election
	node.ResetElectionTimer()

	fmt.Printf("Server %d becomes a candidate\n", node.id)
}

func (node *Raft) CandidateToLeader(){
	node.state = Leader 
	fmt.Printf("Server %d becomes a leader\n", node.id)
}

// once it receives majority votes, it becomes the leader and starts sending heartbeats
func (node *Raft) LeaderToFollower(){
	node.state = Follower
	// reset the votes
	node.votes = 0 

	fmt.Printf("Server %d becomes a follower\n", node.id)
}

func (node *Raft) CandidateToFollower() {
	node.state = Follower

	fmt.Printf("Server %d becomes a follower\n", node.id)
}

