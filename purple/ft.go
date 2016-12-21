package purple

/*
#include <libpurple/purple.h>

extern gssize gopurple_xfer_read_fnc(guchar** buf, PurpleXfer *xfer);
extern gssize gopurple_xfer_write_fnc(guchar *buf, size_t size, PurpleXfer *xfer);
extern void gopurple_xfer_ack_fnc(PurpleXfer *xfer, guchar *buf, size_t size);
extern void gopurple_xfer_request_denied_fnc(PurpleXfer *xfer);
extern void gopurple_xfer_init_fnc(PurpleXfer *xfer);
extern void gopurple_xfer_uninit_fnc(PurpleXfer *xfer);
extern void gopurple_xfer_start_fnc(PurpleXfer *xfer);
extern void gopurple_xfer_end_fnc(PurpleXfer *xfer);
extern void gopurple_xfer_cancel_send_fnc(PurpleXfer *xfer);
extern void gopurple_xfer_cancel_recv_fnc(PurpleXfer *xfer);

*/
import "C"

import "unsafe"

import (
	"log"
)

type XferType int

const (
	XFER_UNKNOWN = XferType(C.PURPLE_XFER_UNKNOWN)
	XFER_SEND    = XferType(C.PURPLE_XFER_SEND)
	XFER_RECEIVE = XferType(C.PURPLE_XFER_RECEIVE)
)

type XferStatusType int

const (
	XFER_STATUS_UNKNOWN       = int(C.PURPLE_XFER_STATUS_UNKNOWN)
	XFER_STATUS_NOT_STARTED   = int(C.PURPLE_XFER_STATUS_NOT_STARTED)
	XFER_STATUS_ACCEPTED      = int(C.PURPLE_XFER_STATUS_ACCEPTED)
	XFER_STATUS_STARTED       = int(C.PURPLE_XFER_STATUS_STARTED)
	XFER_STATUS_DONE          = int(C.PURPLE_XFER_STATUS_DONE)
	XFER_STATUS_CANCEL_LOCAL  = int(C.PURPLE_XFER_STATUS_CANCEL_LOCAL)
	XFER_STATUS_CANCEL_REMOTE = int(C.PURPLE_XFER_STATUS_CANCEL_REMOTE)
)

type XferUiOps struct {

	// private
	fromc bool
	xuo   *C.PurpleXferUiOps
}

func newXferUiOpsFrom(xuo *C.PurpleXferUiOps) *XferUiOps {
	this := &XferUiOps{}
	this.xuo = xuo
	return this
}

type Xfer struct {
	ReadFunc          func(*Xfer) []byte
	WriteFunc         func([]byte, *Xfer) int
	AckFunc           func(*Xfer, []byte)
	RequestDeniedFunc func(*Xfer)
	InitFunc          func(*Xfer)
	UninitFunc        func(*Xfer)
	StartFunc         func(*Xfer)
	EndFunc           func(*Xfer)
	CancelSendFunc    func(*Xfer)
	CancelRecvFunc    func(*Xfer)

	Data interface{}

	// private
	fromc bool
	xfer  *C.PurpleXfer
}

// TODO destructure
func (this *Account) NewXfer(xftype XferType, who string) *Xfer {
	xfer := &Xfer{}
	xfer.xfer = C.purple_xfer_new(this.account, C.PurpleXferType(xftype), CCString(who).Ptr)

	xfer.setCallbakcFuncs()
	xfer.addMap()
	return xfer
}

func (this *Xfer) setCallbakcFuncs() {
	xferc := this.xfer

	C.purple_xfer_set_read_fnc(xferc, go2cfn(C.gopurple_xfer_read_fnc))
	C.purple_xfer_set_write_fnc(xferc, go2cfn(C.gopurple_xfer_write_fnc))
	C.purple_xfer_set_ack_fnc(xferc, go2cfn(C.gopurple_xfer_ack_fnc))
	C.purple_xfer_set_request_denied_fnc(xferc, go2cfn(C.gopurple_xfer_request_denied_fnc))
	C.purple_xfer_set_init_fnc(xferc, go2cfn(C.gopurple_xfer_init_fnc))
	C.purple_xfer_set_start_fnc(xferc, go2cfn(C.gopurple_xfer_start_fnc))
	C.purple_xfer_set_end_fnc(xferc, go2cfn(C.gopurple_xfer_end_fnc))
	C.purple_xfer_set_cancel_send_fnc(xferc, go2cfn(C.gopurple_xfer_cancel_send_fnc))
	C.purple_xfer_set_cancel_recv_fnc(xferc, go2cfn(C.gopurple_xfer_cancel_recv_fnc))

}

func XfersGetAll() *GList {
	r := C.purple_xfers_get_all()
	return newGListFrom(r)
}

func (this *Xfer) Ref()   { C.purple_xfer_ref(this.xfer) }
func (this *Xfer) Unref() { C.purple_xfer_unref(this.xfer) }

func (this *Xfer) RequestAccepted(filename string) {
	C.purple_xfer_request_accepted(this.xfer, CCString(filename).Ptr)
}
func (this *Xfer) Request() {
	C.purple_xfer_request(this.xfer)
}

func (this *Xfer) RequestDenied() { C.purple_xfer_request_denied(this.xfer) }
func (this *Xfer) GetType() XferType {
	return XferType(C.purple_xfer_get_type(this.xfer))
}
func (this *Xfer) GetAccount() *Account {
	return newAccountFrom(C.purple_xfer_get_account(this.xfer))
}
func (this *Xfer) GetRemoteUser() string {
	r := C.purple_xfer_get_remote_user(this.xfer)
	return C.GoString(r)
}
func (this *Xfer) GetStatus() XferStatusType {
	r := C.purple_xfer_get_status(this.xfer)
	return XferStatusType(r)
}
func (this *Xfer) IsCanceled() bool {
	r := C.purple_xfer_is_canceled(this.xfer)
	return c2goBool(r)
}
func (this *Xfer) IsCompleted() bool {
	r := C.purple_xfer_is_completed(this.xfer)
	return c2goBool(r)
}
func (this *Xfer) GetFilename() string {
	r := C.purple_xfer_get_filename(this.xfer)
	return C.GoString(r)
}
func (this *Xfer) GetLocalFilename() string {
	r := C.purple_xfer_get_local_filename(this.xfer)
	return C.GoString(r)
}
func (this *Xfer) GetBytesSent() uint64 {
	r := C.purple_xfer_get_bytes_sent(this.xfer)
	return uint64(r)
}
func (this *Xfer) GetBytesRemaining() uint64 {
	r := C.purple_xfer_get_bytes_remaining(this.xfer)
	return uint64(r)
}
func (this *Xfer) GetSize() uint64 {
	r := C.purple_xfer_get_size(this.xfer)
	return uint64(r)
}
func (this *Xfer) GetProgress() float64 {
	r := C.purple_xfer_get_progress(this.xfer)
	return float64(r)
}
func (this *Xfer) GetLocalPort() uint {
	r := C.purple_xfer_get_local_port(this.xfer)
	return uint(r)
}

func (this *Xfer) GetRemoteIP() string {
	r := C.purple_xfer_get_remote_ip(this.xfer)
	return C.GoString(r)
}
func (this *Xfer) GetRemotePort() uint {
	r := C.purple_xfer_get_remote_port(this.xfer)
	return uint(r)
}
func (this *Xfer) GetStartTime() uint64 {
	r := C.purple_xfer_get_start_time(this.xfer)
	return uint64(r)
}
func (this *Xfer) GetEndTime() uint64 {
	r := C.purple_xfer_get_end_time(this.xfer)
	return uint64(r)
}
func (this *Xfer) SetCompleted(completed bool) {
	C.purple_xfer_set_completed(this.xfer, go2cBool(completed))
}
func (this *Xfer) SetMessage(message string) {
	C.purple_xfer_set_message(this.xfer, CCString(message).Ptr)
}
func (this *Xfer) SetFilename(filename string) {
	C.purple_xfer_set_filename(this.xfer, CCString(filename).Ptr)
}
func (this *Xfer) SetLocalFilename(filename string) {
	C.purple_xfer_set_local_filename(this.xfer, CCString(filename).Ptr)
}
func (this *Xfer) SetSize(size uint64) {
	C.purple_xfer_set_size(this.xfer, C.size_t(size))
}
func (this *Xfer) SetBytesSent(sent uint64) {
	C.purple_xfer_set_bytes_sent(this.xfer, C.size_t(sent))
}

func (this *Xfer) GetUiOps() *XferUiOps {
	r := C.purple_xfer_get_ui_ops(this.xfer)
	return newXferUiOpsFrom(r)
}

// set funcs

//
func (this *Xfer) Read() []byte {
	var buf *C.guchar
	r := C.purple_xfer_read(this.xfer, &buf)
	return C.GoBytes(unsafe.Pointer(buf), C.int(r))
}

func (this *Xfer) Write(buf []byte) int {
	buf_ := (*C.guchar)((unsafe.Pointer)(&buf[0]))
	r := C.purple_xfer_write(this.xfer, buf_, C.gsize(len(buf)))
	return int(r)
}

func (this *Xfer) WriteFile(buf []byte) bool {
	buf_ := (*C.guchar)((unsafe.Pointer)(&buf[0]))
	r := C.purple_xfer_write_file(this.xfer, buf_, C.gsize(len(buf)))
	return c2goBool(r)
}

func (this *Xfer) ReadFile() []byte {
	buf := make([]byte, 256)
	buf_ := (*C.guchar)((unsafe.Pointer)(&buf[0]))
	r := C.purple_xfer_read_file(this.xfer, buf_, 256)
	return buf[0:int(r)]
}

func (this *Xfer) Start(fd int, ip string, port uint) {
	if len(ip) == 0 {
		C.purple_xfer_start(this.xfer, C.int(fd), nil, C.uint(port))
	} else {
		C.purple_xfer_start(this.xfer, C.int(fd), CCString(ip).Ptr, C.uint(port))
	}
}

func (this *Xfer) End() { C.purple_xfer_end(this.xfer) }
func (this *Xfer) Add() { C.purple_xfer_add(this.xfer) }

func (this *Xfer) CancelLocal() {
	C.purple_xfer_cancel_local(this.xfer)
}

func (this *Xfer) CancelRemote() {
	C.purple_xfer_cancel_remote(this.xfer)
}

func XferError(xftype XferType, ac *Account, who string, msg string) {
	C.purple_xfer_error(C.PurpleXferType(xftype), ac.account,
		CCString(who).Ptr, CCString(msg).Ptr)
}

func (this *Xfer) XferError(xftype XferType, ac *Account, who string, msg string) {
	C.purple_xfer_error(C.PurpleXferType(xftype), ac.account,
		CCString(who).Ptr, CCString(msg).Ptr)
}

func (this *Xfer) UpdateProgress() {
	C.purple_xfer_update_progress(this.xfer)
}

func (this *Xfer) ConversationWrite(message string, isError bool) {
	C.purple_xfer_conversation_write(this.xfer, CCString(message).Ptr, go2cBool(isError))
}

func (this *Xfer) UiReady() {
	C.purple_xfer_ui_ready(this.xfer)
}

func (this *Xfer) PrplReady() {
	C.purple_xfer_prpl_ready(this.xfer)
}

func (this *Xfer) GetThumbnail() []byte {
	var ln C.gsize
	r := C.purple_xfer_get_thumbnail(this.xfer, &ln)
	if r == nil {
		return nil
	}
	r_ := (unsafe.Pointer)((*C.char)(r))
	return C.GoBytes(r_, C.int(ln))
}

func (this *Xfer) GetThumbnailMimeype() string {
	r := C.purple_xfer_get_thumbnail_mimetype(this.xfer)
	r_ := (*C.char)(r)
	return C.GoString(r_)
}

func (this *Xfer) SetThumbnail(thumbnail []byte, size uint64, mimeType string) {
	thumbnail_ := (unsafe.Pointer)(&thumbnail[0])
	C.purple_xfer_set_thumbnail(this.xfer, thumbnail_, C.gsize(size),
		(*C.gchar)(CCString(mimeType).Ptr))
}

func (this *Xfer) PrepareThumbnail(formats string) {
	C.purple_xfer_prepare_thumbnail(this.xfer, (*C.gchar)(CCString(formats).Ptr))
}

//
func XfersGetHandle() unsafe.Pointer { return C.purple_xfers_get_handle() }
func XfersInit()                     { C.purple_xfers_init() }
func XfersUninit()                   { C.purple_xfers_uninit() }
func XfersSetUiOps(ops *XferUiOps)   { C.purple_xfers_set_ui_ops(ops.xuo) }
func XfersGetUiOps() *XferUiOps      { return newXferUiOpsFrom(C.purple_xfers_get_ui_ops()) }

//
var xfer_cgo_map = make(map[*C.PurpleXfer]*Xfer, 0)

func (this *Xfer) addMap() {
	xfer_cgo_map[this.xfer] = this
}
func (this *Xfer) delMap() {
	delete(xfer_cgo_map, this.xfer)
}

//export gopurple_xfer_read_fnc
func gopurple_xfer_read_fnc(buf **C.guchar, xfer *C.PurpleXfer) C.gssize {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.ReadFunc != nil {
			r := x.ReadFunc(x)
			*buf = (*C.guchar)(C.calloc(C.size_t(len(r)), 1))
			C.memcpy((unsafe.Pointer)(*buf), (unsafe.Pointer)(&r[0]), C.size_t(len(r)))
			return C.gssize(len(r))
		}
	} else {
		log.Println("not found")
	}
	return 0
}

//export gopurple_xfer_write_fnc
func gopurple_xfer_write_fnc(buf *C.guchar, ln C.size_t, xfer *C.PurpleXfer) C.gssize {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.WriteFunc != nil {
			buf_ := C.GoBytes((unsafe.Pointer)(buf), C.int(ln))
			r := x.WriteFunc(buf_, x)
			return C.gssize(r)
		}
	} else {
		log.Println("not found")
	}
	return 0
}

//export gopurple_xfer_ack_fnc
func gopurple_xfer_ack_fnc(xfer *C.PurpleXfer, buf *C.guchar, size C.size_t) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.AckFunc != nil {
			x.AckFunc(x, C.GoBytes((unsafe.Pointer)(buf), C.int(size)))
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_request_denied_fnc
func gopurple_xfer_request_denied_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.RequestDeniedFunc != nil {
			x.RequestDeniedFunc(x)
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_init_fnc
func gopurple_xfer_init_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.InitFunc != nil {
			x.InitFunc(x)
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_uninit_fnc
func gopurple_xfer_uninit_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.UninitFunc != nil {
			x.UninitFunc(x)
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_start_fnc
func gopurple_xfer_start_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.StartFunc != nil {
			x.StartFunc(x)
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_end_fnc
func gopurple_xfer_end_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.EndFunc != nil {
			x.EndFunc(x)
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_cancel_send_fnc
func gopurple_xfer_cancel_send_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.CancelSendFunc != nil {
			x.CancelSendFunc(x)
		}
	} else {
		log.Println("not found")
	}
}

//export gopurple_xfer_cancel_recv_fnc
func gopurple_xfer_cancel_recv_fnc(xfer *C.PurpleXfer) {
	if x, ok := xfer_cgo_map[xfer]; ok {
		if x.CancelRecvFunc != nil {
			x.CancelRecvFunc(x)
		}
	} else {
		log.Println("not found")
	}
}
