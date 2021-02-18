package bbs

import (
	"reflect"
	"sync"
	"testing"
)

func TestLoadAutoCompleteBoards(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID     UUserID
		startIdxStr string
		nBoards     int
		keyword     string
		isAsc       bool
	}
	tests := []struct {
		name               string
		args               args
		expectedSummaries  []*BoardSummary
		expectedNextIdxStr string
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args:              args{uuserID: "SYSOP", startIdxStr: "", nBoards: 3, keyword: "w", isAsc: true},
			expectedSummaries: []*BoardSummary{testBoardSummary10},
		},
		{
			args:               args{uuserID: "SYSOP", startIdxStr: "", nBoards: 3, keyword: "", isAsc: true},
			expectedSummaries:  []*BoardSummary{testBoardSummary6, testBoardSummary7, testBoardSummary11},
			expectedNextIdxStr: "Note",
		},
		{
			args:               args{uuserID: "SYSOP", startIdxStr: "SYSOP", nBoards: 3, keyword: "", isAsc: false},
			expectedSummaries:  []*BoardSummary{testBoardSummary1, testBoardSummary9, testBoardSummary8},
			expectedNextIdxStr: "EditExp",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotNextIdxStr, err := LoadAutoCompleteBoards(tt.args.uuserID, tt.args.startIdxStr, tt.args.nBoards, tt.args.keyword, tt.args.isAsc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAutoCompleteBoards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummaries, tt.expectedSummaries) {
				t.Errorf("LoadAutoCompleteBoards() gotSummaries = %v, want %v", gotSummaries, tt.expectedSummaries)
			}
			if gotNextIdxStr != tt.expectedNextIdxStr {
				t.Errorf("LoadAutoCompleteBoards() gotNextIdxStr = %v, want %v", gotNextIdxStr, tt.expectedNextIdxStr)
			}
		})
	}
	wg.Wait()
}
