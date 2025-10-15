package store

import (
	"errors"
	"fmt"
	"math"
	"time"
)

func (m innerMap) xadd(key string, streamId StreamIdSpec, fields [][]string) (streamElement StreamElement, err error) {
	sv, ok := m[key]

	fields = cloneStreamFields(fields)

	seqNumber := streamId.Seq
	msTime := streamId.MsTime

	if !ok || sv.isExpired() {
		msTime, seqNumber = getNewStreamId(streamId)

		sEl := StreamElement{
			Id:     storedStreamId{msTime, seqNumber},
			Fields: fields,
		}

		m[key] = newStoreValue(Stream{
			Elements:           []StreamElement{sEl},
			LtsInsertedIdParts: storedStreamId{msTime, seqNumber},
		}, getPossibleEndTime())

		return sEl, nil
	}

	stream, okStream := sv.value.(Stream)
	if !okStream {
		return StreamElement{}, errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if isValidStreamId := validateStreamIdParts(streamId, stream.LtsInsertedIdParts); !isValidStreamId {
		return StreamElement{}, fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if streamId.AutoSeq {
		msTime = streamId.MsTime

		if msTime != stream.LtsInsertedIdParts.MsTime {
			seqNumber = 0
		} else {
			if stream.LtsInsertedIdParts.Seq == math.MaxUint64 {
				return StreamElement{}, fmt.Errorf("ERR sequence overflow for XADD command")
			}

			seqNumber = stream.LtsInsertedIdParts.Seq + 1
		}
	} else if streamId.AutoFull {
		msTime = max(uint64(time.Now().UnixMilli()), stream.LtsInsertedIdParts.MsTime)

		if stream.LtsInsertedIdParts.MsTime == msTime {
			if stream.LtsInsertedIdParts.Seq == math.MaxUint64 {
				return StreamElement{}, fmt.Errorf("ERR sequence overflow for XADD command")
			}

			seqNumber = stream.LtsInsertedIdParts.Seq + 1
		} else {
			seqNumber = 0
		}
	}

	sEl := StreamElement{Id: storedStreamId{msTime, seqNumber}, Fields: fields}

	newElements := append(stream.Elements, sEl)

	m[key] = newStoreValue(Stream{
		Elements:           newElements,
		LtsInsertedIdParts: storedStreamId{msTime, seqNumber},
	}, sv.expiryTime)

	return sEl, nil
}

func cloneStreamFields(fields [][]string) [][]string {
	if len(fields) == 0 {
		return nil
	}

	out := make([][]string, len(fields))
	for i, pair := range fields {
		copied := make([]string, len(pair))
		copy(copied, pair)
		out[i] = copied
	}

	return out
}

func validateStreamIdParts(streamId StreamIdSpec, ltsInsertedIdParts storedStreamId) bool {
	if streamId.AutoFull {
		return true
	}

	if streamId.MsTime < ltsInsertedIdParts.MsTime {
		return false
	}

	if !streamId.AutoSeq && streamId.MsTime == ltsInsertedIdParts.MsTime {
		if streamId.Seq <= ltsInsertedIdParts.Seq {
			return false
		}
	}

	return true
}

func getNewStreamId(streamId StreamIdSpec) (uint64, uint64) {
	if streamId.AutoFull {
		return uint64(time.Now().UnixMilli()), 0
	}

	if streamId.AutoSeq {
		if streamId.MsTime == 0 {
			return streamId.MsTime, 1
		}

		return streamId.MsTime, 0
	}

	return streamId.MsTime, streamId.Seq
}
