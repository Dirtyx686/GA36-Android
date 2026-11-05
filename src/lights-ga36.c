/* Android 5 ARM32 lights ABI, from AOSP android-5.0.2_r1 hardware.h/lights.h.
 * GA36 backlight path confirmed in the device's stock emuelec_autostart.sh.
 */
typedef unsigned int u32; typedef unsigned short u16;
struct module; struct device;
struct methods {int (*open)(const struct module*,const char*,struct device**);};
struct module {u32 tag;u16 api,hal;const char *id,*name,*author;struct methods *methods;void *dso;u32 reserved[25];};
struct device {u32 tag,version;struct module *module;u32 reserved[12];int (*close)(struct device*);};
struct state {unsigned color;int flashMode,flashOnMS,flashOffMS,brightnessMode;};
struct light {struct device common;int (*set_light)(struct light*,const struct state*);};
_Static_assert(sizeof(struct module)==128,"Android ARM32 module ABI");
_Static_assert(sizeof(struct device)==64,"Android ARM32 device ABI");
_Static_assert(sizeof(struct light)==68,"Android ARM32 light ABI");
extern void *calloc(unsigned,unsigned);extern void free(void*);extern int strcmp(const char*,const char*);
extern int snprintf(char*,unsigned,const char*,...);extern int open(const char*,int,...);extern int read(int,void*,unsigned);extern int write(int,const void*,unsigned);extern int close(int);extern int *__errno(void);extern int usleep(unsigned);
struct timespec32 {long sec,nsec;};
extern int clock_gettime(int,struct timespec32*);
extern long lseek(int,long,int);extern int ftruncate(int,long);
static unsigned generation,write_lock;
static int brightness_now(void){
 char buf[16];int fd=open("/sys/devices/virtual/disp/disp/attr/lcdbl",0|0x80000);
 if(fd<0)return -1;int n=read(fd,buf,sizeof(buf));close(fd);
 if(n<=0||buf[0]<'0'||buf[0]>'9')return -1;
 unsigned v=0;for(int i=0;i<n&&buf[i]>='0'&&buf[i]<='9';i++){v=v*10+(unsigned)(buf[i]-'0');if(v>255)return -1;}
 return (int)v;
}
static void trace(unsigned requested,int before,int delayed,int result){
 int fd=open("/dev/ga36-backlight-hal.txt",1|0x400|0x80000);if(fd<0)return;
 if(lseek(fd,0,2)>65536)ftruncate(fd,0);
 struct timespec32 t={0};clock_gettime(1,&t);
 char buf[128];int n=snprintf(buf,sizeof(buf),"uptime=%ld.%03ld request=%u before=%d resume_delay_ms=%d after=%d result=%d\n",t.sec,t.nsec/1000000,requested,before,delayed,brightness_now(),result);
 if(n>0&&n<(int)sizeof(buf))write(fd,buf,(unsigned)n);close(fd);
}
static int close_light(struct device *d){free(d);return 0;}
static int set_backlight(struct light *dev,const struct state *s){
 (void)dev;if(!s)return -22;
 unsigned c=s->color;unsigned value=(77*((c>>16)&255)+150*((c>>8)&255)+29*(c&255))>>8;
 unsigned ticket=__sync_add_and_fetch(&generation,1);
 int before=brightness_now(),delayed=0;
 /* Stock LCD resume finishes asynchronously around430ms after wake. A
  * brightness write before that can be lost when the panel flow finishes.
  * Delay only a transition from physically dark to a requested nonzero level.
  * Newer requests, especially screen-off, invalidate this pending write. */
 if(value>0&&before==0){usleep(650000);delayed=650;}
 while(__sync_lock_test_and_set(&write_lock,1))usleep(1000);
 if(__sync_add_and_fetch(&generation,0)!=ticket){__sync_lock_release(&write_lock);trace(value,before,delayed,1);return 0;}
 char buf[16];int len=snprintf(buf,sizeof(buf),"%u\n",value);
 int fd=open("/sys/devices/virtual/disp/disp/attr/lcdbl",1|0x80000);
 if(fd<0){int result=-*__errno();__sync_lock_release(&write_lock);return result;}
 int n=write(fd,buf,(unsigned)len);int result=n==len?0:n<0?-*__errno():-5;close(fd);trace(value,before,delayed,result);__sync_lock_release(&write_lock);return result;
}
static int open_light(const struct module *m,const char *id,struct device **out){
 if(!out)return -22;*out=0;if(!id||strcmp(id,"backlight"))return -38;
 struct light *l=calloc(1,sizeof(*l));if(!l)return -12;
 l->common.tag=0x48574454;l->common.version=0;l->common.module=(struct module*)m;l->common.close=close_light;l->set_light=set_backlight;*out=&l->common;return 0;
}
static struct methods methods={open_light};
__attribute__((visibility("default"))) struct module HMI={.tag=0x48574d54,.api=1,.hal=0,.id="lights",.name="GA36 backlight with LCD resume ordering",.author="GA36 bring-up",.methods=&methods};
