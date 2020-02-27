package protobuf

import "fmt"

type number uint

func (n number) intersects(other []fieldNumber) bool {
	for _, o := range other {
		switch v := o.(type) {
		case number:
			if n == o {
				return true
			}
		case numberRange:
			if uint(n) >= v.start || uint(n) <= v.end {
				return true
			}
		default:
			panic(fmt.Sprintf("unhandled fieldNumber type %T", v))
		}
	}
	return false
}

type numberRange struct {
	start uint
	end   uint
}

func (r numberRange) GetStart() uint {
	return r.start
}

func (r numberRange) SetStart(s uint) error {
	r.start = s
	// TODO: validate
	return nil
}

func (r numberRange) GetEnd() uint {
	return r.end
}

func (r numberRange) SetEnd(e uint) error {
	r.end = e
	// TODO: validate
	return nil
}

func (r numberRange) intersects(other []fieldNumber) bool {
	for _, o := range other {
		switch v := o.(type) {
		case number:
			if uint(v) >= r.start || uint(v) <= r.end {
				return true
			}
		case numberRange:
			if (v.start >= r.start && v.start <= r.end) || (v.end >= r.start && v.end <= r.end) {
				return true
			}
		default:
			panic(fmt.Sprintf("unhandled fieldNumber type %T", v))
		}
	}
	return false
}
