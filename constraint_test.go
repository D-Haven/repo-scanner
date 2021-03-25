package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestContains_Evaluate(t *testing.T) {
	f, err := ioutil.TempFile("./", "test_")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}

		err = os.Remove(f.Name())
		if err != nil {
			t.Error(err)
		}
	}()

	err = ioutil.WriteFile(f.Name(), []byte("this file contains foo"), 0666)
	if err != nil {
		t.Error(err)
	}

	fi, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Value string
	}
	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Contains 'foo'",
			fields: fields{Value: "foo"},
			args:   args{info: fi},
			want:   true,
		},
		{
			name:   "Does not contain 'bar'",
			fields: fields{Value: "bar"},
			args:   args{info: fi},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Contains{
				Value: tt.fields.Value,
			}
			if got := c.Evaluate(tt.args.info); got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNotBeEmpty_Evaluate(t *testing.T) {
	fText, err := ioutil.TempFile("./", "test_")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := fText.Close()
		if err != nil {
			t.Error(err)
		}

		err = os.Remove(fText.Name())
		if err != nil {
			t.Error(err)
		}
	}()

	err = ioutil.WriteFile(fText.Name(), []byte("this file contains foo"), 0666)
	if err != nil {
		t.Error(err)
	}

	fiText, err := fText.Stat()
	if err != nil {
		t.Error(err)
	}

	fEmpty, err := ioutil.TempFile("./", "test_")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := fEmpty.Close()
		if err != nil {
			t.Error(err)
		}

		err = os.Remove(fEmpty.Name())
		if err != nil {
			t.Error(err)
		}
	}()

	fiEmpty, err := fEmpty.Stat()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty file is flagged",
			args: args{info: fiEmpty},
			want: false,
		},
		{
			name: "Not empty file is successful",
			args: args{info: fiText},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &MustNotBeEmpty{}
			if got := mu.Evaluate(tt.args.info); got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNotContain_Evaluate(t *testing.T) {
	f, err := ioutil.TempFile("./", "test_")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}

		err = os.Remove(f.Name())
		if err != nil {
			t.Error(err)
		}
	}()

	err = ioutil.WriteFile(f.Name(), []byte("this file contains foo"), 0666)
	if err != nil {
		t.Error(err)
	}

	fi, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Value string
	}
	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Must not contain 'foo'",
			fields: fields{Value: "foo"},
			args:   args{info: fi},
			want:   false,
		},
		{
			name:   "Must not contain 'bar'",
			fields: fields{Value: "bar"},
			args:   args{info: fi},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnc := &MustNotContain{
				Value: tt.fields.Value,
			}
			if got := mnc.Evaluate(tt.args.info); got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
