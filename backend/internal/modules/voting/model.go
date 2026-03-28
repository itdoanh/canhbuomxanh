package voting

type castVoteRequest struct {
	TeacherID   uint64 `json:"teacherId"`
	RawScore    uint8  `json:"rawScore"`
	VoteMode    string `json:"voteMode"`
	SemesterKey string `json:"semesterKey"`
}

type releaseVotesRequest struct {
	SemesterKey string `json:"semesterKey"`
}

type forcedVoteAlertRequest struct {
	TeacherID uint64 `json:"teacherId"`
	Detail    string `json:"detail"`
}

type recomputeBadgesRequest struct {
	TeacherID uint64 `json:"teacherId"`
}
