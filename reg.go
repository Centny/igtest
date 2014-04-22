package igtest

import (
	"regexp"
)

var empty = regexp.MustCompile("[\\t\\ ]")

//
var c_reg_BC = regexp.MustCompile("^[\\ \\t]*BC([\\t\\ ]+.*)?$")
var c_reg_BC2 = regexp.MustCompile("^[\\ \\t]*.*[\\>\\+\\-\\<\\*\\/].*[\\ \\t]*")
var c_reg_SET = regexp.MustCompile("^[\\ \\t]*SET([\\t\\ ]+.*)?$")
var c_reg_GET = regexp.MustCompile("^[\\ \\t]*GET([\\t\\ ]+.*)?$")
var c_reg_HR = regexp.MustCompile("^[\\ \\t]*HR([\\t\\ ]+.*)?$")
var c_reg_HG = regexp.MustCompile("^[\\ \\t]*HG([\\t\\ ]+.*)?$")
var c_reg_HP = regexp.MustCompile("^[\\ \\t]*HP([\\t\\ ]+.*)?$")
var c_reg_EX = regexp.MustCompile("^[\\ \\t]*EX([\\t\\ ]+.*)?$")
var c_reg_P = regexp.MustCompile("^[\\ \\t]*P([\\t\\ ]+.*)?$")
var c_reg_SUB = regexp.MustCompile("^[\\ \\t]*SUB([\\t\\ ]+.*)?$")

// var c_reg_INC = regexp.MustCompile("^[\\ \\t]*INC([\\ \t]+.*)?$")
var c_reg_FOR = regexp.MustCompile("^[\\ \\t]*FOR([\\ \\t]+.*)?$")
var c_reg_ROF = regexp.MustCompile("^(.*[\\ \\t]+)?ROF[\\ \\t]*$")
var c_reg_e_ROF = regexp.MustCompile("[\\ \\t]*ROF[\\ \\t]*$")
var c_reg_EXP = regexp.MustCompile("^[\\ \\t]*\\@[^\\@]+.*$")
var c_reg_Y = regexp.MustCompile("^[\\ \\t]*Y([\\t\\ ]+.*)?$")
var c_reg_N = regexp.MustCompile("^[\\ \\t]*N([\\t\\ ]+.*)?$")
var c_reg_R = regexp.MustCompile("^[\\ \\t]*R([\\t\\ ]+.*)?$")
var c_reg_W = regexp.MustCompile("^[\\ \\t]*W([\\t\\ ]+.*)?$")
var c_reg_D = regexp.MustCompile("^[\\ \\t]*D([\\t\\ ]+.*)?$")
var c_reg_M = regexp.MustCompile("^[\\ \\t]*M([\\t\\ ]+.*)?$")
var c_reg_Assign1 = regexp.MustCompile("^[\\ \\t]*\\$[a-zA-Z_]+[_a-zA-Z0-9\\ \\t]*=.*$")
var c_reg_Assign2 = regexp.MustCompile("^[\\ \\t]*\\$\\([^\\(]+\\)[\\ \\t]*=.*$")
var escape = regexp.MustCompile("\\\\.")

//
var b_reg_cmt_beg = regexp.MustCompile("^[\\ \\t]*\\/\\*.*$")
var b_reg_cmt_end = regexp.MustCompile("^.*\\*\\/[\\ \\t]*$")
var l_reg_cmt = regexp.MustCompile("//.*\\n?$")

//
var e_reg_c = regexp.MustCompile("^[\\ \\t]*\\@\\[[^\\]]*\\]$[\\ \\t]*")
var e_reg_m = regexp.MustCompile("^[\\ \\t]*\\@\\{[^\\{]*\\}$[\\ \\t]*")

// var e_reg_m2 = regexp.MustCompile("^[\\ \\t]*.*[\\>\\+\\-\\<\\=\\*\\/][\\=]?.*[\\ \\t]*")
var e_reg_s = regexp.MustCompile("^[\\ \\t]*\\@\\([^\\)]*\\)$[\\ \\t]*")
var e_reg_v = regexp.MustCompile("^[\\ \\t]*\\$\\(?[^\\)\\(\\[\\]\\{\\}]*\\)?$")

// var e_reg_v_p = regexp.MustCompile("^[\\ \\t]*\\$\\([^\\)\\(\\[\\]\\{\\}]*\\)$")
var e_reg_kv = regexp.MustCompile("^[^=]*=[^=]*$")
var reg_http = regexp.MustCompile("^http\\://.*$")
