package utils

import "fmt"

type Segment struct {
	Start, End uint16
}

func (segment Segment) String() string {
	return fmt.Sprintf("[%x - %x]", segment.Start, segment.End)
}

func (segment Segment) Size() uint32 {
	return uint32(segment.End-segment.Start) + 1
}

func (segment Segment) IsIn(value uint16) bool {
	return segment.Start <= value && value <= segment.End
}

func (segment Segment) Intersect(otherSegment *Segment) bool {
	return !(otherSegment.End <= segment.Start || segment.End <= otherSegment.Start)
}

func (segment Segment) OffsetIn(value uint16) uint16 {
	offset := value % uint16(segment.Size())
	return segment.Start + offset
}

func NewSegment(start, end uint16) *Segment {
	return &Segment{
		Start: start,
		End:   end,
	}
}
