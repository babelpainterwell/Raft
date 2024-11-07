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
	heartbeatInterval 		int 
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

// If a follower receives no communication over a period of time, 
// it starts a new election and become the candidate 
func (r *Raft) StartNewElection() {
	// update the term 

	// state transition 

	// initiate RequestVote RPC

	// state transition after RPC, becomes a leader or follower

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
// HeartBeat 
//
func (r *Raft) sendHeartbeats() {
	// send empty AppendEntries to followers periotically based on the heartbeat interval 

	for {
		time.Sleep(time.Duration(r.heartbeatInterval * int(time.Millisecond)))

		appendEntriesArgs := AppendEntriesArgs{
			Term: 			r.currentTerm,
			LeaderId: 		r.id,
			Entries: 		[]LogEntry{}, // empty slice for a heatbeat
		}

		appendEntriesReply := AppendEntriesReply{}


		r.AppendEntries(&appendEntriesArgs, &appendEntriesReply)

		// downgrade a leader if the received term is higher than the leader's current term 
		if appendEntriesReply.Success != true {
			r.LeaderToFollower()
			break
		}

	}
}


//
// RPC HANDLERS
//
func (r *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// deny the request if either args term is lower than current term 
	// or r has voted for another candidate in the same term
	if r.votedFor != -1 || args.Term < r.currentTerm {
		reply.Term = r.currentTerm
		reply.VotedGranted = false 
		return nil 
	}

	// Reset VotedFor if a new term arrives 
	if args.Term > r.currentTerm {
		r.currentTerm = args.Term
		r.votedFor = -1 
	}

	if r.votedFor == -1 {
		reply.Term = r.currentTerm
		reply.VotedGranted = true 
		// set votedFor
		r.votedFor = args.CandidateId
	} else {
		reply.Term = r.currentTerm
		reply.VotedGranted = false 
	}

	return nil

}

func (r *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if args.Term < r.currentTerm {
		reply.Term = r.currentTerm // tell the leader the follower has a higher term 
		reply.Success = false 
		return nil 
	}

	// Get the entries 
	// entries := args.Entries
	reply.Term = r.currentTerm
	reply.Success = true 

	// remember to reset the timer for the follower 
	r.ResetElectionTimer()


	return nil 

}
