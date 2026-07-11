ARCH    ?= x64
PROXY   ?= version

ifeq ($(ARCH),x86)
CC      := i686-w64-mingw32-gcc
else
CC      := x86_64-w64-mingw32-gcc
endif

CFLAGS  := -O2 -Wall -Wextra -Wno-unused-parameter
LDFLAGS := -shared -s -static-libgcc -Wl,--subsystem,windows,--enable-stdcall-fixup

ifeq ($(PROXY),winmm)
CFLAGS  += -DLUCKPROXY_WINMM
endif

PATCH_DIR ?= .
TARGET    := $(PATCH_DIR)/$(PROXY).dll
PATCH_H   := $(PATCH_DIR)/patches.h
DEF_FILE  := $(PROXY).def

all: $(TARGET)

$(TARGET): version.c $(DEF_FILE) $(PATCH_H)
	$(CC) $(CFLAGS) -I$(PATCH_DIR) $(LDFLAGS) -o $@ version.c $(DEF_FILE)

clean:
	rm -f $(TARGET)

.PHONY: all clean
