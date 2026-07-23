package models

type AdminStatistics struct {
	TotalUsers            int
	TotalVoters           int
	TotalAdmins           int
	VerifiedVoters        int
	PendingVerifications  int
	RejectedVerifications int
	ApprovedVerifications int
	TotalVotesCast        int
}
