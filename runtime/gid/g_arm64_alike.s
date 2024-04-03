//go:build arm64 || ppc64 || ppc64le || s390x

#include "textflag.h"

TEXT ·g(SB), NOSPLIT, $0-8
    MOVD g, ret+0(FP)
    RET
