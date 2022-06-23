package middleware

import "testing"

func TestCreateJWT(t *testing.T) {

	got := CreateJWT("user", "123456", []byte("abc"))
	want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoidXNlciIsInJvbGxubyI6IjEyMzQ1NiJ9.Jw307EsDWxVFBhNpzhazEUOLTCGLs19njyvRdopnt70"

	if want != got {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestValidateJWT(t *testing.T) {

	want_rollno := "123456"

	got_rollno, got_ok := validateJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoidXNlciIsInJvbGxubyI6IjEyMzQ1NiJ9.Jw307EsDWxVFBhNpzhazEUOLTCGLs19njyvRdopnt70", []byte("abc"))

	if want_rollno != got_rollno && !got_ok {
		t.Errorf("Expected %s, got %s", want_rollno, got_rollno)
	}

}
