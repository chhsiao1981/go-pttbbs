package ptt

import (
	"sync"
	"testing"
)

func TestAttemptRemoteCaptcha(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	tests := []struct {
		name           string
		expectedUrl    string
		expectedVerify string
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			expectedUrl:    "",
			expectedVerify: "",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUrl, gotVerify, err := AttemptRemoteCaptcha()
			if (err != nil) != tt.wantErr {
				t.Errorf("AttemptRemoteCaptcha() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUrl != tt.expectedUrl {
				t.Errorf("AttemptRemoteCaptcha() gotUrl = %v, want %v", gotUrl, tt.expectedUrl)
			}
			if gotVerify != tt.expectedVerify {
				t.Errorf("AttemptRemoteCaptcha() gotVerify = %v, want %v", gotVerify, tt.expectedVerify)
			}
		})
		wg.Wait()
	}
}
