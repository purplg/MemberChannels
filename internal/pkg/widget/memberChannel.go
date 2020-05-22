package widget

type memberChannel struct {
	ID  string
	ownerID    string
	visitorIDs []string
}

// Returns true if there was an available owner to take over
func (uc *memberChannel) PopToOwner() bool {
	if len(uc.visitorIDs) == 0 {
		return false
	}
	uc.ownerID = uc.visitorIDs[0]
	uc.visitorIDs = uc.visitorIDs[1:]
	return true
}

func (uc *memberChannel) AddVisitor(userID string) {
	uc.visitorIDs = append(uc.visitorIDs, userID)
}

func (uc *memberChannel) RemoveVisitor(userID string) {
	for i := 0; i < len(uc.visitorIDs); i++ {
		if uc.visitorIDs[i] == userID {
			uc.visitorIDs = append(uc.visitorIDs[:i], uc.visitorIDs[i+1:]...)
			return
		}
	}
}
