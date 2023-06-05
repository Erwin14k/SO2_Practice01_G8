obj-m += cpu_202000119.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules  

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean

EXTRA_CFLAGS := -msoft-float -msse -msse2 -msse3 -msse4