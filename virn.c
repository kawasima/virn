#include<stdio.h>
#include<stdlib.h>
#include<curses.h>

typedef __off64_t off64_t;

#include<apr_general.h>
#include<apr_strings.h>
#include<apr_file_io.h>
#include<apr_file_info.h>

/*
 * 元々、APRのライブコーディングをやるつもりだったので
 * 無駄にAPRでかかれている。
 */
int main(int argc, char *argv[])
{
    if (argc <= 1) {
        printf("%s filename", argv[0]);
        return 1;
    }
    apr_pool_t *pool = NULL;
    apr_file_t *file = NULL;
    char buf[1024];
    const char *root_path, *home_dir = NULL;
    const char *path = NULL;
    char *virnfile_path = NULL;
    apr_size_t n;

    apr_initialize();
    apr_pool_create(&pool, NULL);
    
    path = apr_pstrdup(pool, argv[1]);
    apr_filepath_root(&root_path, &path, APR_FILEPATH_NATIVE, pool);
    apr_env_get(&home_dir, "HOME", pool);
    virnfile_path = apr_pstrcat(pool, home_dir, "/.virn/", path, NULL);

    if (apr_file_open(&file, virnfile_path, APR_READ, APR_REG, pool) != APR_SUCCESS) {
        fprintf(stderr, "can't open file %s\n", argv[1]);
        exit(1);
    }
    apr_file_read_full(file, buf, sizeof(buf), &n);
    apr_file_close(file);

    initscr();
    cbreak();
    noecho();
    nl();
    scrollok(stdscr, TRUE);
    keypad(stdscr, TRUE);

    int x, y;
    getmaxyx(stdscr, y, x);
    wresize(stdscr, y-1, x);
    WINDOW *status_win;
    status_win = newwin(1, x, y-1, 0);
    clearok(status_win, true);
    wprintw(status_win, "\"%s\" [New File]", argv[1]);
    wrefresh(status_win);
    getch();
    wdeleteln(status_win);
    wmove(status_win, 0,0);
    wattron(status_win, A_BOLD);
    waddstr(status_win, "--INSERT--");
    wrefresh(status_win);


    int i=0, c;
    while (1) {
        c = getch();
        if (buf[i] == '\n') {
            if (c == '\n') {
                char *indent;
                int j;
                for(j = i; isspace(buf[j]) && buf[j] != '\0'; j++);
                indent = apr_pstrndup(pool, buf+i, j-i);
                addstr(indent);
                i = j;
            }
        } else {
            addch(buf[i++]);
        }
        if (buf[i] == '\0') {
            break;
        }
    }
    getch(); getch();
    endwin();

    apr_file_copy(virnfile_path, argv[1], APR_FILE_SOURCE_PERMS, pool);
    apr_terminate();
}
