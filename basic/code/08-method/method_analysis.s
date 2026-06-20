# command-line-arguments
main.(*Counter).AddPtr STEXT nosplit size=19 args=0x8 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:14)	TEXT	main.(*Counter).AddPtr(SB), NOSPLIT|NOFRAME|ABIInternal, $0-8
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:14)	FUNCDATA	$0, gclocals·wvjpxkknJ4nY1JtrArJJaw==(SB)
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:14)	FUNCDATA	$1, gclocals·J26BEvPExEQhJvjp9E8Whg==(SB)
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:14)	FUNCDATA	$5, main.(*Counter).AddPtr.arginfo1(SB)
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:14)	MOVQ	AX, main.c+8(SP)
	0x0005 00005 (/home/hp/workspace/learn/golang/basic/08-method/method.go:15)	TESTB	AL, (AX)
	0x0007 00007 (/home/hp/workspace/learn/golang/basic/08-method/method.go:15)	TESTB	AL, (AX)
	0x0009 00009 (/home/hp/workspace/learn/golang/basic/08-method/method.go:15)	MOVQ	(AX), CX
	0x000c 00012 (/home/hp/workspace/learn/golang/basic/08-method/method.go:15)	INCQ	CX
	0x000f 00015 (/home/hp/workspace/learn/golang/basic/08-method/method.go:15)	MOVQ	CX, (AX)
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:16)	RET
	0x0000 48 89 44 24 08 84 00 84 00 48 8b 08 48 ff c1 48  H.D$.....H..H..H
	0x0010 89 08 c3                                         ...
main.Counter.AddValue STEXT nosplit size=14 args=0x8 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:19)	TEXT	main.Counter.AddValue(SB), NOSPLIT|NOFRAME|ABIInternal, $0-8
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:19)	FUNCDATA	$0, gclocals·g5+hNtRBP6YXNjfog7aZjQ==(SB)
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:19)	FUNCDATA	$1, gclocals·g5+hNtRBP6YXNjfog7aZjQ==(SB)
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:19)	FUNCDATA	$5, main.Counter.AddValue.arginfo1(SB)
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:19)	MOVQ	AX, main.c+8(SP)
	0x0005 00005 (/home/hp/workspace/learn/golang/basic/08-method/method.go:20)	INCQ	AX
	0x0008 00008 (/home/hp/workspace/learn/golang/basic/08-method/method.go:20)	MOVQ	AX, main.c+8(SP)
	0x000d 00013 (/home/hp/workspace/learn/golang/basic/08-method/method.go:21)	RET
	0x0000 48 89 44 24 08 48 ff c0 48 89 44 24 08 c3        H.D$.H..H.D$..
main.(*Engine).Start STEXT size=124 args=0x8 locals=0x50 funcid=0x0 align=0x0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	TEXT	main.(*Engine).Start(SB), ABIInternal, $80-8
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	CMPQ	SP, 16(R14)
	0x0004 00004 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	PCDATA	$0, $-2
	0x0004 00004 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	JLS	107
	0x0006 00006 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	PCDATA	$0, $-1
	0x0006 00006 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	PUSHQ	BP
	0x0007 00007 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	MOVQ	SP, BP
	0x000a 00010 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	SUBQ	$72, SP
	0x000e 00014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	FUNCDATA	$0, gclocals·wvjpxkknJ4nY1JtrArJJaw==(SB)
	0x000e 00014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	FUNCDATA	$1, gclocals·ceYgNIaaD8ow5EM5cNccoA==(SB)
	0x000e 00014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	FUNCDATA	$2, main.(*Engine).Start.stkobj(SB)
	0x000e 00014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	FUNCDATA	$5, main.(*Engine).Start.arginfo1(SB)
	0x000e 00014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	MOVQ	AX, main.e+88(SP)
	0x0013 00019 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVUPS	X15, main..autotmp_1+56(SP)
	0x0019 00025 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	LEAQ	main..autotmp_1+56(SP), AX
	0x001e 00030 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	AX, main..autotmp_3+24(SP)
	0x0023 00035 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	LEAQ	type:string(SB), DX
	0x002a 00042 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	DX, main..autotmp_1+56(SP)
	0x002f 00047 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	LEAQ	main..stmp_0(SB), DX
	0x0036 00054 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	DX, main..autotmp_1+64(SP)
	0x003b 00059 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	JMP	61
	0x003d 00061 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	AX, main..autotmp_2+32(SP)
	0x0042 00066 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	$1, main..autotmp_2+40(SP)
	0x004b 00075 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	$1, main..autotmp_2+48(SP)
	0x0054 00084 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVL	$1, BX
	0x0059 00089 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	MOVQ	BX, CX
	0x005c 00092 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	PCDATA	$1, $1
	0x005c 00092 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	NOP
	0x0060 00096 (/home/hp/workspace/learn/golang/basic/08-method/method.go:37)	CALL	fmt.Println(SB)
	0x0065 00101 (/home/hp/workspace/learn/golang/basic/08-method/method.go:38)	ADDQ	$72, SP
	0x0069 00105 (/home/hp/workspace/learn/golang/basic/08-method/method.go:38)	POPQ	BP
	0x006a 00106 (/home/hp/workspace/learn/golang/basic/08-method/method.go:38)	RET
	0x006b 00107 (/home/hp/workspace/learn/golang/basic/08-method/method.go:38)	NOP
	0x006b 00107 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	PCDATA	$1, $-1
	0x006b 00107 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	PCDATA	$0, $-2
	0x006b 00107 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	MOVQ	AX, 8(SP)
	0x0070 00112 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	CALL	runtime.morestack_noctxt(SB)
	0x0075 00117 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	PCDATA	$0, $-1
	0x0075 00117 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	MOVQ	8(SP), AX
	0x007a 00122 (/home/hp/workspace/learn/golang/basic/08-method/method.go:36)	JMP	0
	0x0000 49 3b 66 10 76 65 55 48 89 e5 48 83 ec 48 48 89  I;f.veUH..H..HH.
	0x0010 44 24 58 44 0f 11 7c 24 38 48 8d 44 24 38 48 89  D$XD..|$8H.D$8H.
	0x0020 44 24 18 48 8d 15 00 00 00 00 48 89 54 24 38 48  D$.H......H.T$8H
	0x0030 8d 15 00 00 00 00 48 89 54 24 40 eb 00 48 89 44  ......H.T$@..H.D
	0x0040 24 20 48 c7 44 24 28 01 00 00 00 48 c7 44 24 30  $ H.D$(....H.D$0
	0x0050 01 00 00 00 bb 01 00 00 00 48 89 d9 0f 1f 40 00  .........H....@.
	0x0060 e8 00 00 00 00 48 83 c4 48 5d c3 48 89 44 24 08  .....H..H].H.D$.
	0x0070 e8 00 00 00 00 48 8b 44 24 08 eb 84              .....H.D$...
	rel 2+0 t=R_USEIFACE type:string+0
	rel 38+4 t=R_PCREL type:string+0
	rel 50+4 t=R_PCREL main..stmp_0+0
	rel 97+4 t=R_CALL fmt.Println+0
	rel 113+4 t=R_CALL runtime.morestack_noctxt+0
main.(*Car).Start STEXT size=214 args=0x8 locals=0x78 funcid=0x0 align=0x0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	TEXT	main.(*Car).Start(SB), ABIInternal, $120-8
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	CMPQ	SP, 16(R14)
	0x0004 00004 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	PCDATA	$0, $-2
	0x0004 00004 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	JLS	194
	0x000a 00010 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	PCDATA	$0, $-1
	0x000a 00010 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	PUSHQ	BP
	0x000b 00011 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	MOVQ	SP, BP
	0x000e 00014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	SUBQ	$112, SP
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	FUNCDATA	$0, gclocals·Z8zdw/dq+fE82FieA9ctlQ==(SB)
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	FUNCDATA	$1, gclocals·/4KVoIoAUVnbPMDuZquOrw==(SB)
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	FUNCDATA	$2, main.(*Car).Start.stkobj(SB)
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	FUNCDATA	$5, main.(*Car).Start.arginfo1(SB)
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	MOVQ	AX, main.c+128(SP)
	0x001a 00026 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVUPS	X15, main..autotmp_1+96(SP)
	0x0020 00032 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	LEAQ	main..autotmp_1+96(SP), CX
	0x0025 00037 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	CX, main..autotmp_3+64(SP)
	0x002a 00042 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	main.c+128(SP), CX
	0x0032 00050 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	TESTB	AL, (CX)
	0x0034 00052 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	(CX), AX
	0x0037 00055 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	8(CX), BX
	0x003b 00059 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	AX, main..autotmp_4+48(SP)
	0x0040 00064 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	BX, main..autotmp_4+56(SP)
	0x0045 00069 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	PCDATA	$1, $1
	0x0045 00069 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	CALL	runtime.convTstring(SB)
	0x004a 00074 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	AX, main..autotmp_5+40(SP)
	0x004f 00079 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	main..autotmp_3+64(SP), CX
	0x0054 00084 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	TESTB	AL, (CX)
	0x0056 00086 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	LEAQ	type:string(SB), DX
	0x005d 00093 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	DX, (CX)
	0x0060 00096 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	CMPL	runtime.writeBarrier(SB), $0
	0x0067 00103 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	PCDATA	$0, $-2
	0x0067 00103 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	JEQ	107
	0x0069 00105 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	JMP	109
	0x006b 00107 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	JMP	127
	0x006d 00109 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	8(CX), DX
	0x0071 00113 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	CALL	runtime.gcWriteBarrier2(SB)
	0x0076 00118 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	AX, (R11)
	0x0079 00121 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	DX, 8(R11)
	0x007d 00125 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	JMP	127
	0x007f 00127 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	AX, 8(CX)
	0x0083 00131 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	PCDATA	$0, $-1
	0x0083 00131 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	main..autotmp_3+64(SP), CX
	0x0088 00136 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	TESTB	AL, (CX)
	0x008a 00138 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	JMP	140
	0x008c 00140 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	CX, main..autotmp_2+72(SP)
	0x0091 00145 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	$1, main..autotmp_2+80(SP)
	0x009a 00154 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	$1, main..autotmp_2+88(SP)
	0x00a3 00163 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	LEAQ	go:string."%s 跑车正在启动...\n"(SB), AX
	0x00aa 00170 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVL	$25, BX
	0x00af 00175 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVL	$1, DI
	0x00b4 00180 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	MOVQ	DI, SI
	0x00b7 00183 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	PCDATA	$1, $2
	0x00b7 00183 (/home/hp/workspace/learn/golang/basic/08-method/method.go:47)	CALL	fmt.Printf(SB)
	0x00bc 00188 (/home/hp/workspace/learn/golang/basic/08-method/method.go:48)	ADDQ	$112, SP
	0x00c0 00192 (/home/hp/workspace/learn/golang/basic/08-method/method.go:48)	POPQ	BP
	0x00c1 00193 (/home/hp/workspace/learn/golang/basic/08-method/method.go:48)	RET
	0x00c2 00194 (/home/hp/workspace/learn/golang/basic/08-method/method.go:48)	NOP
	0x00c2 00194 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	PCDATA	$1, $-1
	0x00c2 00194 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	PCDATA	$0, $-2
	0x00c2 00194 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	MOVQ	AX, 8(SP)
	0x00c7 00199 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	CALL	runtime.morestack_noctxt(SB)
	0x00cc 00204 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	PCDATA	$0, $-1
	0x00cc 00204 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	MOVQ	8(SP), AX
	0x00d1 00209 (/home/hp/workspace/learn/golang/basic/08-method/method.go:46)	JMP	0
	0x0000 49 3b 66 10 0f 86 b8 00 00 00 55 48 89 e5 48 83  I;f.......UH..H.
	0x0010 ec 70 48 89 84 24 80 00 00 00 44 0f 11 7c 24 60  .pH..$....D..|$`
	0x0020 48 8d 4c 24 60 48 89 4c 24 40 48 8b 8c 24 80 00  H.L$`H.L$@H..$..
	0x0030 00 00 84 01 48 8b 01 48 8b 59 08 48 89 44 24 30  ....H..H.Y.H.D$0
	0x0040 48 89 5c 24 38 e8 00 00 00 00 48 89 44 24 28 48  H.\$8.....H.D$(H
	0x0050 8b 4c 24 40 84 01 48 8d 15 00 00 00 00 48 89 11  .L$@..H......H..
	0x0060 83 3d 00 00 00 00 00 74 02 eb 02 eb 12 48 8b 51  .=.....t.....H.Q
	0x0070 08 e8 00 00 00 00 49 89 03 49 89 53 08 eb 00 48  ......I..I.S...H
	0x0080 89 41 08 48 8b 4c 24 40 84 01 eb 00 48 89 4c 24  .A.H.L$@....H.L$
	0x0090 48 48 c7 44 24 50 01 00 00 00 48 c7 44 24 58 01  HH.D$P....H.D$X.
	0x00a0 00 00 00 48 8d 05 00 00 00 00 bb 19 00 00 00 bf  ...H............
	0x00b0 01 00 00 00 48 89 fe e8 00 00 00 00 48 83 c4 70  ....H.......H..p
	0x00c0 5d c3 48 89 44 24 08 e8 00 00 00 00 48 8b 44 24  ].H.D$......H.D$
	0x00d0 08 e9 2a ff ff ff                                ..*...
	rel 3+0 t=R_USEIFACE type:string+0
	rel 70+4 t=R_CALL runtime.convTstring+0
	rel 89+4 t=R_PCREL type:string+0
	rel 98+4 t=R_PCREL runtime.writeBarrier+-1
	rel 114+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 166+4 t=R_PCREL go:string."%s 跑车正在启动...\n"+0
	rel 184+4 t=R_CALL fmt.Printf+0
	rel 200+4 t=R_CALL runtime.morestack_noctxt+0
main.(*WeChatPay).Pay STEXT size=325 args=0x10 locals=0x88 funcid=0x0 align=0x0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	TEXT	main.(*WeChatPay).Pay(SB), ABIInternal, $136-16
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	LEAQ	-8(SP), R12
	0x0005 00005 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	CMPQ	R12, 16(R14)
	0x0009 00009 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	PCDATA	$0, $-2
	0x0009 00009 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	JLS	290
	0x000f 00015 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	PCDATA	$0, $-1
	0x000f 00015 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	PUSHQ	BP
	0x0010 00016 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	SP, BP
	0x0013 00019 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	ADDQ	$-128, SP
	0x0017 00023 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	FUNCDATA	$0, gclocals·bUB0t99dbIOHL9YDh6V0CA==(SB)
	0x0017 00023 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	FUNCDATA	$1, gclocals·BcfHs9IhtErNtV2JassdQA==(SB)
	0x0017 00023 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	FUNCDATA	$2, main.(*WeChatPay).Pay.stkobj(SB)
	0x0017 00023 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	FUNCDATA	$5, main.(*WeChatPay).Pay.arginfo1(SB)
	0x0017 00023 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	AX, main.w+144(SP)
	0x001f 00031 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	BX, main.amount+152(SP)
	0x0027 00039 (/home/hp/workspace/learn/golang/basic/08-method/method.go:65)	TESTB	AL, (AX)
	0x0029 00041 (/home/hp/workspace/learn/golang/basic/08-method/method.go:65)	TESTB	AL, (AX)
	0x002b 00043 (/home/hp/workspace/learn/golang/basic/08-method/method.go:65)	MOVQ	(AX), CX
	0x002e 00046 (/home/hp/workspace/learn/golang/basic/08-method/method.go:65)	SUBQ	BX, CX
	0x0031 00049 (/home/hp/workspace/learn/golang/basic/08-method/method.go:65)	MOVQ	CX, (AX)
	0x0034 00052 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	LEAQ	main..autotmp_2+96(SP), CX
	0x0039 00057 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVUPS	X15, (CX)
	0x003d 00061 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVUPS	X15, 16(CX)
	0x0042 00066 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	LEAQ	main..autotmp_2+96(SP), CX
	0x0047 00071 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	CX, main..autotmp_4+64(SP)
	0x004c 00076 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	main.amount+152(SP), AX
	0x0054 00084 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$1, $1
	0x0054 00084 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CALL	runtime.convT64(SB)
	0x0059 00089 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, main..autotmp_5+56(SP)
	0x005e 00094 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	main..autotmp_4+64(SP), CX
	0x0063 00099 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	TESTB	AL, (CX)
	0x0065 00101 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	LEAQ	type:int(SB), DX
	0x006c 00108 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	DX, (CX)
	0x006f 00111 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CMPL	runtime.writeBarrier(SB), $0
	0x0076 00118 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$0, $-2
	0x0076 00118 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JEQ	122
	0x0078 00120 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	124
	0x007a 00122 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	142
	0x007c 00124 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	8(CX), DX
	0x0080 00128 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CALL	runtime.gcWriteBarrier2(SB)
	0x0085 00133 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, (R11)
	0x0088 00136 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	DX, 8(R11)
	0x008c 00140 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	142
	0x008e 00142 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, 8(CX)
	0x0092 00146 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$0, $-1
	0x0092 00146 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	main.w+144(SP), CX
	0x009a 00154 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	TESTB	AL, (CX)
	0x009c 00156 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	(CX), AX
	0x009f 00159 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, main..autotmp_6+40(SP)
	0x00a4 00164 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$1, $2
	0x00a4 00164 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CALL	runtime.convT64(SB)
	0x00a9 00169 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, main..autotmp_7+48(SP)
	0x00ae 00174 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	main..autotmp_4+64(SP), CX
	0x00b3 00179 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	TESTB	AL, (CX)
	0x00b5 00181 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	LEAQ	type:int(SB), DX
	0x00bc 00188 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	DX, 16(CX)
	0x00c0 00192 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CMPL	runtime.writeBarrier(SB), $0
	0x00c7 00199 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$0, $-2
	0x00c7 00199 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JEQ	203
	0x00c9 00201 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	205
	0x00cb 00203 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	223
	0x00cd 00205 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	24(CX), DX
	0x00d1 00209 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CALL	runtime.gcWriteBarrier2(SB)
	0x00d6 00214 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, (R11)
	0x00d9 00217 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	DX, 8(R11)
	0x00dd 00221 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	223
	0x00df 00223 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	AX, 24(CX)
	0x00e3 00227 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$0, $-1
	0x00e3 00227 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	main..autotmp_4+64(SP), CX
	0x00e8 00232 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	TESTB	AL, (CX)
	0x00ea 00234 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	JMP	236
	0x00ec 00236 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	CX, main..autotmp_3+72(SP)
	0x00f1 00241 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	$2, main..autotmp_3+80(SP)
	0x00fa 00250 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	$2, main..autotmp_3+88(SP)
	0x0103 00259 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	LEAQ	go:string."微信支付成功，扣除 %d 元， 余额 %d 元\n"(SB), AX
	0x010a 00266 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVL	$52, BX
	0x010f 00271 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVL	$2, DI
	0x0114 00276 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	MOVQ	DI, SI
	0x0117 00279 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	PCDATA	$1, $3
	0x0117 00279 (/home/hp/workspace/learn/golang/basic/08-method/method.go:66)	CALL	fmt.Printf(SB)
	0x011c 00284 (/home/hp/workspace/learn/golang/basic/08-method/method.go:67)	SUBQ	$-128, SP
	0x0120 00288 (/home/hp/workspace/learn/golang/basic/08-method/method.go:67)	POPQ	BP
	0x0121 00289 (/home/hp/workspace/learn/golang/basic/08-method/method.go:67)	RET
	0x0122 00290 (/home/hp/workspace/learn/golang/basic/08-method/method.go:67)	NOP
	0x0122 00290 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	PCDATA	$1, $-1
	0x0122 00290 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	PCDATA	$0, $-2
	0x0122 00290 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	AX, 8(SP)
	0x0127 00295 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	BX, 16(SP)
	0x012c 00300 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	CALL	runtime.morestack_noctxt(SB)
	0x0131 00305 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	PCDATA	$0, $-1
	0x0131 00305 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	8(SP), AX
	0x0136 00310 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	MOVQ	16(SP), BX
	0x013b 00315 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	NOP
	0x0140 00320 (/home/hp/workspace/learn/golang/basic/08-method/method.go:64)	JMP	0
	0x0000 4c 8d 64 24 f8 4d 3b 66 10 0f 86 13 01 00 00 55  L.d$.M;f.......U
	0x0010 48 89 e5 48 83 c4 80 48 89 84 24 90 00 00 00 48  H..H...H..$....H
	0x0020 89 9c 24 98 00 00 00 84 00 84 00 48 8b 08 48 29  ..$........H..H)
	0x0030 d9 48 89 08 48 8d 4c 24 60 44 0f 11 39 44 0f 11  .H..H.L$`D..9D..
	0x0040 79 10 48 8d 4c 24 60 48 89 4c 24 40 48 8b 84 24  y.H.L$`H.L$@H..$
	0x0050 98 00 00 00 e8 00 00 00 00 48 89 44 24 38 48 8b  .........H.D$8H.
	0x0060 4c 24 40 84 01 48 8d 15 00 00 00 00 48 89 11 83  L$@..H......H...
	0x0070 3d 00 00 00 00 00 74 02 eb 02 eb 12 48 8b 51 08  =.....t.....H.Q.
	0x0080 e8 00 00 00 00 49 89 03 49 89 53 08 eb 00 48 89  .....I..I.S...H.
	0x0090 41 08 48 8b 8c 24 90 00 00 00 84 01 48 8b 01 48  A.H..$......H..H
	0x00a0 89 44 24 28 e8 00 00 00 00 48 89 44 24 30 48 8b  .D$(.....H.D$0H.
	0x00b0 4c 24 40 84 01 48 8d 15 00 00 00 00 48 89 51 10  L$@..H......H.Q.
	0x00c0 83 3d 00 00 00 00 00 74 02 eb 02 eb 12 48 8b 51  .=.....t.....H.Q
	0x00d0 18 e8 00 00 00 00 49 89 03 49 89 53 08 eb 00 48  ......I..I.S...H
	0x00e0 89 41 18 48 8b 4c 24 40 84 01 eb 00 48 89 4c 24  .A.H.L$@....H.L$
	0x00f0 48 48 c7 44 24 50 02 00 00 00 48 c7 44 24 58 02  HH.D$P....H.D$X.
	0x0100 00 00 00 48 8d 05 00 00 00 00 bb 34 00 00 00 bf  ...H.......4....
	0x0110 02 00 00 00 48 89 fe e8 00 00 00 00 48 83 ec 80  ....H.......H...
	0x0120 5d c3 48 89 44 24 08 48 89 5c 24 10 e8 00 00 00  ].H.D$.H.\$.....
	0x0130 00 48 8b 44 24 08 48 8b 5c 24 10 0f 1f 44 00 00  .H.D$.H.\$...D..
	0x0140 e9 bb fe ff ff                                   .....
	rel 3+0 t=R_USEIFACE type:int+0
	rel 3+0 t=R_USEIFACE type:int+0
	rel 85+4 t=R_CALL runtime.convT64+0
	rel 104+4 t=R_PCREL type:int+0
	rel 113+4 t=R_PCREL runtime.writeBarrier+-1
	rel 129+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 165+4 t=R_CALL runtime.convT64+0
	rel 184+4 t=R_PCREL type:int+0
	rel 194+4 t=R_PCREL runtime.writeBarrier+-1
	rel 210+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 262+4 t=R_PCREL go:string."微信支付成功，扣除 %d 元， 余额 %d 元\n"+0
	rel 280+4 t=R_CALL fmt.Printf+0
	rel 301+4 t=R_CALL runtime.morestack_noctxt+0
main.main STEXT size=1162 args=0x0 locals=0x198 funcid=0x0 align=0x0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	TEXT	main.main(SB), ABIInternal, $408-0
	0x0000 00000 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	LEAQ	-280(SP), R12
	0x0008 00008 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	CMPQ	R12, 16(R14)
	0x000c 00012 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	PCDATA	$0, $-2
	0x000c 00012 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	JLS	1149
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	PCDATA	$0, $-1
	0x0012 00018 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	PUSHQ	BP
	0x0013 00019 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	MOVQ	SP, BP
	0x0016 00022 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	SUBQ	$400, SP
	0x001d 00029 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	FUNCDATA	$0, gclocals·Dj7m7VTqKq6fxJqfrrXabg==(SB)
	0x001d 00029 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	FUNCDATA	$1, gclocals·MlHovJf8fukPSXh1aq01eg==(SB)
	0x001d 00029 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	FUNCDATA	$2, main.main.stkobj(SB)
	0x001d 00029 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVUPS	X15, main..autotmp_4+168(SP)
	0x0026 00038 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	LEAQ	main..autotmp_4+168(SP), AX
	0x002e 00046 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	AX, main..autotmp_8+96(SP)
	0x0033 00051 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	LEAQ	type:string(SB), DX
	0x003a 00058 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	DX, main..autotmp_4+168(SP)
	0x0042 00066 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	LEAQ	main..stmp_1(SB), DX
	0x0049 00073 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	DX, main..autotmp_4+176(SP)
	0x0051 00081 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	JMP	83
	0x0053 00083 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	AX, main..autotmp_7+104(SP)
	0x0058 00088 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	$1, main..autotmp_7+112(SP)
	0x0061 00097 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	$1, main..autotmp_7+120(SP)
	0x006a 00106 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVL	$1, BX
	0x006f 00111 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	MOVQ	BX, CX
	0x0072 00114 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	PCDATA	$1, $0
	0x0072 00114 (/home/hp/workspace/learn/golang/basic/08-method/method.go:70)	CALL	fmt.Println(SB)
	0x0077 00119 (/home/hp/workspace/learn/golang/basic/08-method/method.go:71)	MOVQ	$0, main..autotmp_9+40(SP)
	0x0080 00128 (/home/hp/workspace/learn/golang/basic/08-method/method.go:71)	MOVQ	$1, main..autotmp_9+40(SP)
	0x0089 00137 (/home/hp/workspace/learn/golang/basic/08-method/method.go:71)	MOVQ	$1, main.c+32(SP)
	0x0092 00146 (/home/hp/workspace/learn/golang/basic/08-method/method.go:73)	LEAQ	main.c+32(SP), AX
	0x0097 00151 (/home/hp/workspace/learn/golang/basic/08-method/method.go:73)	CALL	main.(*Counter).AddPtr(SB)
	0x009c 00156 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	LEAQ	main..autotmp_5+136(SP), DX
	0x00a4 00164 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVUPS	X15, (DX)
	0x00a8 00168 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVUPS	X15, 16(DX)
	0x00ad 00173 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	LEAQ	main..autotmp_5+136(SP), DX
	0x00b5 00181 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	DX, main..autotmp_11+368(SP)
	0x00bd 00189 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	LEAQ	type:string(SB), DX
	0x00c4 00196 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	DX, main..autotmp_5+136(SP)
	0x00cc 00204 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	LEAQ	main..stmp_2(SB), DX
	0x00d3 00211 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	DX, main..autotmp_5+144(SP)
	0x00db 00219 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	main.c+32(SP), AX
	0x00e0 00224 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	PCDATA	$1, $1
	0x00e0 00224 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	CALL	runtime.convT64(SB)
	0x00e5 00229 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	AX, main..autotmp_12+360(SP)
	0x00ed 00237 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	main..autotmp_11+368(SP), DX
	0x00f5 00245 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	TESTB	AL, (DX)
	0x00f7 00247 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	LEAQ	type:int(SB), SI
	0x00fe 00254 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	SI, 16(DX)
	0x0102 00258 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	CMPL	runtime.writeBarrier(SB), $0
	0x0109 00265 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	PCDATA	$0, $-2
	0x0109 00265 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	JEQ	269
	0x010b 00267 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	JMP	271
	0x010d 00269 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	JMP	290
	0x010f 00271 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	24(DX), SI
	0x0113 00275 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	CALL	runtime.gcWriteBarrier2(SB)
	0x0118 00280 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	AX, (R11)
	0x011b 00283 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	SI, 8(R11)
	0x011f 00287 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	NOP
	0x0120 00288 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	JMP	290
	0x0122 00290 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	AX, 24(DX)
	0x0126 00294 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	PCDATA	$0, $-1
	0x0126 00294 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	main..autotmp_11+368(SP), AX
	0x012e 00302 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	TESTB	AL, (AX)
	0x0130 00304 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	JMP	306
	0x0132 00306 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	AX, main..autotmp_10+376(SP)
	0x013a 00314 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	$2, main..autotmp_10+384(SP)
	0x0146 00326 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	$2, main..autotmp_10+392(SP)
	0x0152 00338 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVL	$2, BX
	0x0157 00343 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	MOVQ	BX, CX
	0x015a 00346 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	PCDATA	$1, $0
	0x015a 00346 (/home/hp/workspace/learn/golang/basic/08-method/method.go:74)	CALL	fmt.Println(SB)
	0x015f 00351 (/home/hp/workspace/learn/golang/basic/08-method/method.go:77)	MOVQ	main.c+32(SP), AX
	0x0164 00356 (/home/hp/workspace/learn/golang/basic/08-method/method.go:77)	CALL	main.Counter.AddValue(SB)
	0x0169 00361 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	LEAQ	main..autotmp_5+136(SP), DX
	0x0171 00369 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVUPS	X15, (DX)
	0x0175 00373 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVUPS	X15, 16(DX)
	0x017a 00378 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	LEAQ	main..autotmp_5+136(SP), DX
	0x0182 00386 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	DX, main..autotmp_14+328(SP)
	0x018a 00394 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	LEAQ	type:string(SB), DX
	0x0191 00401 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	DX, main..autotmp_5+136(SP)
	0x0199 00409 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	LEAQ	main..stmp_3(SB), DX
	0x01a0 00416 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	DX, main..autotmp_5+144(SP)
	0x01a8 00424 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	main.c+32(SP), AX
	0x01ad 00429 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	PCDATA	$1, $2
	0x01ad 00429 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	CALL	runtime.convT64(SB)
	0x01b2 00434 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	AX, main..autotmp_15+320(SP)
	0x01ba 00442 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	main..autotmp_14+328(SP), DX
	0x01c2 00450 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	TESTB	AL, (DX)
	0x01c4 00452 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	LEAQ	type:int(SB), SI
	0x01cb 00459 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	SI, 16(DX)
	0x01cf 00463 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	CMPL	runtime.writeBarrier(SB), $0
	0x01d6 00470 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	PCDATA	$0, $-2
	0x01d6 00470 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	JEQ	474
	0x01d8 00472 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	JMP	476
	0x01da 00474 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	JMP	494
	0x01dc 00476 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	24(DX), SI
	0x01e0 00480 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	CALL	runtime.gcWriteBarrier2(SB)
	0x01e5 00485 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	AX, (R11)
	0x01e8 00488 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	SI, 8(R11)
	0x01ec 00492 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	JMP	494
	0x01ee 00494 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	AX, 24(DX)
	0x01f2 00498 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	PCDATA	$0, $-1
	0x01f2 00498 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	main..autotmp_14+328(SP), AX
	0x01fa 00506 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	TESTB	AL, (AX)
	0x01fc 00508 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	JMP	510
	0x01fe 00510 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	AX, main..autotmp_13+336(SP)
	0x0206 00518 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	$2, main..autotmp_13+344(SP)
	0x0212 00530 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	$2, main..autotmp_13+352(SP)
	0x021e 00542 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVL	$2, BX
	0x0223 00547 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	MOVQ	BX, CX
	0x0226 00550 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	PCDATA	$1, $0
	0x0226 00550 (/home/hp/workspace/learn/golang/basic/08-method/method.go:78)	CALL	fmt.Println(SB)
	0x022b 00555 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVUPS	X15, main..autotmp_4+168(SP)
	0x0234 00564 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	LEAQ	main..autotmp_4+168(SP), AX
	0x023c 00572 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	AX, main..autotmp_17+288(SP)
	0x0244 00580 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	LEAQ	type:string(SB), DX
	0x024b 00587 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	DX, main..autotmp_4+168(SP)
	0x0253 00595 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	LEAQ	main..stmp_4(SB), DX
	0x025a 00602 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	DX, main..autotmp_4+176(SP)
	0x0262 00610 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	JMP	612
	0x0264 00612 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	AX, main..autotmp_16+296(SP)
	0x026c 00620 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	$1, main..autotmp_16+304(SP)
	0x0278 00632 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	$1, main..autotmp_16+312(SP)
	0x0284 00644 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVL	$1, BX
	0x0289 00649 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	MOVQ	BX, CX
	0x028c 00652 (/home/hp/workspace/learn/golang/basic/08-method/method.go:80)	CALL	fmt.Println(SB)
	0x0291 00657 (/home/hp/workspace/learn/golang/basic/08-method/method.go:81)	MOVUPS	X15, main..autotmp_18+264(SP)
	0x029a 00666 (/home/hp/workspace/learn/golang/basic/08-method/method.go:81)	MOVQ	$0, main..autotmp_18+280(SP)
	0x02a6 00678 (/home/hp/workspace/learn/golang/basic/08-method/method.go:82)	LEAQ	go:string."Audi"(SB), DX
	0x02ad 00685 (/home/hp/workspace/learn/golang/basic/08-method/method.go:82)	MOVQ	DX, main..autotmp_18+264(SP)
	0x02b5 00693 (/home/hp/workspace/learn/golang/basic/08-method/method.go:82)	MOVQ	$4, main..autotmp_18+272(SP)
	0x02c1 00705 (/home/hp/workspace/learn/golang/basic/08-method/method.go:83)	MOVQ	$500, main..autotmp_18+280(SP)
	0x02cd 00717 (/home/hp/workspace/learn/golang/basic/08-method/method.go:81)	MOVQ	DX, main.myCar+72(SP)
	0x02d2 00722 (/home/hp/workspace/learn/golang/basic/08-method/method.go:81)	MOVQ	$4, main.myCar+80(SP)
	0x02db 00731 (/home/hp/workspace/learn/golang/basic/08-method/method.go:81)	MOVQ	$500, main.myCar+88(SP)
	0x02e4 00740 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	LEAQ	main..autotmp_5+136(SP), DX
	0x02ec 00748 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVUPS	X15, (DX)
	0x02f0 00752 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVUPS	X15, 16(DX)
	0x02f5 00757 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	LEAQ	main..autotmp_5+136(SP), DX
	0x02fd 00765 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	DX, main..autotmp_20+232(SP)
	0x0305 00773 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	LEAQ	type:string(SB), DX
	0x030c 00780 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	DX, main..autotmp_5+136(SP)
	0x0314 00788 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	LEAQ	main..stmp_5(SB), DX
	0x031b 00795 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	DX, main..autotmp_5+144(SP)
	0x0323 00803 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	main.myCar+88(SP), AX
	0x0328 00808 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	PCDATA	$1, $3
	0x0328 00808 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	CALL	runtime.convT64(SB)
	0x032d 00813 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	AX, main..autotmp_21+224(SP)
	0x0335 00821 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	main..autotmp_20+232(SP), DX
	0x033d 00829 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	TESTB	AL, (DX)
	0x033f 00831 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	LEAQ	type:int(SB), SI
	0x0346 00838 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	SI, 16(DX)
	0x034a 00842 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	CMPL	runtime.writeBarrier(SB), $0
	0x0351 00849 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	PCDATA	$0, $-2
	0x0351 00849 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	JEQ	853
	0x0353 00851 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	JMP	855
	0x0355 00853 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	JMP	878
	0x0357 00855 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	24(DX), SI
	0x035b 00859 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	NOP
	0x0360 00864 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	CALL	runtime.gcWriteBarrier2(SB)
	0x0365 00869 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	AX, (R11)
	0x0368 00872 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	SI, 8(R11)
	0x036c 00876 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	JMP	878
	0x036e 00878 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	AX, 24(DX)
	0x0372 00882 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	PCDATA	$0, $-1
	0x0372 00882 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	main..autotmp_20+232(SP), AX
	0x037a 00890 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	TESTB	AL, (AX)
	0x037c 00892 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	JMP	894
	0x037e 00894 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	AX, main..autotmp_19+240(SP)
	0x0386 00902 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	$2, main..autotmp_19+248(SP)
	0x0392 00914 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	$2, main..autotmp_19+256(SP)
	0x039e 00926 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVL	$2, BX
	0x03a3 00931 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	MOVQ	BX, CX
	0x03a6 00934 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	PCDATA	$1, $4
	0x03a6 00934 (/home/hp/workspace/learn/golang/basic/08-method/method.go:86)	CALL	fmt.Println(SB)
	0x03ab 00939 (/home/hp/workspace/learn/golang/basic/08-method/method.go:88)	LEAQ	main.myCar+72(SP), AX
	0x03b0 00944 (/home/hp/workspace/learn/golang/basic/08-method/method.go:88)	CALL	main.(*Car).Start(SB)
	0x03b5 00949 (/home/hp/workspace/learn/golang/basic/08-method/method.go:90)	LEAQ	main.myCar+88(SP), AX
	0x03ba 00954 (/home/hp/workspace/learn/golang/basic/08-method/method.go:90)	PCDATA	$1, $0
	0x03ba 00954 (/home/hp/workspace/learn/golang/basic/08-method/method.go:90)	CALL	main.(*Engine).Start(SB)
	0x03bf 00959 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVUPS	X15, main..autotmp_4+168(SP)
	0x03c8 00968 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	LEAQ	main..autotmp_4+168(SP), AX
	0x03d0 00976 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	AX, main..autotmp_23+192(SP)
	0x03d8 00984 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	LEAQ	type:string(SB), DX
	0x03df 00991 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	DX, main..autotmp_4+168(SP)
	0x03e7 00999 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	LEAQ	main..stmp_6(SB), DX
	0x03ee 01006 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	DX, main..autotmp_4+176(SP)
	0x03f6 01014 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	JMP	1016
	0x03f8 01016 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	AX, main..autotmp_22+200(SP)
	0x0400 01024 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	$1, main..autotmp_22+208(SP)
	0x040c 01036 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	$1, main..autotmp_22+216(SP)
	0x0418 01048 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVL	$1, BX
	0x041d 01053 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	MOVQ	BX, CX
	0x0420 01056 (/home/hp/workspace/learn/golang/basic/08-method/method.go:92)	CALL	fmt.Println(SB)
	0x0425 01061 (/home/hp/workspace/learn/golang/basic/08-method/method.go:93)	MOVQ	$0, main..autotmp_24+48(SP)
	0x042e 01070 (/home/hp/workspace/learn/golang/basic/08-method/method.go:93)	MOVQ	$100, main..autotmp_24+48(SP)
	0x0437 01079 (/home/hp/workspace/learn/golang/basic/08-method/method.go:93)	MOVQ	$100, main.wxWallet+24(SP)
	0x0440 01088 (/home/hp/workspace/learn/golang/basic/08-method/method.go:96)	LEAQ	main.wxWallet+24(SP), AX
	0x0445 01093 (/home/hp/workspace/learn/golang/basic/08-method/method.go:96)	MOVQ	AX, main..autotmp_6+128(SP)
	0x044d 01101 (/home/hp/workspace/learn/golang/basic/08-method/method.go:96)	LEAQ	go:itab.*main.WeChatPay,main.Payer(SB), DX
	0x0454 01108 (/home/hp/workspace/learn/golang/basic/08-method/method.go:96)	MOVQ	DX, main.p2+56(SP)
	0x0459 01113 (/home/hp/workspace/learn/golang/basic/08-method/method.go:96)	MOVQ	AX, main.p2+64(SP)
	0x045e 01118 (/home/hp/workspace/learn/golang/basic/08-method/method.go:96)	NOP
	0x0460 01120 (/home/hp/workspace/learn/golang/basic/08-method/method.go:97)	JMP	1122
	0x0462 01122 (/home/hp/workspace/learn/golang/basic/08-method/method.go:97)	MOVQ	AX, main..autotmp_25+184(SP)
	0x046a 01130 (/home/hp/workspace/learn/golang/basic/08-method/method.go:97)	MOVL	$20, BX
	0x046f 01135 (/home/hp/workspace/learn/golang/basic/08-method/method.go:97)	CALL	main.(*WeChatPay).Pay(SB)
	0x0474 01140 (/home/hp/workspace/learn/golang/basic/08-method/method.go:98)	ADDQ	$400, SP
	0x047b 01147 (/home/hp/workspace/learn/golang/basic/08-method/method.go:98)	POPQ	BP
	0x047c 01148 (/home/hp/workspace/learn/golang/basic/08-method/method.go:98)	RET
	0x047d 01149 (/home/hp/workspace/learn/golang/basic/08-method/method.go:98)	NOP
	0x047d 01149 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	PCDATA	$1, $-1
	0x047d 01149 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	PCDATA	$0, $-2
	0x047d 01149 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	NOP
	0x0480 01152 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	CALL	runtime.morestack_noctxt(SB)
	0x0485 01157 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	PCDATA	$0, $-1
	0x0485 01157 (/home/hp/workspace/learn/golang/basic/08-method/method.go:69)	JMP	0
	0x0000 4c 8d a4 24 e8 fe ff ff 4d 3b 66 10 0f 86 6b 04  L..$....M;f...k.
	0x0010 00 00 55 48 89 e5 48 81 ec 90 01 00 00 44 0f 11  ..UH..H......D..
	0x0020 bc 24 a8 00 00 00 48 8d 84 24 a8 00 00 00 48 89  .$....H..$....H.
	0x0030 44 24 60 48 8d 15 00 00 00 00 48 89 94 24 a8 00  D$`H......H..$..
	0x0040 00 00 48 8d 15 00 00 00 00 48 89 94 24 b0 00 00  ..H......H..$...
	0x0050 00 eb 00 48 89 44 24 68 48 c7 44 24 70 01 00 00  ...H.D$hH.D$p...
	0x0060 00 48 c7 44 24 78 01 00 00 00 bb 01 00 00 00 48  .H.D$x.........H
	0x0070 89 d9 e8 00 00 00 00 48 c7 44 24 28 00 00 00 00  .......H.D$(....
	0x0080 48 c7 44 24 28 01 00 00 00 48 c7 44 24 20 01 00  H.D$(....H.D$ ..
	0x0090 00 00 48 8d 44 24 20 e8 00 00 00 00 48 8d 94 24  ..H.D$ .....H..$
	0x00a0 88 00 00 00 44 0f 11 3a 44 0f 11 7a 10 48 8d 94  ....D..:D..z.H..
	0x00b0 24 88 00 00 00 48 89 94 24 70 01 00 00 48 8d 15  $....H..$p...H..
	0x00c0 00 00 00 00 48 89 94 24 88 00 00 00 48 8d 15 00  ....H..$....H...
	0x00d0 00 00 00 48 89 94 24 90 00 00 00 48 8b 44 24 20  ...H..$....H.D$ 
	0x00e0 e8 00 00 00 00 48 89 84 24 68 01 00 00 48 8b 94  .....H..$h...H..
	0x00f0 24 70 01 00 00 84 02 48 8d 35 00 00 00 00 48 89  $p.....H.5....H.
	0x0100 72 10 83 3d 00 00 00 00 00 74 02 eb 02 eb 13 48  r..=.....t.....H
	0x0110 8b 72 18 e8 00 00 00 00 49 89 03 49 89 73 08 90  .r......I..I.s..
	0x0120 eb 00 48 89 42 18 48 8b 84 24 70 01 00 00 84 00  ..H.B.H..$p.....
	0x0130 eb 00 48 89 84 24 78 01 00 00 48 c7 84 24 80 01  ..H..$x...H..$..
	0x0140 00 00 02 00 00 00 48 c7 84 24 88 01 00 00 02 00  ......H..$......
	0x0150 00 00 bb 02 00 00 00 48 89 d9 e8 00 00 00 00 48  .......H.......H
	0x0160 8b 44 24 20 e8 00 00 00 00 48 8d 94 24 88 00 00  .D$ .....H..$...
	0x0170 00 44 0f 11 3a 44 0f 11 7a 10 48 8d 94 24 88 00  .D..:D..z.H..$..
	0x0180 00 00 48 89 94 24 48 01 00 00 48 8d 15 00 00 00  ..H..$H...H.....
	0x0190 00 48 89 94 24 88 00 00 00 48 8d 15 00 00 00 00  .H..$....H......
	0x01a0 48 89 94 24 90 00 00 00 48 8b 44 24 20 e8 00 00  H..$....H.D$ ...
	0x01b0 00 00 48 89 84 24 40 01 00 00 48 8b 94 24 48 01  ..H..$@...H..$H.
	0x01c0 00 00 84 02 48 8d 35 00 00 00 00 48 89 72 10 83  ....H.5....H.r..
	0x01d0 3d 00 00 00 00 00 74 02 eb 02 eb 12 48 8b 72 18  =.....t.....H.r.
	0x01e0 e8 00 00 00 00 49 89 03 49 89 73 08 eb 00 48 89  .....I..I.s...H.
	0x01f0 42 18 48 8b 84 24 48 01 00 00 84 00 eb 00 48 89  B.H..$H.......H.
	0x0200 84 24 50 01 00 00 48 c7 84 24 58 01 00 00 02 00  .$P...H..$X.....
	0x0210 00 00 48 c7 84 24 60 01 00 00 02 00 00 00 bb 02  ..H..$`.........
	0x0220 00 00 00 48 89 d9 e8 00 00 00 00 44 0f 11 bc 24  ...H.......D...$
	0x0230 a8 00 00 00 48 8d 84 24 a8 00 00 00 48 89 84 24  ....H..$....H..$
	0x0240 20 01 00 00 48 8d 15 00 00 00 00 48 89 94 24 a8   ...H......H..$.
	0x0250 00 00 00 48 8d 15 00 00 00 00 48 89 94 24 b0 00  ...H......H..$..
	0x0260 00 00 eb 00 48 89 84 24 28 01 00 00 48 c7 84 24  ....H..$(...H..$
	0x0270 30 01 00 00 01 00 00 00 48 c7 84 24 38 01 00 00  0.......H..$8...
	0x0280 01 00 00 00 bb 01 00 00 00 48 89 d9 e8 00 00 00  .........H......
	0x0290 00 44 0f 11 bc 24 08 01 00 00 48 c7 84 24 18 01  .D...$....H..$..
	0x02a0 00 00 00 00 00 00 48 8d 15 00 00 00 00 48 89 94  ......H......H..
	0x02b0 24 08 01 00 00 48 c7 84 24 10 01 00 00 04 00 00  $....H..$.......
	0x02c0 00 48 c7 84 24 18 01 00 00 f4 01 00 00 48 89 54  .H..$........H.T
	0x02d0 24 48 48 c7 44 24 50 04 00 00 00 48 c7 44 24 58  $HH.D$P....H.D$X
	0x02e0 f4 01 00 00 48 8d 94 24 88 00 00 00 44 0f 11 3a  ....H..$....D..:
	0x02f0 44 0f 11 7a 10 48 8d 94 24 88 00 00 00 48 89 94  D..z.H..$....H..
	0x0300 24 e8 00 00 00 48 8d 15 00 00 00 00 48 89 94 24  $....H......H..$
	0x0310 88 00 00 00 48 8d 15 00 00 00 00 48 89 94 24 90  ....H......H..$.
	0x0320 00 00 00 48 8b 44 24 58 e8 00 00 00 00 48 89 84  ...H.D$X.....H..
	0x0330 24 e0 00 00 00 48 8b 94 24 e8 00 00 00 84 02 48  $....H..$......H
	0x0340 8d 35 00 00 00 00 48 89 72 10 83 3d 00 00 00 00  .5....H.r..=....
	0x0350 00 74 02 eb 02 eb 17 48 8b 72 18 0f 1f 44 00 00  .t.....H.r...D..
	0x0360 e8 00 00 00 00 49 89 03 49 89 73 08 eb 00 48 89  .....I..I.s...H.
	0x0370 42 18 48 8b 84 24 e8 00 00 00 84 00 eb 00 48 89  B.H..$........H.
	0x0380 84 24 f0 00 00 00 48 c7 84 24 f8 00 00 00 02 00  .$....H..$......
	0x0390 00 00 48 c7 84 24 00 01 00 00 02 00 00 00 bb 02  ..H..$..........
	0x03a0 00 00 00 48 89 d9 e8 00 00 00 00 48 8d 44 24 48  ...H.......H.D$H
	0x03b0 e8 00 00 00 00 48 8d 44 24 58 e8 00 00 00 00 44  .....H.D$X.....D
	0x03c0 0f 11 bc 24 a8 00 00 00 48 8d 84 24 a8 00 00 00  ...$....H..$....
	0x03d0 48 89 84 24 c0 00 00 00 48 8d 15 00 00 00 00 48  H..$....H......H
	0x03e0 89 94 24 a8 00 00 00 48 8d 15 00 00 00 00 48 89  ..$....H......H.
	0x03f0 94 24 b0 00 00 00 eb 00 48 89 84 24 c8 00 00 00  .$......H..$....
	0x0400 48 c7 84 24 d0 00 00 00 01 00 00 00 48 c7 84 24  H..$........H..$
	0x0410 d8 00 00 00 01 00 00 00 bb 01 00 00 00 48 89 d9  .............H..
	0x0420 e8 00 00 00 00 48 c7 44 24 30 00 00 00 00 48 c7  .....H.D$0....H.
	0x0430 44 24 30 64 00 00 00 48 c7 44 24 18 64 00 00 00  D$0d...H.D$.d...
	0x0440 48 8d 44 24 18 48 89 84 24 80 00 00 00 48 8d 15  H.D$.H..$....H..
	0x0450 00 00 00 00 48 89 54 24 38 48 89 44 24 40 66 90  ....H.T$8H.D$@f.
	0x0460 eb 00 48 89 84 24 b8 00 00 00 bb 14 00 00 00 e8  ..H..$..........
	0x0470 00 00 00 00 48 81 c4 90 01 00 00 5d c3 0f 1f 00  ....H......]....
	0x0480 e8 00 00 00 00 e9 76 fb ff ff                    ......v...
	rel 3+0 t=R_USEIFACE type:string+0
	rel 3+0 t=R_USEIFACE type:string+0
	rel 3+0 t=R_USEIFACE type:int+0
	rel 3+0 t=R_USEIFACE type:string+0
	rel 3+0 t=R_USEIFACE type:int+0
	rel 3+0 t=R_USEIFACE type:string+0
	rel 3+0 t=R_USEIFACE type:string+0
	rel 3+0 t=R_USEIFACE type:int+0
	rel 3+0 t=R_USEIFACE type:string+0
	rel 3+0 t=R_USEIFACE type:*main.WeChatPay+0
	rel 54+4 t=R_PCREL type:string+0
	rel 69+4 t=R_PCREL main..stmp_1+0
	rel 115+4 t=R_CALL fmt.Println+0
	rel 152+4 t=R_CALL main.(*Counter).AddPtr+0
	rel 192+4 t=R_PCREL type:string+0
	rel 207+4 t=R_PCREL main..stmp_2+0
	rel 225+4 t=R_CALL runtime.convT64+0
	rel 250+4 t=R_PCREL type:int+0
	rel 260+4 t=R_PCREL runtime.writeBarrier+-1
	rel 276+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 347+4 t=R_CALL fmt.Println+0
	rel 357+4 t=R_CALL main.Counter.AddValue+0
	rel 397+4 t=R_PCREL type:string+0
	rel 412+4 t=R_PCREL main..stmp_3+0
	rel 430+4 t=R_CALL runtime.convT64+0
	rel 455+4 t=R_PCREL type:int+0
	rel 465+4 t=R_PCREL runtime.writeBarrier+-1
	rel 481+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 551+4 t=R_CALL fmt.Println+0
	rel 583+4 t=R_PCREL type:string+0
	rel 598+4 t=R_PCREL main..stmp_4+0
	rel 653+4 t=R_CALL fmt.Println+0
	rel 681+4 t=R_PCREL go:string."Audi"+0
	rel 776+4 t=R_PCREL type:string+0
	rel 791+4 t=R_PCREL main..stmp_5+0
	rel 809+4 t=R_CALL runtime.convT64+0
	rel 834+4 t=R_PCREL type:int+0
	rel 844+4 t=R_PCREL runtime.writeBarrier+-1
	rel 865+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 935+4 t=R_CALL fmt.Println+0
	rel 945+4 t=R_CALL main.(*Car).Start+0
	rel 955+4 t=R_CALL main.(*Engine).Start+0
	rel 987+4 t=R_PCREL type:string+0
	rel 1002+4 t=R_PCREL main..stmp_6+0
	rel 1057+4 t=R_CALL fmt.Println+0
	rel 1104+4 t=R_PCREL go:itab.*main.WeChatPay,main.Payer+0
	rel 1136+4 t=R_CALL main.(*WeChatPay).Pay+0
	rel 1153+4 t=R_CALL runtime.morestack_noctxt+0
main.(*Counter).AddValue STEXT dupok size=70 args=0x8 locals=0x18 funcid=0x17 align=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	main.(*Counter).AddValue(SB), DUPOK|WRAPPER|ABIInternal, $24-8
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	53
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	PUSHQ	BP
	0x0007 00007 (<autogenerated>:1)	MOVQ	SP, BP
	0x000a 00010 (<autogenerated>:1)	SUBQ	$16, SP
	0x000e 00014 (<autogenerated>:1)	FUNCDATA	$0, gclocals·wvjpxkknJ4nY1JtrArJJaw==(SB)
	0x000e 00014 (<autogenerated>:1)	FUNCDATA	$1, gclocals·J26BEvPExEQhJvjp9E8Whg==(SB)
	0x000e 00014 (<autogenerated>:1)	FUNCDATA	$5, main.(*Counter).AddValue.arginfo1(SB)
	0x000e 00014 (<autogenerated>:1)	MOVQ	AX, main.c+32(SP)
	0x0013 00019 (<autogenerated>:1)	TESTQ	AX, AX
	0x0016 00022 (<autogenerated>:1)	JNE	26
	0x0018 00024 (<autogenerated>:1)	JMP	47
	0x001a 00026 (<autogenerated>:1)	TESTB	AL, (AX)
	0x001c 00028 (<autogenerated>:1)	MOVQ	(AX), AX
	0x001f 00031 (<autogenerated>:1)	MOVQ	AX, main..autotmp_1+8(SP)
	0x0024 00036 (<autogenerated>:1)	PCDATA	$1, $1
	0x0024 00036 (<autogenerated>:1)	CALL	main.Counter.AddValue(SB)
	0x0029 00041 (<autogenerated>:1)	ADDQ	$16, SP
	0x002d 00045 (<autogenerated>:1)	POPQ	BP
	0x002e 00046 (<autogenerated>:1)	RET
	0x002f 00047 (<autogenerated>:1)	CALL	runtime.panicwrap(SB)
	0x0034 00052 (<autogenerated>:1)	XCHGL	AX, AX
	0x0035 00053 (<autogenerated>:1)	NOP
	0x0035 00053 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0035 00053 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0035 00053 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x003a 00058 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x003f 00063 (<autogenerated>:1)	PCDATA	$0, $-1
	0x003f 00063 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x0044 00068 (<autogenerated>:1)	JMP	0
	0x0000 49 3b 66 10 76 2f 55 48 89 e5 48 83 ec 10 48 89  I;f.v/UH..H...H.
	0x0010 44 24 20 48 85 c0 75 02 eb 15 84 00 48 8b 00 48  D$ H..u.....H..H
	0x0020 89 44 24 08 e8 00 00 00 00 48 83 c4 10 5d c3 e8  .D$......H...]..
	0x0030 00 00 00 00 90 48 89 44 24 08 e8 00 00 00 00 48  .....H.D$......H
	0x0040 8b 44 24 08 eb ba                                .D$...
	rel 37+4 t=R_CALL main.Counter.AddValue+0
	rel 48+4 t=R_CALL runtime.panicwrap+0
	rel 59+4 t=R_CALL runtime.morestack_noctxt+0
main.Payer.Pay STEXT dupok size=86 args=0x18 locals=0x18 funcid=0x17 align=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	main.Payer.Pay(SB), DUPOK|WRAPPER|ABIInternal, $24-24
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	49
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	PUSHQ	BP
	0x0007 00007 (<autogenerated>:1)	MOVQ	SP, BP
	0x000a 00010 (<autogenerated>:1)	SUBQ	$16, SP
	0x000e 00014 (<autogenerated>:1)	FUNCDATA	$0, gclocals·Ih7UaEzxol4qYE7mnpMg6w==(SB)
	0x000e 00014 (<autogenerated>:1)	FUNCDATA	$1, gclocals·J26BEvPExEQhJvjp9E8Whg==(SB)
	0x000e 00014 (<autogenerated>:1)	FUNCDATA	$5, main.Payer.Pay.arginfo1(SB)
	0x000e 00014 (<autogenerated>:1)	MOVQ	AX, main.~p0+32(SP)
	0x0013 00019 (<autogenerated>:1)	MOVQ	BX, main.~p0+40(SP)
	0x0018 00024 (<autogenerated>:1)	MOVQ	CX, main.amount+48(SP)
	0x001d 00029 (<autogenerated>:1)	TESTB	AL, (AX)
	0x001f 00031 (<autogenerated>:1)	MOVQ	24(AX), DX
	0x0023 00035 (<autogenerated>:1)	MOVQ	BX, AX
	0x0026 00038 (<autogenerated>:1)	MOVQ	CX, BX
	0x0029 00041 (<autogenerated>:1)	PCDATA	$1, $1
	0x0029 00041 (<autogenerated>:1)	CALL	DX
	0x002b 00043 (<autogenerated>:1)	ADDQ	$16, SP
	0x002f 00047 (<autogenerated>:1)	POPQ	BP
	0x0030 00048 (<autogenerated>:1)	RET
	0x0031 00049 (<autogenerated>:1)	NOP
	0x0031 00049 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0031 00049 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0031 00049 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x0036 00054 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x003b 00059 (<autogenerated>:1)	MOVQ	CX, 24(SP)
	0x0040 00064 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x0045 00069 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0045 00069 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x004a 00074 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x004f 00079 (<autogenerated>:1)	MOVQ	24(SP), CX
	0x0054 00084 (<autogenerated>:1)	JMP	0
	0x0000 49 3b 66 10 76 2b 55 48 89 e5 48 83 ec 10 48 89  I;f.v+UH..H...H.
	0x0010 44 24 20 48 89 5c 24 28 48 89 4c 24 30 84 00 48  D$ H.\$(H.L$0..H
	0x0020 8b 50 18 48 89 d8 48 89 cb ff d2 48 83 c4 10 5d  .P.H..H....H...]
	0x0030 c3 48 89 44 24 08 48 89 5c 24 10 48 89 4c 24 18  .H.D$.H.\$.H.L$.
	0x0040 e8 00 00 00 00 48 8b 44 24 08 48 8b 5c 24 10 48  .....H.D$.H.\$.H
	0x0050 8b 4c 24 18 eb aa                                .L$...
	rel 2+0 t=R_USEIFACEMETHOD type:main.Payer+96
	rel 41+0 t=R_CALLIND +0
	rel 65+4 t=R_CALL runtime.morestack_noctxt+0
type:.eq.main.Car STEXT dupok size=191 args=0x10 locals=0x48 funcid=0x0 align=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	type:.eq.main.Car(SB), DUPOK|ABIInternal, $72-16
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	161
	0x000a 00010 (<autogenerated>:1)	PCDATA	$0, $-1
	0x000a 00010 (<autogenerated>:1)	PUSHQ	BP
	0x000b 00011 (<autogenerated>:1)	MOVQ	SP, BP
	0x000e 00014 (<autogenerated>:1)	SUBQ	$64, SP
	0x0012 00018 (<autogenerated>:1)	FUNCDATA	$0, gclocals·TswRR9Pia9Wsluv5u1sUnA==(SB)
	0x0012 00018 (<autogenerated>:1)	FUNCDATA	$1, gclocals·A8pLD7vqL0qgY87/mhUKyA==(SB)
	0x0012 00018 (<autogenerated>:1)	FUNCDATA	$5, type:.eq.main.Car.arginfo1(SB)
	0x0012 00018 (<autogenerated>:1)	MOVQ	AX, main.p+80(SP)
	0x0017 00023 (<autogenerated>:1)	MOVQ	BX, main.q+88(SP)
	0x001c 00028 (<autogenerated>:1)	MOVB	$0, main.r+31(SP)
	0x0021 00033 (<autogenerated>:1)	MOVQ	8(AX), DX
	0x0025 00037 (<autogenerated>:1)	MOVQ	DX, main..autotmp_3+40(SP)
	0x002a 00042 (<autogenerated>:1)	MOVQ	8(BX), SI
	0x002e 00046 (<autogenerated>:1)	MOVQ	SI, main..autotmp_4+32(SP)
	0x0033 00051 (<autogenerated>:1)	CMPQ	DX, SI
	0x0036 00054 (<autogenerated>:1)	JEQ	58
	0x0038 00056 (<autogenerated>:1)	JMP	141
	0x003a 00058 (<autogenerated>:1)	MOVQ	main.p+80(SP), DX
	0x003f 00063 (<autogenerated>:1)	MOVQ	16(DX), DX
	0x0043 00067 (<autogenerated>:1)	CMPQ	16(BX), DX
	0x0047 00071 (<autogenerated>:1)	JEQ	75
	0x0049 00073 (<autogenerated>:1)	JMP	139
	0x004b 00075 (<autogenerated>:1)	MOVQ	main.p+80(SP), DX
	0x0050 00080 (<autogenerated>:1)	MOVQ	8(DX), DX
	0x0054 00084 (<autogenerated>:1)	MOVQ	DX, main..autotmp_4+32(SP)
	0x0059 00089 (<autogenerated>:1)	MOVQ	main.p+80(SP), DX
	0x005e 00094 (<autogenerated>:1)	MOVQ	(DX), DX
	0x0061 00097 (<autogenerated>:1)	MOVQ	DX, main..autotmp_5+56(SP)
	0x0066 00102 (<autogenerated>:1)	MOVQ	main.q+88(SP), DX
	0x006b 00107 (<autogenerated>:1)	MOVQ	(DX), BX
	0x006e 00110 (<autogenerated>:1)	MOVQ	BX, main..autotmp_6+48(SP)
	0x0073 00115 (<autogenerated>:1)	MOVQ	main..autotmp_4+32(SP), CX
	0x0078 00120 (<autogenerated>:1)	MOVQ	main..autotmp_5+56(SP), AX
	0x007d 00125 (<autogenerated>:1)	PCDATA	$1, $1
	0x007d 00125 (<autogenerated>:1)	NOP
	0x0080 00128 (<autogenerated>:1)	CALL	runtime.memequal(SB)
	0x0085 00133 (<autogenerated>:1)	MOVB	AL, main.r+31(SP)
	0x0089 00137 (<autogenerated>:1)	JMP	150
	0x008b 00139 (<autogenerated>:1)	JMP	143
	0x008d 00141 (<autogenerated>:1)	JMP	143
	0x008f 00143 (<autogenerated>:1)	MOVB	$0, main.r+31(SP)
	0x0094 00148 (<autogenerated>:1)	JMP	150
	0x0096 00150 (<autogenerated>:1)	MOVBLZX	main.r+31(SP), AX
	0x009b 00155 (<autogenerated>:1)	ADDQ	$64, SP
	0x009f 00159 (<autogenerated>:1)	POPQ	BP
	0x00a0 00160 (<autogenerated>:1)	RET
	0x00a1 00161 (<autogenerated>:1)	NOP
	0x00a1 00161 (<autogenerated>:1)	PCDATA	$1, $-1
	0x00a1 00161 (<autogenerated>:1)	PCDATA	$0, $-2
	0x00a1 00161 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x00a6 00166 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x00ab 00171 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x00b0 00176 (<autogenerated>:1)	PCDATA	$0, $-1
	0x00b0 00176 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x00b5 00181 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x00ba 00186 (<autogenerated>:1)	JMP	0
	0x0000 49 3b 66 10 0f 86 97 00 00 00 55 48 89 e5 48 83  I;f.......UH..H.
	0x0010 ec 40 48 89 44 24 50 48 89 5c 24 58 c6 44 24 1f  .@H.D$PH.\$X.D$.
	0x0020 00 48 8b 50 08 48 89 54 24 28 48 8b 73 08 48 89  .H.P.H.T$(H.s.H.
	0x0030 74 24 20 48 39 f2 74 02 eb 53 48 8b 54 24 50 48  t$ H9.t..SH.T$PH
	0x0040 8b 52 10 48 39 53 10 74 02 eb 40 48 8b 54 24 50  .R.H9S.t..@H.T$P
	0x0050 48 8b 52 08 48 89 54 24 20 48 8b 54 24 50 48 8b  H.R.H.T$ H.T$PH.
	0x0060 12 48 89 54 24 38 48 8b 54 24 58 48 8b 1a 48 89  .H.T$8H.T$XH..H.
	0x0070 5c 24 30 48 8b 4c 24 20 48 8b 44 24 38 0f 1f 00  \$0H.L$ H.D$8...
	0x0080 e8 00 00 00 00 88 44 24 1f eb 0b eb 02 eb 00 c6  ......D$........
	0x0090 44 24 1f 00 eb 00 0f b6 44 24 1f 48 83 c4 40 5d  D$......D$.H..@]
	0x00a0 c3 48 89 44 24 08 48 89 5c 24 10 e8 00 00 00 00  .H.D$.H.\$......
	0x00b0 48 8b 44 24 08 48 8b 5c 24 10 e9 41 ff ff ff     H.D$.H.\$..A...
	rel 129+4 t=R_CALL runtime.memequal+0
	rel 172+4 t=R_CALL runtime.morestack_noctxt+0
type:.eq.[2]interface {} STEXT dupok size=212 args=0x10 locals=0x50 funcid=0x0 align=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	type:.eq.[2]interface {}(SB), DUPOK|ABIInternal, $80-16
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	181
	0x000a 00010 (<autogenerated>:1)	PCDATA	$0, $-1
	0x000a 00010 (<autogenerated>:1)	PUSHQ	BP
	0x000b 00011 (<autogenerated>:1)	MOVQ	SP, BP
	0x000e 00014 (<autogenerated>:1)	SUBQ	$72, SP
	0x0012 00018 (<autogenerated>:1)	FUNCDATA	$0, gclocals·TswRR9Pia9Wsluv5u1sUnA==(SB)
	0x0012 00018 (<autogenerated>:1)	FUNCDATA	$1, gclocals·EYsUeQHkIlelPup/TMZjqA==(SB)
	0x0012 00018 (<autogenerated>:1)	FUNCDATA	$5, type:.eq.[2]interface {}.arginfo1(SB)
	0x0012 00018 (<autogenerated>:1)	MOVQ	AX, main.p+88(SP)
	0x0017 00023 (<autogenerated>:1)	MOVQ	BX, main.q+96(SP)
	0x001c 00028 (<autogenerated>:1)	MOVB	$0, main.r+31(SP)
	0x0021 00033 (<autogenerated>:1)	MOVQ	$0, main..autotmp_3+32(SP)
	0x002a 00042 (<autogenerated>:1)	JMP	44
	0x002c 00044 (<autogenerated>:1)	CMPQ	main..autotmp_3+32(SP), $2
	0x0032 00050 (<autogenerated>:1)	JLT	54
	0x0034 00052 (<autogenerated>:1)	JMP	163
	0x0036 00054 (<autogenerated>:1)	MOVQ	main..autotmp_3+32(SP), DX
	0x003b 00059 (<autogenerated>:1)	SHLQ	$4, DX
	0x003f 00063 (<autogenerated>:1)	ADDQ	main.q+96(SP), DX
	0x0044 00068 (<autogenerated>:1)	MOVQ	(DX), SI
	0x0047 00071 (<autogenerated>:1)	MOVQ	8(DX), DX
	0x004b 00075 (<autogenerated>:1)	MOVQ	SI, main..autotmp_4+56(SP)
	0x0050 00080 (<autogenerated>:1)	MOVQ	DX, main..autotmp_4+64(SP)
	0x0055 00085 (<autogenerated>:1)	MOVQ	main..autotmp_3+32(SP), DX
	0x005a 00090 (<autogenerated>:1)	SHLQ	$4, DX
	0x005e 00094 (<autogenerated>:1)	ADDQ	main.p+88(SP), DX
	0x0063 00099 (<autogenerated>:1)	MOVQ	(DX), AX
	0x0066 00102 (<autogenerated>:1)	MOVQ	8(DX), BX
	0x006a 00106 (<autogenerated>:1)	MOVQ	AX, main..autotmp_5+40(SP)
	0x006f 00111 (<autogenerated>:1)	MOVQ	BX, main..autotmp_5+48(SP)
	0x0074 00116 (<autogenerated>:1)	CMPQ	main..autotmp_4+56(SP), AX
	0x0079 00121 (<autogenerated>:1)	JEQ	125
	0x007b 00123 (<autogenerated>:1)	JMP	152
	0x007d 00125 (<autogenerated>:1)	MOVQ	main..autotmp_4+64(SP), CX
	0x0082 00130 (<autogenerated>:1)	PCDATA	$1, $0
	0x0082 00130 (<autogenerated>:1)	CALL	runtime.efaceeq(SB)
	0x0087 00135 (<autogenerated>:1)	TESTB	AL, AL
	0x0089 00137 (<autogenerated>:1)	JNE	141
	0x008b 00139 (<autogenerated>:1)	JMP	150
	0x008d 00141 (<autogenerated>:1)	INCQ	main..autotmp_3+32(SP)
	0x0092 00146 (<autogenerated>:1)	JMP	148
	0x0094 00148 (<autogenerated>:1)	JMP	44
	0x0096 00150 (<autogenerated>:1)	JMP	154
	0x0098 00152 (<autogenerated>:1)	JMP	154
	0x009a 00154 (<autogenerated>:1)	JMP	156
	0x009c 00156 (<autogenerated>:1)	MOVB	$0, main.r+31(SP)
	0x00a1 00161 (<autogenerated>:1)	JMP	170
	0x00a3 00163 (<autogenerated>:1)	MOVB	$1, main.r+31(SP)
	0x00a8 00168 (<autogenerated>:1)	JMP	170
	0x00aa 00170 (<autogenerated>:1)	MOVBLZX	main.r+31(SP), AX
	0x00af 00175 (<autogenerated>:1)	ADDQ	$72, SP
	0x00b3 00179 (<autogenerated>:1)	POPQ	BP
	0x00b4 00180 (<autogenerated>:1)	RET
	0x00b5 00181 (<autogenerated>:1)	NOP
	0x00b5 00181 (<autogenerated>:1)	PCDATA	$1, $-1
	0x00b5 00181 (<autogenerated>:1)	PCDATA	$0, $-2
	0x00b5 00181 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x00ba 00186 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x00bf 00191 (<autogenerated>:1)	NOP
	0x00c0 00192 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x00c5 00197 (<autogenerated>:1)	PCDATA	$0, $-1
	0x00c5 00197 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x00ca 00202 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x00cf 00207 (<autogenerated>:1)	JMP	0
	0x0000 49 3b 66 10 0f 86 ab 00 00 00 55 48 89 e5 48 83  I;f.......UH..H.
	0x0010 ec 48 48 89 44 24 58 48 89 5c 24 60 c6 44 24 1f  .HH.D$XH.\$`.D$.
	0x0020 00 48 c7 44 24 20 00 00 00 00 eb 00 48 83 7c 24  .H.D$ ......H.|$
	0x0030 20 02 7c 02 eb 6d 48 8b 54 24 20 48 c1 e2 04 48   .|..mH.T$ H...H
	0x0040 03 54 24 60 48 8b 32 48 8b 52 08 48 89 74 24 38  .T$`H.2H.R.H.t$8
	0x0050 48 89 54 24 40 48 8b 54 24 20 48 c1 e2 04 48 03  H.T$@H.T$ H...H.
	0x0060 54 24 58 48 8b 02 48 8b 5a 08 48 89 44 24 28 48  T$XH..H.Z.H.D$(H
	0x0070 89 5c 24 30 48 39 44 24 38 74 02 eb 1b 48 8b 4c  .\$0H9D$8t...H.L
	0x0080 24 40 e8 00 00 00 00 84 c0 75 02 eb 09 48 ff 44  $@.......u...H.D
	0x0090 24 20 eb 00 eb 96 eb 02 eb 00 eb 00 c6 44 24 1f  $ ...........D$.
	0x00a0 00 eb 07 c6 44 24 1f 01 eb 00 0f b6 44 24 1f 48  ....D$......D$.H
	0x00b0 83 c4 48 5d c3 48 89 44 24 08 48 89 5c 24 10 90  ..H].H.D$.H.\$..
	0x00c0 e8 00 00 00 00 48 8b 44 24 08 48 8b 5c 24 10 e9  .....H.D$.H.\$..
	0x00d0 2c ff ff ff                                      ,...
	rel 131+4 t=R_CALL runtime.efaceeq+0
	rel 193+4 t=R_CALL runtime.morestack_noctxt+0
go:cuinfo.producer.main SDWARFCUINFO dupok size=0
	0x0000 2d 4e 20 2d 6c 20 72 65 67 61 62 69              -N -l regabi
runtime.interequal·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR runtime.interequal+0
runtime.memequal64·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR runtime.memequal64+0
runtime.gcbits.0100000000000000 SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
type:.namedata.*main.Payer. SRODATA dupok size=13
	0x0000 01 0b 2a 6d 61 69 6e 2e 50 61 79 65 72           ..*main.Payer
type:*main.Payer SRODATA size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 20 ef 7e b6 28 08 08 16 00 00 00 00 00 00 00 00   .~.(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Payer.+0
	rel 48+8 t=R_ADDR type:main.Payer+0
runtime.gcbits.0200000000000000 SRODATA dupok size=8
	0x0000 02 00 00 00 00 00 00 00                          ........
type:.namedata.*func(int)- SRODATA dupok size=12
	0x0000 00 0a 2a 66 75 6e 63 28 69 6e 74 29              ..*func(int)
type:*func(int) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 5e eb a8 b3 28 08 08 16 00 00 00 00 00 00 00 00  ^...(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(int)-+0
	rel 48+8 t=R_ADDR type:func(int)+0
type:func(int) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 f2 a3 a2 b7 22 08 08 13 00 00 00 00 00 00 00 00  ...."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(int)-+0
	rel 44+4 t=RelocType(-32763) type:*func(int)+0
	rel 56+8 t=R_ADDR type:int+0
type:.importpath.main. SRODATA dupok size=6
	0x0000 00 04 6d 61 69 6e                                ..main
type:.namedata.Pay. SRODATA dupok size=5
	0x0000 01 03 50 61 79                                   ..Pay
type:main.Payer SRODATA size=104
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 aa 68 29 30 07 08 08 14 00 00 00 00 00 00 00 00  .h)0............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 18 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.interequal·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0200000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Payer.+0
	rel 44+4 t=R_ADDROFF type:*main.Payer+0
	rel 48+8 t=R_ADDR type:.importpath.main.+0
	rel 56+8 t=R_ADDR type:main.Payer+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+4 t=R_ADDROFF type:.namedata.Pay.+0
	rel 100+4 t=R_ADDROFF type:func(int)+0
type:.namedata.*main.WeChatPay. SRODATA dupok size=17
	0x0000 01 0f 2a 6d 61 69 6e 2e 57 65 43 68 61 74 50 61  ..*main.WeChatPa
	0x0010 79                                               y
runtime.gcbits. SRODATA dupok size=0
type:.namedata.Balance. SRODATA dupok size=9
	0x0000 01 07 42 61 6c 61 6e 63 65                       ..Balance
type:main.WeChatPay SRODATA size=120
	0x0000 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 5a 50 86 09 0f 08 08 19 00 00 00 00 00 00 00 00  ZP..............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 28 00 00 00 00 00 00 00  ........(.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.WeChatPay.+0
	rel 44+4 t=R_ADDROFF type:*main.WeChatPay+0
	rel 56+8 t=R_ADDR type:main.WeChatPay+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+8 t=R_ADDR type:.namedata.Balance.+0
	rel 104+8 t=R_ADDR type:int+0
type:.namedata.*func(*main.WeChatPay, int)- SRODATA dupok size=29
	0x0000 00 1b 2a 66 75 6e 63 28 2a 6d 61 69 6e 2e 57 65  ..*func(*main.We
	0x0010 43 68 61 74 50 61 79 2c 20 69 6e 74 29           ChatPay, int)
type:*func(*main.WeChatPay, int) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 d1 82 ab c8 28 08 08 16 00 00 00 00 00 00 00 00  ....(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.WeChatPay, int)-+0
	rel 48+8 t=R_ADDR type:func(*main.WeChatPay, int)+0
type:func(*main.WeChatPay, int) SRODATA dupok size=72
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 23 c9 cf 0b 22 08 08 13 00 00 00 00 00 00 00 00  #..."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00                          ........
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.WeChatPay, int)-+0
	rel 44+4 t=RelocType(-32763) type:*func(*main.WeChatPay, int)+0
	rel 56+8 t=R_ADDR type:*main.WeChatPay+0
	rel 64+8 t=R_ADDR type:int+0
type:*main.WeChatPay SRODATA size=88
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 53 26 9a 91 29 08 08 16 00 00 00 00 00 00 00 00  S&..)...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 01 00 01 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.WeChatPay.+0
	rel 48+8 t=R_ADDR type:main.WeChatPay+0
	rel 56+4 t=R_ADDROFF type:.importpath.main.+0
	rel 72+4 t=R_ADDROFF type:.namedata.Pay.+0
	rel 76+4 t=R_METHODOFF type:func(int)+0
	rel 80+4 t=R_METHODOFF main.(*WeChatPay).Pay+0
	rel 84+4 t=R_METHODOFF main.(*WeChatPay).Pay+0
go:itab.*main.WeChatPay,main.Payer SRODATA dupok size=32
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 53 26 9a 91 00 00 00 00 00 00 00 00 00 00 00 00  S&..............
	rel 0+8 t=R_ADDR type:main.Payer+0
	rel 8+8 t=R_ADDR type:*main.WeChatPay+0
	rel 24+8 t=RelocType(-32767) main.(*WeChatPay).Pay+0
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
main..inittask SNOPTRDATA size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+0 t=R_INITORDER fmt..inittask+0
type:.namedata.*main.Counter. SRODATA dupok size=15
	0x0000 01 0d 2a 6d 61 69 6e 2e 43 6f 75 6e 74 65 72     ..*main.Counter
type:.namedata.*func(*main.Counter)- SRODATA dupok size=22
	0x0000 00 14 2a 66 75 6e 63 28 2a 6d 61 69 6e 2e 43 6f  ..*func(*main.Co
	0x0010 75 6e 74 65 72 29                                unter)
type:*func(*main.Counter) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 35 17 62 fd 28 08 08 16 00 00 00 00 00 00 00 00  5.b.(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Counter)-+0
	rel 48+8 t=R_ADDR type:func(*main.Counter)+0
type:func(*main.Counter) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 2b 08 ea df 22 08 08 13 00 00 00 00 00 00 00 00  +..."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Counter)-+0
	rel 44+4 t=RelocType(-32763) type:*func(*main.Counter)+0
	rel 56+8 t=R_ADDR type:*main.Counter+0
type:.namedata.AddPtr. SRODATA dupok size=8
	0x0000 01 06 41 64 64 50 74 72                          ..AddPtr
type:.namedata.*func()- SRODATA dupok size=9
	0x0000 00 07 2a 66 75 6e 63 28 29                       ..*func()
type:*func() SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 75 53 d6 d8 28 08 08 16 00 00 00 00 00 00 00 00  uS..(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func()-+0
	rel 48+8 t=R_ADDR type:func()+0
type:func() SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 fe 05 46 7f 22 08 08 13 00 00 00 00 00 00 00 00  ..F."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00                                      ....
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func()-+0
	rel 44+4 t=RelocType(-32763) type:*func()+0
type:.namedata.AddValue. SRODATA dupok size=10
	0x0000 01 08 41 64 64 56 61 6c 75 65                    ..AddValue
type:*main.Counter SRODATA size=104
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 d4 11 9e c8 29 08 08 16 00 00 00 00 00 00 00 00  ....)...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 02 00 02 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Counter.+0
	rel 48+8 t=R_ADDR type:main.Counter+0
	rel 56+4 t=R_ADDROFF type:.importpath.main.+0
	rel 72+4 t=R_ADDROFF type:.namedata.AddPtr.+0
	rel 76+4 t=R_METHODOFF type:func()+0
	rel 80+4 t=R_METHODOFF main.(*Counter).AddPtr+0
	rel 84+4 t=R_METHODOFF main.(*Counter).AddPtr+0
	rel 88+4 t=R_ADDROFF type:.namedata.AddValue.+0
	rel 92+4 t=R_METHODOFF type:func()+0
	rel 96+4 t=R_METHODOFF main.(*Counter).AddValue+0
	rel 100+4 t=R_METHODOFF main.(*Counter).AddValue+0
type:.namedata.Count. SRODATA dupok size=7
	0x0000 01 05 43 6f 75 6e 74                             ..Count
type:.namedata.*func(main.Counter)- SRODATA dupok size=21
	0x0000 00 13 2a 66 75 6e 63 28 6d 61 69 6e 2e 43 6f 75  ..*func(main.Cou
	0x0010 6e 74 65 72 29                                   nter)
type:*func(main.Counter) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 33 fb da ce 28 08 08 16 00 00 00 00 00 00 00 00  3...(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(main.Counter)-+0
	rel 48+8 t=R_ADDR type:func(main.Counter)+0
type:func(main.Counter) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 b2 5a 27 95 22 08 08 13 00 00 00 00 00 00 00 00  .Z'."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(main.Counter)-+0
	rel 44+4 t=RelocType(-32763) type:*func(main.Counter)+0
	rel 56+8 t=R_ADDR type:main.Counter+0
type:main.Counter SRODATA size=136
	0x0000 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 c8 46 48 7d 0f 08 08 19 00 00 00 00 00 00 00 00  .FH}............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 01 00 01 00 28 00 00 00 00 00 00 00  ........(.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0080 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Counter.+0
	rel 44+4 t=R_ADDROFF type:*main.Counter+0
	rel 56+8 t=R_ADDR type:main.Counter+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+8 t=R_ADDR type:.namedata.Count.+0
	rel 104+8 t=R_ADDR type:int+0
	rel 120+4 t=R_ADDROFF type:.namedata.AddValue.+0
	rel 124+4 t=R_METHODOFF type:func()+0
	rel 128+4 t=R_METHODOFF main.(*Counter).AddValue+0
	rel 132+4 t=R_METHODOFF main.Counter.AddValue+0
type:.namedata.*main.Engine. SRODATA dupok size=14
	0x0000 01 0c 2a 6d 61 69 6e 2e 45 6e 67 69 6e 65        ..*main.Engine
type:.namedata.*func(*main.Engine)- SRODATA dupok size=21
	0x0000 00 13 2a 66 75 6e 63 28 2a 6d 61 69 6e 2e 45 6e  ..*func(*main.En
	0x0010 67 69 6e 65 29                                   gine)
type:*func(*main.Engine) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 5e af 55 9d 28 08 08 16 00 00 00 00 00 00 00 00  ^.U.(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Engine)-+0
	rel 48+8 t=R_ADDR type:func(*main.Engine)+0
type:func(*main.Engine) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 0a 84 fc c9 22 08 08 13 00 00 00 00 00 00 00 00  ...."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Engine)-+0
	rel 44+4 t=RelocType(-32763) type:*func(*main.Engine)+0
	rel 56+8 t=R_ADDR type:*main.Engine+0
type:.namedata.Start. SRODATA dupok size=7
	0x0000 01 05 53 74 61 72 74                             ..Start
type:*main.Engine SRODATA size=88
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 97 46 2d 67 29 08 08 16 00 00 00 00 00 00 00 00  .F-g)...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 01 00 01 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Engine.+0
	rel 48+8 t=R_ADDR type:main.Engine+0
	rel 56+4 t=R_ADDROFF type:.importpath.main.+0
	rel 72+4 t=R_ADDROFF type:.namedata.Start.+0
	rel 76+4 t=R_METHODOFF type:func()+0
	rel 80+4 t=R_METHODOFF main.(*Engine).Start+0
	rel 84+4 t=R_METHODOFF main.(*Engine).Start+0
type:.namedata.Power. SRODATA dupok size=7
	0x0000 01 05 50 6f 77 65 72                             ..Power
type:main.Engine SRODATA size=120
	0x0000 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 67 e9 54 6e 0f 08 08 19 00 00 00 00 00 00 00 00  g.Tn............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 28 00 00 00 00 00 00 00  ........(.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Engine.+0
	rel 44+4 t=R_ADDROFF type:*main.Engine+0
	rel 56+8 t=R_ADDR type:main.Engine+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+8 t=R_ADDR type:.namedata.Power.+0
	rel 104+8 t=R_ADDR type:int+0
type:.eqfunc.main.Car SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR type:.eq.main.Car+0
type:.namedata.*main.Car. SRODATA dupok size=11
	0x0000 01 09 2a 6d 61 69 6e 2e 43 61 72                 ..*main.Car
type:.namedata.*func(*main.Car)- SRODATA dupok size=18
	0x0000 00 10 2a 66 75 6e 63 28 2a 6d 61 69 6e 2e 43 61  ..*func(*main.Ca
	0x0010 72 29                                            r)
type:*func(*main.Car) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 bf ab 99 eb 28 08 08 16 00 00 00 00 00 00 00 00  ....(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Car)-+0
	rel 48+8 t=R_ADDR type:func(*main.Car)+0
type:func(*main.Car) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 35 1c 9e 1e 22 08 08 13 00 00 00 00 00 00 00 00  5..."...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Car)-+0
	rel 44+4 t=RelocType(-32763) type:*func(*main.Car)+0
	rel 56+8 t=R_ADDR type:*main.Car+0
type:*main.Car SRODATA size=88
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 9c 23 4e 38 29 08 08 16 00 00 00 00 00 00 00 00  .#N8)...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 01 00 01 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Car.+0
	rel 48+8 t=R_ADDR type:main.Car+0
	rel 56+4 t=R_ADDROFF type:.importpath.main.+0
	rel 72+4 t=R_ADDROFF type:.namedata.Start.+0
	rel 76+4 t=R_METHODOFF type:func()+0
	rel 80+4 t=R_METHODOFF main.(*Car).Start+0
	rel 84+4 t=R_METHODOFF main.(*Car).Start+0
type:.namedata.Name. SRODATA dupok size=6
	0x0000 01 04 4e 61 6d 65                                ..Name
type:.namedata.Engine..embedded SRODATA dupok size=8
	0x0000 09 06 45 6e 67 69 6e 65                          ..Engine
type:main.Car SRODATA size=144
	0x0000 18 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 37 c0 ff 9a 07 08 08 19 00 00 00 00 00 00 00 00  7...............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 02 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 40 00 00 00 00 00 00 00  ........@.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0080 00 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	rel 24+8 t=R_ADDR type:.eqfunc.main.Car+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Car.+0
	rel 44+4 t=R_ADDROFF type:*main.Car+0
	rel 56+8 t=R_ADDR type:main.Car+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+8 t=R_ADDR type:.namedata.Name.+0
	rel 104+8 t=R_ADDR type:string+0
	rel 120+8 t=R_ADDR type:.namedata.Engine..embedded+0
	rel 128+8 t=R_ADDR type:main.Engine+0
go:string."引擎启动..." SRODATA dupok size=15
	0x0000 e5 bc 95 e6 93 8e e5 90 af e5 8a a8 2e 2e 2e     ...............
main..stmp_0 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 0f 00 00 00 00 00 00 00  ................
	rel 0+8 t=R_ADDR go:string."引擎启动..."+0
go:string."%s 跑车正在启动...\n" SRODATA dupok size=25
	0x0000 25 73 20 e8 b7 91 e8 bd a6 e6 ad a3 e5 9c a8 e5  %s .............
	0x0010 90 af e5 8a a8 2e 2e 2e 0a                       .........
go:string."微信支付成功，扣除 %d 元， 余额 %d 元\n" SRODATA dupok size=52
	0x0000 e5 be ae e4 bf a1 e6 94 af e4 bb 98 e6 88 90 e5  ................
	0x0010 8a 9f ef bc 8c e6 89 a3 e9 99 a4 20 25 64 20 e5  ........... %d .
	0x0020 85 83 ef bc 8c 20 e4 bd 99 e9 a2 9d 20 25 64 20  ..... ...... %d 
	0x0030 e5 85 83 0a                                      ....
go:string."----------- 方法基本使用 ----------" SRODATA dupok size=41
	0x0000 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d 20 e6 96 b9 e6  ----------- ....
	0x0010 b3 95 e5 9f ba e6 9c ac e4 bd bf e7 94 a8 20 2d  .............. -
	0x0020 2d 2d 2d 2d 2d 2d 2d 2d 2d                       ---------
go:string."指针方法执行后: " SRODATA dupok size=23
	0x0000 e6 8c 87 e9 92 88 e6 96 b9 e6 b3 95 e6 89 a7 e8  ................
	0x0010 a1 8c e5 90 8e 3a 20                             .....: 
go:string."值方法执行后: " SRODATA dupok size=20
	0x0000 e5 80 bc e6 96 b9 e6 b3 95 e6 89 a7 e8 a1 8c e5  ................
	0x0010 90 8e 3a 20                                      ..: 
go:string."----------- 组合 ----------" SRODATA dupok size=29
	0x0000 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d 20 e7 bb 84 e5  ----------- ....
	0x0010 90 88 20 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d           .. ----------
go:string."马力：" SRODATA dupok size=9
	0x0000 e9 a9 ac e5 8a 9b ef bc 9a                       .........
go:string."----------- 方法集 ----------" SRODATA dupok size=32
	0x0000 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d 20 e6 96 b9 e6  ----------- ....
	0x0010 b3 95 e9 9b 86 20 2d 2d 2d 2d 2d 2d 2d 2d 2d 2d  ..... ----------
go:string."Audi" SRODATA dupok size=4
	0x0000 41 75 64 69                                      Audi
main..stmp_1 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 29 00 00 00 00 00 00 00  ........).......
	rel 0+8 t=R_ADDR go:string."----------- 方法基本使用 ----------"+0
main..stmp_2 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 17 00 00 00 00 00 00 00  ................
	rel 0+8 t=R_ADDR go:string."指针方法执行后: "+0
main..stmp_3 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 14 00 00 00 00 00 00 00  ................
	rel 0+8 t=R_ADDR go:string."值方法执行后: "+0
main..stmp_4 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 1d 00 00 00 00 00 00 00  ................
	rel 0+8 t=R_ADDR go:string."----------- 组合 ----------"+0
main..stmp_5 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 09 00 00 00 00 00 00 00  ................
	rel 0+8 t=R_ADDR go:string."马力："+0
main..stmp_6 SRODATA static size=16
	0x0000 00 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00  ........ .......
	rel 0+8 t=R_ADDR go:string."----------- 方法集 ----------"+0
type:.namedata.*[1]interface {}- SRODATA dupok size=18
	0x0000 00 10 2a 5b 31 5d 69 6e 74 65 72 66 61 63 65 20  ..*[1]interface 
	0x0010 7b 7d                                            {}
runtime.nilinterequal·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR runtime.nilinterequal+0
type:[1]interface {} SRODATA dupok size=72
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 6e df 95 c2 02 08 08 11 00 00 00 00 00 00 00 00  n...............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.nilinterequal·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0200000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*[1]interface {}-+0
	rel 44+4 t=RelocType(-32763) type:*[1]interface {}+0
	rel 48+8 t=R_ADDR type:interface {}+0
	rel 56+8 t=R_ADDR type:[]interface {}+0
type:*[1]interface {} SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 a8 f1 a8 c9 28 08 08 16 00 00 00 00 00 00 00 00  ....(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*[1]interface {}-+0
	rel 48+8 t=R_ADDR type:[1]interface {}+0
type:.namedata.*[2]interface {}- SRODATA dupok size=18
	0x0000 00 10 2a 5b 32 5d 69 6e 74 65 72 66 61 63 65 20  ..*[2]interface 
	0x0010 7b 7d                                            {}
type:.eqfunc.[2]interface {} SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR type:.eq.[2]interface {}+0
runtime.gcbits.0a00000000000000 SRODATA dupok size=8
	0x0000 0a 00 00 00 00 00 00 00                          ........
type:[2]interface {} SRODATA dupok size=72
	0x0000 20 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00   ....... .......
	0x0010 0a 0c 4b 4b 02 08 08 11 00 00 00 00 00 00 00 00  ..KK............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 02 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR type:.eqfunc.[2]interface {}+0
	rel 32+8 t=R_ADDR runtime.gcbits.0a00000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*[2]interface {}-+0
	rel 44+4 t=RelocType(-32763) type:*[2]interface {}+0
	rel 48+8 t=R_ADDR type:interface {}+0
	rel 56+8 t=R_ADDR type:[]interface {}+0
type:*[2]interface {} SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 53 23 94 ff 28 08 08 16 00 00 00 00 00 00 00 00  S#..(...........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*[2]interface {}-+0
	rel 48+8 t=R_ADDR type:[2]interface {}+0
gclocals·wvjpxkknJ4nY1JtrArJJaw== SRODATA dupok size=10
	0x0000 02 00 00 00 01 00 00 00 01 00                    ..........
gclocals·J26BEvPExEQhJvjp9E8Whg== SRODATA dupok size=8
	0x0000 02 00 00 00 00 00 00 00                          ........
main.(*Counter).AddPtr.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
gclocals·g5+hNtRBP6YXNjfog7aZjQ== SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
main.Counter.AddValue.arginfo1 SRODATA static dupok size=5
	0x0000 fe 00 08 fd ff                                   .....
gclocals·ceYgNIaaD8ow5EM5cNccoA== SRODATA dupok size=10
	0x0000 02 00 00 00 06 00 00 00 00 00                    ..........
main.(*Engine).Start.stkobj SRODATA static size=24
	0x0000 01 00 00 00 00 00 00 00 f0 ff ff ff 10 00 00 00  ................
	0x0010 10 00 00 00 00 00 00 00                          ........
	rel 20+4 t=R_ADDROFF runtime.gcbits.0200000000000000+0
main.(*Engine).Start.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
gclocals·Z8zdw/dq+fE82FieA9ctlQ== SRODATA dupok size=11
	0x0000 03 00 00 00 01 00 00 00 01 00 00                 ...........
gclocals·/4KVoIoAUVnbPMDuZquOrw== SRODATA dupok size=14
	0x0000 03 00 00 00 09 00 00 00 00 00 08 00 00 00        ..............
main.(*Car).Start.stkobj SRODATA static size=24
	0x0000 01 00 00 00 00 00 00 00 f0 ff ff ff 10 00 00 00  ................
	0x0010 10 00 00 00 00 00 00 00                          ........
	rel 20+4 t=R_ADDROFF runtime.gcbits.0200000000000000+0
main.(*Car).Start.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
gclocals·bUB0t99dbIOHL9YDh6V0CA== SRODATA dupok size=12
	0x0000 04 00 00 00 01 00 00 00 01 01 00 00              ............
gclocals·BcfHs9IhtErNtV2JassdQA== SRODATA dupok size=16
	0x0000 04 00 00 00 0a 00 00 00 00 00 04 00 04 00 00 00  ................
main.(*WeChatPay).Pay.stkobj SRODATA static size=24
	0x0000 01 00 00 00 00 00 00 00 e0 ff ff ff 20 00 00 00  ............ ...
	0x0010 20 00 00 00 00 00 00 00                           .......
	rel 20+4 t=R_ADDROFF runtime.gcbits.0a00000000000000+0
main.(*WeChatPay).Pay.arginfo1 SRODATA static dupok size=5
	0x0000 00 08 08 08 ff                                   .....
gclocals·Dj7m7VTqKq6fxJqfrrXabg== SRODATA dupok size=8
	0x0000 05 00 00 00 00 00 00 00                          ........
gclocals·MlHovJf8fukPSXh1aq01eg== SRODATA dupok size=38
	0x0000 05 00 00 00 2b 00 00 00 00 00 00 00 00 00 00 00  ....+...........
	0x0010 00 00 80 00 00 00 00 00 04 00 04 00 40 00 00 00  ............@...
	0x0020 04 00 00 00 00 00                                ......
main.main.stkobj SRODATA static size=56
	0x0000 03 00 00 00 00 00 00 00 b8 fe ff ff 18 00 00 00  ................
	0x0010 08 00 00 00 00 00 00 00 f8 fe ff ff 20 00 00 00  ............ ...
	0x0020 20 00 00 00 00 00 00 00 18 ff ff ff 10 00 00 00   ...............
	0x0030 10 00 00 00 00 00 00 00                          ........
	rel 20+4 t=R_ADDROFF runtime.gcbits.0100000000000000+0
	rel 36+4 t=R_ADDROFF runtime.gcbits.0a00000000000000+0
	rel 52+4 t=R_ADDROFF runtime.gcbits.0200000000000000+0
main.(*Counter).AddValue.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
gclocals·Ih7UaEzxol4qYE7mnpMg6w== SRODATA dupok size=10
	0x0000 02 00 00 00 02 00 00 00 02 00                    ..........
main.Payer.Pay.arginfo1 SRODATA static dupok size=9
	0x0000 fe 00 08 08 08 fd 10 08 ff                       .........
gclocals·TswRR9Pia9Wsluv5u1sUnA== SRODATA dupok size=10
	0x0000 02 00 00 00 02 00 00 00 03 00                    ..........
gclocals·A8pLD7vqL0qgY87/mhUKyA== SRODATA dupok size=10
	0x0000 02 00 00 00 02 00 00 00 00 00                    ..........
type:.eq.main.Car.arginfo1 SRODATA static dupok size=5
	0x0000 00 08 08 08 ff                                   .....
gclocals·EYsUeQHkIlelPup/TMZjqA== SRODATA dupok size=10
	0x0000 02 00 00 00 04 00 00 00 00 00                    ..........
type:.eq.[2]interface {}.arginfo1 SRODATA static dupok size=3
	0x0000 08 08 ff                                         ...
