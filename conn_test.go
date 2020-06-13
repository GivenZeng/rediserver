package rediserver

import (
	"testing"
)

func TestRsfp2Cmd(t *testing.T) {
	commands := []string{
		"*2\r\n$3\r\ndel\r\n$3\r\nkey\r\n",
		"*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
		"*2\r\n$6\r\nexpire\r\n$3\r\nkey\r\n",

		"*3\r\n$4\r\nhget\r\n$3\r\nkey\r\n$5\r\nfield\r\n",
		"*4\r\n$4\r\nhset\r\n$3\r\nkey\r\n$5\r\nfield\r\n$3\r\nval\r\n",
		"*3\r\n$4\r\nhdel\r\n$3\r\nkey\r\n$5\r\nfield\r\n",
		"*2\r\n$7\r\nhgetall\r\n$3\r\nkey\r\n",
	}
	for _, cmd := range commands {
		rsfp2Cmd(t, []byte(cmd))
	}
}

func rsfp2Cmd(t *testing.T, buf []byte) {
	cmd, err := Rsfp2Cmd(buf)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("cmd = %s", cmd.String())
}
