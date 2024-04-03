#include "textflag.h"

TEXT ·g(SB), NOSPLIT, $0-8
    MOVQ (TLS), AX
    MOVQ AX, ret+0(FP)
    RET
