#define _GNU_SOURCE
#include <errno.h>
#include <fcntl.h>
#include <linux/input.h>
#include <linux/uinput.h>
#include <poll.h>
#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>

_Static_assert(sizeof(void*)==4 && sizeof(struct input_event)==16,
               "This binary requires the stock ARM32 evdev ABI");

/* Only the GA36 physical controller is read. It is never grabbed: Android
 * continues to receive its gamepad and D-pad events. Right stick drives a
 * separate, conventional relative mouse; R1 holds its left mouse button. */
static volatile sig_atomic_t stop;
static void stopped(int sig) { (void)sig; stop=1; }
static long long milliseconds(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC,&ts);
    return (long long)ts.tv_sec*1000+ts.tv_nsec/1000000;
}
static int velocity(int value, int minimum, int maximum) {
    int span=maximum-minimum;
    if (span<4 || value<minimum || value>maximum) return 0;
    int offset=2*value-minimum-maximum;
    int sign=offset<0?-1:1;
    int distance=abs(offset), dead=span/5;
    if (distance<=dead) return 0;
    return sign*(1+(distance-dead)*7/(span-dead));
}
static int emit(int fd, unsigned short type, unsigned short code, int value) {
    struct input_event e;
    memset(&e,0,sizeof(e)); e.type=type; e.code=code; e.value=value;
    return write(fd,&e,sizeof(e))==(ssize_t)sizeof(e)?0:-1;
}
static int controller(struct input_absinfo *x,struct input_absinfo *y,char *path) {
    for (int i=0;i<32;i++) {
        snprintf(path,64,"/dev/input/event%d",i);
        int fd=open(path,O_RDONLY|O_NONBLOCK|O_CLOEXEC);
        if(fd<0) continue;
        char name[128]={0}; struct stat st;
        if(fstat(fd,&st)==0 && S_ISCHR(st.st_mode) &&
           ioctl(fd,EVIOCGNAME(sizeof(name)-1),name)>=0 &&
           strcmp(name,"micro_gamepad")==0 &&
           ioctl(fd,EVIOCGABS(ABS_RX),x)==0 && ioctl(fd,EVIOCGABS(ABS_RY),y)==0 &&
           x->maximum-x->minimum>=4 && x->maximum-x->minimum<=65535 &&
           y->maximum-y->minimum>=4 && y->maximum-y->minimum<=65535) return fd;
        close(fd);
    }
    return -1;
}
static int mouse(void) {
    int fd=open("/dev/uinput",O_WRONLY|O_NONBLOCK|O_CLOEXEC);
    if(fd<0) fd=open("/dev/input/uinput",O_WRONLY|O_NONBLOCK|O_CLOEXEC);
    if(fd<0) return -1;
    if(ioctl(fd,UI_SET_EVBIT,EV_SYN)<0 || ioctl(fd,UI_SET_EVBIT,EV_KEY)<0 ||
       ioctl(fd,UI_SET_EVBIT,EV_REL)<0 || ioctl(fd,UI_SET_KEYBIT,BTN_LEFT)<0 ||
       ioctl(fd,UI_SET_RELBIT,REL_X)<0 || ioctl(fd,UI_SET_RELBIT,REL_Y)<0) {
        close(fd); return -1;
    }
    struct uinput_user_dev dev;
    memset(&dev,0,sizeof(dev));
    snprintf(dev.name,sizeof(dev.name),"GA36 UI Pointer");
    dev.id.bustype=BUS_VIRTUAL; dev.id.vendor=0x1; dev.id.product=0x36;
    dev.id.version=1;
    if(write(fd,&dev,sizeof(dev))!=(ssize_t)sizeof(dev) || ioctl(fd,UI_DEV_CREATE)<0) {
        close(fd); return -1;
    }
    return fd;
}
static int selftest(void) {
    for(int v=0;v<=18;v++) {
        if(velocity(v,0,18)!=-velocity(18-v,0,18)) return 1;
        if(abs(velocity(v,0,18))>8) return 2;
    }
    if(velocity(9,0,18)||velocity(8,0,18)||velocity(10,0,18)||
       velocity(0,0,18)!=-8||velocity(18,0,18)!=8||velocity(-1,0,18)) return 3;
    int p[2]; if(pipe(p)<0) return 4;
    if(emit(p[1],EV_KEY,BTN_LEFT,1)<0||emit(p[1],EV_SYN,SYN_REPORT,0)<0) return 5;
    struct input_event e[2];
    if(read(p[0],e,sizeof(e))!=(ssize_t)sizeof(e)||e[0].type!=EV_KEY||
       e[0].code!=BTN_LEFT||e[0].value!=1||e[1].type!=EV_SYN) return 6;
    close(p[0]);close(p[1]);
    printf("GA36_POINTER_SELFTEST_OK pointer_bytes=%zu input_event_bytes=%zu right_stick_deadzone_and_event_encoding=passed\n",sizeof(void*),sizeof(struct input_event));
    return 0;
}
int main(int argc,char **argv) {
    if(argc==2 && strcmp(argv[1],"--selftest")==0) return selftest();
    if(argc!=2 || strcmp(argv[1],"--run")!=0) return 64;
    FILE *log=fopen("/ga36logs/input-pointer.txt","a");
    if(!log) return 1;
    setvbuf(log,NULL,_IOLBF,0);
    signal(SIGTERM,stopped); signal(SIGINT,stopped);
    struct input_absinfo x,y; char path[64]; int input=-1;
    for(int n=0;n<100 && !stop;n++) {
        input=controller(&x,&y,path);
        if(input>=0) break;
        usleep(100000);
    }
    if(input<0) {fprintf(log,"ERROR physical micro_gamepad not available\n");fclose(log);return 2;}
    int output=mouse();
    if(output<0) {fprintf(log,"ERROR uinput mouse setup: %s\n",strerror(errno));close(input);fclose(log);return 3;}
    fprintf(log,"READY physical=%s virtual=GA36 UI Pointer right_axes=%d..%d,%d..%d R1=left_click\n",path,x.minimum,x.maximum,y.minimum,y.maximum);
    int vx=x.value,vy=y.value,clicked=0,records=0,frames=0;
    long long next=milliseconds()+20;
    struct pollfd pollfd={.fd=input,.events=POLLIN};
    while(!stop) {
        long long now=milliseconds();
        int wait=next>now?(int)(next-now):0;
        int ready=poll(&pollfd,1,wait);
        if(ready<0 && errno!=EINTR) break;
        if(pollfd.revents&(POLLERR|POLLHUP|POLLNVAL)) break;
        if(ready>0 && (pollfd.revents&POLLIN)) {
            struct input_event e[64]; ssize_t bytes;
            while((bytes=read(input,e,sizeof(e)))>0) {
                if(bytes%sizeof(e[0])) {stop=1;break;}
                for(size_t i=0;i<(size_t)bytes/sizeof(e[0]);i++) {
                    if(e[i].type==EV_SYN && e[i].code==SYN_DROPPED) {
                        if(ioctl(input,EVIOCGABS(ABS_RX),&x)==0) vx=x.value;
                        if(ioctl(input,EVIOCGABS(ABS_RY),&y)==0) vy=y.value;
                        emit(output,EV_KEY,BTN_LEFT,0);emit(output,EV_SYN,SYN_REPORT,0);clicked=0;
                    }
                    if(e[i].type==EV_ABS && e[i].code==ABS_RX) vx=e[i].value;
                    if(e[i].type==EV_ABS && e[i].code==ABS_RY) vy=e[i].value;
                    if(e[i].type==EV_KEY) {
                        if(records++<1024) fprintf(log,"KEY code=%u value=%d\n",e[i].code,e[i].value);
                        if(e[i].code==BTN_TR && e[i].value!=2) {
                            clicked=e[i].value!=0;
                            if(emit(output,EV_KEY,BTN_LEFT,clicked)<0||emit(output,EV_SYN,SYN_REPORT,0)<0) stop=1;
                        }
                    }
                }
            }
        }
        now=milliseconds();
        if(now>=next) {
            int dx=velocity(vx,x.minimum,x.maximum),dy=velocity(vy,y.minimum,y.maximum);
            if(dx||dy) {
                if(emit(output,EV_REL,REL_X,dx)<0||emit(output,EV_REL,REL_Y,dy)<0||emit(output,EV_SYN,SYN_REPORT,0)<0) break;
                if(frames++%100==0 && records++<1024) fprintf(log,"POINTER right_axes=%d,%d relative=%d,%d\n",vx,vy,dx,dy);
            }
            next=now+20;
        }
    }
    if(clicked) {emit(output,EV_KEY,BTN_LEFT,0);emit(output,EV_SYN,SYN_REPORT,0);}
    ioctl(output,UI_DEV_DESTROY);close(output);close(input);
    fprintf(log,"STOP frames=%d key_records=%d\n",frames,records);fclose(log);
    return 0;
}
