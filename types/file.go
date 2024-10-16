package types

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
)

//IsDir
//
//dashd in cmsys/file.c
func IsDir(path string) bool {
	theState, err := os.Stat(path)
	if err != nil {
		return false
	}
	return theState.IsDir()
}

func CopyFileToFile(src string, dst string) (err error) {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, SYS_BUFFER_SIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

func CopyFile(src string, dst string) (err error) {
	if IsDir(dst) {
		statSrc, err := os.Stat(src)
		if err != nil {
			return err
		}

		modeSrc := statSrc.Mode()

		if statSrc.IsDir() {
			return CopyDirToDir(src, dst)
		} else if modeSrc.IsRegular() {
			return CopyFileToDir(src, dst)
		} else {
			return ErrInvalidFile
		}
	} else if IsDir(src) {
		return CopyDirToDir(src, dst)
	} else {
		return CopyFileToFile(src, dst)
	}
}

func Mkdir(path string) error {
	return os.Mkdir(path, DEFAULT_FOLDER_CREATE_PERM)
}

func MkdirAll(path string) error {
	return os.MkdirAll(path, DEFAULT_FOLDER_CREATE_PERM)
}

func CopyDirToDir(src string, dst string) (err error) {
	_, err = os.Stat(src)
	if err != nil {
		return err
	}

	_, err = os.Stat(dst)
	if err != nil {
		err = MkdirAll(dst)
		if err != nil {
			return err
		}
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if entryName == "." || entryName == ".." {
			continue
		}

		childSrc := strings.Join([]string{src, entryName}, string(os.PathSeparator))
		childDst := strings.Join([]string{dst, entryName}, string(os.PathSeparator))

		if IsDir(childSrc) {
			err = MkdirAll(childDst)
			if err != nil {
				return err
			}
		}

		err = CopyFile(childSrc, childDst)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyFileToDir(src string, dst string) (err error) {
	basename := path.Base(src)
	dstFilename := strings.Join([]string{dst, basename}, string(os.PathSeparator))

	return CopyFileToFile(src, dstFilename)
}

func Unlink(filename string) (err error) {
	return os.Remove(filename)
}

//Rename
//
//Force rename src to dst by recursively-deleting old dst.
func Rename(src string, dst string) (err error) {
	_, err = os.Stat(dst)
	if err == nil {
		os.RemoveAll(dst)
	}

	dirname := path.Dir(dst)
	_, err = os.Stat(dirname)
	if err != nil {
		err = os.MkdirAll(dirname, DEFAULT_FOLDER_CREATE_PERM)
		if err != nil {
			return err
		}
	}

	return os.Rename(src, dst)
}

func DashT(filename string) (t Time4) {
	info, err := os.Stat(filename)
	if err != nil {
		return -1
	}

	return TimeToTime4(info.ModTime())
}

func Symlink(src string, dst string) (err error) {
	return os.Symlink(src, dst)
}

func DashS(filename string) (theSize int64) {
	info, err := os.Stat(filename)
	if err != nil {
		return -1
	}

	return info.Size()
}

func DashD(filename string) (isDir bool) {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func DashF(filename string) (isExists bool, err error) {
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func OpenCreate(filename string, flags int) (file *os.File, err error) {
	return os.OpenFile(filename, flags|os.O_CREATE, DEFAULT_FILE_CREATE_PERM)
}

func BinaryRead(reader io.ReadSeeker, order binary.ByteOrder, data interface{}) (err error) {
	origPos, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	defer func() {
		err2 := recover()
		if err2 == nil {
			return
		}

		err = ErrRecover(err2)

		newPos, err3 := reader.Seek(0, io.SeekCurrent)
		if err3 != nil {
			log.Errorf("BinaryRead (recover): unable to seek cur")
			return
		}
		if newPos < origPos {
			log.Errorf("BinaryRead (recover): newPos < currentPos: newPos: %v currentPos: %v", newPos, origPos)
			return
		}
		_, err3 = reader.Seek(origPos, io.SeekStart)
		if err3 != nil {
			log.Errorf("BinaryRead (recover): unable to seek orig-pos")
			return
		}
		theBytes := make([]byte, newPos-origPos)
		n, err3 := reader.Read(theBytes)
		log.Warnf("BinaryRead (recover): err: %v origPos: %v newPos: %v sz: %v bytes: %v, n: %v err3: %v", err, origPos, newPos, newPos-origPos, theBytes, n, err3)
		debug.PrintStack()
	}()

	return binary.Read(reader, order, data)
}

func BinaryWrite(writer io.Writer, order binary.ByteOrder, data interface{}) (err error) {
	defer func() {
		err2 := recover()
		if err2 == nil {
			return
		}

		err = ErrRecover(err2)

		log.Warnf("BinaryWrite (recover): err: %v data: %v", err, data)
		debug.PrintStack()
	}()

	return binary.Write(writer, order, data)
}
