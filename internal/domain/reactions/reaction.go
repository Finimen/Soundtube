package reactions

type Reaction struct {
	targetID int
	reaction string
}

func (r *Reaction) GetTargetID() int { return r.targetID }
func (r *Reaction) GetType() string  { return r.reaction }

func NewReaction(targetID int, reaction string) *Reaction {
	return &Reaction{
		targetID: targetID,
		reaction: reaction,
	}
}
