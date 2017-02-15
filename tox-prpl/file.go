/*
  all implemention code except must code.
*/
package main

import (
	"log"
	"os"
	/*
		"fmt"
		"strings"
		"time"
	*/
	"io/ioutil"

	"go-purple/purple"

	"github.com/kitech/go-toxcore"
)

type xferSendData struct {
	fileNumber   uint32
	fileId       string
	friendNumber uint32
	who          string
	fp           *os.File
	fileKind     uint32
	fileSize     uint64
	fileName     string
}

// embed into ToxPlugin struct
// used for avatar only
type fileTransferFields struct {
	ftqueue map[uint32]*xferSendData // fileNumber => *xferSendData
}

func (this *ToxPlugin) setupModuleFields() {
	this.ftqueue = make(map[uint32]*xferSendData)
}

func (this *ToxPlugin) setupFileCallbacks(ac *purple.Account) {
	this._tox.CallbackFileChunkRequest(this.onFileChunkRequest, ac)
	this._tox.CallbackFileRecv(this.onFileRecv, ac)
	this._tox.CallbackFileRecvChunk(this.onFileRecvChunk, ac)
	this._tox.CallbackFileRecvControl(this.onFileRecvControl, ac)
}

func (this *ToxPlugin) onFileRecv(t *tox.Tox, friendNumber uint32, fileNumber uint32, kind uint32, fileSize uint64, fileName string, ud interface{}) {
	log.Println("herere", fileName)
	ac := ud.(*purple.Account)

	pubkey, err := this._tox.FriendGetPublicKey(friendNumber)
	if err != nil {
		log.Println(err)
	}

	xd := &xferSendData{}
	xd.friendNumber = friendNumber
	xd.fileNumber = fileNumber
	xd.who = pubkey
	xd.fileKind = kind
	xd.fileSize = fileSize

	if kind == tox.FILE_KIND_AVATAR {
		log.Println("what can I do?")
		// get avatar dir
		// recieve avatar data and save to avatar dir, with friendId as name
		// notify purple peer's icon changed
		this.ftqueue[fileNumber] = xd
		fname := os.Getenv("HOME") + "/" + purple.UserDir() + "/icons/" + pubkey + ".jpg"
		log.Println(fname)
		xd.fileName = fname
		xd.fp, err = os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Println(err)
		} else {
			// go on
			_, err = this._tox.FileControl(friendNumber, fileNumber, tox.FILE_CONTROL_RESUME)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		xfer := ac.NewXfer(purple.XFER_RECEIVE, pubkey)
		this.xferSetCallbacks(xfer)
		xfer.Data = xd
		xfer.SetFilename(fileName)
		xfer.SetSize(fileSize)
		xferSendMap[fileNumber] = xfer

		xfer.Request()
	}
}

func (this *ToxPlugin) onFileRecvChunk(t *tox.Tox, friendNumber uint32, fileNumber uint32, position uint64, data []byte, ud interface{}) {
	ac := ud.(*purple.Account)

	if position == 0 || len(data) == 0 {
		log.Println("start or finished:", position, len(data))
	}

	// process avatar
	if xd, ok := this.ftqueue[fileNumber]; ok {
		if position == xd.fileSize {
			log.Println("avatar recv finished.")
			xd.fp.Close()
			delete(this.ftqueue, fileNumber)
			// TODO change purple's buddy icon
			fname := xd.fileName
			log.Println("change buddy's icon to:", fname)
			icon_data, err := ioutil.ReadFile(fname)
			if err != nil {
				log.Println(err)
			} else {
				ac.BuddyIconsSetForUser(xd.who, icon_data)
			}
			return
		}
		n, err := xd.fp.WriteAt(data, int64(position))
		if err != nil {
			log.Println(err, n, len(data), xd.fp.Fd())
		} else if n != len(data) {
			log.Println(n, len(data))
		} else {
			// write ok
		}
		return
	}

	// process normal file
	x, ok := xferSendMap[fileNumber]
	if !ok {
		log.Println("wtf", ok, x)
		return
	}
	sd := x.Data.(*xferSendData)

	if position == x.GetSize() {
		log.Println("finished")
		x.SetBytesSent(position)
		x.SetCompleted(true)
		if err := sd.fp.Truncate(int64(position)); err != nil {
			log.Println(err)
		}
		sd.fp.Close()
		delete(xferSendMap, fileNumber)
		return
	}

	n, err := sd.fp.WriteAt(data, int64(position))
	if err != nil {
		log.Println(err, n, len(data), sd.fp.Fd())
	} else if n != len(data) {
		log.Println(n, len(data))
	} else {
		x.SetBytesSent(position)
		x.UpdateProgress()
	}

}

func (this *ToxPlugin) onFileRecvControl(t *tox.Tox, friendNumber uint32, fileNumber uint32,
	control int, ud interface{}) {
	log.Println("herere", control)
	// ac := ud.(*purple.Account)

	x, ok := xferSendMap[fileNumber]
	if !ok {
		log.Println("wtf", ok, x)
		return
	}
	sd := x.Data.(*xferSendData)

	var err error
	switch control {
	case tox.FILE_CONTROL_RESUME:
		log.Println(x.GetLocalFilename())
		if sd.fp != nil {
			log.Println("already opened, skip:", sd.fp)
			break
		}
		sd.fp, err = os.Open(x.GetLocalFilename())
		if err != nil {
			log.Println(err)
		} else {
			x.Start(-1, "", 0)
		}
	case tox.FILE_CONTROL_CANCEL:
		if sd.fp != nil {
			sd.fp.Close()
		}
		x.End()
		delete(xferSendMap, fileNumber)
	case tox.FILE_CONTROL_PAUSE:
	}
}

func (this *ToxPlugin) onFileChunkRequest(t *tox.Tox, friend_number uint32, file_number uint32, position uint64, length int, ud interface{}) {
	ac := ud.(*purple.Account)

	if position == 0 || length == 0 {
		log.Println("herere", position, length)
	}

	x, ok := xferSendMap[file_number]
	if !ok {
		log.Println("wtf", ok, x)
		return
	}
	sd := x.Data.(*xferSendData)

	if length == 0 {
		sd.fp.Close()
		log.Println("finished")
		x.SetBytesSent(position)
		x.UpdateProgress()
		x.SetCompleted(true)
		// x.End() // end is cancel
		delete(xferSendMap, file_number)
		return
	}

	buf := make([]byte, length)
	n, err := sd.fp.ReadAt(buf, int64(position))
	if err != nil {
		log.Println(err)
		x.XferError(x.GetType(), ac, sd.who, err.Error())
	} else if n != length {
		log.Println(n, length)
	} else {
		_, err := this._tox.FileSendChunk(friend_number, sd.fileNumber, position, buf)
		if err != nil {
			log.Println(err)
			x.XferError(x.GetType(), ac, sd.who, err.Error())
		} else {
			x.SetBytesSent(position)
			x.UpdateProgress()
		}
	}

}

// purples
func (this *ToxPlugin) CanReceiveFile(gc *purple.Connection, who string) bool {
	log.Println(who)
	if true {
		return true // let's ui show send file, but is friend offline, just show error
	}
	friendNumber, err := this._tox.FriendByPublicKey(who)
	if err != nil {
		log.Println(err)
	} else {
		st, err := this._tox.FriendGetConnectionStatus(friendNumber)
		if err != nil {
			log.Println(err)
		} else if st > tox.CONNECTION_NONE {
			return true
		}
	}
	return false
}

func (this *ToxPlugin) SendFile(gc *purple.Connection, who string, filename string) {
	log.Println(who, filename)
	ac := gc.ConnGetAccount()
	xfer := ac.NewXfer(purple.XFER_SEND, who)
	this.xferSetCallbacks(xfer)
	xfer.Data = &xferSendData{}
	(xfer.Data.(*xferSendData)).friendNumber, _ = this._tox.FriendByPublicKey(who)

	/* Perform the request */
	if len(filename) > 0 {
		xfer.RequestAccepted(filename)
	} else {
		xfer.Request()
	}
}

func (this *ToxPlugin) NewXfer(gc *purple.Connection, who string) *purple.Xfer {
	log.Println(who)
	return nil
}

// xfer funcs
func (this *ToxPlugin) xferSetCallbacks(x *purple.Xfer) {
	x.ReadFunc = this.xferRead
	x.WriteFunc = this.xferWrite
	x.AckFunc = this.xferAck
	x.RequestDeniedFunc = this.xferRequestDenied
	x.InitFunc = this.xferInit
	x.UninitFunc = this.xferUninit
	x.StartFunc = this.xferStart
	x.EndFunc = this.xferEnd
	x.CancelSendFunc = this.xferCancelSend
	x.CancelRecvFunc = this.xferCancelRecv
}

var xferSendMap = make(map[uint32]*purple.Xfer)

func (this *ToxPlugin) xferInit(x *purple.Xfer) {
	log.Println(x.GetType(), x.GetFilename(), x.GetLocalFilename(), x.GetSize(), x)
	switch x.GetType() {
	case purple.XFER_SEND:
		this.xferInitSend(x)
	case purple.XFER_RECEIVE:
		this.xferInitRecv(x)
	default:
		log.Println("unknown:", x.GetType())
	}
}
func (this *ToxPlugin) xferInitSend(x *purple.Xfer) {
	sd := x.Data.(*xferSendData)
	fi, err := os.Stat(x.GetLocalFilename())
	if err != nil {
		log.Println(err)
		return
	}
	filename := fi.Name()
	log.Println(filename, x.GetFilename(), err)

	fileNumber, err := this._tox.FileSend(sd.friendNumber, tox.FILE_KIND_DATA,
		x.GetSize(), "hehhe", filename)
	if err != nil {
		log.Println(err)
	} else {
		sd.fileNumber = fileNumber
		xferSendMap[fileNumber] = x
	}
}
func (this *ToxPlugin) xferInitRecv(x *purple.Xfer) {
	sd := x.Data.(*xferSendData)
	if sd.fp != nil {
		log.Println("already opened, skiop:", sd.fp)
	} else {
		var err error
		sd.fp, err = os.OpenFile(x.GetLocalFilename(), os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Println(err)
		}
	}
	_, err := this._tox.FileControl(sd.friendNumber, sd.fileNumber, tox.FILE_CONTROL_RESUME)
	if err != nil {
		log.Println(err, sd)
	}

	x.Start(-1, "", 0)
}
func (this *ToxPlugin) xferUninit(x *purple.Xfer) {
	log.Println(x)
}
func (this *ToxPlugin) xferRead(x *purple.Xfer) []byte {
	log.Println(x)
	return nil
}
func (this *ToxPlugin) xferWrite(buf []byte, x *purple.Xfer) int {
	log.Println(x)
	return 0
}
func (this *ToxPlugin) xferAck(x *purple.Xfer, buf []byte) {
	log.Println(x)
}
func (this *ToxPlugin) xferRequestDenied(x *purple.Xfer) {
	log.Println(x)
	sd := x.Data.(*xferSendData)

	var err error
	_, err = this._tox.FileControl(sd.friendNumber, sd.fileNumber, tox.FILE_CONTROL_CANCEL)
	if err != nil {
		log.Println(err)
	}
	delete(xferSendMap, sd.fileNumber)
}
func (this *ToxPlugin) xferStart(x *purple.Xfer) {
	log.Println(x)
}
func (this *ToxPlugin) xferEnd(x *purple.Xfer) {
	log.Println(x)
}
func (this *ToxPlugin) xferCancelSend(x *purple.Xfer) {
	log.Println(x)
	sd := x.Data.(*xferSendData)
	_, err := this._tox.FileControl(sd.friendNumber, sd.fileNumber, tox.FILE_CONTROL_CANCEL)
	if err != nil {
		log.Println(err)
	}
	delete(xferSendMap, sd.fileNumber)
}
func (this *ToxPlugin) xferCancelRecv(x *purple.Xfer) {
	log.Println(x)

	sd := x.Data.(*xferSendData)
	_, err := this._tox.FileControl(sd.friendNumber, sd.fileNumber, tox.FILE_CONTROL_CANCEL)
	if err != nil {
		log.Println(err)
	}
	delete(xferSendMap, sd.fileNumber)
	sd.fp.Close()
	err = os.Remove(x.GetLocalFilename())
	if err != nil {
		log.Println(err)
	}
}
