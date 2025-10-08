package reactions

type Reaction struct {
	targetID      int
	totalLikes    int
	totalDislikes int
}

func (r *Reaction) GetTargetID() int      { return r.targetID }
func (r *Reaction) GetTotalLikes() int    { return r.totalLikes }
func (r *Reaction) GetTotalDislikes() int { return r.totalDislikes }

func NewReaction(targetID int) *Reaction {
	return &Reaction{
		targetID:      targetID,
		totalLikes:    0,
		totalDislikes: 0,
	}
}

func RestoreReactionFromStorage(targetID, totalLikes, totalDislikes int) *Reaction {
	return &Reaction{
		targetID:      targetID,
		totalLikes:    totalLikes,
		totalDislikes: totalDislikes,
	}
}
