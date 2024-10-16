package ptttype

import (
	"bufio"
	"os"
	"regexp"

	configutil "github.com/Ptt-official-app/go-pttbbs/configutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

const configPrefix = "go-pttbbs:ptttype"

func InitConfig() (err error) {
	config()

	err = postInitConfig()
	if err != nil {
		return err
	}

	err = checkTypes()
	if err != nil {
		return err
	}

	return initVars()
}

func setStringConfig(idx string, orig string) string {
	return configutil.SetStringConfig(configPrefix, idx, orig)
}

func setBoolConfig(idx string, orig bool) bool {
	return configutil.SetBoolConfig(configPrefix, idx, orig)
}

func setColorConfig(idx string, orig string) string {
	return configutil.SetColorConfig(configPrefix, idx, orig)
}

func setIntConfig(idx string, orig int) int {
	return configutil.SetIntConfig(configPrefix, idx, orig)
}

func setDoubleConfig(idx string, orig float64) float64 {
	return configutil.SetDoubleConfig(configPrefix, idx, orig)
}

func setServiceMode(serviceMode ServiceMode) ServiceMode {
	switch serviceMode {
	case DEV:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	return serviceMode
}

//SetBBSHOME
//
//This is to safely set BBSHOME
//Public to be used in the tests of other modules.
//
//Params
//	bbshome: new bbshome
//
//Return
//	string: original bbshome
func SetBBSHOME(bbshome string) string {
	origBBSHome := BBSHOME
	log.Debugf("SetBBSHOME: %v", bbshome)

	// config.go
	BBSHOME = bbshome
	BBSPROG = BBSHOME + /* 主程式 */
		string(os.PathSeparator) +
		BBSPROGPOSTFIX

	HAVE_USERAGREEMENT = BBSHOME +
		string(os.PathSeparator) +
		HAVE_USERAGREEMENT_POSTFIX
	HAVE_USERAGREEMENT_VERSION = BBSHOME +
		string(os.PathSeparator) +
		HAVE_USERAGREEMENT_VERSION_POSTFIX
	HAVE_USERAGREEMENT_ACCEPTABLE = BBSHOME +
		string(os.PathSeparator) +
		HAVE_USERAGREEMENT_ACCEPTABLE_POSTFIX

	// common.go
	FN_CONF_BANIP = BBSHOME + // 禁止連線的 IP 列表
		string(os.PathSeparator) +
		FN_CONF_BANIP_POSTFIX
	FN_PASSWD = BBSHOME + /* User records */
		string(os.PathSeparator) +
		FN_PASSWD_POSTFIX
	FN_BOARD = BBSHOME + /* board list */
		string(os.PathSeparator) +
		FN_BOARD_POSTFIX

	// const.go
	FN_FRESH = BBSHOME + string(os.PathSeparator) + FN_FRESH_POSTFIX /* mbbsd/register.c line: 381 */

	FN_ALLOW_EMAIL_LIST = BBSHOME + string(os.PathSeparator) + FN_ALLOW_EMAIL_LIST_POSTFIX
	FN_REJECT_EMAIL_LIST = BBSHOME + string(os.PathSeparator) + FN_REJECT_EMAIL_LIST_POSTFIX

	FN_DEFAULT_FAVS = BBSHOME + string(os.PathSeparator) + FN_DEFAULT_FAVS_POSTFIX

	return origBBSHome
}

func setBBSName(bbsname string) (origBBSName string) {
	origBBSName = BBSNAME
	BBSNAME = bbsname

	BBSNAME_BIG5 = types.Utf8ToBig5(BBSNAME)

	return origBBSName
}

//setBBSMName
//
//This is to safely set BBSMNAME
//
//Params
//	bbsmname: new bbsmname
//
//Return
//	string: original bbsmname
func setBBSMName(bbsmname string) string {
	origBBSMName := BBSMNAME
	log.Debugf("SetBBSMNAME: %v", bbsmname)

	BBSMNAME = bbsmname

	// regex-replace

	BN_SECURITY_s = regexReplace(BN_SECURITY_s, "BBSMNAME", BBSMNAME)
	BN_NOTE_s = regexReplace(BN_NOTE_s, "BBSMNAME", BBSMNAME)
	BN_RECORD_s = regexReplace(BN_RECORD_s, "BBSMNAME", BBSMNAME)
	BN_SYSOP_s = regexReplace(BN_SYSOP_s, "BBSMNAME", BBSMNAME)
	BN_TEST_s = regexReplace(BN_SECURITY_s, "BBSMNAME", BBSMNAME)
	BN_BUGREPORT_s = regexReplace(BN_BUGREPORT_s, "BBSMNAME", BBSMNAME)
	BN_LAW_s = regexReplace(BN_LAW_s, "BBSMNAME", BBSMNAME)
	BN_NEWBIE_s = regexReplace(BN_NEWBIE_s, "BBSMNAME", BBSMNAME)
	BN_ASKBOARD_s = regexReplace(BN_ASKBOARD_s, "BBSMNAME", BBSMNAME)
	BN_FOREIGN_s = regexReplace(BN_FOREIGN_s, "BBSMNAME", BBSMNAME)

	// config.go
	if IS_BN_FIVECHESS_LOG_INFERRED {
		BN_FIVECHESS_LOG_s = BBSMNAME + "Five"
	}
	if IS_BN_CCHESS_LOG_INFERRED {
		BN_CCHESS_LOG_s = BBSMNAME + "CChess"
	}
	if IS_MONEYNAME_INFFERRED {
		MONEYNAME = BBSMNAME + "幣"
	}

	BN_BUGREPORT_s = BBSMNAME + "Bug"
	BN_LAW_s = BBSMNAME + "Law"
	BN_NEWBIE_s = BBSMNAME + "NewHand"
	BN_FOREIGN_s = BBSMNAME + "Foreign"

	return origBBSMName
}

func regexReplace(str string, substr string, repl string) string {
	theRe := regexp.MustCompile("\\s*" + substr + "\\s*")
	if theRe == nil {
		return str
	}

	return theRe.ReplaceAllString(str, repl)
}

func setCAPTCHAInsertServerAddr(captchaInsertServerAddr string) string {
	origCAPTCHAInsertServerAddr := CAPTCHA_INSERT_SERVER_ADDR

	CAPTCHA_INSERT_SERVER_ADDR = captchaInsertServerAddr

	if IS_CAPTCHA_INSERT_HOST_INFERRED {
		CAPTCHA_INSERT_HOST = CAPTCHA_INSERT_SERVER_ADDR
	}

	return origCAPTCHAInsertServerAddr
}

//setMyHostname
//
//Params
//	myHostName: new my hostname
//
//Return
//	string: orig my hostname
func setMyHostname(myHostname string) string {
	origMyHostname := MYHOSTNAME

	MYHOSTNAME = myHostname

	if IS_AID_HOSTNAME_INFERRED {
		AID_HOSTNAME = MYHOSTNAME
	}

	return origMyHostname
}

//setRecycleBinName
//
//Params
//	recycleBinName: new recycle bin name
//
//Return
//	string: orig recycle bin name
func setRecycleBinName(recycleBinName string) string {
	origRecycleBinName := recycleBinName

	RECYCLE_BIN_NAME = recycleBinName
	RECYCLE_BIN_OWNER = "[" + RECYCLE_BIN_NAME + "]"

	return origRecycleBinName
}

func setFNSafeDel(fnSafeDel string) (origFNSafeDel string) {
	origFNSafeDel = FN_SAFEDEL

	FN_SAFEDEL = fnSafeDel
	FN_SAFEDEL_b = []byte(FN_SAFEDEL)
	FN_SAFEDEL_PREFIX_LEN = len(fnSafeDel)

	return origFNSafeDel
}

func setMaxPostMoney(maxPostMoney int) (origMaxPostMoney int) {
	origMaxPostMoney = MAX_POST_MONEY

	MAX_POST_MONEY = maxPostMoney
	ENTROPY_MAX = int(float64(MAX_POST_MONEY) * ENTROPY_RATIO)

	return origMaxPostMoney
}

func setEntropyRatio(entropyRatio float64) (origEntropyRatio float64) {
	origEntropyRatio = ENTROPY_RATIO

	ENTROPY_RATIO = entropyRatio
	ENTROPY_MAX = int(float64(MAX_POST_MONEY) * ENTROPY_RATIO)

	return origEntropyRatio
}

func postInitConfig() error {
	_ = setServiceMode(SERVICE_MODE)
	_ = SetBBSHOME(BBSHOME)
	_ = setBBSName(BBSNAME)
	_ = setBBSMName(BBSMNAME)
	_ = setCAPTCHAInsertServerAddr(CAPTCHA_INSERT_SERVER_ADDR)
	_ = setMyHostname(MYHOSTNAME)
	_ = setRecycleBinName(RECYCLE_BIN_NAME)
	_ = setBoards()
	_ = setFNSafeDel(FN_SAFEDEL)
	_ = setMaxPostMoney(MAX_POST_MONEY)
	_ = setEntropyRatio(ENTROPY_RATIO)

	return nil
}

func checkTypes() (err error) {
	if USEREC2_RAW_SZ != DEFAULT_USEREC2_RAW_SZ {
		log.Errorf("userec2 is not aligned: userec2: %v default-userec2: %v", USEREC2_RAW_SZ, DEFAULT_USEREC2_RAW_SZ)

		return ErrInvalidType
	}

	return nil
}

func initVars() (err error) {
	initReservedUserIDs()

	return nil
}

func initReservedUserIDs() {
	ReservedUserIDs = []types.Cstr{}

	filename := SetBBSHomePath(FN_CONF_RESERVED_ID)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for theBytes, err := types.ReadLine(reader); err == nil; theBytes, err = types.ReadLine(reader) {
		first, _ := types.CstrTokenR(theBytes, BYTES_SPACE)
		if len(first) == 0 {
			continue
		}
		log.Infof("initReservedUserIDs: theBytes: %v first: %v", string(theBytes), string(first))

		firstBytes := types.CstrToBytes(first)
		ReservedUserIDs = append(ReservedUserIDs, firstBytes)
	}
}

func setBoards() (err error) {
	BN_SECURITY = ToBoardID([]byte(BN_SECURITY_s))
	BN_NOTE = ToBoardID([]byte(BN_NOTE_s))
	BN_RECORD = ToBoardID([]byte(BN_RECORD_s))
	BN_SYSOP = ToBoardID([]byte(BN_SYSOP_s))
	BN_TEST = ToBoardID([]byte(BN_TEST_s))
	BN_BUGREPORT = ToBoardID([]byte(BN_BUGREPORT_s))
	BN_LAW = ToBoardID([]byte(BN_LAW_s))
	BN_NEWBIE = ToBoardID([]byte(BN_NEWBIE_s))
	BN_ASKBOARD = ToBoardID([]byte(BN_ASKBOARD_s))
	BN_FOREIGN = ToBoardID([]byte(BN_FOREIGN_s))
	BN_ARTDSN = ToBoardID([]byte(BN_ARTDSN_s))
	BN_BBSMOVIE = ToBoardID([]byte(BN_BBSMOVIE_s))
	BN_WHOAMI = ToBoardID([]byte(BN_WHOAMI_s))
	BN_FIVECHESS_LOG = ToBoardID([]byte(BN_FIVECHESS_LOG_s))
	BN_CCHESS_LOG = ToBoardID([]byte(BN_CCHESS_LOG_s))

	BN_ID_PROBLEM = ToBoardID([]byte(BN_ID_PROBLEM_s))
	BN_DELETED = ToBoardID([]byte(BN_DELETED_s))
	BN_JUNK = ToBoardID([]byte(BN_JUNK_s))
	BN_POLICELOG = ToBoardID([]byte(BN_POLICELOG_s))
	BN_UNANONYMOUS = ToBoardID([]byte(BN_UNANONYMOUS_s))
	BN_NEWIDPOST = ToBoardID([]byte(BN_NEWIDPOST_s))
	BN_ALLPOST = ToBoardID([]byte(BN_ALLPOST_s))
	BN_ALLHIDPOST = ToBoardID([]byte(BN_ALLHIDPOST_s))

	return nil
}
