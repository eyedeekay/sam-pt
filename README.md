sam-pt
======

It's a Tor pluggable transport over I2P's SAM API because YOLO.

Client Configuration
--------------------

To set up a client, add the following to your torrc

        UseBridges 1
        Bridge sam <base64 address, base32 address, or i2p domain in address book>

        ClientTransportPlugin sam exec /usr/bin/samclient <base64 address, base32 address, or i2p domain in address book>

Server Configuration
--------------------

To set up a server, add the following to your torrc

        BridgeRelay 1
        ORPort 9001
        ExtORPort 9002

        ServerTransportPlugin sam exec /usr/bin/samserver <optional path to client torrc fragment>
        ServerTransportListenAddr sam

