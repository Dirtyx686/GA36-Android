#define _GNU_SOURCE
#include <sys/mman.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
_Static_assert(sizeof(void*)==4,"ARM32 required");
int main(void){
 uint32_t *p=mmap(0,4096,PROT_READ|PROT_WRITE|PROT_EXEC,MAP_PRIVATE|MAP_ANONYMOUS,-1,0);
 if(p==MAP_FAILED)return 1;
 p[0]=0xe1a03000; /* mov r3,r0 */
 p[1]=0xe3530001; /* real kernel cmp r3,1 */
 p[2]=0x0a000019; /* real vendor suppression branch */
 p[3]=0xe3a00001;p[4]=0xe12fff1e;
 p[29]=0xe3a00000;p[30]=0xe12fff1e;
 __builtin___clear_cache((char*)p,(char*)p+128);
 int(*f)(int)=(int(*)(int))p;
 if(f(0)!=1||f(1)!=0||f(2)!=1)return 2;
 p[2]=0xe320f000; /* V19 exactly guarded replacement */
 __builtin___clear_cache((char*)p,(char*)p+128);
 if(f(0)!=1||f(1)!=1||f(2)!=1)return 3;
 munmap(p,4096);
 puts("ARM_POWER_BRANCH_TEST_OK original_flag1_suppressed=true patched_flag0_1_2_delivered=true");return 0;
}
