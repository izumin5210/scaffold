package fs

import (
	"os"

	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/pkg/errors"
)

func (f *fs) GetEntries(root string, compact bool) ([]Entry, error) {
	dir, err := f.afs.DirExists(root)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read %q", root)
	}
	if !dir {
		return nil, errors.Errorf("GetEntries() requires a directory, but %q is a file", root)
	}
	if compact {
		return f.getEntriesCompact(root)
	}
	return f.getEntries(root)
}

func (f *fs) getEntries(root string) ([]Entry, error) {
	tree, err := f.getTree(root)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build tree %q", root)
	}

	entries := []Entry{}

	tree.Root().Walk(func(k []byte, v interface{}) bool {
		entries = append(entries, v.(Entry))
		return false
	})

	return entries, nil
}

func (f *fs) getEntriesCompact(root string) ([]Entry, error) {
	tree, err := f.getTree(root)

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build tree %q", root)
	}

	entries := []Entry{}

	tree.Root().Walk(func(ki []byte, v interface{}) bool {
		if v.(Entry).IsDir() {
			skip := false

			tree.Root().WalkPrefix(ki, func(kj []byte, _ interface{}) bool {
				if string(ki) != string(kj) {
					skip = true
				}
				return skip
			})

			if skip {
				return false
			}
		}
		entries = append(entries, v.(Entry))
		return false
	})

	return entries, nil
}

func (f *fs) getTree(root string) (*iradix.Tree, error) {
	txn := iradix.New().Txn()

	err := f.afs.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if root == path {
			return nil
		}

		e, err := NewEntry(path, info.IsDir())
		if err != nil {
			return err
		}

		txn.Insert([]byte(path), e)
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to walk %q", root)
	}

	return txn.Commit(), nil
}
