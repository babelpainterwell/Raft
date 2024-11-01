package raft


type RequestVoteArgs struct {
	Term 			int 
	CandidateId 	int 
}

type RequestVoteReply struct {
	Term			int 
	VotedGranted 	bool 
}

type AppendEntriesArgs struct {

}

type AppendEntriesReply struct {
	
}