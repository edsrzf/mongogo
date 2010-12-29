# Copyright 2010, Evan Shaw.  All rights reserved.
# Use of this source code is governed by a BSD-style License
# that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=github.com/edsrzf/mongogo
GOFILES=\
	collection.go\
	conn.go\
	cursor.go\
	database.go\
	error.go\
	query.go\

include $(GOROOT)/src/Make.pkg

