include $(GOROOT)/src/Make.inc

TARG=htmlsoap
GOFILES=\
	suds.go\
	parser.go\
	node.go\
	command.go

include $(GOROOT)/src/Make.cmd