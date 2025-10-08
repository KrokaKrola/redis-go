package store

import (
	"errors"
	"fmt"
	"math"
	"time"
)

func (m innerMap) xadd(key string, streamId StreamIdSpec, fields [][]string) (newEntryId string, err error) {
	sv, ok := m[key]

	fields = cloneStreamFields(fields)

	seqNumber := streamId.Seq
	msTime := streamId.MsTime

	if !ok || sv.isExpired() {
		msTime, seqNumber = getNewStreamId(streamId)

		m[key] = newStoreValue(Stream{
			Elements:           []streamElement{{id: storedStreamId{msTime, seqNumber}, fields: fields}},
			LtsInsertedIdParts: storedStreamId{msTime, seqNumber},
		}, getPossibleEndTime())

		return fmt.Sprintf("%d-%d", msTime, seqNumber), nil
	}

	stream, okStream := sv.value.(Stream)
	if !okStream {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if isValidStreamId := validateStreamIdParts(streamId, stream.LtsInsertedIdParts); !isValidStreamId {
		return "", fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	if streamId.AutoSeq {
		msTime = streamId.MsTime

		if msTime != stream.LtsInsertedIdParts.msTime {
			seqNumber = 0
		} else {
			if stream.LtsInsertedIdParts.seq == math.MaxUint64 {
				return "", fmt.Errorf("ERR sequence overflow for XADD command")
			}

			seqNumber = stream.LtsInsertedIdParts.seq + 1
		}
	} else if streamId.AutoFull {
		msTime = max(uint64(time.Now().UnixMilli()), stream.LtsInsertedIdParts.msTime)

		if stream.LtsInsertedIdParts.msTime == msTime {
			if stream.LtsInsertedIdParts.seq == math.MaxUint64 {
				return "", fmt.Errorf("ERR sequence overflow for XADD command")
			}

			seqNumber = stream.LtsInsertedIdParts.seq + 1
		} else {
			seqNumber = 0
		}
	}

	newElements := append(stream.Elements, streamElement{id: storedStreamId{msTime, seqNumber}, fields: fields})

	m[key] = newStoreValue(Stream{
		Elements:           newElements,
		LtsInsertedIdParts: storedStreamId{msTime, seqNumber},
	}, sv.expiryTime)

	return fmt.Sprintf("%d-%d", msTime, seqNumber), nil
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

	if streamId.MsTime < ltsInsertedIdParts.msTime {
		return false
	}

	if !streamId.AutoSeq && streamId.MsTime == ltsInsertedIdParts.msTime {
		if streamId.Seq <= ltsInsertedIdParts.seq {
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
