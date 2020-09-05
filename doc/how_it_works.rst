# Components



## Transport

* a generic interface for SIPProvider - receives and sends messages
sending messages all transport layers include the WriteTo function

## Director

sorts incoming messages and sends them through channels to registered goroutinges

## cleaner

periodically sends nil to all registered channels, and if blocked removes them
