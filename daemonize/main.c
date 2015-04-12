#include <sys/types.h>
#include <sys/stat.h>
#include <stdio.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
    if (argc < 2) {
        return 1;
    }

    pid_t pid = fork();
    if (pid < 0) {
        return 2;
    } else if (pid > 0) {
        // Parent
        printf("%d", pid);
        return 0;
    } else {
        // Child
        
        umask(0);
        pid_t sid = setsid();
        if (sid == 0) {
            return 4;
        }
        chdir("/");
        close(STDIN_FILENO);
        close(STDOUT_FILENO);
        close(STDERR_FILENO);

        // Setup the array of args to pass to the new process.
        char *newArgs[argc];
        int i;
        for (i = 0; i < argc-1; i++) {
            newArgs[i] = argv[i+1];
        }
        newArgs[argc-1] = NULL;

        // Start the process.
        execv(argv[1], newArgs);

        // This should never be reached.
        return 3;
    }
}
