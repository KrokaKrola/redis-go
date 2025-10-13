package store

import (
	"fmt"
	"strconv"
	"strings"
)

func (m innerMap) xread(keys [][]string) ([]Stream, error) {
	streams := []Stream{}

	for _, key := range keys {
		id := key[1]
		streamIdSpec, ok := parseXreadStreamId(id)

		if !ok {
			return streams, fmt.Errorf("ERR invalid stream id")
		}

		keyValue := key[0]
		storeStreamRawValue, ok := m[keyValue]

		if !ok {
			streams = append(streams, Stream{})
			continue
		}

		if storeStreamRawValue.isExpired() {
			m.delete(keyValue)
			streams = append(streams, Stream{})
			continue
		}

		storeStream, ok := storeStreamRawValue.value.(Stream)
		if !ok {
			return streams, fmt.Errorf("MISSTYPE of the element in the underlying stream")
		}

		elements := []StreamElement{}

		for _, streamElement := range storeStream.Elements {
			if streamElement.Id.MsTime < streamIdSpec.MsTime {
				continue
			}

			if streamElement.Id.MsTime == streamIdSpec.MsTime {
				if streamElement.Id.Seq <= streamIdSpec.Seq {
					continue
				}
			}
			elements = append(elements, streamElement)
		}

		streams = append(streams, Stream{Elements: elements})
	}

	return streams, nil
}

func parseXreadStreamId(id string) (spec StreamIdSpec, ok bool) {
	before, after, found := strings.Cut(id, "-")

	if !found {
		return StreamIdSpec{}, false
	}

	msTime, err := strconv.ParseUint(before, 10, 64)
	if err != nil {
		return StreamIdSpec{}, false
	}

	seq, err := strconv.ParseUint(after, 10, 64)

	if err != nil {
		return StreamIdSpec{}, false
	}

	return StreamIdSpec{
		MsTime: msTime,
		Seq:    seq,
	}, true
}
