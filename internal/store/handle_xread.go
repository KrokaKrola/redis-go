package store

import (
	"fmt"
	"strconv"
	"strings"
)

func (m innerMap) xread(keys []string, id string) ([]Stream, error) {
	streams := []Stream{}

	streamIdSpec, ok := parseXreadStreamId(id)
	if !ok {
		return streams, fmt.Errorf("ERR invalid stream id")
	}

	for _, key := range keys {
		storeStreamRawValue, ok := m[key]

		if !ok {
			continue
		}

		if storeStreamRawValue.isExpired() {
			m.delete(key)
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
