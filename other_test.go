package spvalidator

import (
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

type methodValidator struct {
	ok bool
}

func (m methodValidator) Validate() error {
	if m.ok {
		return nil
	}
	return errors.New("validation failed")
}

func (m methodValidator) Check() bool {
	return m.ok
}

func (m methodValidator) NoReturn() {}

func (m methodValidator) Panic() error {
	panic("boom")
}

func (m methodValidator) WithArg(int) error {
	return nil
}

func TestFilesystemValidators(t *testing.T) {
	dir := t.TempDir()
	file := writeTempFile(t, dir, "sample.txt", []byte("hello"))

	expectNoError(t, Dir(dir))
	validateErr(t, Dir(filepath.Join(dir, "missing")), "dir")
	expectNoError(t, File(file))
	validateErr(t, File(dir), "file")
	expectNoError(t, DirPath(dir))
	validateErr(t, DirPath("."), "dirpath")
	expectNoError(t, FilePath(file))
	validateErr(t, FilePath("."), "filepath")

	imgPath := filepath.Join(dir, "sample.png")
	imgFile, err := os.Create(imgPath)
	if err != nil {
		t.Fatalf("create image: %v", err)
	}
	if err := png.Encode(imgFile, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		imgFile.Close()
		t.Fatalf("encode image: %v", err)
	}
	if err := imgFile.Close(); err != nil {
		t.Fatalf("close image: %v", err)
	}
	expectNoError(t, Image(imgPath))
	validateErr(t, Image(filepath.Join(dir, "missing.png")), "image")
}

func TestCollectionValidators(t *testing.T) {
	expectNoError(t, MIMEType("text/plain; charset=utf-8"))
	validateErr(t, MIMEType("not a mime type"), "mimetype")

	expectNoError(t, IsDefault(0))
	validateErr(t, IsDefault(1), "isdefault")

	expectNoError(t, Len("é", 1))
	expectNoError(t, Len([]int{1, 2}, 2))
	expectNoError(t, Len(3, 3))
	validateErr(t, Len("abc", 2), "len")

	expectNoError(t, Max("abc", 3))
	validateErr(t, Max("abc", 2), "max")
	expectNoError(t, Min(3, 2))
	validateErr(t, Min(2, 3), "min")

	expectNoError(t, OneOf("b", "a", "b"))
	validateErr(t, OneOf("c", "a", "b"), "oneof")
	expectNoError(t, NoneOf("c", "a", "b"))
	validateErr(t, NoneOf("a", "a", "b"), "noneof")

	expectNoError(t, Unique("abc"))
	validateErr(t, Unique("abca"), "unique")
	expectNoError(t, Unique([]int{1, 2, 3}))
	validateErr(t, Unique([]int{1, 2, 1}), "unique")
	expectNoError(t, Unique(map[string]int{"a": 1, "b": 1}))
}

func TestValidateFn(t *testing.T) {
	expectNoError(t, ValidateFn(methodValidator{ok: true}))
	validateErr(t, ValidateFn(methodValidator{ok: false}), "validateFn")

	expectNoError(t, ValidateFn(methodValidator{ok: true}, "Check"))
	validateErr(t, ValidateFn(methodValidator{ok: false}, "Check"), "validateFn")

	expectNoError(t, ValidateFn(methodValidator{ok: true}, "NoReturn"))
	validateErr(t, ValidateFn(methodValidator{ok: true}, "Panic"), "validateFn")
	validateErr(t, ValidateFn(methodValidator{ok: true}, "WithArg"), "validateFn")
	validateErr(t, ValidateFn(methodValidator{ok: true}, "Missing"), "validateFn")
}
