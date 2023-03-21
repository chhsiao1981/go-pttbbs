package ptt

import (
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func TestReadPost(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("WhoAmI"))

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))

	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	filename2 := &ptttype.Filename_t{}
	copy(filename2[:], []byte("M.1607202239.A.31D"))
	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.UID
		boardID    *ptttype.BoardID_t
		bid        ptttype.Bid
		filename   *ptttype.Filename_t
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
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				boardID:  boardID1,
				bid:      10,
				filename: filename1,
			},
			expectedContent: testContent1,
			expectedMtime:   1607209066,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				filename:   filename1,
				retrieveTS: 1607209066,
			},
			expectedContent: nil,
			expectedMtime:   1607209066,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				filename:   filename2,
				retrieveTS: 1607209066,
			},
			expectedContent: nil,
			expectedMtime:   0,
			wantErr:         true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, gotMtime, _, err := ReadPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.filename, tt.args.retrieveTS, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("ReadPost() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
			if !reflect.DeepEqual(gotMtime, tt.expectedMtime) {
				t.Errorf("ReadPost() gotMtime = %v, want %v", gotMtime, tt.expectedMtime)
			}
		})
		wg.Wait()
	}
}

func TestNewPost(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	SetupNewUser(testNewPostUser1)

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))
	ip0 := &ptttype.IPv4_t{}
	copy(ip0[:], []byte("127.0.0.1"))

	class0 := []byte("test")
	title0 := []byte("this is a test")
	fullTitle0 := ptttype.Title_t{}
	copy(fullTitle0[:], []byte("[test] this is a test"))
	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("A1"))
	expectedSummary0 := &ptttype.ArticleSummaryRaw{
		Aid:     3,
		BoardID: boardID0,
		FileHeaderRaw: &ptttype.FileHeaderRaw{
			Title: fullTitle0,
			Owner: owner0,
		},
	}

	expectedSummary1 := &ptttype.ArticleSummaryRaw{
		Aid:     4,
		BoardID: boardID0,
		FileHeaderRaw: &ptttype.FileHeaderRaw{
			Title: fullTitle0,
			Owner: owner0,
		},
	}

	cache.Shm.Shm.BCache[9].NUser = 40

	uid0, _ := cache.DoSearchUserRaw(&testNewPostUser1.UserID, nil)

	content0 := [][]byte{[]byte("test1"), []byte("test2")}
	expected0 := []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'A', '1', ' ', // 作者: A1
		0x28, 0xaf, 0xab, //(神
		0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, //) 看板
		0x3a, 0x20, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x0a, //: WhoAmI
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, // 標題: [tes
		0x74, 0x5d, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, // t] this is
		0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // a test
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, 0x00, 0x00, 0x00, 0x20, // 時間: 000
		0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, // 000 00 00:
		0x00, 0x00, 0x3a, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, // 00:00 0000
		0x0a,
		0x74, 0x65, 0x73, 0x74, 0x31, 0x0a, // test1
		0x74, 0x65, 0x73, 0x74, 0x32, 0x0a, // test2
		0x0a,
		0x2d, 0x2d, 0x0a, //--
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, //※ 發信站:
		0x20, 0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, // 新批踢踢(
		0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // ptt2.cc),
		0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, // 來自: 127.
		0x30, 0x2e, 0x30, 0x2e, 0x31, 0x0a, // 0.0.1
		0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, //※ 文章網
		0xa7, 0x7d, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, // 址: http:/
		0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, ///localhost
		0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, ///bbs/WhoAm
		0x49, 0x2f, 0x4d, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // I/M.000000
		0x00, 0x00, 0x00, 0x00, 0x2e, 0x41, 0x2e, 0x00, 0x00, 0x00, // 0000.A.000
		0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x0a, //.html
	}

	removeIdxes := []int{
		61, 62, 63, 65, 66, 67, 69, 70, 72, 73, 75, 76, 78, 79, 81, 82, 83, 84, // 時間
		192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 205, 206, 207, // 文章網址
	}

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		boardID  *ptttype.BoardID_t
		bid      ptttype.Bid
		posttype []byte
		title    []byte
		content  [][]byte
		ip       *ptttype.IPv4_t
		from     []byte
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary *ptttype.ArticleSummaryRaw
		expected        []byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "post0",
			args: args{
				user:     testNewPostUser1,
				uid:      uid0,
				boardID:  boardID0,
				bid:      10,
				posttype: class0,
				title:    title0,
				content:  content0,
				ip:       ip0,
			},
			expectedSummary: expectedSummary0,
			expected:        expected0,
		},
		{
			name: "post1",
			args: args{
				user:     testNewPostUser1,
				uid:      uid0,
				boardID:  boardID0,
				bid:      10,
				posttype: class0,
				title:    title0,
				content:  content0,
				ip:       ip0,
			},
			expectedSummary: expectedSummary1,
			expected:        expected0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := NewPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.posttype, tt.args.title, tt.args.content, tt.args.ip, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			content, mtime, _, err := ReadPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, &gotSummary.FileHeaderRaw.Filename, 0, false)
			if err != nil {
				t.Errorf("NewPost() unable to ReadPost: e: %v", err)
				return
			}

			if mtime != gotSummary.Modified {
				t.Errorf("NewPost() mtime: %v expected: %v", mtime, gotSummary.Modified)
			}
			gotSummary.Filename = ptttype.Filename_t{}
			gotSummary.Modified = 0
			gotSummary.Date = ptttype.Date_t{}
			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)

			for _, idx := range removeIdxes {
				if idx >= len(content) {
					break
				}
				content[idx] = 0x00
			}

			testutil.TDeepEqual(t, "content", content, tt.expected)
		})
		wg.Wait()
	}
}

func TestCrossPost(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	testBoardID0 := &ptttype.BoardID_t{}
	copy(testBoardID0[:], []byte("mewboard0"))
	testBrdClass0 := []byte("CPBL")
	testBrdTitle0 := []byte("new-board")

	testForwardBoardID0 := &ptttype.BoardID_t{}
	copy(testForwardBoardID0[:], []byte("fwboard0"))
	testForwardBrdTitle0 := []byte("fw-board")

	boardSummary, err := NewBoard(testUserecRaw3, 6, 2, testBoardID0, testBrdClass0, testBrdTitle0, nil, ptttype.BRD_CPLOG, 0, ptttype.CHESSCODE_NONE, false)
	logrus.Infof("TestCrossPost: after NewBoard: e: %v", err)

	forwardBoardSummary, _ := NewBoard(testUserecRaw3, 6, 2, testForwardBoardID0, testBrdClass0, testForwardBrdTitle0, nil, ptttype.BRD_CPLOG, 0, ptttype.CHESSCODE_NONE, false)

	SetupNewUser(testNewPostUser1)

	ip0 := &ptttype.IPv4_t{}
	copy(ip0[:], []byte("127.0.0.1"))

	class0 := []byte("test")
	title0 := []byte("this is a test")
	fullTitle0 := ptttype.Title_t{}
	copy(fullTitle0[:], []byte("[test] this is a test"))
	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("A1"))

	fullForwardTitle0 := ptttype.Title_t{}
	copy(fullForwardTitle0[:], []byte("Fw: [test] this is a test"))

	cache.Shm.Shm.BCache[9].NUser = 40

	uid0, _ := cache.DoSearchUserRaw(&testNewPostUser1.UserID, nil)

	content0 := [][]byte{[]byte("test1"), []byte("test2")}

	ownerForward0 := ptttype.Owner_t{}
	copy(ownerForward0[:], []byte("SYSOP"))

	expectedSummary0 := &ptttype.ArticleSummaryRaw{
		Aid:     1,
		BoardID: testForwardBoardID0,
		FileHeaderRaw: &ptttype.FileHeaderRaw{
			Title: fullForwardTitle0,
			Owner: ownerForward0,
		},
	}

	expectedContent0 := []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'A', '1', ' ', // 作者: A1
		0x28, 0xaf, 0xab, //(神
		0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, //) 看板
		0x3a, 0x20, 0x6d, 0x65, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, 0x0a, //: newboard0
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, // 標題: [tes (30~39)
		0x74, 0x5d, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, // t] this is (40~49)
		0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // a test (50~57)
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, 0x00, 0x00, 0x00, 0x20, // 時間: 000 (58~67)
		0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, // 000 00 00: (68~77)
		0x00, 0x00, 0x3a, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, // 00:00 0000 (78~88)
		0x0a,
		0x74, 0x65, 0x73, 0x74, 0x31, 0x0a, // test1 (79~84)
		0x74, 0x65, 0x73, 0x74, 0x32, 0x0a, // test2 (85~90)
		0x0a,
		0x2d, 0x2d, 0x0a, //--
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, //※ 發信站:
		0x20, 0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, // 新批踢踢(
		0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // ptt2.cc),
		0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, // 來自: 127.
		0x30, 0x2e, 0x30, 0x2e, 0x31, 0x0a, // 0.0.1
		0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, //※ 文章網
		0xa7, 0x7d, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, // 址: http:/
		0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, // /localhost
		0x2f, 0x62, 0x62, 0x73, 0x2f, 0x6d, 0x65, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, // /bbs/newboard0
		0x2f, 0x4d, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // /M.000000
		0x00, 0x00, 0x00, 0x00, 0x2e, 0x41, 0x2e, 0x00, 0x00, 0x00, // 0000.A.000
		0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x0a, //.html
		0xa1, 0xb0, 0x20, // ※
		0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x32, 0x6d, // [1;32m
		0x53, 0x59, 0x53, 0x4f, 0x50, // SYSOP
		0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x32, 0x6d, // [0;32m:
		0x3a, 0xc2, 0xe0, 0xbf, 0xfd, 0xa6, 0xdc, 0xac, 0xdd, 0xaa, 0x4f, 0x20, // 轉錄至看板
		0x66, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, // fwboard0
		0x1b, 0x5b, 0x6d, //
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, //
		0x00, 0x00, 0x2f, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x0a, // 00/00 00:00
	}

	removeIdxes := []int{
		64, 65, 66, 68, 69, 70, 72, 73, 75, 76, 78, 79, 81, 82, 84, 85, 86, 87,
		198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 211, 212, 213,
		304, 305, 307, 308, 310, 311, 313, 314,
	}

	expectedComment0 := []byte{
		0xa1, 0xb0, 0x20, // ※
		0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x32, 0x6d, // [1;32m
		0x53, 0x59, 0x53, 0x4f, 0x50, // SYSOP
		0x1b, 0x5b, 0x30, 0x3b, 0x33, 0x32, 0x6d, // [0;32m:
		0x3a, 0xc2, 0xe0, 0xbf, 0xfd, 0xa6, 0xdc, 0xac, 0xdd, 0xaa, 0x4f, 0x20, // 轉錄至看板
		0x66, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, // fwboard0
		0x1b, 0x5b, 0x6d, //
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, //
		0x00, 0x00, 0x2f, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x0a, // 00/00 00:00

	}

	removeCommentIdxes := []int{
		84, 85, 87, 88, 90, 91, 93, 94,
	}

	expectedContentForward0 := []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'S', 'Y', 'S', 'O', 'P', ' ', // 作者: SYSOP
		0x28, 0xaf, 0xab, 0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, //(神) 看板
		0x3a, 0x20, 0x66, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, 0x0a, //: fwboard0
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x46, 0x77, 0x3a, 0x20, // 標題: Fw:
		0x5b, 0x74, 0x65, 0x73, 0x74, 0x5d, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, //[test] this
		0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // is a test
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, 0x00, 0x00, 0x00, 0x20, // 時間: 000
		0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, // 000 00 00:
		0x00, 0x00, 0x3a, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, // 00:00 0000 (78~88)
		0x0a,
		0xa1, 0xb0, 0x20, 0x5b, 0xa5, 0xbb, 0xa4, 0xe5, 0xc2, 0xe0, 0xbf, 0xfd, 0xa6, 0xdb, 0x20, //※ [本文轉錄自
		0x6d, 0x65, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, 0x20, // newboard0
		0xac, 0xdd, 0xaa, 0x4f, 0x20, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5d, 0x0a, // 看板 #1X87_A66]
		0x0a,
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'A', '1', ' ', // 作者: A1
		0x28, 0xaf, 0xab, 0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, 0x3a, 0x20, //(神) 看板:
		0x6d, 0x65, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, 0x0a, // newboard0
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, 0x74, 0x5d, 0x20, // 標題: [test]
		0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // this is a test
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x20, // 時間: Sat Aug 21
		0x00, 0x00, 0x3a, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, // 12:23:37 2021
		0x0a,
		0x74, 0x65, 0x73, 0x74, 0x31, 0x0a, // test1
		0x74, 0x65, 0x73, 0x74, 0x32, 0x0a, // test2
		0x0a,
		0x2d, 0x2d, 0x0a, //--
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, 0x20, //※ 發信站:
		0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, 0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // 新批踢踢(ptt2.cc),
		0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x0a, // 來自: 127.0.0.1
		0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, 0xa7, 0x7d, 0x3a, 0x20, //※ 文章網址:
		0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x2f, 0x62, 0x62, 0x73, 0x2f, 0x6d, 0x65, 0x77, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x30, 0x2f, 0x4d, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x41, 0x2e, 0x00, 0x00, 0x00, 0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x0a, // http://localhost/bbs/newboard0/M.1629519818.A.186.html
		0x0a,
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, 0x20, //※ 發信站:
		0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, 0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x0a, // 新批踢踢(ptt2.cc)
		0xa1, 0xb0, 0x20, 0xc2, 0xe0, 0xbf, 0xfd, 0xaa, 0xcc, 0x3a, 0x20, //※ 轉錄者:
		'S', 'Y', 'S', 'O', 'P', ' ', 0x28, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x29, 0x2c, 0x20, // SYSOP (127.0.0.1),
		0x00, 0x00, 0x2f, 0x00, 0x00, 0x2f, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x0a, // 08/21/2021 12:23:37
	}

	removeIdxesForward := []int{
		70, 71, 72, 74, 75, 76, 78, 79, 81, 82, 84, 85, 87, 88, 90, 91, 92, 93,
		127, 128, 129, 130, 131, 132, 133, 134,
		202, 203, 204, 206, 207, 208, 210, 211, 213, 214, 216, 217, 219, 220, 222, 223, 224, 225,
		336, 337, 338, 339, 340, 341, 342, 343, 344, 345, 349, 350, 351,
		418, 419, 421, 422, 424, 425, 426, 427, 429, 430, 432, 433, 435, 436,
	}

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		boardID  *ptttype.BoardID_t
		bid      ptttype.Bid
		xBoardID *ptttype.BoardID_t
		xBid     ptttype.Bid
		filemode ptttype.FileMode
		ip       *ptttype.IPv4_t
		from     []byte

		contentUser  *ptttype.UserecRaw
		contentUID   ptttype.UID
		contentTitle []byte
		posttype     []byte
		content      [][]byte
	}
	tests := []struct {
		name                   string
		args                   args
		expectedArticleSummary *ptttype.ArticleSummaryRaw
		expectedComment        []byte
		expectedCommentMTime   types.Time4
		expectedContent        []byte
		expectedContentForward []byte
		wantErr                bool
	}{
		// TODO: Add test cases.
		{
			name: "post0",
			args: args{
				user:         testUserecRaw1,
				uid:          1,
				contentUser:  testNewPostUser1,
				contentUID:   uid0,
				contentTitle: title0,
				boardID:      boardSummary.Brdname,
				bid:          boardSummary.Bid,
				posttype:     class0,
				content:      content0,
				xBoardID:     forwardBoardSummary.Brdname,
				xBid:         forwardBoardSummary.Bid,
				ip:           ip0,
			},
			expectedComment:        expectedComment0,
			expectedArticleSummary: expectedSummary0,
			expectedContent:        expectedContent0,
			expectedContentForward: expectedContentForward0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			summary, err := NewPost(tt.args.contentUser, tt.args.contentUID, tt.args.boardID, tt.args.bid, tt.args.posttype, tt.args.contentTitle, tt.args.content, tt.args.ip, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrossPost() unable to NewPost: error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotArticleSummary, gotComment, gotCommentMTime, err := CrossPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, &summary.Filename, tt.args.xBoardID, tt.args.xBid, tt.args.filemode, tt.args.ip, tt.args.from)

			if (err != nil) != tt.wantErr {
				t.Errorf("CrossPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotArticleSummary.Date = ptttype.Date_t{}
			gotArticleSummaryFilenameStr := gotArticleSummary.Filename.String()
			gotArticleSummary.Filename = ptttype.Filename_t{}
			testutil.TDeepEqual(t, "summary", gotArticleSummary, tt.expectedArticleSummary)

			for _, idx := range removeCommentIdxes {
				if idx >= len(gotComment) {
					break
				}
				gotComment[idx] = 0x00
			}
			testutil.TDeepEqual(t, "comment", gotComment, tt.expectedComment)

			if gotCommentMTime == 0 {
				t.Errorf("CrossPost() gotCommentMTime == 0")
			}

			// content
			filename, _ := path.SetBFile(tt.args.boardID, summary.Filename.String())
			logrus.Infof("CrossPost: filename: %v", filename)
			file, err := os.Open(filename)
			if err != nil {
				t.Errorf("CrossPost(): unable to open file: %v e: %v", filename, err)
				return
			}
			defer file.Close()

			content, _ := io.ReadAll(file)
			logrus.Infof("CrossPost: content: %v", content)

			for _, idx := range removeIdxes {
				if idx >= len(content) {
					break
				}
				content[idx] = 0x00
			}

			testutil.TDeepEqual(t, "content", content, tt.expectedContent)

			// content-forward
			filenameForward, _ := path.SetBFile(tt.args.xBoardID, gotArticleSummaryFilenameStr)
			logrus.Infof("CrossPost: filenameForward: %v", filenameForward)
			fileForward, err := os.Open(filenameForward)
			if err != nil {
				t.Errorf("CrossPost(): unable to open file: %v e: %v", filenameForward, err)
				return
			}
			defer fileForward.Close()

			contentForward, _ := io.ReadAll(fileForward)
			logrus.Infof("CrossPost: contentForward: %v", contentForward)
			for _, idx := range removeIdxesForward {
				if idx >= len(contentForward) {
					break
				}
				contentForward[idx] = 0x00
			}

			testutil.TDeepEqual(t, "contentForward", contentForward, tt.expectedContentForward)
		})
		wg.Wait()
	}
}

func TestEditPost(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	SetupNewUser(testNewPostUser1)

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))
	ip0 := &ptttype.IPv4_t{}
	copy(ip0[:], []byte("127.0.0.1"))

	class0 := []byte("test")
	title0 := []byte("this is a test")
	fullTitle0 := ptttype.Title_t{}
	copy(fullTitle0[:], []byte("[test] this is a test"))
	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("A1"))

	uid0, _ := cache.DoSearchUserRaw(&testNewPostUser1.UserID, nil)
	content0 := [][]byte{[]byte("test1"), []byte("test2")}

	articleSummary, _ := NewPost(testNewPostUser1, uid0, boardID0, 10, class0, title0, content0, ip0, nil)

	filename0, _ := path.SetBFile(boardID0, articleSummary.Filename.String())
	file0, _ := os.Open(filename0)
	defer file0.Close()
	postContent0, _ := io.ReadAll(file0)

	oldSZ0 := len(postContent0)
	oldSum0 := cmsys.FNV1_64_INIT
	oldSum0 = cmsys.Fnv64Buf(postContent0, oldSZ0, oldSum0)

	logrus.Infof("oldSZ0: %v oldSum0: %v", oldSZ0, oldSum0)

	editContent0 := [][]byte{
		{
			0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'A', '1', ' ', // 作者: A1
			0x28, 0xaf, 0xab, //(神
			0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, //) 看板
			0x3a, 0x20, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, //: WhoAmI
		},
		{
			0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, // 標題: [tes
			0x74, 0x5d, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, // t] this is
			0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, // a test
		},
		{
			0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, // 時間:
		},
		{},
		{
			0x74, 0x65, 0x73, 0x74, 0x38, // test8
		},
		{
			0x74, 0x65, 0x73, 0x74, 0x39, // test9
		},
		{},
		{
			0x2d, 0x2d, //--
		},
		{
			0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, //※ 發信站:
			0x20, 0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, // 新批踢踢(
			0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // ptt2.cc),
			0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, // 來自: 127.
			0x30, 0x2e, 0x30, 0x2e, 0x31, // 0.0.1
		},
		{
			0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, //※ 文章網
			0xa7, 0x7d, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, // 址: http:/
			0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, // /localhost
			0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x2f, 0x4d, 0x2e, // /bbs/WhoAmI/M.
		},
	}

	expectedContent0 := []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 0x41, 0x31, 0x20, // 作者: A1
		0x28, 0xaf, 0xab, 0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, 0x3a, 0x20, //(神) 看板:
		0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x0a, // WhoAmI
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, 0x74, 0x5d, 0x20, // 標題: [test]
		0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // this is a test
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x0a, // 時間:
		0x0a,
		0x74, 0x65, 0x73, 0x74, 0x38, 0x0a, // test8
		0x74, 0x65, 0x73, 0x74, 0x39, 0x0a, // test9
		0x0a,
		0x2d, 0x2d, 0x0a, //--
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, 0x20, // ※ 發信站:
		0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, // 新批踢踢
		0x28, 0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // (ptt2.cc),
		0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x0a, // 來自: 127.0.0.1
		0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, 0xa7, 0x7d, 0x3a, 0x20, // ※ 文章網址:
		0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, // http://localhost
		0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x2f, 0x4d, 0x2e, 0x0a, // /bbs/WhoAmI/M.
		0xa1, 0xb0, 0x20, 0xbd, 0x73, 0xbf, 0xe8, 0x3a, 0x20, 0x41, 0x31, 0x20, // ※ 編輯: A1
		0x28, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x29, 0x2c, 0x20, // (127.0.0.1)
		0x00, 0x00, 0x2f, 0x00, 0x00, 0x2f, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x0a, // 00/00/0000 00:00:00
	}

	removeIdxes0 := []int{
		193, 194, 196, 197, 199, 200, 201, 202, 204, 205, 207, 208, 210, 211, // 時間
	}

	expectedTitle0 := &fullTitle0

	editContent1 := [][]byte{
		{
			0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 0x41, 0x31, 0x20, // 作者: A1
			0x28, 0xaf, 0xab, 0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, 0x3a, 0x20, //(神) 看板:
			0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, // WhoAmI
		},
		{
			0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, 0x74, 0x5d, 0x20, // 標題: [test]
			0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, // this is a test
		},
		{
			0xae, 0xc9, 0xb6, 0xa1, 0x3a, // 時間:
		},
		{},
		{
			0x74, 0x65, 0x73, 0x74, 0x31, 0x30, // test10
		},
		{
			0x74, 0x65, 0x73, 0x74, 0x31, 0x31, // test11
		},
		{},
		{
			0x2d, 0x2d, //--
		},
		{
			0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, 0x20, // ※ 發信站:
			0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, // 新批踢踢
			0x28, 0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // (ptt2.cc),
			0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, // 來自: 127.0.0.1
		},
		{
			0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, 0xa7, 0x7d, 0x3a, 0x20, // ※ 文章網址:
			0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, // http://localhost
			0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x2f, 0x4d, 0x2e, // /bbs/WhoAmI/M.
		},
		{
			0xa1, 0xb0, 0x20, 0xbd, 0x73, 0xbf, 0xe8, 0x3a, 0x20, 0x41, 0x31, 0x20, // ※ 編輯: A1
			0x28, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x29, 0x2c, 0x20, // (127.0.0.1)
			0x00, 0x00, 0x2f, 0x00, 0x00, 0x2f,
		},
	}

	expectedContent1 := []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 0x41, 0x31, 0x20, // 作者: A1
		0x28, 0xaf, 0xab, 0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, 0x3a, 0x20, //(神) 看板:
		0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x0a, // WhoAmI
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, 0x74, 0x5d, 0x20, // 標題: [test]
		0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // this is a test
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x0a, // 時間:
		0x0a,
		0x74, 0x65, 0x73, 0x74, 0x31, 0x30, 0x0a, // test10
		0x74, 0x65, 0x73, 0x74, 0x31, 0x31, 0x0a, // test11
		0x0a,
		0x2d, 0x2d, 0x0a, //--
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, 0x20, // ※ 發信站:
		0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, // 新批踢踢
		0x28, 0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // (ptt2.cc),
		0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x0a, // 來自: 127.0.0.1
		0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, 0xa7, 0x7d, 0x3a, 0x20, // ※ 文章網址:
		0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, // http://localhost
		0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x2f, 0x4d, 0x2e, 0x0a, // /bbs/WhoAmI/M.
		0xa1, 0xb0, 0x20, 0xbd, 0x73, 0xbf, 0xe8, 0x3a, 0x20, 0x41, 0x31, 0x20, // ※ 編輯: A1
		0x28, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x29, 0x2c, 0x0a, // (127.0.0.1)
		0xa1, 0xb0, 0x20, 0xbd, 0x73, 0xbf, 0xe8, 0x3a, 0x20, 0x53, 0x59, 0x53, 0x4f, 0x50, 0x20, // ※ 編輯: SYSOP
		0x28, 0x31, 0x32, 0x37, 0x2e, 0x30, 0x2e, 0x30, 0x2e, 0x31, 0x29, 0x2c, 0x20, // (127.0.0.1)
		0x00, 0x00, 0x2f, 0x00, 0x00, 0x2f, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x3a, 0x00, 0x00, 0x0a, // 00/00/0000 00:00:00
	}

	removeIdxes1 := []int{
		223, 224, 226, 227, 229, 230, 231, 232, 234, 235, 237, 238, 240, 241, // 時間
	}

	expectedTitle1 := &fullTitle0

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.UID
		boardID  *ptttype.BoardID_t
		bid      ptttype.Bid
		filename *ptttype.Filename_t
		posttype []byte
		title    []byte
		content  [][]byte
		ip       *ptttype.IPv4_t
		from     []byte
	}
	tests := []struct {
		name               string
		args               args
		expectedNewContent []byte
		expectedMtime      types.Time4
		expectedTitle      *ptttype.Title_t
		removeIdxes        []int
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:     testNewPostUser1,
				uid:      uid0,
				boardID:  boardID0,
				bid:      10,
				filename: &articleSummary.Filename,
				content:  editContent0,
				ip:       ip0,
			},
			removeIdxes:        removeIdxes0,
			expectedNewContent: expectedContent0,
			expectedTitle:      expectedTitle0,
		},
		{
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				boardID:  boardID0,
				bid:      10,
				filename: &articleSummary.Filename,
				content:  editContent1,
				ip:       ip0,
			},
			removeIdxes:        removeIdxes1,
			expectedNewContent: expectedContent1,
			expectedTitle:      expectedTitle1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			filename, _ := path.SetBFile(boardID0, articleSummary.Filename.String())
			file, _ := os.Open(filename)
			defer file.Close()
			postContent, _ := io.ReadAll(file)

			oldSZ := len(postContent)
			oldSum := cmsys.FNV1_64_INIT
			oldSum = cmsys.Fnv64Buf(postContent, oldSZ, oldSum)

			gotNewContent, _, gotTitle, err := EditPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.filename, tt.args.posttype, tt.args.title, tt.args.content, oldSZ, oldSum, tt.args.ip, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("EditPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			logrus.Infof("gotNewContent: %x", gotNewContent)

			for _, idx := range tt.removeIdxes {
				if idx >= len(gotNewContent) {
					break
				}
				gotNewContent[idx] = 0x00
			}

			testutil.TDeepEqual(t, "got", gotNewContent, tt.expectedNewContent)

			testutil.TDeepEqual(t, "title", gotTitle, tt.expectedTitle)
		})
		wg.Wait()
	}
}

func TestReadPostTemplate(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	filename := "testcase/boards/W/WhoAmI/postsample.0"
	mtime := time.Unix(1679419504, 0)
	os.Chtimes(filename, mtime, mtime)

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("WhoAmI"))

	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.UID
		boardID    *ptttype.BoardID_t
		bid        ptttype.Bid
		templateID ptttype.SortIdx
		retrieveTS types.Time4
		isHash     bool
	}
	tests := []struct {
		name         string
		args         args
		wantContent  []byte
		wantMtime    types.Time4
		wantChecksum cmsys.Fnv64_t
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				templateID: 1,
			},
			wantContent: testContent1,
			wantMtime:   1679419504,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				templateID: 1,
				isHash:     true,
			},
			wantContent:  testContent1,
			wantMtime:    1679419504,
			wantChecksum: 708027644387720498,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, gotMtime, gotChecksum, err := ReadPostTemplate(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.templateID, tt.args.retrieveTS, tt.args.isHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPostTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("ReadPostTemplate() gotContent = %v, want %v", gotContent, tt.wantContent)
			}
			if !reflect.DeepEqual(gotMtime, tt.wantMtime) {
				t.Errorf("ReadPostTemplate() gotMtime = %v, want %v", gotMtime, tt.wantMtime)
			}
			if !reflect.DeepEqual(gotChecksum, tt.wantChecksum) {
				t.Errorf("ReadPostTemplate() gotChecksum = %v, want %v", gotChecksum, tt.wantChecksum)
			}
		})
		wg.Wait()
	}
}
