#include "textflag.h"

TEXT ·g(SB), NOSPLIT, $0-4
    MOVL (TLS), AX
    MOVL AX, ret+0(FP)
    RET
