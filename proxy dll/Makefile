CC      := x86_64-w64-mingw32-gcc
CFLAGS  := -O2 -Wall -Wextra -Wno-unused-parameter
LDFLAGS := -shared -s -static-libgcc -Wl,--subsystem,windows,--enable-stdcall-fixup

PATCH_DIR ?= .
TARGET    := $(PATCH_DIR)/version.dll
PATCH_H   := $(PATCH_DIR)/patches.h

all: $(TARGET)

$(TARGET): version.c version.def $(PATCH_H)
	$(CC) $(CFLAGS) -I$(PATCH_DIR) $(LDFLAGS) -o $@ version.c version.def

clean:
	rm -f $(TARGET)

.PHONY: all clean
