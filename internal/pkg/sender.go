package pkg

type Network interface {
	SendPkg(connectionId string, data []byte)
}

type SenderOpts struct {
	Network    Network
	Recipients []string
}

type Sender struct {
	network    Network
	recipients []string
}

func NewSender(opts SenderOpts) *Sender {
	return &Sender{
		network:    opts.Network,
		recipients: opts.Recipients,
	}
}

func (s *Sender) ToAll() *Protocol {
	return &Protocol{
		n:          s.network,
		recipients: s.recipients,
	}
}

func (s *Sender) To(recipients ...string) Protocol {
	return Protocol{
		n:          s.network,
		recipients: recipients,
	}
}
