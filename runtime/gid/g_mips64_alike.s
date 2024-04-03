//go:build mips64 || mips64le

#include "textflag.h"

TEXT ·g(SB), NOSPLIT, $0-8
    MOVV g, ret+0(FP)
    RET
