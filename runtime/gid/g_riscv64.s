#include "textflag.h"

TEXT ·g(SB), NOSPLIT, $0-8
    MOV g, ret+0(FP)
    RET
