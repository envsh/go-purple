import sys, os
import signal
import asyncio


from PyQt5.QtCore import *
from PyQt5.QtDBus import *
from quamash import *


app = QCoreApplication(sys.argv)
loop = QEventLoop(app)
asyncio.set_event_loop(loop)



def sigint_handler():
    app.exit()
    sys.exit()
    return
signal.signal(signal.SIGINT, sigint_handler)  # no effect
app.aboutToQuit.connect(sigint_handler, Qt.QueuedConnection)  # no effect
loop.add_signal_handler(signal.SIGINT, sigint_handler)


service = 'im.pidgin.gopurple.PurpleService'
path = '/im/pidgin/gopurple/PurpleObject'
iface = 'im.pidgin.gopurple.PurpleInterface'

sesbus = QDBusConnection.sessionBus()
oiface = QDBusInterface(service, path, iface, sesbus)

class Controller(QObject):
    def __init__(self):
        super(Controller, self).__init__()
        return

    @pyqtSlot(QDBusMessage)
    def dbus_message_handler(self, msg):
        print(msg, msg.arguments(), msg.member(), msg.signature(), msg.type())
        args = msg.arguments()
        if msg.member() == 'ReceivedChatMsg':
            proto = oiface.call('PurpleAccountGetProtocolName', args[0])
            print('protocol name:', proto.arguments())
        return

ctrl = Controller()
if qVersion() >= '5.6' and qVersion() <= '5.7.9':
    sesbus.registerObject('/hotfixidontknowwhy_' + ctrl.__class__.__name__, ctrl)

bret = sesbus.connect(service, path, iface, '', ctrl.dbus_message_handler)
# bret = sesbus.connect(service, path, iface, 'ReceivedChatMsg', ctrl.dbus_message_handler)
qDebug('connected server message bus: {}'.format(bret))


qDebug('looping...')
with loop:
    loop.run_forever()


