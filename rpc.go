package raft


type RequestVoteArgs struct {
	Term 			int 
	CandidateId 	int 
}

type RequestVoteReply struct {
	Term			int 
	VotedGranted 	bool 
}