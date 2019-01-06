package cmd

import (
	"testing"
	"time"
)

func TestCheckConflictSchedule(t *testing.T) {
	s10_11 := buildSchedule(10, 11)
	s11_12 := buildSchedule(11, 12)
	s10_13 := buildSchedule(10, 13)
	s12_14 := buildSchedule(12, 14)
	checkConflict(false, s10_11, s11_12, t)
	checkConflict(false, s11_12, s10_11, t)
	checkConflict(false, s10_11, s12_14, t)
	checkConflict(true, s10_11, s10_11, t)
	checkConflict(true, s11_12, s10_13, t)
	checkConflict(true, s10_13, s12_14, t)
}

func checkConflict(expect bool, src Schedule, tar Schedule, t *testing.T) {
	var msg string
	if expect {
		msg = "this case should conflict"
	} else {
		msg = "this case should not conflict"
	}
	if CheckConflictSchedule(src, tar) != expect {
		t.Errorf("%v ====> %v %v", msg, src, tar)
	}
}

func TestPreviousSchedule(t *testing.T) {
	s10_11 := buildSchedule(10, 11)
	s11_12 := buildSchedule(11, 12)
	if !PreviousSchedule(s10_11, s11_12) {
		t.Errorf("this case is previous case")
	}
	if PreviousSchedule(s11_12, s10_11) {
		t.Errorf("this case is not previous case")
	}
}

func buildSchedule(s int, e int) Schedule {
	now := time.Now()
	sD := time.Date(now.Year(), now.Month(), now.Day(), s, 0, 0, 0, now.Location())
	eD := time.Date(now.Year(), now.Month(), now.Day(), e, 0, 0, 0, now.Location())
	return Schedule{sD, eD}

}
