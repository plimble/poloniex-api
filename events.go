package poloniex

import "github.com/chuckpreslar/emission"

func (p *Poloniex) On(event interface{}, listener interface{}) *emission.Emitter {
	return p.emitter.On(event, listener)
}

func (p *Poloniex) Emit(event interface{}, arguments ...interface{}) *emission.Emitter {
	return p.emitter.Emit(event, arguments...)
}

func (p *Poloniex) Off(event interface{}, listener interface{}) *emission.Emitter {
	return p.emitter.Off(event, listener)
}
