
TOR_PT_MANAGED_TRANSPORT_VER=1 \
TOR_PT_STATE_LOCATION=/tmp \
TOR_PT_SERVER_TRANSPORTS=sam \
TOR_PT_SERVER_BINDADDR=sam-0.0.0.0:1231 \
TOR_PT_ORPORT=127.0.0.1:1080 ./samserver -client-config sam.torrc -i2p-keys sam.torrc.i2pkeys
