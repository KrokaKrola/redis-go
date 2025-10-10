package store

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (m innerMap) xrange(key, start, end string) (Stream, error) {
	v, ok := m[key]

	if !ok {
		return Stream{}, nil
	}

	if v.isExpired() {
		return Stream{}, nil
	}

	stream, ok := v.value.(Stream)
	if !ok {
		return Stream{}, fmt.Errorf("MISSTYPE of the element in the underlying stream")
	}

	startId, ok := parseXrangeStreamId(start)
	if !ok {
		return Stream{}, fmt.Errorf("ERR invalid start value for XRANGE command")
	}
	startBound := storedStreamId{MsTime: startId.MsTime, Seq: startId.Seq}
	if startId.AutoSeq {
		startBound.Seq = 0
	}

	endId, ok := parseXrangeStreamId(end)
	if !ok {
		return Stream{}, fmt.Errorf("ERR invalid end value for XRANGE command")
	}
	endBound := storedStreamId{MsTime: endId.MsTime, Seq: endId.Seq}
	if endId.AutoSeq {
		endBound.Seq = math.MaxUint64
	}

	result := Stream{}

	for _, streamElement := range stream.Elements {
		if less(streamElement.Id, startBound) {
			continue
		}

		if greater(streamElement.Id, endBound) {
			break
		}

		result.Elements = append(result.Elements, streamElement)
	}

	return result, nil
}

func parseXrangeStreamId(id string) (streamId StreamIdSpec, ok bool) {
	before, after, found := strings.Cut(id, "-")

	msTime, err := strconv.ParseUint(before, 10, 64)
	if err != nil {
		return StreamIdSpec{}, false
	}

	if found {
		seq, err := strconv.ParseUint(after, 10, 64)

		if err != nil {
			return StreamIdSpec{}, false
		}

		return StreamIdSpec{
			MsTime: msTime,
			Seq:    seq,
		}, true
	}

	return StreamIdSpec{
		MsTime:  msTime,
		AutoSeq: true,
	}, true
}

func less(a, b storedStreamId) bool {
	if a.MsTime != b.MsTime {
		return a.MsTime < b.MsTime
	}

	return a.Seq < b.Seq
}

func greater(a, b storedStreamId) bool {
	return less(b, a)
}
