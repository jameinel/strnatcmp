#include <stdio.h>
#include "strnatcmp.h"

int main(int argc, char **argv) {
	int res;
	if (argc != 3) {
		fprintf(stderr, "need 2 arguments\n");
	}
	res = strnatcmp(argv[1], argv[2]);
	fprintf(stdout, "Result: %d\n", res);
	strnatprintcalls();
}
