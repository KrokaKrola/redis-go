package resp

import "fmt"

type parserState string

const (
	parserInitialState    parserState = "initial"
	parserProcessingState parserState = "processing"
	parserDoneState       parserState = "done"
)

type parser struct {
	state  parserState
	input  []byte
	result *Resp
}

func NewParser(input []byte) *parser {
	return &parser{
		state: parserInitialState,
		input: input,
	}
}

func (p *parser) Parse() (*Resp, error) {
	idx := 0

	p.state = parserProcessingState

	for idx <= len(p.input) {
		firstByte := p.input[idx]

		switch firstByte {
		case byte('+'):
			n, err := p.processSimpleString(p.input[idx:])

			if err != nil {
				return nil, err
			}

			idx += n
		case byte('*'):
			n, err := p.processArray(p.input[idx:])

			if err != nil {
				return nil, err
			}

			idx += n
		case byte('$'):
			n, err := p.processBulkString(p.input[idx:])

			if err != nil {
				return nil, err
			}

			idx += n
		default:
			return nil, fmt.Errorf("unknown data type: %s", string(p.input))
		}
	}

	return p.result, nil
}

func (p *parser) processSimpleString(input []byte) (int, error) {
	bytesProcessed := 0

	return bytesProcessed, nil
}

func (p *parser) processArray(input []byte) (int, error) {
	bytesProcessed := 0

	return bytesProcessed, nil
}

func (p *parser) processBulkString(input []byte) (int, error) {
	bytesProcessed := 0

	return bytesProcessed, nil
}
