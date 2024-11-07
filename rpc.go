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
	Term 			int 
	LeaderId 		int  // for client redirection use as the client might contact any server in the cluster
	Entries 		[]LogEntry
}

type AppendEntriesReply struct {
	Term 			int  // for the candidate to update itself 
	Success 		bool
}