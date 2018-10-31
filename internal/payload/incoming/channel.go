package incoming

type (
	Channels struct{}

	ChannelJoin struct {
		ChannelID string `json:"id"`
	}

	ChannelPart struct {
		ChannelID string `json:"id"`
	}

	ChannelViewRecord struct {
		ChannelID     uint64 `json:"channelID,string,omitempty"`
		LastMessageID uint64 `json:"lastMessageID,string,omitempty"`
	}

	ChannelCreate struct {
		Name  *string `json:"name"`
		Topic *string `json:"topic"`
		Type  *string `json:"type"`
	}

	ChannelUpdate struct {
		ID    string  `json:"id"`
		Name  *string `json:"name"`
		Topic *string `json:"topic"`
		Type  *string `json:"type"`
	}

	// ChannelActivity is sent from the client when there is an activity on the channel...
	ChannelActivity struct {
		ChannelID uint64 `json:"ID,string"`
		Kind      string `json:"kind,omitempty"`
	}
)
