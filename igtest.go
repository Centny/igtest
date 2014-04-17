package igtest

const Version string = "Beta v0.1.0"

func Exec(fname string) error {
	return ExecCtx(fname, NewCtx(nil))
}

func ExecCtx(fname string, ctx *Ctx) error {
	c := Compiler{}
	err := c.Load(fname)
	if err != nil {
		return err
	}
	err = c.CompileAndExec(ctx, YesOnExeced)
	if err != nil {
		return err
	}
	return nil
}
