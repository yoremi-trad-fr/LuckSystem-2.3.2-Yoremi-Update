CC      := x86_64-w64-mingw32-gcc
CFLAGS  := -O2 -Wall -Wextra -Wno-unused-parameter
LDFLAGS := -shared -s -static-libgcc -Wl,--subsystem,windows,--enable-stdcall-fixup

TARGET  := version.dll

all: $(TARGET)

$(TARGET): version.c version.def
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ version.c version.def

clean:
	rm -f $(TARGET)

.PHONY: all clean
