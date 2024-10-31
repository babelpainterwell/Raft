# Raft

## Q & A

1. How is message considered as comitted?

- An entry is considered committed if it is safe for that entry to be applied to state machines.
- Logs that have been appended to a local server are not necessarily considered committed until they are confirmed by a majority of servers in the cluster.

2. Any timeout could lead to leader re-election?

- If a follower receives no communication over a period of time called the election timeout, then it assumes there is no vi- able leader and begins an election to choose a new leader.

3. When a Server Receives a RequestVote RPC, Must It Vote Yes?

- Not necessarily. a.The requesting candidate's term must be at least as high as the receiver's current term. b. The candidate’s log must be at least as up-to-date as the receiver’s. If these conditions are met and the server hasn’t already voted in the current term, it will vote for the candidate.
- If it has already voted for iteself, it cannot vote for another candidate in that term.

4. A mechnism for a candidate to become a follower after split votes, otherwise indefinite split votes?

- Yes, if a candidate sees a higher term in a RequestVote or AppendEntries message from another server, it immediately steps down (from candidate or leader) to the follower state, recognizing the other server’s term as more up-to-date.
