package log

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

func escape(x interface{}) string {
	r := fmt.Sprint(x)
	r = strings.Replace(r, "\"", "\\\"", -1)
	r = strings.Replace(r, "'", "\\'", -1)
	r = strings.Replace(r, "\t", "    ", -1)
	return r
}

func format(sub, pred interface{}, props ...interface{}) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\t%s", escape(sub), escape(pred))
	for i := 0; i < len(props); i += 2 {
		fmt.Fprintf(&buf, "\t%s=%s",
			escape(props[i]),
			escape(props[i+1]))
	}
	return buf.String()
}

func Print(sub, pred interface{}, props ...interface{}) {
	log.Print(format(sub, pred, props...))
}

func Fatal(sub, pred interface{}, props ...interface{}) {
	log.Fatal(format(sub, pred, props...))
}

func Panic(sub, pred interface{}, props ...interface{}) {
	log.Panic(format(sub, pred, props...))
}
