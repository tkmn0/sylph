package engine

type MessageParcer struct {
	buffer []byte
}

func NewMessageParcer() *MessageParcer {
	return &MessageParcer{}
}

func (p *MessageParcer) Parce(buff []byte) (MessageType, []byte) {
	switch buff[0] {
	case uint8(MessageTypeHeartBeat):
		return MessageTypeHeartBeat, nil
	case uint8(MessageTypeBody):
		if p.buffer != nil {
			data := append(p.buffer, buff[1:]...)
			p.buffer = nil
			return MessageTypeBody, data
		} else {
			return MessageTypeBody, buff[1:]
		}
	case uint8(MessageTypeChunk):
		if p.buffer == nil {
			p.buffer = []byte{}
		}
		p.buffer = append(p.buffer, buff[1:]...)
		return MessageTypeChunk, nil
	case uint8(MessageTypeInitialize):
		return MessageTypeInitialize, buff[1:]
	}
	return MessageTypeUnknown, nil
}
