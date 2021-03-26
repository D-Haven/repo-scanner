package main

import "testing"

func TestContains_Evaluate(t *testing.T) {
	type fields struct {
		Value string
	}
	type args struct {
		filename string
		b        []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{
			name:   "Contains 'foo' is true",
			fields: fields{Value: "foo"},
			args: args{
				filename: "testfile.bat",
				b:        []byte("This string is foo-rific"),
			},
			want:  true,
			want1: "",
		},
		{
			name:   "Contains 'bar' is false",
			fields: fields{Value: "bar"},
			args: args{
				filename: "testfile.bat",
				b:        []byte("This string is foo-rific"),
			},
			want:  false,
			want1: "(testfile.bat) does not contain text: bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Contains{
				Value: tt.fields.Value,
			}
			got, got1 := c.Evaluate(tt.args.filename, tt.args.b)
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMultiConstraint_Evaluate(t *testing.T) {
	type fields struct {
		Constraints []Constraint
	}
	type args struct {
		filename string
		b        []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{
			name: "Good path: Multi-constraint must not be empty, must contain 'foo', must not contain 'bar'",
			fields: fields{
				Constraints: []Constraint{
					&MustNotBeEmpty{},
					&Contains{Value: "foo"},
					&MustNotContain{Value: "bar"},
				},
			},
			args: args{
				filename: "Jenkinsfile",
				b:        []byte("{ build foo }"),
			},
			want:  true,
			want1: "",
		},
		{
			name: "empty file: Multi-constraint must not be empty, must contain 'foo', must not contain 'bar'",
			fields: fields{
				Constraints: []Constraint{
					&MustNotBeEmpty{},
					&Contains{Value: "foo"},
					&MustNotContain{Value: "bar"},
				},
			},
			args: args{
				filename: "Jenkinsfile",
				b:        []byte{},
			},
			want:  false,
			want1: "empty: Jenkinsfile",
		},
		{
			name: "has 'bar': Multi-constraint must not be empty, must contain 'foo', must not contain 'bar'",
			fields: fields{
				Constraints: []Constraint{
					&MustNotBeEmpty{},
					&Contains{Value: "foo"},
					&MustNotContain{Value: "bar"},
				},
			},
			args: args{
				filename: "Jenkinsfile",
				b:        []byte("{ build foo bar }"),
			},
			want:  false,
			want1: "(Jenkinsfile) contains text: bar",
		},
		{
			name: "missing foo: Multi-constraint must not be empty, must contain 'foo', must not contain 'bar'",
			fields: fields{
				Constraints: []Constraint{
					&MustNotBeEmpty{},
					&Contains{Value: "foo"},
					&MustNotContain{Value: "bar"},
				},
			},
			args: args{
				filename: "Jenkinsfile",
				b:        []byte("{ build bar }"),
			},
			want:  false,
			want1: "(Jenkinsfile) does not contain text: foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &MultiConstraint{
				Constraints: tt.fields.Constraints,
			}
			got, got1 := mc.Evaluate(tt.args.filename, tt.args.b)
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMustNotBeEmpty_Evaluate(t *testing.T) {
	type args struct {
		filename string
		b        []byte
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			name: "Empty file returns false with error message",
			args: args{
				filename: "another-test.txt",
				b:        []byte{},
			},
			want:  false,
			want1: "empty: another-test.txt",
		},
		{
			name: "File with content returns true",
			args: args{
				filename: "another-test.txt",
				b:        []byte{13, 122, 55, 3, 0},
			},
			want:  true,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &MustNotBeEmpty{}
			got, got1 := mu.Evaluate(tt.args.filename, tt.args.b)
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMustNotContain_Evaluate(t *testing.T) {
	type fields struct {
		Value string
	}
	type args struct {
		filename string
		b        []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{
			name:   "File containing 'foo' will return false",
			fields: fields{Value: "foo"},
			args: args{
				filename: "example.go",
				b:        []byte("This file contains foo in it."),
			},
			want:  false,
			want1: "(example.go) contains text: foo",
		},
		{
			name:   "File containing 'bar' will return false",
			fields: fields{Value: "bar"},
			args: args{
				filename: "example.go",
				b:        []byte("This file contains foo in it."),
			},
			want:  true,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnc := &MustNotContain{
				Value: tt.fields.Value,
			}
			got, got1 := mnc.Evaluate(tt.args.filename, tt.args.b)
			if got != tt.want {
				t.Errorf("Evaluate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
