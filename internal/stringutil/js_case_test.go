package stringutil

import "testing"

func TestJSCasing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "ascii lowercase", got: ToLowerJS("HELLO"), want: "hello"},
		{name: "ascii uppercase", got: ToUpperJS("hello"), want: "HELLO"},
		{name: "lowercase dotted i", got: ToLowerJS("İSPANYOL"), want: "i̇spanyol"},
		{name: "lowercase final sigma", got: ToLowerJS("ΟΣ"), want: "ος"},
		{name: "uppercase sharp s", got: ToUpperJS("ßfoo"), want: "SSFOO"},
		{name: "uppercase ligature", got: ToUpperJS("ﬁoo"), want: "FIOO"},
		{name: "capitalize-style uppercase", got: ToUpperJS("ß") + "foo", want: "SSfoo"},
		{name: "uncapitalize-style lowercase", got: ToLowerJS("İ") + "foo", want: "i̇foo"},
		{name: "lowercase final sigma after lowercase letter without uppercase mapping", got: ToLowerJS("ʕΣ"), want: "ʕς"},
		{name: "lowercase sigma after modifier letter", got: ToLowerJS("ʰΣ"), want: "ʰσ"},
		{name: "lowercase sigma after case ignorable ypogegrammeni", got: ToLowerJS("ͅΣ"), want: "ͅσ"},
		{name: "lowercase sigma after uncased extended latin letter", got: ToLowerJS("ꟋΣ"), want: "Ɤσ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}
