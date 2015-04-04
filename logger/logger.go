package logger

import (
	"io"
	"log"
)

type Logger struct {
	r io.Reader
	w io.Writer
}

func New(r io.Reader, w io.Writer) (*Logger) {
	return &Logger{r, w}
}

func (l *Logger) Read(buf []byte) (n int, err error) {
	n, err = l.r.Read(buf)
	if err == nil {
		log.Printf("IN: %s\n", buf[:n])
	}
	return n, err
}

func (l* Logger) Write(buf []byte) (n int, err error) {
	log.Printf("OUT: %s\n", buf)
	n, err = l.w.Write(buf)
	return n, err
}

/*
func (l *Logger) Close() error {
	log.Printf("CLOSE\n")
l.w.Close()
	return l.r.Close() 
}
*/
