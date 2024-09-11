package eco

type Guild struct {
	TempChannels map[string]bool
}

func NewGuild() *Guild {
	return &Guild{
		TempChannels: map[string]bool{},
	}
}
