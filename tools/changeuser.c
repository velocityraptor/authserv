#include <stdio.h>
#include <stdlib.h>

void
usage()
{
	printf("Velocityraptor Authserver\n");
	printf("authserv-changeuser [username]\n");
}

int
main(int argc, char *argv[])
{
	if(argc < 2) {
		usage();
		return 0;
	}
}

