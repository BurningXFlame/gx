//go:build arm || mips || mipsle

#include "textflag.h"

TEXT ·g(SB), NOSPLIT, $0-4
    MOVW g, ret+0(FP)
    RET
