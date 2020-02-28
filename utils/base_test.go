package utils

import "testing"

func BenchmarkUnsafeBytesToString(b *testing.B) {
	type args struct {
		bytes []byte
	}
	tests := struct {
		name string
		args args
		want bool
	}{name:"BenchmarkUnsafeBytesToString", args:args{}}
	for i:=0;i<10000;i++ {
		tests.args.bytes = append(tests.args.bytes, 'a')
	}
	b.StopTimer()
	b.StartTimer()
	for i:=0;i<b.N;i++ {
		b.Run(tests.name, func(b *testing.B) {
			UnsafeBytesToString(tests.args.bytes)
		})
	}
}

func BenchmarkBytesToString(b *testing.B) {
	type args struct {
		bytes []byte
	}
	tests := struct {
		name string
		args args
		want bool
	}{name:"BenchmarkUnsafeBytesToString", args:args{}}
	for i:=0;i<10000;i++ {
		tests.args.bytes = append(tests.args.bytes, 'a')
	}
	b.StopTimer()
	b.StartTimer()
	for i:=0;i<b.N;i++ {
		b.Run(tests.name, func(b *testing.B) {
			BytesToString(tests.args.bytes)
		})
	}
}
