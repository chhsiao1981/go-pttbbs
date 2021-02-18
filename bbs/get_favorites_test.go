package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestGetFavorites(t *testing.T) {
	setupTest()
	defer teardownTest()

	expectedContent := []byte{
		0x23, 0x0d, 0x03, 0x00, 0x02, 0x01, 0x01, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x03, 0x01, 0x01, 0x02,
		0x01, 0x01, 0xb7, 0x73, 0xaa, 0xba, 0xa5, 0xd8,
		0xbf, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x01, 0x01,
		0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	type args struct {
		uuserID    UUserID
		retrieveTS types.Time4
	}
	tests := []struct {
		name            string
		args            args
		expectedContent []byte
		expectedMtime   types.Time4
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{uuserID: "CodingMan"},
			expectedContent: expectedContent,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, _, err := GetFavorites(tt.args.uuserID, tt.args.retrieveTS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFavorites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("GetFavorites() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
		})
	}
	wg.Wait()
}
